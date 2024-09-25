package types

import (
	"C2S/internal/models"
	"context"

	"github.com/gofiber/fiber/v2"
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

type AnswerPayload struct {
	QuestionId int `json:"questionId" validate:"required"`
	Question string `json:"question" validate:"required"`
	Room     string `json:"room" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
}

type UserStore interface {
	GetUserByUserName(userName string) (*models.User, error)
	GetUserByID(id primitive.ObjectID) (*models.User, error)
	CreateUser(user *models.User) error
	SeedQuestionsForUser(ctx context.Context,userID primitive.ObjectID) error
	GetAllUsers(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx,) error
	UpdateUser(c *fiber.Ctx) error
	UpdateScore(c *fiber.Ctx) error
	GetUserByUserNameHandler(c* fiber.Ctx) error

}
type RoomStore interface {
	EnterRoom(ctx context.Context, userID primitive.ObjectID, roomEntered string) error
	EscapeRoom(ctx context.Context, userID primitive.ObjectID, roomEntered string) error
}

type QuestionStore interface {
	GetNextQuestion(ctx context.Context, userID primitive.ObjectID) (models.Question, error)
	QuestionAnswered(ctx context.Context, userID primitive.ObjectID, question models.Question) error
}