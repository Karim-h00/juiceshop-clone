package main

import (
	"net/http"
)

func (cfg *config) handlerAdminTest(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, "Welcome")
}
