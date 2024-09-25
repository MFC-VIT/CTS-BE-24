package utils

import (
	"fmt"
	"log"
	"os"
	//"strings"

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
	tokenAuth := c.Get("Authorization")
	tokenQuery := c.Query("token")

	if tokenAuth != "" {
		return tokenAuth
	}
	if tokenQuery != "" {
		return tokenQuery
	}
	return ""
}

type Question struct {
	Question   string `yaml:"question"`   
	QuestionId int    `yaml:"question_id"`
	Answer     string `yaml:"answer"`     
	Room       string `yaml:"room"`       
	Answered   bool   `yaml:"answered"`    
}

type AnswerData struct {
	Questions []Question `yaml:"questions"`
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

	for _, q := range data.Questions {
		log.Printf("Loaded answer from file - Question ID: %d, Question: %v, Answer: %s", q.QuestionId,q.Question, q.Answer)
	}

	return data,nil
}