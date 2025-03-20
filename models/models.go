package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Platform  string    `json:"platform" db:"platform"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatSession struct {
	ID        int       `json:"id" db:"id"`
	SessionID uuid.UUID `json:"session_id" db:"session_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	EndedAt   time.Time `json:"ended_at,omitempty" db:"ended_at"`
}

type Message struct {
	ID        int       `json:"id" db:"id"`
	SessionID uuid.UUID `json:"session_id" db:"session_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Message   string    `json:"message" db:"message"`
	BotReply  string    `json:"bot_reply" db:"bot_reply"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type ChatLog struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	Message   string    `json:"message"`
	BotReply  string    `json:"bot_reply"`
	Timestamp time.Time `json:"timestamp"`
} 