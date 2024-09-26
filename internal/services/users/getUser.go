package users

import (
	"C2S/internal/models"
	"C2S/internal/utils"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Store) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.usersCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (s *Store) GetUserByUserName(UserName string) (*models.User, error) {
	var user models.User
	err := s.usersCollection.FindOne(context.TODO(), bson.M{"username": UserName}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Store) GetUserByUserNameHandler(c *fiber.Ctx) error {

    userName := c.Params("username")


    pipeline := mongo.Pipeline{

        bson.D{
            {Key: "$match", Value: bson.M{"username": userName}},
        },

        bson.D{
            {Key: "$lookup",
                Value: bson.D{
                    {Key: "from", Value: "rooms"},
                    {Key: "localField", Value: "_id"},
                    {Key: "foreignField", Value: "user_id"},
                    {Key: "as", Value: "room_details"},
                },
            },
        },

        bson.D{
            {Key: "$lookup",
                Value: bson.D{
                    {Key: "from", Value: "questions"},
                    {Key: "localField", Value: "_id"},
                    {Key: "foreignField", Value: "user_id"},
                    {Key: "as", Value: "question_details"},
                },
            },
        },

        bson.D{
            {Key: "$project",
                Value: bson.D{
                    {Key: "room_details._id", Value: 0},
                    {Key: "room_details.username", Value: 0},
                    {Key: "room_details.user_id", Value: 0},
                    {Key: "question_details.user_id", Value: 0},
                    {Key: "question_details._id", Value: 0},
                    {Key: "question_details.username", Value: 0},
                    {Key: "password", Value: 0},
                },
            },
        },
    }

    cursor, err := r.usersCollection.Aggregate(c.Context(), pipeline)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch user data")
    }
    defer cursor.Close(context.TODO())

    var users []bson.M
    if err := cursor.All(c.Context(), &users); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user data")
    }

    if len(users) == 0 {
        return c.Status(400).SendString("User Not Found")
    }

    return c.Status(fiber.StatusOK).JSON(users[0])
}


func (s * Store) GetAllUsers(c* fiber.Ctx) error{
	pipeline := mongo.Pipeline{
		bson.D{
			{Key:"$lookup",
				Value:bson.D{
					{Key:"from", Value:"rooms"},
					{Key:"localField", Value:"_id"},
					{Key:"foreignField", Value:"user_id"},
					{Key:"as", Value:"room_details"},
				},
			},
		},

		bson.D{
			{Key:"$lookup",
			Value:bson.D{
					{Key:"from", Value:"questions"},
					{Key:"localField", Value:"_id"},
					{Key:"foreignField", Value:"user_id"},
					{Key:"as", Value:"question_details"},
				},
			},
		},

		bson.D{
			{Key:"$project",
			Value:bson.D{
					{Key:"room_details._id",Value: 0},
					{Key:"room_details.username",Value: 0},
					{Key:"room_details.user_id", Value:0},
					{Key:"question_details.user_id",Value: 0},
					{Key:"question_details._id", Value:0},
					{Key:"question_details.username", Value:0},
					{Key:"password",Value: 0},
				},
			},
		},
	}
	cursor, err := s.usersCollection.Aggregate(c.Context(), pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get users")
	}
	defer cursor.Close(context.TODO())

	var users []bson.M
	if err := cursor.All(c.Context(), &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse users")
	}

	return c.Status(fiber.StatusOK).JSON(users)	
}

