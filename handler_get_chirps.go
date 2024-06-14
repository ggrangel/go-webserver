package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ggrangel/go-webserver/database"
)

func handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	log.Printf("here")

	path_id := strings.TrimPrefix(r.URL.Path, "/api/chirps/")

	log.Printf("Path ID: %s", path_id)

	if path_id == "" {
		handlerGetAllChirps(w, r)
		return
	}

	id, err := strconv.Atoi(path_id)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	handlerGetChirpById(w, r, id)
}

func handlerGetChirpById(w http.ResponseWriter, r *http.Request, id int) {
	db, err := database.NewDB("database.json")

	fmt.Println("ID: ", id)

	if err != nil {
		w.WriteHeader(500)
		return
	}

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

func handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
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
