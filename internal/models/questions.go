package models

import "go.mongodb.org/mongo-driver/bson/primitive"


type Question struct {
	Question string `bson:"question"`
	QuestionId int `bson:"question_id"`
	Answer   string `bson:"answer"`
	Room     string `bson:"room"`
	Answered string   `bson:"answered"` 
}


type Room struct {
	Questions []Question `bson:"questions"` 
}

type Questions struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	UserName   string             `bson:"username"`
	RoomA Room               `bson:"room_a"`
	RoomB Room               `bson:"room_b"`
	RoomC Room               `bson:"room_c"`
	RoomD Room               `bson:"room_d"`
}
