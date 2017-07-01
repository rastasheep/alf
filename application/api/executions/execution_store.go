package executions

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ExecutionStore struct {
	DbStore *sqlx.DB
	DbData  *sqlx.DB
}

type Execution struct {
	ID        int       `json:"id,omitempty"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (store ExecutionStore) CreateExecution(e Execution) (Execution, error) {
	err := store.DbStore.QueryRow(`insert into executions (query) values ($1) returning id, query, created_at`, e.Query).Scan(&e.ID, &e.Query, &e.CreatedAt)

	return e, err
}

func (store ExecutionStore) GetExecution(id int) (Execution, error) {
	var e Execution
	err := store.DbStore.QueryRow(`select id, query, created_at from executions where id = $1`, id).Scan(&e.ID, &e.Query, &e.CreatedAt)

	return e, err
}

func (store ExecutionStore) DeleteExecution(e Execution) error {
	_, err := store.DbStore.Exec(`delete from executions where id = $1`, e.ID)

	return err
}

func (store ExecutionStore) ListExecutions(perPage int, lastId int) ([]Execution, error) {
	executions := make([]Execution, 0)

	rows, err := store.DbStore.Query(`select id, query, created_at from executions where id < $1 order by id desc limit $2`, lastId, perPage)
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

func (store ExecutionStore) Execute(id int) ([]map[string]interface{}, error) {
	e, err := store.GetExecution(id)
	if err != nil {
		return nil, fmt.Errorf("execution not found: %v", err)
	}

	tx, err := store.DbData.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}

	_, err = tx.Exec("SET TRANSACTION READ ONLY;")
	if err != nil {
		return nil, err
	}

	rows, err := tx.Queryx(e.Query)
	if err != nil {
		return nil, fmt.Errorf("error executing user query: %v", err)
	}

	data := make([]map[string]interface{}, 0)

	defer rows.Close()
	for rows.Next() {
		entry := make(map[string]interface{})

		err := rows.MapScan(entry)
		if err != nil {
			return nil, fmt.Errorf("error serializing user query results: %v", err)
		}

		data = append(data, entry)
	}

	tx.Commit()

	return data, nil
}
