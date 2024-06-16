package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ggrangel/go-webserver/auth"
)

func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete chirp")
	chirpIdString := r.PathValue("chirpID")
	chirpId, err := strconv.Atoi(chirpIdString)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	claims, err := auth.ParseToken(apiCfg.jwtSecret, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	userId, err := claims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	chirp, err := apiCfg.DB.GetChirp(chirpId)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	if chirp.AuthorId != userIdInt {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = apiCfg.DB.DeleteChirp(chirpId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
