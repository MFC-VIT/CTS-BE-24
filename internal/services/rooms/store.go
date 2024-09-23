package rooms

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomControllerStore struct {
	db *mongo.Database
}


func NewRoomStore(db *mongo.Database) *RoomControllerStore {
	return &RoomControllerStore{db: db}
}