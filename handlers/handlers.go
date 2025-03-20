package handlers

import (
	"net/http"
	"time"

	"chatbot/kafka"
	"chatbot/models"
	"chatbot/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	repo     *repository.Repository
	producer *kafka.Producer
}

func NewHandler(repo *repository.Repository, producer *kafka.Producer) *Handler {
	return &Handler{
		repo:     repo,
		producer: producer,
	}
}

type SendMessageRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Platform string `json:"platform" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

type SendMessageResponse struct {
	Reply string `json:"reply"`
}

func (h *Handler) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create or get user
	user := &models.User{
		UserID:   req.UserID,
		Platform: req.Platform,
	}
	if err := h.repo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// TODO: Process message with NLP engine (Dialogflow/Rasa/LLM)
	// For now, return a mock response
	reply := "This is a mock response. In production, this would be processed by an NLP engine."

	// Save message
	message := &models.Message{
		SessionID: uuid.New(),
		UserID:    req.UserID,
		Message:   req.Message,
		BotReply:  reply,
		Timestamp: time.Now(),
	}
	if err := h.repo.SaveMessage(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Log to Kafka
	chatLog := &models.ChatLog{
		UserID:    req.UserID,
		SessionID: message.SessionID.String(),
		Message:   req.Message,
		BotReply:  reply,
		Timestamp: time.Now(),
	}
	if err := h.producer.LogMessage(chatLog); err != nil {
		// Log error but don't fail the request
		c.Error(err)
	}

	c.JSON(http.StatusOK, SendMessageResponse{Reply: reply})
}

func (h *Handler) GetConversationHistory(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	messages, err := h.repo.GetUserConversation(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":      userID,
		"conversation": messages,
	})
}

type StartSessionRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Platform string `json:"platform" binding:"required"`
}

type StartSessionResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

func (h *Handler) StartSession(c *gin.Context) {
	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := &models.ChatSession{
		SessionID: uuid.New(),
		UserID:    req.UserID,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, StartSessionResponse{
		SessionID: session.SessionID.String(),
		Message:   "Session started",
	})
}

type EndSessionRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	SessionID string `json:"session_id" binding:"required"`
}

func (h *Handler) EndSession(c *gin.Context) {
	var req EndSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.EndSession(req.SessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session ended"})
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
