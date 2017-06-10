package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"gopkg.in/matryer/respond.v1"
	"net/http"
	"regexp"
)

type Execution struct {
	ID    int    `json:"id"`
	Query string `json:"query"`
	Name  string `json:"name"`
}

func (s *Server) createExecution(w http.ResponseWriter, r *http.Request) {
	var e Execution
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respond.With(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	re := regexp.MustCompile("(?i)(SET.*TRANSACTION)|(SET.*SESSION.*CHARACTERISTICS)")
	matched := re.MatchString(e.Query)
	if matched {
		s.logger.Printf("Blocked execution of query: %v", e.Query)
		respond.With(w, r, http.StatusBadRequest, "pq: you are not allowed to change the characteristics of the current transaction")
		return
	}

	s.logger.Printf("Executing query: %v", e.Query)

	tx, err := s.db.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	_, err = tx.Exec("SET TRANSACTION READ ONLY;")
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	rows, err := tx.Queryx(e.Query)

	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	data := make([]map[string]interface{}, 0)

	defer rows.Close()
	for rows.Next() {
		entry := make(map[string]interface{})

		err := rows.MapScan(entry)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		data = append(data, entry)
	}
	tx.Commit()

	respond.With(w, r, http.StatusOK, data)
}

func (s *Server) listExecutions(w http.ResponseWriter, r *http.Request) {
	respond.With(w, r, http.StatusInternalServerError, errors.New("This is the RESTful api"))
}
