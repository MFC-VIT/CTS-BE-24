package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomsDone struct {
	RoomA bool `bson:"room_a"`
	RoomB bool `bson:"room_b"`
	RoomC bool `bson:"room_c"`
	RoomD bool `bson:"room_d"`
}

type RoomsGiveUp struct{
	RoomA bool `bson:"room_a"`
	RoomB bool `bson:"room_b"`
	RoomC bool `bson:"room_c"`
	RoomD bool `bson:"room_d"`
}

type Rooms struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"` 
	UserID primitive.ObjectID `bson:"user_id"`
	IsRoomsDone     RoomsDone           `bson:"is_rooms_done"` 
	IsRoomsGiveUp     RoomsDone           `bson:"is_rooms_giveup"` 
}