package users

import (
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
    ID                  primitive.ObjectID 
    Location           string
	usersCollection *mongo.Collection
    roomsCollection *mongo.Collection
	questionsCollection *mongo.Collection
    db         *mongo.Database 
}

func NewUserStore(db *mongo.Database) *Store {
    return &Store{
        db:db,
        usersCollection: db.Collection(os.Getenv("MONGO_USER_COLLECTION")),
        roomsCollection: db.Collection(os.Getenv("MONGO_ROOMS_COLLECTION")),
        questionsCollection: db.Collection(os.Getenv("MONGO_QUESTIONS_COLLECTION")),
    }
}
