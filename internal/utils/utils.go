package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

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

type LocationData struct {
	Locations []Location `yaml:"locations"`
}

type Location struct {
	Location string `yaml:"location"`
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

func LoadLocations(filePath string) ([]string, error) {
	var locData LocationData

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &locData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML file: %v", err)
	}

	var locations []string
	for _, loc := range locData.Locations {
		locations = append(locations, loc.Location)
	}

	return locations, nil
}



func GetRandomLocation(locations []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(locations), func(i, j int) {
		locations[i], locations[j] = locations[j], locations[i]
	})
	if len(locations) > 0 {
		return locations[0]
	}
	return "" 
}