package handlers

import (
	"crmcandidate/models"
	"crmcandidate/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CandidateHandler struct {
	service *services.CandidateService
}

func NewCandidateHandler(service *services.CandidateService) *CandidateHandler {
	return &CandidateHandler{service: service}
}

func (h *CandidateHandler) CreateCandidate(w http.ResponseWriter, r *http.Request) {
	var candidate models.Candidate
	if err := json.NewDecoder(r.Body).Decode(&candidate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCandidate(&candidate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(candidate)
}

func (h *CandidateHandler) GetCandidate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	candidate, err := h.service.GetCandidate(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(candidate)
}
