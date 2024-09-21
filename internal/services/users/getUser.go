package users

import (
	"C2S/internal/models"
	"context"

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