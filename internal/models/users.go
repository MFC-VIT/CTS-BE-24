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
	IsRoomADone bool 			 `bson:"is_room_a"`
	IsRoomBDone bool 			 `bson:"is_room_b"`
	IsRoomCDone bool 			 `bson:"is_room_c"`
	IsRoomDDone bool 			 `bson:"is_room_d"`
	IsQuestionSeeded bool 		 `bson:"is_seeded"`
	RoomEntered     map[string]time.Time      `bson:"room_entered,omitempty"`
}