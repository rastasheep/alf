package templates

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rastasheep/alf/respond"
)

const (
	maxSerial = 2147483647
)

type TemplateHandler struct {
	Store    *TemplateStore
	Logger   *log.Logger
	PageSize int
}

func (h *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.HandlerFunc
	r.URL.Path = strings.Trim(r.URL.Path, "/")
	id, _ := strconv.Atoi(r.URL.Path)
	ctx := context.WithValue(r.Context(), "templateId", id)

	switch {
	case r.URL.Path == "" && r.Method == "GET":
		handler = h.listTemplates

	case r.URL.Path == "" && r.Method == "POST":
		handler = h.createTemplate

	case id != 0 && r.Method == "GET":
		handler = h.getTemplate

	case id != 0 && r.Method == "PUT":
		handler = h.updateTemplate

	case id != 0 && r.Method == "DELETE":
		handler = h.deleteTemplate

	default:
		respond.With(w, r, http.StatusNotFound, errors.New("not found"))
		return
	}

	http.HandlerFunc(handler).ServeHTTP(w, r.WithContext(ctx))
}

func (h *TemplateHandler) listTemplates(w http.ResponseWriter, r *http.Request) {
	lastId, err := strconv.Atoi(r.FormValue("lastId"))
	if err != nil {
		lastId = maxSerial
	}

	t, err := h.Store.ListTemplates(h.PageSize, lastId)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, t)
}

func (h *TemplateHandler) getTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("templateId").(int)

	t, err := h.Store.GetTemplate(id)
	if err != nil {
		respond.With(w, r, http.StatusNotFound, errors.New("template not found"))
		return
	}

	respond.With(w, r, http.StatusOK, t)
}

func (h *TemplateHandler) createTemplate(w http.ResponseWriter, r *http.Request) {
	t := &Template{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(t); err != nil {
		respond.With(w, r, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}

	if err := t.Valid(); err != nil {
		respond.With(w, r, http.StatusBadRequest, fmt.Errorf("invalid request payload: %v", err))
		return
	}

	t, err := h.Store.CreateTemplate(t)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusCreated, t)
}

func (h *TemplateHandler) updateTemplate(w http.ResponseWriter, r *http.Request) {
	t := &Template{}
	id := r.Context().Value("templateId").(int)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(t); err != nil {
		respond.With(w, r, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	if err := t.Valid(); err != nil {
		respond.With(w, r, http.StatusBadRequest, fmt.Errorf("invalid request payload: %v", err))
		return
	}

	t.ID = id
	t, err := h.Store.UpdateTemplate(t)
	if err != nil {
		respond.With(w, r, http.StatusNotFound, errors.New("template not found"))
		return
	}

	respond.With(w, r, http.StatusCreated, t)
}

func (h *TemplateHandler) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("templateId").(int)

	if err := h.Store.DeleteTemplate(id); err != nil {
		respond.With(w, r, http.StatusNotFound, errors.New("template not found"))
		return
	}

	respond.With(w, r, http.StatusNoContent, nil)
}
