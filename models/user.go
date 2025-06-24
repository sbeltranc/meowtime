package models

import (
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
	query := `SELECT id, email, name, picture, team_id, team_name FROM users WHERE id = ? LIMIT 1`
	if err := session.Query(query, id).Consistency(gocql.One).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture, &user.TeamID, &user.TeamName,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserBySUB(session *gocql.Session, sub string) (*User, error) {
	var user User
	query := `SELECT id, email, name, picture, team_id, team_name FROM users WHERE sub = ? LIMIT 1`
	if err := session.Query(query, sub).Consistency(gocql.One).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture, &user.TeamID, &user.TeamName,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(session *gocql.Session, user *User) error {
	query := `INSERT INTO users (id, sub, email, name, picture, team_id, team_name) VALUES (?, ?, ?, ?, ?, ?, ?)`
	return session.Query(query, user.ID, user.SUB, user.Email, user.Name, user.Picture, user.TeamID, user.TeamName).Exec()
}

func InitSchema(session *gocql.Session) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		sub TEXT,
		email TEXT,
		name TEXT,
		picture TEXT,
		team_id TEXT,
		team_name TEXT
	);
	`
	return session.Query(query).Exec()
}

func UpdateUser(session *gocql.Session, user *User) error {
	query := `UPDATE users SET email = ?, name = ?, picture = ?, team_id = ?, team_name = ? WHERE id = ?`
	return session.Query(query, user.Email, user.Name, user.Picture, user.TeamID, user.TeamName, user.ID).Exec()
}
