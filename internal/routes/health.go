package routes

import (
	"C2S/internal/controllers/health"

	"github.com/gofiber/fiber/v2"
)

type HealthRoutes struct{}

func NewHealthRoutes() *HealthRoutes {
	return &HealthRoutes{}
}

func (h *HealthRoutes) RegisterRoutes(router fiber.Router) {
	router.Get("/health", health.HandleHealth) 
}

