package question

import (
	"C2S/internal/models"
	"C2S/internal/types"
	"C2S/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	store types.QuestionStore
}

func NewHandler(store types.QuestionStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) HandleGetQuestion(c *fiber.Ctx) error {
	userIDParam := c.Params("userID")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("invalid userID format"))
	}
	fmt.Printf("Received request to get next question for user ID: %s\n", userID.Hex())
	question, err := h.store.GetNextQuestion(c.Context(), userID)
	if err != nil {
		fmt.Printf("Error fetching next question: %v\n", err)
		if err.Error() == "user is not in any room" {
			return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("user is not in any room"))
		}
		return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("failed to fetch the next question"))
	}
	fmt.Printf("Successfully fetched question: %s\n", question.Question)
	return utils.WriteJSON(c, fiber.StatusOK, question)
}

func (h *Handler) HandlePostAnswer(c *fiber.Ctx) error {
		var payload types.AnswerPayload
		userIDParam := c.Params("userID")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("invalid userID format"))
		}
		if err := c.BodyParser(&payload); err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("invalid request format"))
		}
	
		fmt.Printf("Received request to submit answer for user ID: %s, Question: %s\n", userID.Hex(), payload.Question)
		questionData := models.Question{
			Room:     payload.Room,
			Question: payload.Question,
			Answer:   payload.Answer,
		}
	
		err = h.store.QuestionAnswered(c.Context(), userID, questionData)
		if err != nil {
			fmt.Printf("Error submitting answer: %v\n", err)
			if err.Error() == "incorrect answer" {
				return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("incorrect answer"))
			}
			if err.Error() == "question already answered or not found" {
				return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("question already answered or not found"))
			}
			return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("failed to submit answer"))
		}
	
		fmt.Printf("Successfully submitted answer for user ID: %s\n", userID.Hex())
		return utils.WriteJSON(c, fiber.StatusOK, fiber.Map{"message": "Answer submitted successfully"})
	}