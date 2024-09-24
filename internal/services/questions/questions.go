package questions

import (
	"C2S/internal/models"
	"C2S/internal/utils"
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (qs *QuestionControllerStore) GetNextQuestion(ctx context.Context, userID primitive.ObjectID) (models.Question, error) {
	
	var user models.User

	err := qs.usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return models.Question{}, fmt.Errorf("user not found: %v", err)
	}

	if user.RoomEntered == "" {
		return models.Question{}, fmt.Errorf("user is not in any room")
	}

	var questionData models.Questions
	err = qs.questionsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&questionData)
	if err != nil {
		return models.Question{}, fmt.Errorf("failed to fetch questions: %v", err)
	}

	var questions []models.Question
	switch user.RoomEntered {
	case "A":
		questions = questionData.RoomA.Questions
	case "B":
		questions = questionData.RoomB.Questions
	case "C":
		questions = questionData.RoomC.Questions
	case "D":
		questions = questionData.RoomD.Questions
	default:
		return models.Question{}, fmt.Errorf("unknown room: %s", user.RoomEntered)
	}

	allAnswered := true
	for _, question := range questions {
		if question.Answered == "false" {
			allAnswered = false
			break
		}
	}

	if allAnswered {
		err = qs.markRoomAsDone(ctx, userID, user.RoomEntered, qs.roomsCollection)
		if err != nil {
			return models.Question{}, fmt.Errorf("failed to update room status: %v", err)
		}

		switch user.RoomEntered {
		case "A":
			return models.Question{}, fmt.Errorf("clue: A")
		case "B":
			return models.Question{}, fmt.Errorf("clue: B")
		case "C":
			return models.Question{}, fmt.Errorf("clue: C")
		case "D":
			return models.Question{}, fmt.Errorf("clue: D")
		}
	}

	for _, question := range questions {
		if question.Answered == "false" {
			return question, nil
		}
	}

	err = qs.markRoomAsDone(ctx, userID, user.RoomEntered,qs.roomsCollection)
	if err != nil {
		return models.Question{}, fmt.Errorf("failed to update room status: %v", err)
	}

	return models.Question{}, fmt.Errorf("all questions answered in room: %s", user.RoomEntered)
}


func (qs *QuestionControllerStore) markRoomAsDone(ctx context.Context, userID primitive.ObjectID, room string, roomsCollection *mongo.Collection) error {
	var roomStatus models.Rooms
	err := roomsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&roomStatus)
	if err != nil {
		return fmt.Errorf("failed to fetch room status for user: %v", err)
	}

	switch room {
	case "A":
		if roomStatus.IsRoomsDone.RoomA {
			return nil
		}
	case "B":
		if roomStatus.IsRoomsDone.RoomB {
			return nil
		}
	case "C":
		if roomStatus.IsRoomsDone.RoomC {
			return nil
		}
	case "D":
		if roomStatus.IsRoomsDone.RoomD {
			return nil
		}
	default:
		return fmt.Errorf("unknown room: %s", room)
	}

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("is_rooms_done.room_%s", strings.ToLower(room)): true,
		},
	}

	_, err = roomsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update room status in rooms collection: %v", err)
	}

	return nil
}

func (qs *QuestionControllerStore) QuestionAnswered(ctx context.Context, userID primitive.ObjectID, question models.Question) error {
	var questionData models.Questions
	err := qs.questionsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&questionData)
	if err != nil {
		return fmt.Errorf("failed to fetch questions for user: %v", err)
	}

	var questions []models.Question
	switch question.Room {
	case "A":
		questions = questionData.RoomA.Questions
	case "B":
		questions = questionData.RoomB.Questions
	case "C":
		questions = questionData.RoomC.Questions
	case "D":
		questions = questionData.RoomD.Questions
	default:
		return fmt.Errorf("unknown room: %s", question.Room)
	}

	answerFilePath := "internal/seeders/answer.yaml" 
	answerData, err := utils.LoadAnswers(answerFilePath)
	if err != nil {
		return fmt.Errorf("failed to load answers: %v", err)
	}

	var correctAnswer string
	for _, ansQuestion := range answerData.Questions {
		if ansQuestion.Question == question.Question && ansQuestion.Room == question.Room {
			correctAnswer = ansQuestion.Answer
			break
		}
	}

	if correctAnswer == "" {
		return fmt.Errorf("question not found in answer file")
	}

	if question.Answer != correctAnswer {
		return fmt.Errorf("incorrect answer")
	}

	for i, q := range questions {
		if q.Question == question.Question && q.Answered == "false" {
			filter := bson.M{
				"user_id":userID,
				fmt.Sprintf("room_%s.questions.%d.question", strings.ToLower(question.Room),i): q.Question,
			}
			update := bson.M{
				"$set": bson.M{
					fmt.Sprintf("room_%s.questions.%d.answered", strings.ToLower(question.Room), i): "true",
				},
			}

			_, err := qs.questionsCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("failed to mark question as answered: %v", err)
			}
			return nil 
		}
		if i == 3{
			allAnswered := true
			for _, q := range questions {
				if q.Answered == "false" {
					allAnswered = false
					break
				}
			}

			if allAnswered {
				err = qs.markRoomAsDone(ctx, userID, question.Room, qs.roomsCollection)
				if err != nil {
					return fmt.Errorf("failed to update room status: %v", err)
				}

				updateUser := bson.M{
					"$set": bson.M{
						"room_entered": "",
					},
				}
			
				_, err = qs.usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, updateUser)
				if err != nil {
					return fmt.Errorf("failed to clear RoomEntered field for user: %v", err)
				}
			}
		}
	}
	return fmt.Errorf("question already answered or not found")
}

