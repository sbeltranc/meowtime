package models

import (
	"fmt"

	"github.com/gocql/gocql"
)

type User struct {
	ID       string `json:"id" cql:"id"`
	SUB      string `json:"sub" cql:"sub"`
	Email    string `json:"email" cql:"email"`
	Name     string `json:"name" cql:"name"`
	Picture  string `json:"picture" cql:"picture"`
	TeamID   string `json:"team_id,omitempty" cql:"team_id"`
	TeamName string `json:"team_name,omitempty" cql:"team_name"`

	CreatedAt string `json:"created_at,omitempty" cql:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" cql:"updated_at"`
}

func GetUserByID(session *gocql.Session, id string) (*User, error) {
	var user User
	query := `SELECT id, sub, email, name, picture, team_id, team_name FROM users WHERE id = ? LIMIT 1`
	if err := session.Query(query, id).Consistency(gocql.One).Scan(
		&user.ID, &user.SUB, &user.Email, &user.Name, &user.Picture, &user.TeamID, &user.TeamName,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserBySUB(session *gocql.Session, sub string) (*User, error) {
	var user User
	query := `SELECT sub, id, email, name, picture, team_id, team_name FROM users WHERE sub = ? LIMIT 1`
	if err := session.Query(query, sub).Consistency(gocql.One).Scan(
		&user.SUB, &user.ID, &user.Email, &user.Name, &user.Picture, &user.TeamID, &user.TeamName,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(session *gocql.Session, user *User) error {
	query := `INSERT INTO users (id, sub, email, name, picture, team_id, team_name) VALUES (?, ?, ?, ?, ?, ?, ?)`
	return session.Query(query, user.ID, user.SUB, user.Email, user.Name, user.Picture, user.TeamID, user.TeamName).Exec()
}

func InitUserSchema(db *gocql.Session) error {
	// creating the table IF it doesn't exist
	err := db.Query(`
		CREATE TABLE IF NOT EXISTS users (
			id uuid,
			sub text,
			email text,
			name text,
			picture text,
			team_id text,
			team_name text,
			PRIMARY KEY (id, sub)
		)
	`).Exec()

	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	// now checking if something is missing so we can add it up
	columns := []struct {
		name     string
		datatype string
	}{
		{"sub", "text"},
		{"email", "text"},
		{"name", "text"},
		{"picture", "text"},
		{"team_id", "text"},
		{"team_name", "text"},
	}

	for _, col := range columns {
		query := "ALTER TABLE users ADD " + col.name + " " + col.datatype
		db.Query(query).Exec()
	}

	return nil
}

func UpdateUser(session *gocql.Session, user *User) error {
	query := `UPDATE users SET email = ?, name = ?, picture = ?, team_id = ?, team_name = ? WHERE id = ?`
	return session.Query(query, user.Email, user.Name, user.Picture, user.TeamID, user.TeamName, user.ID).Exec()
}
