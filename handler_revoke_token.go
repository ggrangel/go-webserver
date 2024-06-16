package main

import (
	"net/http"

	"github.com/ggrangel/go-webserver/auth"
)

func (apiCfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	user, err := apiCfg.DB.GetUserByRefreshToken(token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = apiCfg.DB.RevokeUserRefreshToken(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
