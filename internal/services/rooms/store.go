package rooms

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type RoomControllerStore struct {
	db *mongo.Database
	usersCollection *mongo.Collection
    roomsCollection *mongo.Collection
	questionsCollection *mongo.Collection
}


func NewRoomStore(db *mongo.Database) *RoomControllerStore {
	return &RoomControllerStore{
		db: db,
		usersCollection: db.Collection(os.Getenv("MONGO_USER_COLLECTION")),
        roomsCollection: db.Collection(os.Getenv("MONGO_ROOMS_COLLECTION")),
		questionsCollection: db.Collection(os.Getenv("MONGO_QUESTIONS_COLLECTION")),
	}
}