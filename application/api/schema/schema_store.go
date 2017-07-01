package schema

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SchemaStore struct {
	*sqlx.DB
}

type SchemaItem struct {
	TableName  string `json:"table_name"`
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
}

func (store *SchemaStore) GetSchema() (*[]SchemaItem, error) {
	var schema = []SchemaItem{}

	rows, err := store.Query("select table_name, column_name, data_type from information_schema.columns where table_schema='public';")
	if err != nil {
		return nil, fmt.Errorf("error querying database schema: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var item SchemaItem
		err = rows.Scan(&item.TableName, &item.ColumnName, &item.DataType)

		if err != nil {
			return nil, fmt.Errorf("could not read schema query results: %v", err)
		}

		schema = append(schema, item)
	}

	return &schema, nil
}
