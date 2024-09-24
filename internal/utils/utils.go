package utils

import (
	"C2S/internal/models"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v2"
)

var Validate = validator.New()

func ParseJSON(r *fiber.Ctx,payload any)error{
	if len(r.Body())==0{
		return fmt.Errorf("missing request body")
	}
	return r.BodyParser(payload)
}

func WriteJSON(r *fiber.Ctx, status int, v any) error {
	return r.Status(status).JSON(v)
}


func WriteError(r *fiber.Ctx, status int, err error) error {
	// WriteJSON(r, status, map[string]string{"error": err.Error()})
	return r.Status(status).JSON(map[string]string{"error": err.Error()})
}

func GetTokenFromRequest(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	log.Printf("Authorization Header: %s", authHeader)

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("User is not authorized or token missing")
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		log.Println("Authorization header format must be Bearer {token}")
		return ""
	}

	log.Printf("Token Part: %s", parts[1])
	return parts[1]
}


type AnswerData struct {
	Questions []models.Question `yaml:"questions"`
}

func LoadAnswers(filePath string) (AnswerData,error){
	var data AnswerData

	file, err := os.ReadFile(filePath)

	if err !=nil {
		return data, fmt.Errorf("error reading file; %v",err)
	}

	if err := yaml.Unmarshal(file,&data);err!=nil{
		return data, fmt.Errorf("error unmarshalling yaml: %v",err)
	}
	return data,nil
}