package questions

import "go.mongodb.org/mongo-driver/mongo"

type QuestionControllerStore struct {
	db *mongo.Database
}

func NewQuestionStore(db *mongo.Database) *QuestionControllerStore {
	return &QuestionControllerStore{db: db}
}