package execution

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rastasheep/alf/respond"
)

const (
	maxSerial = 2147483647
)

type ExecutionHandler struct {
	logger  *log.Logger
	store   *ExecutionStore
	perPage int
}

func NewExecutionHandler(logger *log.Logger, db *sqlx.DB, perPage int) *ExecutionHandler {
	return &ExecutionHandler{
		logger:  logger,
		store:   &ExecutionStore{db},
		perPage: perPage,
	}
}

func (h *ExecutionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.HandlerFunc
	r.URL.Path = strings.Trim(r.URL.Path, "/")
	id, _ := strconv.Atoi(r.URL.Path)
	ctx := context.WithValue(r.Context(), "executionId", id)

	switch {
	case r.URL.Path == "" && r.Method == "GET":
		handler = h.listExecutions

	case r.URL.Path == "" && r.Method == "POST":
		handler = h.createExecution

	case id != 0 && r.Method == "GET":
		handler = h.getExecution

	case id != 0 && r.Method == "DELETE":
		handler = h.deleteExecution

	default:
		respond.With(w, r, http.StatusNotFound, errors.New("not found"))
		return
	}

	http.HandlerFunc(handler).ServeHTTP(w, r.WithContext(ctx))
}

func (h *ExecutionHandler) listExecutions(w http.ResponseWriter, r *http.Request) {
	lastId, err := strconv.Atoi(r.FormValue("lastId"))
	if err != nil {
		lastId = maxSerial
	}

	executions, err := h.store.ListExecutions(h.perPage, lastId)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, executions)
}

func (h *ExecutionHandler) createExecution(w http.ResponseWriter, r *http.Request) {
	var e Execution
	logPrefix := r.Context().Value("logPrefix").(string)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&e); err != nil {
		respond.With(w, r, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	re := regexp.MustCompile("(?i)(SET.*TRANSACTION)|(SET.*SESSION.*CHARACTERISTICS)")
	matched := re.MatchString(e.Query)
	if matched {
		h.logger.Printf("%s blocked execution of query: %s", logPrefix, e.Query)
		respond.With(w, r, http.StatusBadRequest, errors.New("you are not allowed to change the characteristics of transaction"))
		return
	}

	e, err := h.store.CreateExecution(e)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	// go s.resultsCache.GetResults(strconv.Itoa(e.ID), nil)

	respond.With(w, r, http.StatusCreated, e)
}

func (h *ExecutionHandler) getExecution(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("executionId").(int)
	e := Execution{ID: id}

	e, err := h.store.GetExecution(e)
	if err != nil {
		respond.With(w, r, http.StatusNotFound, errors.New("execution not found"))
		return
	}

	respond.With(w, r, http.StatusOK, e)
}

func (h *ExecutionHandler) deleteExecution(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("executionId").(int)
	e := Execution{ID: id}

	if err := h.store.DeleteExecution(e); err != nil {
		respond.With(w, r, http.StatusNotFound, errors.New("execution not found"))
		return
	}

	respond.With(w, r, http.StatusNoContent, nil)
}
