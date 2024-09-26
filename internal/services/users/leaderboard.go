package users

import (
	"C2S/internal/utils"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Store) GetLeaderBoardHandler(c *fiber.Ctx) error {
	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.D{{Key: "score", Value: -1}}}}, 
		{{Key: "$limit", Value: 10}},                               
	}

	cursor, err := r.usersCollection.Aggregate(c.Context(), pipeline)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("failed to retrieve leaderboard: %v", err))
	}
	defer cursor.Close(context.TODO())

	var leaderboardData []map[string]interface{}
	for cursor.Next(c.Context()) {
		var user bson.M
		if err := cursor.Decode(&user); err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("failed to decode leaderboard user: %v", err))
		}

		leaderboardData = append(leaderboardData, map[string]interface{}{
			"username": user["username"],
			"score":    user["score"],
		})
	}

	if err := cursor.Err(); err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("cursor error: %v", err))
	}

	return utils.WriteJSON(c, fiber.StatusOK, leaderboardData)
}