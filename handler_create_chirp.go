package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ggrangel/go-webserver/database"
)

func handlerCreateChirp(dbPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Body string `json:"body"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		db, err := database.NewDB(dbPath)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error creating databse")
			return
		}

		chirp, err := db.CreateChirp(request.Body)
		if err != nil {
			fmt.Println("Error creating chirp:", err)
		}

		response, err := json.Marshal(chirp)
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}
