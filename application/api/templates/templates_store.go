package templates

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type TemplateStore struct {
	*sqlx.DB
}

type Template struct {
	ID        int       `json:"id,omitempty"`
	Query     string    `json:"query,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (store TemplateStore) GetTemplate(id int) (*Template, error) {
	var t Template
	err := store.QueryRow(`select id, query, updated_at, created_at from templates where id = $1`, id).Scan(&t.ID, &t.Query, &t.UpdatedAt, &t.CreatedAt)

	return &t, err
}

func (store TemplateStore) CreateTemplate(t *Template) (*Template, error) {
	err := store.QueryRow(`insert into templates (query) values ($1) returning id, query, updated_at, created_at`, t.Query).Scan(&t.ID, &t.Query, &t.UpdatedAt, &t.CreatedAt)

	return t, err
}

func (store TemplateStore) UpdateTemplate(t *Template) (*Template, error) {
	err := store.QueryRow(`update templates set query = $1, updated_at = default returning id, query, updated_at, created_at`, t.Query).Scan(&t.ID, &t.Query, &t.UpdatedAt, &t.CreatedAt)

	return t, err
}

func (store TemplateStore) DeleteTemplate(id int) error {
	_, err := store.Exec(`delete from templates where id = $1`, id)

	return err
}

func (store TemplateStore) ListTemplates(perPage int, lastId int) (*[]Template, error) {
	templates := make([]Template, 0)

	rows, err := store.Query(`select id, query, updated_at, created_at from templates where id < $1 order by id desc limit $2`, lastId, perPage)
	if err != nil {
		return nil, fmt.Errorf("error querying template list: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var t Template
		err := rows.Scan(&t.ID, &t.Query, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("could not serialize template list results: %v", err)
		}
		templates = append(templates, t)
	}

	return &templates, err
}

func (t *Template) Valid() error {
	if t.Query == "" {
		return fmt.Errorf("query is required")
	}

	return nil
}
