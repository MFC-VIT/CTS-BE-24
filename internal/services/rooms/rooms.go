package rooms

import (
	"C2S/internal/models"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func (rc *RoomControllerStore) EnterRoom(ctx context.Context, userID primitive.ObjectID, roomEntered string) error {
	usersCollection := rc.db.Collection(os.Getenv("MONGO_USER_COLLECTION"))

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	fmt.Printf("User ID: %s, Current Room: %s, Requested Room: %s\n", userID.Hex(), user.RoomEntered, roomEntered)

	roomStatus, err := rc.collectUserRoomsStatus(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to collect room status: %v", err)
	}
	fmt.Printf("Room status for user %s: %v\n", userID.Hex(), roomStatus)
	if user.RoomEntered != "" && user.RoomEntered != roomEntered {
		var canEnter bool
		for _, status := range roomStatus {
			if status == fmt.Sprintf("%sD", roomEntered) || status == fmt.Sprintf("%sG", roomEntered) {
				canEnter = true
				break
			}
		}

		if !canEnter {
			fmt.Printf("User is already in room: %s and cannot enter room: %s\n", user.RoomEntered, roomEntered)
			return fmt.Errorf("user is already in room: %s and cannot enter room: %s", user.RoomEntered, roomEntered)
		}

		fmt.Printf("Allowing user to switch rooms from %s to %s as room is done or giveup.\n", user.RoomEntered, roomEntered)
		return nil
	}

	update := bson.M{"$set": bson.M{"room_entered": roomEntered}}
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("failed to update room entered: %v", err)
	}
	fmt.Printf("User %s has entered room: %s\n", userID.Hex(), roomEntered)

	return nil
}

func (rc *RoomControllerStore) EscapeRoom(ctx context.Context, userID primitive.ObjectID, roomEntered string) error {

	var user models.User
	err := rc.usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if user.RoomEntered != roomEntered {
		return fmt.Errorf("user is not in room: %s", roomEntered)
	}

	roomDoneField := fmt.Sprintf("is_rooms_done.room_%s", strings.ToLower(roomEntered))
	roomGiveUpField := fmt.Sprintf("is_rooms_giveup.room_%s", strings.ToLower(roomEntered))
	update := bson.M{"$set": bson.M{roomGiveUpField: true,roomDoneField:true}}

	_, err = rc.roomsCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		return fmt.Errorf("failed to update room give up: %v", err)
	}

	_, err = rc.usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"room_entered": ""}})
	if err != nil {
		return fmt.Errorf("failed to update user's room status: %v", err)
	}

	if err := rc.CheckUnansweredQuestionsAndUpdateScore(ctx, userID, roomEntered); err != nil {
		return err
	}

	fmt.Printf("User %s has successfully escaped room %s. Room marked as done and giveup.\n", userID.Hex(), roomEntered)

	return nil
}

//Go concurrency
func (rc *RoomControllerStore) collectUserRoomsStatus(ctx context.Context, userID primitive.ObjectID) ([]string, error) {

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

func (rc *RoomControllerStore) CheckUnansweredQuestionsAndUpdateScore(ctx context.Context, userID primitive.ObjectID, roomEntered string) error {

	var questionData models.Questions
	err := rc.questionsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&questionData)
	if err != nil {
		return fmt.Errorf("failed to fetch questions for user: %v", err)
	}

	var questions []models.Question
	switch strings.ToUpper(roomEntered) {
	case "A":
		questions = questionData.RoomA.Questions
	case "B":
		questions = questionData.RoomB.Questions
	case "C":
		questions = questionData.RoomC.Questions
	case "D":
		questions = questionData.RoomD.Questions
	default:
		return fmt.Errorf("unknown room: %s", roomEntered)
	}

	unansweredCount := 0
	for _, question := range questions {
		if question.Answered == "false" {
			unansweredCount++
		}
	}

	scoreDeduction := unansweredCount * 10

	_, err = rc.usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$inc": bson.M{"score": -scoreDeduction}})
	if err != nil {
		return fmt.Errorf("failed to update user's score: %v", err)
	}

	fmt.Printf("User %s has %d unanswered questions in room %s. Score deducted by %d.\n", userID.Hex(), unansweredCount, roomEntered, scoreDeduction)

	return nil
}


// func (rc *RoomControllerStore) collectUserRoomsStatus(ctx context.Context, userID primitive.ObjectID) ([]string, error) {
// 	var roomStatus models.Rooms
// 	err := rc.roomsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&roomStatus)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch room status: %v", err)
// 	}

// 	var roomStatusArray []string

// 	rooms := []struct {
// 		done    bool
// 		giveUp  bool
// 		room    string
// 	}{
// 		{roomStatus.IsRoomsDone.RoomA, roomStatus.IsRoomsGiveUp.RoomA, "A"},
// 		{roomStatus.IsRoomsDone.RoomB, roomStatus.IsRoomsGiveUp.RoomB, "B"},
// 		{roomStatus.IsRoomsDone.RoomC, roomStatus.IsRoomsGiveUp.RoomC, "C"},
// 		{roomStatus.IsRoomsDone.RoomD, roomStatus.IsRoomsGiveUp.RoomD, "D"},
// 	}

// 	for _, r := range rooms {
// 		if r.done {
// 			roomStatusArray = append(roomStatusArray, fmt.Sprintf("%sD", r.room)) 
// 		} else if r.giveUp {
// 			roomStatusArray = append(roomStatusArray, fmt.Sprintf("%sG", r.room)) 
// 		} else {
// 			roomStatusArray = append(roomStatusArray, fmt.Sprintf("%s-", r.room)) 
// 		}
// 	}
// 	return roomStatusArray, nil
// }












