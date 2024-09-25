package middleware

import (
	"C2S/internal/types"
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
var objectIdRegex = regexp.MustCompile(`^ObjectID\("([0-9a-fA-F]{24})"\)$`)

func IsAdmin(store types.UserStore) fiber.Handler{
	return func (c* fiber.Ctx) error{
		userID := c.Locals(UserKey)
		if(userID==""){
			return permissionDenied(c)
		}
		log.Printf("%T",userID)
		log.Println(userID)
		matches := objectIdRegex.FindStringSubmatch(userID.(string))
		if len(matches) != 2 {
			log.Println("Invalid ObjectID format:", userID)
			return permissionDenied(c)
		}
		hexID := matches[1]
		log.Println("Extracted hex string:", hexID)
		usid, err := primitive.ObjectIDFromHex(hexID)
		if err != nil {
			log.Println("Invalid ObjectID format:", err)
			return permissionDenied(c)
		}
		log.Println("UserID:", usid)
		user,err := store.GetUserByID(usid)
		if err!=nil{
			log.Println("Cannot find user")
			return permissionDenied(c)
		}
		if !user.IsAdmin{
			permissionDenied(c)
		}
		return c.Next()
	}
}

