package execution

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ExecutionStore struct {
	*sqlx.DB
}

type Execution struct {
	ID        int       `json:"id,omitempty"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (store ExecutionStore) CreateExecution(e Execution) (Execution, error) {
	err := store.QueryRow(`insert into executions (query) values ($1) returning id, query, created_at`, e.Query).Scan(&e.ID, &e.Query, &e.CreatedAt)

	return e, err
}

func (store ExecutionStore) GetExecution(e Execution) (Execution, error) {
	err := store.QueryRow(`select id, query, created_at from executions where id = $1`, e.ID).Scan(&e.ID, &e.Query, &e.CreatedAt)

	return e, err
}

func (store ExecutionStore) DeleteExecution(e Execution) error {
	_, err := store.Exec(`delete from executions where id = $1`, e.ID)

	return err
}

func (store ExecutionStore) ListExecutions(perPage int, lastId int) ([]Execution, error) {
	executions := make([]Execution, 0)

	rows, err := store.Query(`select id, query, created_at from executions where id < $1 order by id desc limit $2`, lastId, perPage)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var e Execution
		err := rows.Scan(&e.ID, &e.Query, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("could not read execution results: %v", err)
		}
		executions = append(executions, e)
	}

	return executions, err
}
