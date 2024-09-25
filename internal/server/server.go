package server

import (
	"C2S/internal/controllers/question"
	room "C2S/internal/controllers/rooms"
	"C2S/internal/controllers/user"
	"C2S/internal/routes"
	"C2S/internal/services/questions"
	rooms "C2S/internal/services/rooms"
	"C2S/internal/services/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type FiberServer struct {
	App *fiber.App
	db  *mongo.Database
}


func New(app *fiber.App, client *mongo.Client, db *mongo.Database) *FiberServer {
	server := &FiberServer{
		App: app,
		db:  db,
	}

	server.App.Use(logger.New())

	return server
}


func (s *FiberServer) RegisterFiberRoutes() {
	api := s.App.Group("/api/v1")

	userStore := users.NewUserStore(s.db)
	roomStore := rooms.NewRoomStore(s.db)
	questionStore := questions.NewQuestionStore(s.db)
	adminStore := users.NewUserStore(s.db)

	userHandler := user.NewHandler(userStore)
	userRoutes := routes.NewUserRoutes(userHandler)
	userRoutes.RegisterRoutes(api)

	healthRoutes := routes.NewHealthRoutes()
	healthRoutes.RegisterRoutes(api)

	roomHandler := room.NewHandler(userStore, roomStore) 
	roomRoutes := routes.NewRoomRoutes(roomHandler,userStore)
	roomRoutes.RegisterRoutes(api)

	questionHandler := question.NewHandler(questionStore)
	questionRoutes := routes.NewQuestionRoutes(questionHandler,userStore)
	questionRoutes.RegisterRoutes(api)

	adminRoutes := routes.NewAdminRoutes(adminStore)
	adminRoutes.RegisterRoutes(api)
}
