package questions

import (
	"C2S/internal/models"
	"C2S/internal/utils"
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (qs *QuestionControllerStore) GetNextQuestion(ctx context.Context, userID primitive.ObjectID) (models.Question, error) {
	usersCollection := qs.db.Collection(os.Getenv("MONGO_USER_COLLECTION"))
	questionsCollection := qs.db.Collection(os.Getenv("MONGO_QUESTIONS_COLLECTION"))

	var user models.User

	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return models.Question{}, fmt.Errorf("user not found: %v", err)
	}

	if user.RoomEntered == "" {
		return models.Question{}, fmt.Errorf("user is not in any room")
	}

	var questionData models.Questions
	err = questionsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&questionData)
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

	for _, question := range questions {
		if question.Answered == "false" {
			return question, nil
		}
	}

	return models.Question{}, fmt.Errorf("all questions answered in room: %s", user.RoomEntered)
}

func (qs *QuestionControllerStore) QuestionAnswered(ctx context.Context, userID primitive.ObjectID, question models.Question) error {
	questionsCollection := qs.db.Collection(os.Getenv("MONGO_QUESTIONS_COLLECTION"))

	var questionData models.Questions
	err := questionsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&questionData)
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

			_, err := questionsCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("failed to mark question as answered: %v", err)
			}

			return nil 
		}
	}

	return fmt.Errorf("question already answered or not found")
}

