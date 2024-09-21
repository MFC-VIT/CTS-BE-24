package routes

import (
	"C2S/internal/controllers/user"

	"github.com/gofiber/fiber/v2"
)

type UserRoutes struct {
	handler *user.Handler
}

func NewUserRoutes(handler *user.Handler) *UserRoutes {
	return &UserRoutes{handler: handler}
}

func (u *UserRoutes) RegisterRoutes(router fiber.Router) {

	router.Post("/user/login", u.handler.HandleLogin)
	router.Post("/user/register", u.handler.HandleRegister)
}
