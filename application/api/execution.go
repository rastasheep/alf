package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/rastasheep/alf/respond"
)

type Execution struct {
	ID        int       `json:"id,omitempty"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type ExecutionStore struct {
	*sqlx.DB
}

func (db ExecutionStore) CreateExecution(execution Execution) (Execution, error) {
	err := db.QueryRow(`
    insert into executions (query)
    values ($1)
    returning id, query, created_at
  `, execution.Query).Scan(&execution.ID, &execution.Query, &execution.CreatedAt)

	return execution, err
}

func (db ExecutionStore) GetExecution(execution Execution) (Execution, error) {
	err := db.QueryRow(`
    select id, query, created_at from executions
    where id = $1
  `, execution.ID).Scan(&execution.ID, &execution.Query, &execution.CreatedAt)
	return execution, err
}

func (db ExecutionStore) DeleteExecution(execution Execution) error {
	_, err := db.Exec(`
    delete from executions
    where id = $1
  `, execution.ID)
	return err
}

func (db ExecutionStore) ListExecutions(perPage int, lastId int) ([]Execution, error) {
	executions := make([]Execution, 0)

	rows, err := db.Query(`
    select id, query, created_at from executions
    where id < $1
    order by id desc
    limit $2
  `, lastId, perPage)

	defer rows.Close()
	for rows.Next() {
		var e Execution
		err := rows.Scan(&e.ID, &e.Query, &e.CreatedAt)
		if err != nil {
			return make([]Execution, 0), err
		}
		executions = append(executions, e)
	}

	return executions, err
}

func (s *Server) createExecution(w http.ResponseWriter, r *http.Request) {
	var e Execution
	logPrefix := r.Context().Value("logPrefix").(string)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respond.With(w, r, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}
	defer r.Body.Close()

	re := regexp.MustCompile("(?i)(SET.*TRANSACTION)|(SET.*SESSION.*CHARACTERISTICS)")
	matched := re.MatchString(e.Query)
	if matched {
		s.logger.Printf("%s blocked execution of query: %s", logPrefix, e.Query)
		respond.With(w, r, http.StatusBadRequest, errors.New("you are not allowed to change the characteristics of transaction"))
		return
	}

	db := ExecutionStore{s.dbStore}

	e, err := db.CreateExecution(e)

	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	go s.resultsCache.GetResults(strconv.Itoa(e.ID), nil)

	respond.With(w, r, http.StatusCreated, e)
}

func (s *Server) getExecution(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	e := Execution{ID: id}
	db := ExecutionStore{s.dbStore}

	e, err = db.GetExecution(e)
	if err != nil {
		respond.With(w, r, http.StatusNotFound, errors.New("Execution not found"))
		return
	}

	respond.With(w, r, http.StatusOK, e)
}

func (s *Server) deleteExecution(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	e := Execution{ID: id}
	db := ExecutionStore{s.dbStore}

	err = db.DeleteExecution(e)
	if err != nil {
		respond.With(w, r, http.StatusNotFound, errors.New("Execution not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listExecutions(w http.ResponseWriter, r *http.Request) {
	lastId, err := strconv.Atoi(r.FormValue("lastId"))
	if err != nil {
		lastId = maxSerial
	}

	db := ExecutionStore{s.dbStore}

	executions, err := db.ListExecutions(s.perPage, lastId)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, executions)
}
