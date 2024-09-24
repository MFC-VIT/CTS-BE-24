package rooms

import (
	"C2S/internal/middleware"
	"C2S/internal/types"
	"C2S/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	store types.UserStore
	roomcontrollerstore types.RoomStore
}

func NewHandler(store types.UserStore, roomcontrollerstore types.RoomStore) *Handler {
	return &Handler{store: store, roomcontrollerstore: roomcontrollerstore}
}

func (h *Handler) HandleEnterRoom(c *fiber.Ctx) error {
	userIDParam := c.Params("userID")

	userIDFromToken := c.Locals(middleware.UserKey).(string)

	if userIDParam != userIDFromToken {
		return utils.WriteError(c, fiber.StatusForbidden, fmt.Errorf("permission denied: user ID mismatch"))
	}

	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("invalid userid format"))
	}

	var payload struct {
		RoomEntered string `bson:"room_entered"`
	}

		fmt.Printf("input: %s", payload.RoomEntered)

	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest,fmt.Errorf("invalid payload"))
	}

	err = h.roomcontrollerstore.EnterRoom(c.Context(), userID, payload.RoomEntered)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("failed to enter the room"))
	}

	return utils.WriteJSON(c, fiber.StatusOK, "Successfully entered the room")
}

func (h *Handler) HandleEscapeRoom(c *fiber.Ctx) error {
	userIDParam := c.Params("userID")
	
	userIDFromToken := c.Locals(middleware.UserKey).(string)

	if userIDParam != userIDFromToken {
		return utils.WriteError(c, fiber.StatusForbidden, fmt.Errorf("permission denied: user ID mismatch"))
	}
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest,  fmt.Errorf("invalid userid format"))
	}
	var payload struct {
		RoomEntered string `bson:"room_entered"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest,  fmt.Errorf("invalid payload"))
	}

	err = h.roomcontrollerstore.EscapeRoom(c.Context(), userID, payload.RoomEntered)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError,  fmt.Errorf("failed to escape the room"))
	}

	return utils.WriteJSON(c, fiber.StatusOK, "Successfully escaped the room")
}
