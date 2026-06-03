package main

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/google/uuid"
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

func (cfg *config) handlerAddJuice(w http.ResponseWriter, r *http.Request) {
	type juice_params struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int    `json:"price"`
		Stock       int    `json:"stock"`
	}

	decoder := json.NewDecoder(r.Body)
	params := juice_params{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding params")
		return
	}
	if params.Name == "" {
		respondWithError(w, 400, "Name is required")
		return
	}
	if params.Price <= 0 {
		respondWithError(w, 400, "price must be positive")
		return
	}
	if params.Stock < 0 {
		respondWithError(w, 400, "stock must not be negative")
		return
	}
	juice, err := cfg.queries.AddJuice(r.Context(), database.AddJuiceParams{
		Name:        params.Name,
		Description: params.Description,
		Price:       int32(params.Price),
		Stock:       int32(params.Stock),
	})
	if err != nil {
		respondWithError(w, 500, "Error creating juice")
		return
	}
	respondWithJSON(w, 201, juice)
}

func (cfg *config) handlerUpdateJuice(w http.ResponseWriter, r *http.Request) {
	type update_juice_params struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int    `json:"price"`
		Stock       int    `json:"stock"`
	}

	juiceID := r.PathValue("juiceID")
	parsedID, err := uuid.Parse(juiceID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := update_juice_params{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding params")
		return
	}

	if params.Name == "" {
		respondWithError(w, 400, "Name is required")
		return
	}
	if params.Price <= 0 {
		respondWithError(w, 400, "price must be positive")
		return
	}
	if params.Stock < 0 {
		respondWithError(w, 400, "stock must not be negative")
		return
	}
	juice, err := cfg.queries.UpdateJuice(r.Context(), database.UpdateJuiceParams{
		ID:          parsedID,
		Name:        params.Name,
		Description: params.Description,
		Price:       int32(params.Price),
		Stock:       int32(params.Stock),
	})
	if err != nil {
		respondWithError(w, 500, "Error updating juice")
		return
	}
	respondWithJSON(w, 200, juice)
}

func (cfg *config) handlerUpdateJuiceImage(w http.ResponseWriter, r *http.Request) {
	juiceID := r.PathValue("juiceID")
	parsedID, err := uuid.Parse(juiceID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	_, err = cfg.queries.GetJuiceByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, 404, "juice doesn't exist")
		return
	}

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("juiceImage")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file")
		return
	}
	defer file.Close()
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type")
		return
	}
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type")
		return
	}

	assetPath := getAssetPath(mediaType)
	assetDiskPath := cfg.getAssetDiskPath(assetPath)

	dst, err := os.Create(assetDiskPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving file")
		return
	}

	url := cfg.getAssetURL(assetPath)

	err = cfg.queries.UpdateJuiceImage(r.Context(), database.UpdateJuiceImageParams{
		ImageUrl: url,
		ID:       parsedID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update juice")
		return
	}
	w.WriteHeader(http.StatusOK)
}
