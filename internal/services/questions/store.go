package questions

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionControllerStore struct {
	db *mongo.Database
	usersCollection *mongo.Collection
    roomsCollection *mongo.Collection
	questionsCollection *mongo.Collection
}

func NewQuestionStore(db *mongo.Database) *QuestionControllerStore {
	return &QuestionControllerStore{
		db: db,
		usersCollection: db.Collection(os.Getenv("MONGO_USER_COLLECTION")),
        roomsCollection: db.Collection(os.Getenv("MONGO_ROOMS_COLLECTION")),
		questionsCollection: db.Collection(os.Getenv("MONGO_QUESTIONS_COLLECTION")),
	}
}