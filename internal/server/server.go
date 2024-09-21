package server

import (
	"C2S/internal/controllers/user"
	"C2S/internal/routes"
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

	userHandler := user.NewHandler(users.NewUserStore(s.db))
	userRoutes := routes.NewUserRoutes(userHandler)
	userRoutes.RegisterRoutes(api)

	healthRoutes := routes.NewHealthRoutes()
	healthRoutes.RegisterRoutes(api)
}
