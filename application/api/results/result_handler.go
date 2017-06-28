package results

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/rastasheep/alf/respond"
)

type ResultHandler struct {
	Logger      *log.Logger
	ResultCache *ResultCache
}

func (h *ResultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" || r.URL.Path != "" {
		respond.With(w, r, http.StatusNotFound, errors.New("not found"))
		return
	}

	executionId := r.FormValue("executionId")
	if _, err := strconv.Atoi(executionId); err != nil {
		respond.With(w, r, http.StatusBadRequest, errors.New("invalid execution id"))
		return
	}

	var results []map[string]interface{}

	err := h.ResultCache.GetResults(executionId, &results)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, results)
}
