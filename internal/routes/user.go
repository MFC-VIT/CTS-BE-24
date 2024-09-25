package routes

import (
	"C2S/internal/controllers/user"
	"C2S/internal/middleware"
	"C2S/internal/types"

	"github.com/gofiber/fiber/v2"
)

type UserRoutes struct {
	handler *user.Handler
	userStore types.UserStore
}

func NewUserRoutes(handler *user.Handler,userStore types.UserStore) *UserRoutes {
	return &UserRoutes{handler: handler,userStore: userStore}
}

func (u *UserRoutes) RegisterRoutes(router fiber.Router) {

	router.Post("/user/login", u.handler.HandleLogin)
	router.Post("/user/register", u.handler.HandleRegister)
	router.Get("/user/:userID/location", middleware.WithJWTAuth(u.userStore),u.handler.HandleGetRandomLocation)
}
