package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

func (cfg *config) handlerGetAuditLogs(w http.ResponseWriter, r *http.Request) {

	type auditLogResponse struct {
		ID         uuid.UUID  `json:"id"`
		UserID     *uuid.UUID `json:"user_id"`
		Action     string     `json:"action"`
		TargetType string     `json:"target_type"`
		TargetID   *uuid.UUID `json:"target_id"`
		TargetName *string    `json:"target_name"`
		CreatedAt  time.Time  `json:"created_at"`
	}

	page := 1

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			cfg.logger.Warn("invalid page number", "page", pageStr, "ip", getClientIP(r))
			respondWithError(w, 400, "invalid page number")
			return
		}
		page = parsedPage
	}
	offset := (page - 1) * 10

	data, err := cfg.queries.GetAuditLogs(r.Context(), database.GetAuditLogsParams{
		Limit:  20,
		Offset: int32(offset),
	})
	if err != nil {
		cfg.logger.Error("get audit logs", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get audit logs")
		return
	}

	logs := make([]auditLogResponse, len(data))
	for i, log := range data {
		var userID *uuid.UUID
		if log.UserID.Valid {
			userID = &log.UserID.UUID
		}
		var targetID *uuid.UUID
		if log.TargetID.Valid {
			targetID = &log.TargetID.UUID
		}
		var targetName *string
		if log.TargetName.Valid {
			targetName = &log.TargetName.String
		}
		logs[i] = auditLogResponse{
			ID:         log.ID,
			UserID:     userID,
			Action:     log.Action,
			TargetType: log.TargetType,
			TargetID:   targetID,
			TargetName: targetName,
			CreatedAt:  log.CreatedAt,
		}
	}

	respondWithJSON(w, http.StatusOK, logs)
}
