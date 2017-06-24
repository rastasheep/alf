package main

import (
	"errors"
	"net/http"

	"github.com/rastasheep/alf/respond"
)

func (s *Server) listResults(w http.ResponseWriter, r *http.Request) {
	executionId := r.FormValue("executionId")
	if executionId == "" {
		respond.With(w, r, http.StatusInternalServerError, errors.New("Execution not found"))
		return
	}

	var results []map[string]interface{}

	err := s.resultsCache.GetResults(executionId, &results)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, results)
}
