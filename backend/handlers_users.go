package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

func (cfg *config) handlerGetAllUsers(w http.ResponseWriter, r *http.Request) {
	type userResponse struct {
		ID        uuid.UUID `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	data, err := cfg.queries.GetAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "no users found")
		return
	}

	users := make([]userResponse, len(data))
	for i, u := range data {
		users[i] = userResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
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
		respondWithError(w, 400, "Error decoding params")
		return
	}
	user, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Username: params.Username,
		Email:    params.Email,
		ID:       userID,
	})
	if err != nil {
		respondWithError(w, 500, "Error creating user")
		return
	}
	respondWithJSON(w, 200, user)
}

func (cfg *config) handlerAdminUpdate(w http.ResponseWriter, r *http.Request) {
	type Body_params struct {
		Role string `json:"role"`
	}
	role := r.Context().Value(contextKeyRole).(string)
	if role != "admin" {
		respondWithError(w, 403, "Forbidden")
		return
	}

	callerID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	userID := r.PathValue("userID")
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := Body_params{}

	err = decoder.Decode(&params)
	if err != nil {
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
		respondWithError(w, http.StatusInternalServerError, "couldn't update user")
		return
	}
	w.WriteHeader(200)
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
		respondWithError(w, 400, "Error decoding params")
		return
	}

	db_password, err := cfg.queries.GetPasswordByUserID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 400, "something went wrong")
		return
	}

	_, err = auth.CheckPasswordHash(params.Password, db_password)
	if err != nil {
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
		respondWithError(w, 500, "could not hash password")
		return
	}
	err = cfg.queries.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, 500, "could not update password")
		return
	}
	respondWithJSON(w, 200, struct{}{})
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
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	w.WriteHeader(204)
}
