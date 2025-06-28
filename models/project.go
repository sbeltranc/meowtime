package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Project struct {
	ID          string `json:"id" cql:"id"`
	Name        string `json:"name" cql:"name"`
	Description string `json:"description" cql:"description"`

	RepositoryName    string `json:"repository_name" cql:"repository_name"`
	RepositoryAuthor  string `json:"repository_author" cql:"repository_author"`
	RepositoryService string `json:"repository_service" cql:"repository_service"`

	OwnerID string `json:"owner_id" cql:"owner_id"`

	CreatedAt string `json:"created_at" cql:"created_at"`
	UpdatedAt string `json:"updated_at" cql:"updated_at"`
}

func InitProjectSchema(scylla *gocql.Session) error {
	err := scylla.Query(`CREATE TABLE IF NOT EXISTS projects (
		id text,
		name text,
		description text,
		repository_name text,
		repository_author text,
		repository_service text,
		owner_id text,
		created_at text,
		updated_at text,
		PRIMARY KEY (id, name)
	)`).Exec()
	if err != nil {
		return err
	}

	columns := []struct {
		name     string
		datatype string
	}{
		{"id", "text"},
		{"name", "text"},
		{"description", "text"},
		{"repository_name", "text"},
		{"repository_author", "text"},
		{"repository_service", "text"},
		{"owner_id", "text"},
		{"created_at", "text"},
		{"updated_at", "text"},
	}

	for _, col := range columns {
		query := "ALTER TABLE projects ADD " + col.name + " " + col.datatype
		scylla.Query(query).Exec()
	}

	return nil
}

func CreateProject(name string, ownerID string, scylla *gocql.Session) (*Project, error) {
	project := &Project{
		ID:          gocql.TimeUUID().String(),
		Name:        name,
		Description: "",

		RepositoryName:    "",
		RepositoryAuthor:  "",
		RepositoryService: "",

		OwnerID: ownerID,

		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	err := scylla.Query(`INSERT INTO projects (id, name, description, repository_name, repository_author, repository_service, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		project.ID, project.Name, project.Description, project.RepositoryName, project.RepositoryAuthor, project.RepositoryService, project.OwnerID, project.CreatedAt, project.UpdatedAt).Exec()
	if err != nil {
		return nil, err
	}

	return project, nil
}

func FindProjectByNameAndOwner(name string, ownerID string, scylla *gocql.Session) (*Project, error) {
	var project Project
	err := scylla.Query(`SELECT id, name, description, repository_name, repository_author, repository_service, owner_id, created_at, updated_at FROM projects WHERE name = ? AND owner_id = ?`, name, ownerID).Scan(&project.ID, &project.Name, &project.Description, &project.RepositoryName, &project.RepositoryAuthor, &project.RepositoryService, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func FindProjectByID(id string, scylla *gocql.Session) (*Project, error) {
	var project Project
	err := scylla.Query(`SELECT id, name, description, repository_name, repository_author, repository_service, owner_id, created_at, updated_at FROM projects WHERE id = ?`, id).Scan(&project.ID, &project.Name, &project.Description, &project.RepositoryName, &project.RepositoryAuthor, &project.RepositoryService, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &project, nil
}
