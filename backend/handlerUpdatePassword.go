package main

import (
	"encoding/json"
	"net/http"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

type update_password struct {
	Password     string `json:"password"`
	New_password string `json:"new_password"`
}

func (cfg *config) handlerUpdatePassword(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userID, _, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 400, "Could not make session")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := update_password{}
	err = decoder.Decode(&params)
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

	hashedPassword, err := auth.HashPassword(params.New_password)
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
