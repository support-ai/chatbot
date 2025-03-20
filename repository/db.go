package repository

import (
	"database/sql"
	"fmt"
	"os"

	"chatbot/models"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Close() {
	panic("unimplemented")
}

func NewRepository() (*Repository, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (user_id, platform)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING
	`
	_, err := r.db.Exec(query, user.UserID, user.Platform)
	return err
}

func (r *Repository) CreateSession(session *models.ChatSession) error {
	query := `
		INSERT INTO chat_sessions (session_id, user_id, status)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, session.SessionID, session.UserID, session.Status)
	return err
}

func (r *Repository) EndSession(sessionID string) error {
	query := `
		UPDATE chat_sessions
		SET status = 'ended', ended_at = NOW()
		WHERE session_id = $1
	`
	_, err := r.db.Exec(query, sessionID)
	return err
}

func (r *Repository) SaveMessage(message *models.Message) error {
	query := `
		INSERT INTO messages (session_id, user_id, message, bot_reply)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, message.SessionID, message.UserID, message.Message, message.BotReply)
	return err
}

func (r *Repository) GetUserConversation(userID string) ([]models.Message, error) {
	query := `
		SELECT id, session_id, user_id, message, bot_reply, timestamp
		FROM messages
		WHERE user_id = $1
		ORDER BY timestamp ASC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.SessionID, &msg.UserID, &msg.Message, &msg.BotReply, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
