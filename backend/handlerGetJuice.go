package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *config) handlerGetJuice(w http.ResponseWriter, r *http.Request) {

	data, err := cfg.queries.GetAllJuice(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error retrieving chirp")
		return
	}
	respondWithJSON(w, 200, data)

}

func (cfg *config) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	juiceID := r.PathValue("juiceID")
	parsedID, err := uuid.Parse(juiceID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}
	data, err := cfg.queries.GetJuiceByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, 500, "Error retrieving chirp")
		return
	}
	respondWithJSON(w, 200, data)
}
