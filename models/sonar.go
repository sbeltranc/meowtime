package models

import (
	"encoding/json"
	"fmt"

	"github.com/gocql/gocql"
)

type Sonar struct {
	ID        string          `json:"id" cql:"id"`
	UserID    string          `json:"user_id" cql:"user_id"`
	Metadata  json.RawMessage `json:"metadata" cql:"metadata"`
	IPAddress string          `json:"ip_address" cql:"ip_address"`
	ProjectID string          `json:"project_id" cql:"project_id"`
	TotalTime int64           `json:"total_time" cql:"total_time"`
	Software  string          `json:"software" cql:"software"`
}

func InitSonarSchema(db *gocql.Session) error {
	// creating the table IF it doesn't exist
	err := db.Query(`
		CREATE TABLE IF NOT EXISTS sonar (
			id uuid,
			user_id text,
			metadata blob,
			ip_address text,
			project_id text,
			total_time bigint,
			software text,
			PRIMARY KEY (user_id, project_id)
		)
	`).Exec()

	if err != nil {
		return fmt.Errorf("error creating sonar table: %w", err)
	}

	// now checking if something is missing so we can add it up
	columns := []struct {
		name     string
		datatype string
	}{
		{"user_id", "text"},
		{"metadata", "blob"},
		{"ip_address", "text"},
		{"project_id", "text"},
		{"total_time", "bigint"},
		{"software", "text"},
	}

	for _, col := range columns {
		query := "ALTER TABLE sonar ADD " + col.name + " " + col.datatype
		db.Query(query).Exec()
	}

	return nil
}

func CreateSonar(sonar *Sonar, db *gocql.Session) (*Sonar, error) {
	sonar = &Sonar{
		ID:        sonar.ID,
		UserID:    sonar.UserID,
		Metadata:  sonar.Metadata,
		IPAddress: sonar.IPAddress,
		ProjectID: sonar.ProjectID,
		TotalTime: sonar.TotalTime,
		Software:  sonar.Software,
	}

	err := db.Query(`
		INSERT INTO sonar (id, user_id, metadata, ip_address, project_id, total_time, software)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, sonar.ID, sonar.UserID, sonar.Metadata, sonar.IPAddress, sonar.ProjectID, sonar.TotalTime, sonar.Software).Exec()

	if err != nil {
		fmt.Println("Error creating sonar:", err)
		return nil, err
	}

	return sonar, nil
}
