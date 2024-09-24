package routes

import (
	"C2S/internal/controllers/rooms"
	"C2S/internal/middleware"
	"C2S/internal/types"

	"github.com/gofiber/fiber/v2"
)

type RoomRoutes struct {
	handler *rooms.Handler
	userStore types.UserStore
}

func NewRoomRoutes(handler *rooms.Handler,userStore types.UserStore) *RoomRoutes {
	return &RoomRoutes{handler: handler,userStore: userStore}
}

func (r *RoomRoutes) RegisterRoutes(router fiber.Router) {
	router.Post("/room/:userID/enter",middleware.WithJWTAuth(r.userStore),  r.handler.HandleEnterRoom)
	router.Post("/room/:userID/escape", middleware.WithJWTAuth(r.userStore), r.handler.HandleEscapeRoom)
}