package models

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

type Session struct {
	ID           string `json:"id" cql:"id"`
	UserID       string `json:"user_id" cql:"user_id"`
	SessionToken string `json:"session_token" cql:"session_token"`
	ExpiresAt    int64  `json:"expires_at" cql:"expires_at"`
}

func InitSessionSchema(db *gocql.Session) error {
	// creating the table IF it doesn't exist
	err := db.Query(`
		CREATE TABLE IF NOT EXISTS session (
			id uuid,
			user_id text,
			session_token text,
			expires_at bigint,
			PRIMARY KEY (user_id, session_token)
		)
	`).Exec()

	if err != nil {
		return fmt.Errorf("error creating session table: %w", err)
	}

	// now checking if something is missing so we can add it up
	columns := []struct {
		name     string
		datatype string
	}{
		{"user_id", "text"},
		{"session_token", "text"},
		{"expires_at", "bigint"},
	}

	for _, col := range columns {
		query := "ALTER TABLE session ADD " + col.name + " " + col.datatype
		db.Query(query).Exec()
	}

	return nil
}

func CreateSession(session *Session, db *gocql.Session) (*Session, error) {
	session = &Session{
		ID:           session.ID,
		UserID:       session.UserID,
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
	}

	err := db.Query(`
		INSERT INTO session (id, user_id, session_token, expires_at)
		VALUES (?, ?, ?, ?)
	`, session.ID, session.UserID, session.SessionToken, session.ExpiresAt).Exec()

	if err != nil {
		fmt.Println("Error creating session:", err)
	}

	return session, err
}

func FindSession(sessionToken string, db *gocql.Session) (*User, error) {
	var session Session

	// finding the session by session token
	query := `SELECT id, user_id, session_token, expires_at FROM session WHERE session_token = ? LIMIT 1`
	if err := db.Query(query, sessionToken).Consistency(gocql.One).Scan(
		&session.ID, &session.UserID, &session.SessionToken, &session.ExpiresAt,
	); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// checking if session is expired
	if time.Now().Unix() > session.ExpiresAt {
		deleteQuery := `DELETE FROM session WHERE session_token = ?`

		if err := db.Query(deleteQuery, sessionToken).Exec(); err != nil {
			return nil, fmt.Errorf("error deleting expired session: %w", err)
		}

		return nil, fmt.Errorf("session expired")
	}

	// session was found, get the associated user info
	user, err := GetUserByID(db, session.UserID)

	if err != nil {
		return nil, fmt.Errorf("user not found for session: %w", err)
	}

	// returning the user info
	return user, nil
}
