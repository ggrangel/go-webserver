package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ggrangel/go-webserver/auth"
	"github.com/ggrangel/go-webserver/database"
)

type LoginResponse struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := apiCfg.DB.GetUser(request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var expiresIn time.Duration
	expiresIn = time.Hour * 1
	// if request.Expires_in_seconds == 0 || request.Expires_in_seconds > 3600 {
	// 	expiresIn = time.Hour * 1
	// } else {
	// 	expiresIn = time.Second * time.Duration(request.Expires_in_seconds)
	// }

	token, err := auth.GenerateToken(apiCfg.jwtSecret, user.Id, expiresIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Token = token

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sixtyDays := time.Hour * 24 * 60
	user.RefreshToken = database.RefreshToken{
		Token:  refreshToken,
		Expiry: int(time.Now().Add(sixtyDays).Unix()),
	}

	err = apiCfg.DB.SaveUserRefreshToken(user.Id, user.RefreshToken.Token, user.RefreshToken.Expiry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResponse := LoginResponse{
		Id:           user.Id,
		Email:        user.Email,
		Token:        user.Token,
		RefreshToken: user.RefreshToken.Token,
	}

	json.NewEncoder(w).Encode(loginResponse)
	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
}
