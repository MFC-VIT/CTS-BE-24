package middleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"C2S/internal/types"
	"C2S/internal/utils"
)

func CreateJWT(secret []byte, userID string) (string, error) {
	expirationString := os.Getenv("JWTEXPINSEC")
	expirationSeconds, _ := strconv.Atoi(expirationString)
	expiration := time.Second * time.Duration(expirationSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	log.Printf("Created Token: %s", tokenString) // Log the created token for debugging
	return tokenString, nil
}

const UserKey = "userID"

func WithJWTAuth(store types.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := utils.GetTokenFromRequest(c)
		log.Printf("Token from Request: %s", tokenString)

		token, err := validateJWT(tokenString)
		log.Println(token)
		if err != nil {
			log.Printf("Failed to validate token: %v", err)
			return permissionDenied(c)
		}
		if !token.Valid {
			log.Println("Token is invalid")
			return permissionDenied(c)
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)
		log.Printf("Extracted UserID from Token: %s", userID)

		c.Locals(UserKey, userID)

		return c.Next()
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := []byte(os.Getenv("JWTSECRET"))
	log.Printf("Using JWT Secret for validation: %s", secret) // Log the secret (avoid in production)
	log.Printf("Token %s", tokenString )
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}

func permissionDenied(c *fiber.Ctx) error {
	return utils.WriteError(c, fiber.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserKey).(string)
	if !ok {
		return "nil"
	}
	return userID
}
