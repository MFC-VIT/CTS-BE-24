package health

import "github.com/gofiber/fiber/v2"

func HandleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "up",
		"message": "Service is healthy",
	})
}
