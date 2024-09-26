package middleware

import (
	"C2S/internal/types"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IsAdmin middleware checks if the user is an admin based on the JWT-stored userID
func IsAdmin(store types.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals(UserKey).(string) // Extract userID from the JWT stored in Locals

		// Check if userID is empty
		if userID == "" {
			return permissionDenied(c)
		}

		log.Printf("Extracted userID: %s", userID)

		// Convert userID (Hex string) to ObjectID
		usid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Printf("Invalid ObjectID format: %v", err)
			return permissionDenied(c)
		}

		// Fetch the user from the database using the ObjectID
		user, err := store.GetUserByID(usid)
		if err != nil {
			log.Println("Cannot find user")
			return permissionDenied(c)
		}

		// Check if the user has admin privileges
		if !user.IsAdmin {
			log.Println("User is not an admin")
			return permissionDenied(c)
		}

		// Allow the request to proceed
		return c.Next()
	}
}
