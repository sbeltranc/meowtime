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
	ProjectID string          `json:"project_id" cql:"project_id"`
	TotalTime int64           `json:"total_time" cql:"total_time"`
}

func CreateSonar(sonar *Sonar, db *gocql.Session) {
	sonar = &Sonar{
		ID:        sonar.ID,
		UserID:    sonar.UserID,
		Metadata:  sonar.Metadata,
		ProjectID: sonar.ProjectID,
		TotalTime: sonar.TotalTime,
	}

	err := db.Query(`
		INSERT INTO sonar (id, user_id, metadata, project_id, total_time)
		VALUES (?, ?, ?, ?, ?)
	`, sonar.ID, sonar.UserID, sonar.Metadata, sonar.ProjectID, sonar.TotalTime).Exec()

	if err != nil {
		fmt.Println("Error creating sonar:", err)
	}
}
