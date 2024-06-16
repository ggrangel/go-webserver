package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ggrangel/go-webserver/database"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	path_id := strings.TrimPrefix(r.URL.Path, "/api/chirps/")

	if path_id == "" {
		handlerGetAllChirps(w, apiCfg.DB)
		return
	}

	id, err := strconv.Atoi(path_id)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	handlerGetChirpById(w, apiCfg.DB, id)
}

func handlerGetChirpById(w http.ResponseWriter, db *database.DB, id int) {
	chirp, err := db.GetChirp(id)

	fmt.Println(chirp)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	response, err := json.Marshal(chirp)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func handlerGetAllChirps(w http.ResponseWriter, db *database.DB) {
	db, err := database.NewDB("database.json")

	if err != nil {
		w.WriteHeader(500)
		return
	}

	chirps, err := db.GetChirps()

	response, err := json.Marshal(chirps)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}
