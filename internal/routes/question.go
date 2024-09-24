package routes

import (
	"C2S/internal/controllers/question"
	"C2S/internal/middleware"
	"C2S/internal/types"

	"github.com/gofiber/fiber/v2"
)

type QuestionRoutes struct {
	handler *question.Handler
	userStore types.UserStore
}

func NewQuestionRoutes(handler *question.Handler,userStore types.UserStore) *QuestionRoutes{
	return &QuestionRoutes{handler:handler,userStore: userStore}
}

func (q *QuestionRoutes) RegisterRoutes(router fiber.Router) {
	router.Get("/question/:userID/getQuestions",middleware.WithJWTAuth(q.userStore),  q.handler.HandleGetQuestion)
	router.Post("/question/:userID/postAnswer",middleware.WithJWTAuth(q.userStore),  q.handler.HandlePostAnswer)	
}