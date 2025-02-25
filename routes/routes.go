package routes

import (
	"crmcandidate/handlers"
	"crmcandidate/services"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router *chi.Mux, candidateService *services.CandidateService) {
	candidateHandler := handlers.NewCandidateHandler(candidateService)

	router.Route("/candidates", func(r chi.Router) {
		r.Post("/", candidateHandler.CreateCandidate)
		r.Get("/{id}", candidateHandler.GetCandidate)
	})
}
