package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ggrangel/go-webserver/auth"
)

func (apiCfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	claims, err := auth.ParseToken(apiCfg.jwtSecret, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	encryptedPassword, err := auth.HashPassword(request.Password)
	user, err := apiCfg.DB.UpdateUser(id, request.Email, encryptedPassword)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
