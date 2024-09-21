package seeders

import (
	"C2S/internal/models"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v2"
)

type QuestionYAML struct {
	Question string `yaml:"question"`
	Answer   string `yaml:"answer"`
	Room     string `yaml:"room"`
	Answered string   `yaml:"answered"`
}

type QuestionsYAML struct {
	Questions []QuestionYAML `yaml:"questions"`
}


func getUsernameByID(ctx context.Context, db *mongo.Database, userID primitive.ObjectID) (string, error) {
	usersCollection := db.Collection(os.Getenv("MONGO_USER_COLLECTION")) 
	var user models.User 
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %v", err)
	}
	return user.UserName, nil 
}

func SeedQuestions(ctx context.Context, db *mongo.Database, userID primitive.ObjectID, yamlFile string) error {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %v", err)
	}

	var questionsYAML QuestionsYAML
	if err := yaml.Unmarshal(data, &questionsYAML); err != nil {
		return fmt.Errorf("failed to unmarshal YAML data: %v", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(questionsYAML.Questions), func(i, j int) {
		questionsYAML.Questions[i], questionsYAML.Questions[j] = questionsYAML.Questions[j], questionsYAML.Questions[i]
	})

	var roomAQuestions, roomBQuestions, roomCQuestions, roomDQuestions []models.Question

	for _, q := range questionsYAML.Questions {
		question := models.Question{
			Question: q.Question,
			Answer:   q.Answer,
			Room:     q.Room,
			Answered: q.Answered,
		}

		switch q.Room {
		case "A":
			if len(roomAQuestions) < 4 {
				roomAQuestions = append(roomAQuestions, question)
			}
		case "B":
			if len(roomBQuestions) < 4 {
				roomBQuestions = append(roomBQuestions, question)
			}
		case "C":
			if len(roomCQuestions) < 4 {
				roomCQuestions = append(roomCQuestions, question)
			}
		case "D":
			if len(roomDQuestions) < 4 {
				roomDQuestions = append(roomDQuestions, question)
			}
		}
	}

	roomA := models.Room{Questions: roomAQuestions}
	roomB := models.Room{Questions: roomBQuestions}
	roomC := models.Room{Questions: roomCQuestions}
	roomD := models.Room{Questions: roomDQuestions}

	userName, err := getUsernameByID(ctx, db, userID)
	if err != nil {
		return fmt.Errorf("failed to get username: %v", err)
	}

	collectionName := os.Getenv("MONGO_QUESTIONS_COLLECTION")
	fmt.Printf("Collection: %s\n", collectionName)

	questionCollection := db.Collection(collectionName)

	questions := models.Questions{
		ID:       primitive.NewObjectID(),
		UserID:   userID,
		UserName: userName, 
		RoomA:    roomA,
		RoomB:    roomB,
		RoomC:    roomC,
		RoomD:    roomD,
	}

	_, err = questionCollection.InsertOne(ctx, questions)
	if err != nil {
		return fmt.Errorf("failed to insert seeded questions: %v", err)
	}

	fmt.Println("Seeded questions successfully!")
	return nil
}
