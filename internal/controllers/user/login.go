package user

import (
	"C2S/internal/middleware"
	"C2S/internal/services/auth"
	"C2S/internal/types"
	"C2S/internal/utils"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler{
	return &Handler{store:store}
}

func (h *Handler) HandleLogin(c *fiber.Ctx) error{
	var user types.LoginUserPayload
	if err:= utils.ParseJSON(c,&user);
	err !=nil{
		return utils.WriteError(c, fiber.StatusBadRequest,err)
	}

	if err := utils.Validate.Struct(user); 
	err != nil {
		errors := err.(validator.ValidationErrors)
		return utils.WriteError(c,fiber.StatusBadRequest,fmt.Errorf("invalid username or user not found: %v", errors))
	}
	
	u,err:= h.store.GetUserByUserName(user.UserName)
	if err!=nil{
		return utils.WriteError(c,fiber.StatusBadRequest,fmt.Errorf("invalid username or user not found: %v", err))
	}
	if !auth.Comparepasswords(u.Password,[]byte(user.Password)){
		return utils.WriteError(c,fiber.StatusBadRequest,fmt.Errorf("invalid username or password: %v", err))
	}
	secret := []byte(os.Getenv("JWTSecret"))
	userIDString := u.ID.Hex()
	token, err := middleware.CreateJWT(secret, userIDString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create token",
		})
	}
	response := struct {
		UserID primitive.ObjectID `json:"id"`
		UserName  string             `json:"username"`
		Token  string             `json:"token"`
	}{
		UserID: u.ID,
		UserName:  u.UserName,
		Token:  token,
	}

	return c.Status(fiber.StatusOK).JSON(response)

}

func (h *Handler) HandleGetRandomLocation(c *fiber.Ctx) error {
	userIDParam := c.Params("userID") 
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("invalid user ID: %v", err))
	}

	locationsFilePath := "internal/files/location.yaml" 
	
	randomLocation, err := h.store.GetRandomLocation(c.Context(), userID, locationsFilePath)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("failed to get random location: %v", err))
	}

	if randomLocation == "" {
		return utils.WriteError(c, fiber.StatusOK, fmt.Errorf("solve all rooms to get location"))
	}

	response := struct {
		Location string `json:"location"`
	}{
		Location: randomLocation,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) HandleGetRoomStatus(c *fiber.Ctx) error {
	userIDParam := c.Params("userID")
	userIDFromToken := c.Locals(middleware.UserKey).(string)

	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("invalid user ID: %v", err))
	}
	if userIDParam != userIDFromToken {
		return utils.WriteError(c, fiber.StatusForbidden, fmt.Errorf("permission denied: user ID mismatch"))
	}
	roomStatus, err := h.store.GetUserRoomStatus(c.Context(), userID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("failed to get room status: %v", err))
	}

	return c.Status(fiber.StatusOK).JSON(roomStatus)
}
