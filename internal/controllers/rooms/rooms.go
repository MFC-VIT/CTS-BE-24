package rooms

import (
	"C2S/internal/middleware"
	"C2S/internal/types"
	"C2S/internal/utils"
	"fmt"
	"log"
	"regexp"

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
var objectIdRegex = regexp.MustCompile(`^ObjectID\("([0-9a-fA-F]{24})"\)$`)
func (h *Handler) HandleEnterRoom(c *fiber.Ctx) error {
	userIDParam := c.Params("userID")

	userIDFromToken := c.Locals(middleware.UserKey).(string)
	log.Printf("User ID from token: %s", userIDFromToken)
	log.Printf("User ID from params: %s", userIDParam)
	matches := objectIdRegex.FindStringSubmatch(userIDFromToken)
		if len(matches) != 2 {
			log.Println("Invalid ObjectID format:", userIDParam)
			return utils.WriteError(c,fiber.StatusForbidden, fmt.Errorf("invalid token"))
		}
		hexID := matches[1]
		log.Println("Extracted hex string:", hexID)
	if userIDParam != hexID {
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

	//serIDFromToken := c.Locals(middleware.UserKey).(string)
	log.Printf("User ID from token: %s", userIDFromToken)
	log.Printf("User ID from params: %s", userIDParam)
	matches := objectIdRegex.FindStringSubmatch(userIDFromToken)
		if len(matches) != 2 {
			log.Println("Invalid ObjectID format:", userIDParam)
			return utils.WriteError(c,fiber.StatusForbidden, fmt.Errorf("invalid token"))
		}
		hexID := matches[1]
		log.Println("Extracted hex string:", hexID)
	if userIDParam != hexID {
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
	//log.Println(Room)
	err = h.roomcontrollerstore.EscapeRoom(c.Context(), userID, payload.RoomEntered)
	if err != nil {
		log.Println(err)
		return utils.WriteError(c, fiber.StatusInternalServerError,  fmt.Errorf("failed to escape the room"))
	}

	return utils.WriteJSON(c, fiber.StatusOK, "Successfully escaped the room")
}
