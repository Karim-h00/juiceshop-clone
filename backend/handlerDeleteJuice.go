package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *config) handlerDeleteJuice(w http.ResponseWriter, r *http.Request) {

	juiceID := r.PathValue("juiceID")
	parsedID, err := uuid.Parse(juiceID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	err = cfg.queries.DeleteJuice(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, 500, "could not delete juice")
		return
	}
	w.WriteHeader(204)
}
