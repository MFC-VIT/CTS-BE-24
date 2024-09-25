package users

import (
	"C2S/internal/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Store) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (s *Store) GetUserByUserName(UserName string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(context.TODO(), bson.M{"username": UserName}).Decode(&user)
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

    cursor, err := r.collection.Aggregate(c.Context(), pipeline)
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
	cursor, err := s.collection.Aggregate(c.Context(), pipeline)
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
