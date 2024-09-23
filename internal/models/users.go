package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserName string             `bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	IsAdmin   bool               `bson:"is_admin"`
	IsQuestionSeeded bool 		 `bson:"is_seeded"`
	RoomEntered     string      `bson:"room_entered"`
	Score            int          `bson:"score"`
}