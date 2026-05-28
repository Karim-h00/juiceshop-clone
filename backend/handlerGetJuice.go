package main

import (
	"net/http"

	"github.com/karim-h00/juiceshop-clone/internal/database"
)

func (cfg *config) handlerGetJuice(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query().Get("q")
	var data = []database.Juice{}
	var err error

	if q != "" {
		data, err = cfg.queries.GetJuiceByName(r.Context(), "%"+q+"%")
		if err != nil {
			respondWithError(w, 500, "Error retrieving Juices")
			return
		}
	} else {
		data, err = cfg.queries.GetAllJuice(r.Context())
		if err != nil {
			respondWithError(w, 500, "Error retrieving Juices")
			return
		}
	}
	respondWithJSON(w, 200, data)

}

func (cfg *config) handlerGetJuiceByName(w http.ResponseWriter, r *http.Request) {
	juiceName := r.PathValue("juiceName")

	data, err := cfg.queries.GetJuiceByName(r.Context(), juiceName)
	if err != nil {
		respondWithError(w, 500, "Error retrieving Juice")
		return
	}
	respondWithJSON(w, 200, data)
}
