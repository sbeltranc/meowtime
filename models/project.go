package models

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
