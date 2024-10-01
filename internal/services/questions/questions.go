package questions

import (
	"C2S/internal/models"
	"C2S/internal/utils"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

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

	roomStatusArray, err := qs.collectStatus(ctx, userID)
	if err != nil {
		return models.Question{}, fmt.Errorf("failed to fetch room status: %v", err)
	}

	for _, status := range roomStatusArray {
		if status == fmt.Sprintf("%sD", user.RoomEntered) {
			switch user.RoomEntered {
			case "A":
				return models.Question{}, fmt.Errorf("clue: del")
			case "B":
				return models.Question{}, fmt.Errorf("clue: apto")
			case "C":
				return models.Question{}, fmt.Errorf("clue: mas")
			case "D":
				return models.Question{}, fmt.Errorf("clue: supervivencia")
			default:
				return models.Question{}, fmt.Errorf("unknown room: %s", user.RoomEntered)
			}
		}
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
			return models.Question{}, fmt.Errorf("clue: del")
		case "B":
			return models.Question{}, fmt.Errorf("clue: apto")
		case "C":
			return models.Question{}, fmt.Errorf("clue: mas")
		case "D":
			return models.Question{}, fmt.Errorf("clue: supervivencia")
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


//Go concurrency
func (rc *QuestionControllerStore) collectStatus(ctx context.Context, userID primitive.ObjectID) ([]string, error) {

	var roomStatus models.Rooms
	err := rc.roomsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&roomStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room status: %v", err)
	}
	roomStatusChan := make(chan string, 4)

	var wg sync.WaitGroup

	checkRoomStatus := func(roomName string, done bool, giveUp bool) {
		defer wg.Done()
		var result string
		if done {
			result = fmt.Sprintf("%sD", roomName)
		} else if giveUp {
			result = fmt.Sprintf("%sG", roomName) 
		} else {
			result = fmt.Sprintf("%s-", roomName) 
		}
		roomStatusChan <- result
	}

	wg.Add(4)
	go checkRoomStatus("A", roomStatus.IsRoomsDone.RoomA, roomStatus.IsRoomsGiveUp.RoomA)
	go checkRoomStatus("B", roomStatus.IsRoomsDone.RoomB, roomStatus.IsRoomsGiveUp.RoomB)
	go checkRoomStatus("C", roomStatus.IsRoomsDone.RoomC, roomStatus.IsRoomsGiveUp.RoomC)
	go checkRoomStatus("D", roomStatus.IsRoomsDone.RoomD, roomStatus.IsRoomsGiveUp.RoomD)

	go func() {
		wg.Wait()
		close(roomStatusChan)
	}()

	var roomStatusArray []string
	for status := range roomStatusChan {
		roomStatusArray = append(roomStatusArray, status)
	}
	log.Printf("room status array: %+v", roomStatusArray)
	return roomStatusArray, nil
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
	log.Println("Starting QuestionAnswered function")

	var questionData models.Questions
	err := qs.questionsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&questionData)
	if err != nil {
		log.Printf("Failed to fetch questions for user %v: %v\n", userID, err)
		return fmt.Errorf("failed to fetch questions for user: %v", err)
	}

	var questions []models.Question
	log.Printf("Room: %s\n", question.Room)

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
		log.Printf("Unknown room: %s\n", question.Room)
		return fmt.Errorf("unknown room: %s", question.Room)
	}

	answerFilePath := "internal/files/answer.yaml"
	log.Println("Loading answers from file:", answerFilePath)
	answerData, err := utils.LoadAnswers(answerFilePath)
	if err != nil {
		log.Printf("Failed to load answers: %v\n", err)
		return fmt.Errorf("failed to load answers: %v", err)
	}

	var correctAnswer string
	for _, ansQuestion := range answerData.Questions {
		if ansQuestion.QuestionId == question.QuestionId && ansQuestion.Room == question.Room {
			correctAnswer = ansQuestion.Answer
			break
		}
	}

	if correctAnswer == "" {
		log.Printf("Question ID %d not found in answer file for room %s\n", question.QuestionId, question.Room)
		return fmt.Errorf("question not found in answer file")
	}

	if question.Answer != correctAnswer {
		log.Printf("Incorrect answer: got %s, expected %s\n", question.Answer, correctAnswer)
		return fmt.Errorf("incorrect answer")
	}

	for i, q := range questions {
		if q.QuestionId == question.QuestionId && q.Answered == "false" {
			log.Printf("Marking question %v as answered for user %v\n", q.QuestionId, userID)

			filter := bson.M{
				"user_id": userID,
				fmt.Sprintf("room_%s.questions.%d.question", strings.ToLower(question.Room), i): q.Question,
			}
			update := bson.M{
				"$set": bson.M{
					fmt.Sprintf("room_%s.questions.%d.answered", strings.ToLower(question.Room), i): "true",
				},
			}

			_, err := qs.questionsCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				log.Printf("Failed to mark question as answered: %v\n", err)
				return fmt.Errorf("failed to mark question as answered: %v", err)
			}
			return nil
		}
		if i == 3 {
			allAnswered := true
			for _, q := range questions {
				if q.Answered == "false" {
					allAnswered = false
					break
				}
			}

			if allAnswered {
				log.Printf("All questions answered in room %s for user %v\n", question.Room, userID)
				err = qs.markRoomAsDone(ctx, userID, question.Room, qs.roomsCollection)
				if err != nil {
					log.Printf("Failed to update room status: %v\n", err)
					return fmt.Errorf("failed to update room status: %v", err)
				}

				updateUser := bson.M{
					"$set": bson.M{
						"room_entered": "",
					},
				}

				_, err = qs.usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, updateUser)
				if err != nil {
					log.Printf("Failed to clear RoomEntered field for user %v: %v\n", userID, err)
					return fmt.Errorf("failed to clear RoomEntered field for user: %v", err)
				}
			}
		}
	}
	log.Println("Question already answered or not found")
	return fmt.Errorf("question already answered or not found")
}

