package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

type update_password struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

func (cfg *config) handlerUpdatePassword(w http.ResponseWriter, r *http.Request) {
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
