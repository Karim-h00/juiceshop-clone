package main

import (
	"net/http"
)

func (cfg *config) handlerGetJuice(w http.ResponseWriter, r *http.Request) {

	data, err := cfg.queries.GetAllJuice(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error retrieving chirp")
		return
	}
	respondWithJSON(w, 200, data)

}
