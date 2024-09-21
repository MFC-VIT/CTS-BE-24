package types

import (
	"C2S/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUserPayload struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserStore interface {
	GetUserByUserName(userName string) (*models.User, error)
	GetUserByID(id primitive.ObjectID) (*models.User, error)
	CreateUser(user *models.User) error
	SeedQuestionsForUser(ctx context.Context,userID primitive.ObjectID) error
}