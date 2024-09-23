package routes

import (
	"C2S/internal/controllers/question"

	"github.com/gofiber/fiber/v2"
)

type QuestionRoutes struct {
	handler *question.Handler
}

func NewQuestionRoutes(handler *question.Handler) *QuestionRoutes{
	return &QuestionRoutes{handler:handler}
}

func (q *QuestionRoutes) RegisterRoutes(router fiber.Router) {
	router.Get("/question/:userID/getQuestions",q.handler.HandleGetQuestion)
	router.Post("/question/:userID/postAnswer",q.handler.HandlePostAnswer)	
}