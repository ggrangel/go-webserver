package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ggrangel/go-webserver/database"
)

func handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.NewDB("database.json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Creating user with email: ", request.Email)

	user, err := db.CreateUser(request.Email)
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
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
