package main

import (
	"gopkg.in/matryer/respond.v1"
	"net/http"
)

type SchemaItem struct {
	TableName  string `json:"table_name"`
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
}

func (s *Server) getSchema(w http.ResponseWriter, r *http.Request) {
	var schema = []SchemaItem{}

	rows, err := s.db.Query("select table_name, column_name, data_type from information_schema.columns where table_schema='public';")
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var item SchemaItem
		err = rows.Scan(&item.TableName, &item.ColumnName, &item.DataType)

		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		schema = append(schema, item)
	}

	respond.With(w, r, http.StatusOK, schema)
}
