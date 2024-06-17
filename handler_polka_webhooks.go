package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ggrangel/go-webserver/auth"
)

func (apiCfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("request.Event: ", request.Event)

	if request.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	polkaApiKey, err := auth.GetPolkaApiKey(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("polkaApiKey: ", polkaApiKey)
	fmt.Println("apiCfg.polkaApiKey: ", apiCfg.polkaApiKey)
	if polkaApiKey != apiCfg.polkaApiKey {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	user, err := apiCfg.DB.GetUserById(request.Data.UserId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	apiCfg.DB.SetUserToRedMember(user.Id)
	w.WriteHeader(http.StatusNoContent)
}
