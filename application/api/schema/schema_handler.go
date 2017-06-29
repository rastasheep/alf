package schema

import (
	"errors"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rastasheep/alf/respond"
)

type SchemaHandler struct {
	logger *log.Logger
	store  *SchemaStore
}

func NewSchemaHandler(logger *log.Logger, db *sqlx.DB) *SchemaHandler {
	return &SchemaHandler{
		logger: logger,
		store:  &SchemaStore{db},
	}
}

func (h *SchemaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" || r.URL.Path != "" {
		respond.With(w, r, http.StatusNotFound, errors.New("not found"))
		return
	}

	schema, err := h.store.GetSchema()
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, schema)
}