func (s *Store) DeleteUser(c *fiber.Ctx) error {
	userID,err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		log.Println("Error converting user ID to ObjectID:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	query := bson.D{{Key: "_id", Value: userID}}

	result, err := s.db.Collection("users").DeleteOne(c.Context(), query)
	if err != nil {
		log.Println("Error deleting user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Error during deletion"})
	}

	if result.DeletedCount == 0 {
		log.Println("No user found with ID:",userID.Hex())
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	log.Println("User deleted successfully:", userID.Hex())
	return c.Status(200).JSON(fiber.Map{"message": "User deleted successfully"})
}



func (s *Store) UpdateUser(c* fiber.Ctx) error{
	userId,err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		return c.Status(400).SendString("Error")
	}
	user := new(models.User)
	if err := c.BodyParser(user);err !=nil{
		return c.Status(400).SendString(err.Error())
	}
	query:=bson.D{{Key:"_id",Value:userId}}
	update:=bson.D{
		{
			Key:"$set",
			Value:bson.D{
				{Key:"username",Value:user.UserName},
				{Key:"password",Value:user.Password},
			},
		},
	}
	result,err := s.db.Collection("user").UpdateOne(c.Context(),query,update)
	if err!=nil{
		return c.Status(400).SendString("Cannot update the user")
	}

	if result.MatchedCount == 0{
		return c.Status(400).SendString("User Not Found")
	}

	return c.Status(200).SendString("User updated")
}
func (s* Store) UpdateScore(c *fiber.Ctx) error{
	userID,err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		c.Status(400).SendString("Could not retrive userid")
	}
	user := new(models.User)
	if err := c.BodyParser(user);err !=nil{
		return c.Status(400).SendString("Error")
	}
	query :=bson.D{{Key:"_id",Value:userID}}
	update := bson.D{
		{
			Key:"$set",
			Value:bson.D{
				{Key:"score",Value:user.Score},
			},
		},
	}
	result,err := s.db.Collection("user").UpdateOne(c.Context(),query,update)
	if err!=nil{
		c.SendStatus(400)
	}
	if result.MatchedCount == 0{
		return c.Status(400).SendString("User Not Found")
	}

	return c.Status(200).SendString("User updated")
}
func (s *Store) GetRandomLocation(ctx context.Context, userID primitive.ObjectID, locationsFilePath string) (string, error) {

	user, err := s.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %v", err)
	}
	if user.Location != "" {
		return user.Location, nil
	}


	roomStatusArray, err := s.collectUserRoomsStatus(ctx, userID) 
	if err != nil {
		return "", fmt.Errorf("failed to collect user rooms status: %v", err)
	}

	allRoomsDone := true
	for _, status := range roomStatusArray {
		if !strings.HasSuffix(status, "D") {
			allRoomsDone = false
			break
		}
	}

	if allRoomsDone {
		locations, err := utils.LoadLocations(locationsFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to load locations: %v", err)
		}
		randomLocation := utils.GetRandomLocation(locations)
		s.Location = randomLocation
		err = s.UpdateUserLocation(ctx, userID, randomLocation) 
		if err != nil {
			return "", fmt.Errorf("failed to update user location: %v", err)
		}

		return randomLocation, nil
	}

	return "", nil 
}


func (s *Store) UpdateUserLocation(ctx context.Context, userID primitive.ObjectID, location string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"location": location}}

	_, err := s.usersCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) collectUserRoomsStatus(ctx context.Context, userID primitive.ObjectID) ([]string, error) {

	var roomStatus models.Rooms
	err := s.roomsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&roomStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room status: %v", err)
	}
	roomStatusChan := make(chan string, 4)

	var wg sync.WaitGroup

	checkRoomStatus := func(roomName string, done bool, giveUp bool) {
		defer wg.Done()
		var result string
		if done {
			result = fmt.Sprintf("%sD", roomName)
		} else if giveUp {
			result = fmt.Sprintf("%sG", roomName) 
		} else {
			result = fmt.Sprintf("%s-", roomName) 
		}
		roomStatusChan <- result
	}

	wg.Add(4)
	go checkRoomStatus("A", roomStatus.IsRoomsDone.RoomA, roomStatus.IsRoomsGiveUp.RoomA)
	go checkRoomStatus("B", roomStatus.IsRoomsDone.RoomB, roomStatus.IsRoomsGiveUp.RoomB)
	go checkRoomStatus("C", roomStatus.IsRoomsDone.RoomC, roomStatus.IsRoomsGiveUp.RoomC)
	go checkRoomStatus("D", roomStatus.IsRoomsDone.RoomD, roomStatus.IsRoomsGiveUp.RoomD)

	go func() {
		wg.Wait()
		close(roomStatusChan)
	}()

	var roomStatusArray []string
	for status := range roomStatusChan {
		roomStatusArray = append(roomStatusArray, status)
	}
	log.Printf("room status array: %+v", roomStatusArray)
	return roomStatusArray, nil
}


func (s *Store) GetUserRoomStatus(ctx context.Context, userID primitive.ObjectID) (*models.Rooms, error) {
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{{Key: "_id", Value: userID}}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Rooms"}, 
				{Key: "localField", Value: "_id"}, 
				{Key: "foreignField", Value: "user_id"}, 
				{Key: "as", Value: "room_details"}, 
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$room_details"}, {Key: "preserveNullAndEmptyArrays", Value: true}}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "username", Value: 1}, 
				{Key: "is_rooms_done", Value: bson.D{
					{Key: "room_a", Value: "$room_details.is_rooms_done.room_a"},
					{Key: "room_b", Value: "$room_details.is_rooms_done.room_b"},
					{Key: "room_c", Value: "$room_details.is_rooms_done.room_c"},
					{Key: "room_d", Value: "$room_details.is_rooms_done.room_d"},
				}},
				{Key: "is_rooms_giveup", Value: bson.D{
					{Key: "room_a", Value: "$room_details.is_rooms_giveup.room_a"},
					{Key: "room_b", Value: "$room_details.is_rooms_giveup.room_b"},
					{Key: "room_c", Value: "$room_details.is_rooms_giveup.room_c"},
					{Key: "room_d", Value: "$room_details.is_rooms_giveup.room_d"},
				}},
			}},
		},
	}

	cursor, err := s.usersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to perform aggregation: %v", err)
	}
	defer cursor.Close(ctx)

	var result models.Rooms
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode result: %v", err)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &result, nil
}
