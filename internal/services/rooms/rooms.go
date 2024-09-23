package rooms

import (
	"C2S/internal/models"
	"context"
	"fmt"
	"os"

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

	if user.RoomEntered != "" && user.RoomEntered != roomEntered {
		fmt.Printf("User is already in room: %s\n", user.RoomEntered)
		return fmt.Errorf("user is already in room: %s", user.RoomEntered)
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
	roomsCollection := rc.db.Collection("rooms")
	usersCollection := rc.db.Collection("users")

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if user.RoomEntered != roomEntered {
		return fmt.Errorf("user is not in room: %s", roomEntered)
	}

	roomGiveUpField := fmt.Sprintf("is_rooms_giveup.room_%s", roomEntered)
	update := bson.M{"$set": bson.M{roomGiveUpField: true}}
	_, err = roomsCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		return fmt.Errorf("failed to update room give up: %v", err)
	}

	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"room_entered": ""}})
	if err != nil {
		return fmt.Errorf("failed to update user's room status: %v", err)
	}

	return nil
}

