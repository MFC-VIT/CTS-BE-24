package routes

import (
	"C2S/internal/controllers/rooms"

	"github.com/gofiber/fiber/v2"
)

type RoomRoutes struct {
	handler *rooms.Handler
}

func NewRoomRoutes(handler *rooms.Handler) *RoomRoutes {
	return &RoomRoutes{handler: handler}
}

func (r *RoomRoutes) RegisterRoutes(router fiber.Router) {
	router.Post("/room/:userID/enter", r.handler.HandleEnterRoom)
	router.Post("/room/:userID/escape", r.handler.HandleEscapeRoom)
}
