package main

import (
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

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
