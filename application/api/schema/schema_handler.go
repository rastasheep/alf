package schema

import (
	"errors"
	"log"
	"net/http"

	"github.com/rastasheep/alf/respond"
)

type SchemaHandler struct {
	Logger *log.Logger
	Store  *SchemaStore
}

func (h *SchemaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" || r.URL.Path != "" {
		respond.With(w, r, http.StatusNotFound, errors.New("not found"))
		return
	}

	schema, err := h.Store.GetSchema()
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, schema)
}
