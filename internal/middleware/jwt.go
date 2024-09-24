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
	"go.mongodb.org/mongo-driver/bson/primitive"

	"C2S/internal/types"
	"C2S/internal/utils"
)

func CreateJWT(secret []byte, userID string) (string, error) {
	expirationString := os.Getenv("JWTEXPINSEC")
	expirationSeconds, _ := strconv.Atoi(expirationString)
	expiration := time.Second * time.Duration(expirationSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


const UserKey  = "userID"

func WithJWTAuth(store types.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := utils.GetTokenFromRequest(c)
		log.Printf("Token: %s", tokenString)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token(1): %v", err)
			return permissionDenied(c)
		}
		if !token.Valid {
			log.Println("invalid token")
			return permissionDenied(c)
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := primitive.ObjectIDFromHex(str)
		if err != nil {
			log.Printf("failed to convert userID to ObjectID: %v", err)
			return permissionDenied(c)
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get userID: %v", err)
			return permissionDenied(c)
		}

		c.Locals(UserKey, u.ID.Hex())


		return c.Next()
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWTSECRET")), nil
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
