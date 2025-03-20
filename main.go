package main

import (
	"log"
	"os"

	"chatbot/handlers"
	"chatbot/kafka"
	"chatbot/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize repository
	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize Kafka producer
	producer, err := kafka.NewProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize handler
	handler := handlers.NewHandler(repo, producer)

	// Set up Gin router
	router := gin.Default()

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Messages
		v1.POST("/messages", handler.SendMessage)
		v1.GET("/messages/:user_id", handler.GetConversationHistory)

		// Sessions
		v1.POST("/session/start", handler.StartSession)
		v1.POST("/session/end", handler.EndSession)
	}

	// Health check
	router.GET("/health", handler.HealthCheck)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
