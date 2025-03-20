package repository

import (
	"database/sql"
	"fmt"
	"os"

	"chatbot/models"

	"github.com/google/uuid"
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
		VALUES ($1::uuid, $2, $3)
	`
	_, err := r.db.Exec(query, session.SessionID.String(), session.UserID, session.Status)
	return err
}

func (r *Repository) EndSession(sessionID string) error {
	query := `
		UPDATE chat_sessions
		SET status = 'ended', ended_at = NOW()
		WHERE session_id = $1::uuid
	`
	_, err := r.db.Exec(query, sessionID)
	return err
}

func (r *Repository) SaveMessage(message *models.Message) error {
	query := `
		INSERT INTO messages (session_id, user_id, message, bot_reply)
		VALUES ($1::uuid, $2, $3, $4)
	`
	_, err := r.db.Exec(query, message.SessionID.String(), message.UserID, message.Message, message.BotReply)
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
		var sessionIDStr string
		err := rows.Scan(&msg.ID, &sessionIDStr, &msg.UserID, &msg.Message, &msg.BotReply, &msg.Timestamp)
		if err != nil {
			return nil, err
		}

		msg.SessionID, err = uuid.Parse(sessionIDStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing session ID: %v", err)
		}

		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *Repository) GetActiveSession(userID string) (*models.ChatSession, error) {
	query := `
		SELECT id, session_id, user_id, status, created_at, ended_at
		FROM chat_sessions
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var session models.ChatSession
	var sessionIDStr string
	err := r.db.QueryRow(query, userID).Scan(
		&session.ID,
		&sessionIDStr,
		&session.UserID,
		&session.Status,
		&session.CreatedAt,
		&session.EndedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Convert string to UUID
	session.SessionID, err = uuid.Parse(sessionIDStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing session ID: %v", err)
	}

	return &session, nil
}

func (r *Repository) GetSession(sessionID uuid.UUID) (*models.ChatSession, error) {
	var session models.ChatSession
	var sessionIDStr string

	err := r.db.QueryRow(`
		SELECT session_id::text, user_id, status, created_at
		FROM chat_sessions
		WHERE session_id = $1::uuid
	`, sessionID.String()).Scan(&sessionIDStr, &session.UserID, &session.Status, &session.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	session.SessionID, err = uuid.Parse(sessionIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID: %v", err)
	}

	return &session, nil
}
