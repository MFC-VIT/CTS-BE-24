package users

import (
	"C2S/internal/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s * Store) GetAllUsers(c* fiber.Ctx) error{
	query := bson.D{{}}
	var  users []  models.User = make([]models.User,0 )
	cursor,err := s.collection.Find(c.Context(),query)
	if err!=nil{
		return err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()){
		var user models.User
		if err:= cursor.Decode(&user);err!=nil{
			return err
		}
		users = append(users,user)
	}
	return c.JSON(users)	
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