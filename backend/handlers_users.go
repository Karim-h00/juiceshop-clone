package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

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
