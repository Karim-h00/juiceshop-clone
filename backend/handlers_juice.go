package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

type JuiceResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        int32     `json:"price"`
	ImageUrl     string    `json:"image_url"`
	Stock        int32     `json:"stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	AvgRating    float64   `json:"avg_rating"`
	ReviewsCount int64     `json:"reviews_count"`
}

func (cfg *config) handlerGetJuice(w http.ResponseWriter, r *http.Request) {

	rows, err := cfg.queries.GetAllJuice(r.Context())
	if err != nil {
		cfg.logger.Error("Error retrieving all Juices", "error", err)
		respondWithError(w, http.StatusInternalServerError, "error retrieving juicess")
		return
	}

	data := make([]JuiceResponse, len(rows))
	for i, row := range rows {
		avgRating := 0.0
		if b, ok := row.AvgRating.([]byte); ok {
			parsed, err := strconv.ParseFloat(string(b), 64)
			if err == nil {
				avgRating = parsed
			}
		}

		data[i] = JuiceResponse{
			ID:           row.ID,
			Name:         row.Name,
			Description:  row.Description,
			Price:        row.Price,
			ImageUrl:     row.ImageUrl,
			Stock:        row.Stock,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
			AvgRating:    avgRating,
			ReviewsCount: row.ReviewsCount,
		}
	}
	respondWithJSON(w, http.StatusOK, data)
}

func (cfg *config) handlerGetJuiceByName(w http.ResponseWriter, r *http.Request) {
	juiceName := r.PathValue("juiceName")

	row, err := cfg.queries.GetJuiceDetails(r.Context(), juiceName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.logger.Warn("get juice details", "reason", "juice not found", "error", err, "ip", r.RemoteAddr)
			respondWithError(w, http.StatusNotFound, "Juice not found")
			return
		}
		cfg.logger.Error("get juice details", "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	avgRating := 0.0
	if b, ok := row.AvgRating.([]byte); ok {
		parsed, err := strconv.ParseFloat(string(b), 64)
		if err == nil {
			avgRating = parsed
		} else {
			cfg.logger.Warn("get juice details", "reason", "failed to parse avg rating", "error", err)
		}
	}

	respondWithJSON(w, http.StatusOK, JuiceResponse{
		ID:           row.ID,
		Name:         row.Name,
		Description:  row.Description,
		Price:        row.Price,
		ImageUrl:     row.ImageUrl,
		Stock:        row.Stock,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		AvgRating:    avgRating,
		ReviewsCount: row.ReviewsCount,
	})
}

func (cfg *config) handlerDeleteJuice(w http.ResponseWriter, r *http.Request) {
	adminID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	juiceID := r.PathValue("juiceID")
	parsedID, err := uuid.Parse(juiceID)
	if err != nil {
		cfg.logger.Warn("delete juice", "reason", "failed to parse juice id", "error", err)
		respondWithError(w, http.StatusBadRequest, "failed to parse ID")
		return
	}

	err = cfg.queries.DeleteJuice(r.Context(), parsedID)
	if err != nil {
		cfg.logger.Error("delete juice", "admin_id", adminID, "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "could not delete juice")
		return
	}
	cfg.logger.Info("delete juice", "admin_id", adminID, "ip", r.RemoteAddr)
	w.WriteHeader(204)
}

func (cfg *config) handlerAddJuice(w http.ResponseWriter, r *http.Request) {
	type juice_params struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int    `json:"price"`
		Stock       int    `json:"stock"`
	}

	adminID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	decoder := json.NewDecoder(r.Body)
	params := juice_params{}

	err := decoder.Decode(&params)
	if err != nil {
		cfg.logger.Error("add juice", "admin_id", adminID, "reason", "error decoding params", "error", err)
		respondWithError(w, 400, "Error decoding params")
		return
	}
	if params.Name == "" {
		cfg.logger.Warn("add juice", "reason", "name parameter empty", "ip", r.RemoteAddr)
		respondWithError(w, 400, "Name is required")
		return
	}
	if params.Price <= 0 {
		cfg.logger.Warn("add juice", "reason", "price parameter 0 or negative", "ip", r.RemoteAddr)
		respondWithError(w, 400, "price must be positive")
		return
	}
	if params.Stock < 0 {
		cfg.logger.Warn("add juice", "reason", "stock parameter negative", "ip", r.RemoteAddr)
		respondWithError(w, 400, "stock must not be negative")
		return
	}
	if len(params.Description) > 500 {
		cfg.logger.Warn("add juice", "reason", "description param too big", "ip", r.RemoteAddr)
		respondWithError(w, 400, "description too big")
		return
	}
	juice, err := cfg.queries.AddJuice(r.Context(), database.AddJuiceParams{
		Name:        params.Name,
		Description: params.Description,
		Price:       int32(params.Price),
		Stock:       int32(params.Stock),
	})
	if err != nil {
		cfg.logger.Error("add juice", "admin_id", adminID, "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "Error creating juice")
		return
	}
	cfg.logger.Info("add juice successful", "admin_id", adminID, "ip", r.RemoteAddr, "juice_name", params.Name)
	respondWithJSON(w, 201, juice)
}

func (cfg *config) handlerUpdateJuice(w http.ResponseWriter, r *http.Request) {
	type update_juice_params struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int    `json:"price"`
		Stock       int    `json:"stock"`
	}

	adminID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	juiceID := r.PathValue("juiceID")
	parsedID, err := uuid.Parse(juiceID)
	if err != nil {
		cfg.logger.Error("update juice", "admin_id", adminID, "reason", "error parsing juice id", "error", err)
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := update_juice_params{}

	err = decoder.Decode(&params)
	if err != nil {
		cfg.logger.Error("update juice", "admin_id", adminID, "reason", "error decoding params", "error", err)
		respondWithError(w, 400, "Error decoding params")
		return
	}

	if params.Name == "" {
		cfg.logger.Warn("update juice", "reason", "name parameter empty", "ip", r.RemoteAddr)
		respondWithError(w, 400, "Name is required")
		return
	}
	if params.Price <= 0 {
		cfg.logger.Warn("update juice", "reason", "price parameter 0 or negative", "ip", r.RemoteAddr)
		respondWithError(w, 400, "price must be positive")
		return
	}
	if params.Stock < 0 {
		cfg.logger.Warn("update juice", "reason", "stock parameter negative", "ip", r.RemoteAddr)
		respondWithError(w, 400, "stock must not be negative")
		return
	}
	if len(params.Description) > 500 {
		cfg.logger.Warn("update juice", "reason", "description param too big", "ip", r.RemoteAddr)
		respondWithError(w, 400, "description too big")
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
		cfg.logger.Error("update juice", "admin_id", adminID, "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "Error updating juice")
		return
	}
	cfg.logger.Info("update juice successful", "admin_id", adminID, "ip", r.RemoteAddr, "juice_name", params.Name)
	respondWithJSON(w, http.StatusOK, juice)
}

func (cfg *config) handlerUpdateJuiceImage(w http.ResponseWriter, r *http.Request) {
	adminID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	juiceID := r.PathValue("juiceID")
	parsedID, err := uuid.Parse(juiceID)
	if err != nil {
		cfg.logger.Error("update juice image", "admin_id", adminID, "reason", "error parsing juice id", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	_, err = cfg.queries.GetJuiceByID(r.Context(), parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.logger.Warn("update juice img", "reason", "juice not found", "error", err, "ip", r.RemoteAddr)
			respondWithError(w, http.StatusNotFound, "Juice not found")
			return
		}
		cfg.logger.Error("update juice img", "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("juiceImage")
	if err != nil {
		cfg.logger.Error("update juice img", "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file")
		return
	}
	defer file.Close()
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		cfg.logger.Error("update juice img", "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type")
		return
	}
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		cfg.logger.Warn("update juice img", "reason", "invalid media type", "media_type", mediaType, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusBadRequest, "Invalid file type")
		return
	}

	assetPath := getAssetPath(mediaType)
	assetDiskPath := cfg.getAssetDiskPath(assetPath)

	dst, err := os.Create(assetDiskPath)
	if err != nil {
		cfg.logger.Error("update juice img", "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "Unable to create file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		cfg.logger.Error("update juice img", "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "Error saving file")
		return
	}

	url := cfg.getAssetURL(assetPath)

	err = cfg.queries.UpdateJuiceImage(r.Context(), database.UpdateJuiceImageParams{
		ImageUrl: url,
		ID:       parsedID,
	})
	if err != nil {
		cfg.logger.Error("update juice img", "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update juice")
		return
	}
	cfg.logger.Info("update juice img successful", "admin_id", adminID, "juice_id", juiceID, "ip", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
