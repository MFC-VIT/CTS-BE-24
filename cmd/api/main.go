package main

import (
	"C2S/internal/db"
	"C2S/internal/server"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", 
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DBNAME")

	client, db := db.ConnectMongoDB(mongoURI, dbName)
	if client == nil || db == nil {
		log.Fatal("Failed to connect to MongoDB")
	}

	fiberServer := server.New(app, client, db)
	fiberServer.RegisterFiberRoutes()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Pong",
		})
	})

	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %s", portStr)
	}


	log.Printf("Server is starting on PORT:%d", port)
	if err := app.Listen(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
}
