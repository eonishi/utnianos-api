package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"utnianos-api/internal/models"
	"utnianos-api/internal/scraper"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	scraper *scraper.Scraper
}

func New() *Handler {
	return &Handler{
		scraper: scraper.New(),
	}
}

func (h *Handler) GetForums(w http.ResponseWriter, r *http.Request) {
	result, err := h.scraper.GetForums()
	if err != nil {
		h.sendError(w, r, http.StatusInternalServerError, "Error fetching forums", err.Error())
		return
	}

	h.sendJSON(w, r, http.StatusOK, result)
}

func (h *Handler) GetForum(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		h.sendError(w, r, http.StatusBadRequest, "Missing forum slug", "")
		return
	}

	filters := make(map[string][]string)
	for key, values := range r.URL.Query() {
		switch key {
		case "aporte":
			filters["filtertf_tipo_aporte[]"] = values
		case "materia":
			filters["filtertf_materia[]"] = values
		default:
			filters[key] = values
		}
	}

	result, err := h.scraper.GetForum(slug, filters)
	if err != nil {
		h.sendError(w, r, http.StatusInternalServerError, "Error fetching forum", err.Error())
		return
	}

	h.sendJSON(w, r, http.StatusOK, result)
}

func (h *Handler) GetTopic(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		h.sendError(w, r, http.StatusBadRequest, "Missing topic slug", "")
		return
	}

	result, err := h.scraper.GetTopic(slug)
	if err != nil {
		h.sendError(w, r, http.StatusNotFound, "Topic not found", err.Error())
		return
	}

	h.sendJSON(w, r, http.StatusOK, result)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		query = r.URL.Query().Get("keywords")
	}

	action := r.URL.Query().Get("action")

	forumID := 0
	if fids := r.URL.Query().Get("fids"); fids != "" {
		forumID, _ = strconv.Atoi(fids)
	}

	result, err := h.scraper.Search(query, forumID, action)
	if err != nil {
		h.sendError(w, r, http.StatusInternalServerError, "Error searching", err.Error())
		return
	}

	h.sendJSON(w, r, http.StatusOK, result)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		h.sendError(w, r, http.StatusBadRequest, "Missing username", "")
		return
	}

	result, err := h.scraper.GetUser(username)
	if err != nil {
		h.sendError(w, r, http.StatusNotFound, "User not found", err.Error())
		return
	}

	h.sendJSON(w, r, http.StatusOK, result)
}

func (h *Handler) sendJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(w, `{"error": "Encoding error"}`)
	}
}

func (h *Handler) sendError(w http.ResponseWriter, r *http.Request, status int, err, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)

	response := models.ErrorResponse{
		Error:   err,
		Message: details,
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(response)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)
}
