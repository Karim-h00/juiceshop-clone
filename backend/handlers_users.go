package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
	"github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Token     string    `json:"token"`
}

func (cfg *config) handlerGetAllUsers(w http.ResponseWriter, r *http.Request) {
	type userResponse struct {
		ID        uuid.UUID `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	page := 1

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			respondWithError(w, 400, "invalid page number")
			return
		}
		page = parsedPage
	}
	offset := (page - 1) * 10

	search := r.URL.Query().Get("q")
	var users []userResponse

	if search == "" {
		data, err := cfg.queries.GetAllUsers(r.Context(), database.GetAllUsersParams{
			Limit:  50,
			Offset: int32(offset),
		})
		if err != nil {
			cfg.logger.Error("get users", "error", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
			return
		}
		users = make([]userResponse, len(data))
		for i, u := range data {
			users[i] = userResponse{
				ID:        u.ID,
				Username:  u.Username,
				Email:     u.Email,
				Role:      u.Role,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			}
		}
	} else {
		data, err := cfg.queries.SearchUsers(r.Context(), database.SearchUsersParams{
			Column1: sql.NullString{String: search, Valid: true},
			Limit:   50,
			Offset:  int32(offset),
		})
		if err != nil {
			cfg.logger.Error("search users", "error", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
			return
		}
		users = make([]userResponse, len(data))
		for i, u := range data {
			users[i] = userResponse{
				ID:        u.ID,
				Username:  u.Username,
				Email:     u.Email,
				Role:      u.Role,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			}
		}
	}

	respondWithJSON(w, 200, users)
}

func (cfg *config) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type User_parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := User_parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding params")
		return
	}
	if params.Password == "" {
		respondWithError(w, 400, "password is required")
		return
	}
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error creating user")
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		cfg.logger.Error("create user", "error", err)
		respondWithError(w, 500, "Error creating user")
		return
	}
	respondWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Username:  user.Username,
	})
}

func (cfg *config) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type user_update_params struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	decoder := json.NewDecoder(r.Body)
	params := user_update_params{}
	err := decoder.Decode(&params)
	if err != nil {
		cfg.logger.Warn("update user", "user_id", userID, "reason", "failed to decode body", "error", err, "ip", getClientIP(r))
		respondWithError(w, 400, "Error decoding params")
		return
	}
	user, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Username: params.Username,
		Email:    params.Email,
		ID:       userID,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			respondWithError(w, http.StatusConflict, "username or email already taken")
			return
		}
		cfg.logger.Error("update user", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *config) handlerAdminUpdate(w http.ResponseWriter, r *http.Request) {
	type Body_params struct {
		Role string `json:"role"`
	}
	role := r.Context().Value(contextKeyRole).(string)
	if role != "admin" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	callerID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	userID := r.PathValue("userID")
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		cfg.logger.Warn("admin update user", "user_id", callerID, "reason", "failed to parse user id", "error", err, "ip", getClientIP(r))
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := Body_params{}

	err = decoder.Decode(&params)
	if err != nil {
		cfg.logger.Warn("admin update user", "user_id", callerID, "reason", "failed to decode body", "error", err, "ip", getClientIP(r))
		respondWithError(w, 400, "Error decoding params")
		return
	}
	if params.Role == "" {
		respondWithError(w, 400, "role is required")
		return
	}
	if params.Role != "admin" && params.Role != "user" {
		respondWithError(w, 400, "incorrect role")
		return
	}
	if callerID == parsedID {
		respondWithError(w, 403, "cannot change your own role")
		return
	}
	err = cfg.queries.UpdateUserRole(r.Context(), database.UpdateUserRoleParams{
		Role: params.Role,
		ID:   parsedID,
	})
	if err != nil {
		cfg.logger.Error("admin update user", "error", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't update user")
		return
	}
	err = cfg.queries.AddLog(r.Context(), database.AddLogParams{
		UserID:     uuid.NullUUID{UUID: callerID, Valid: true},
		Action:     "update_role",
		TargetType: "user",
		TargetID:   uuid.NullUUID{UUID: parsedID, Valid: true},
		TargetName: sql.NullString{String: params.Role, Valid: true},
		CreatedAt:  time.Now().UTC(),
	})
	if err != nil {
		cfg.logger.Error("add audit log", "error", err)
	}
	cfg.logger.Info("admin update user", "admin_id", callerID, "target_id", parsedID, "ip", getClientIP(r))
	w.WriteHeader(http.StatusOK)
}

func (cfg *config) handlerUpdatePassword(w http.ResponseWriter, r *http.Request) {
	type update_password struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}
	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	decoder := json.NewDecoder(r.Body)
	params := update_password{}
	err := decoder.Decode(&params)
	if err != nil {
		cfg.logger.Warn("update user password", "user_id", userID, "reason", "failed to decode body", "error", err, "ip", getClientIP(r))
		respondWithError(w, 400, "Error decoding params")
		return
	}

	db_password, err := cfg.queries.GetPasswordByUserID(r.Context(), userID)
	if err != nil {
		cfg.logger.Error("update password", "user_id", userID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	_, err = auth.CheckPasswordHash(params.Password, db_password)
	if err != nil {
		cfg.logger.Warn("update password", "user_id", userID, "reason", "wrong password", "ip", getClientIP(r))
		respondWithError(w, 401, "wrong password")
		return
	}
	if params.NewPassword == "" {
		respondWithError(w, 400, "new password is required")
		return
	}
	if params.NewPassword == params.Password {
		respondWithError(w, 400, "passwords must be different")
		return
	}
	hashedPassword, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}
	err = cfg.queries.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		cfg.logger.Error("update user password", "error", err)
		respondWithError(w, http.StatusInternalServerError, "could not update password")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (cfg *config) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {

	role := r.Context().Value(contextKeyRole).(string)
	if role != "admin" {
		respondWithError(w, 403, "Forbidden")
		return
	}

	userID := r.PathValue("userID")
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}
	callerID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	if callerID == parsedID {
		respondWithError(w, 403, "cannot delete your own account")
		return
	}

	err = cfg.queries.DeleteUserByID(r.Context(), parsedID)
	if err != nil {
		cfg.logger.Error("delete user", "admin_id", callerID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	cfg.queries.AddLog(r.Context(), database.AddLogParams{
		UserID:     uuid.NullUUID{UUID: callerID, Valid: true},
		Action:     "delete",
		TargetType: "user",
		TargetID:   uuid.NullUUID{UUID: parsedID, Valid: true},
		TargetName: sql.NullString{Valid: false},
		CreatedAt:  time.Now().UTC(),
	})
	w.WriteHeader(204)
}
