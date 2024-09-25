package users

import (
	"C2S/internal/models"
	seed "C2S/internal/seeders"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) CreateUser(user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	_, err := s.usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SeedQuestionsForUser(ctx context.Context,userID primitive.ObjectID) error {
	yamlFile := "internal/files/questions.yaml" 
	return seed.SeedQuestions(ctx, s.db, userID, yamlFile) 
}

