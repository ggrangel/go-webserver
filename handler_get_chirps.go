package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ggrangel/go-webserver/database"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	path_id := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	authorIdString := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")

	if path_id == "" && authorIdString == "" {
		apiCfg.handlerGetAllChirps(w, sort)
		return
	}
	if authorIdString != "" {
		apiCfg.handlerGetChirpsByAuthorId(w, authorIdString, sort)
	}

	id, err := strconv.Atoi(path_id)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	apiCfg.handlerGetChirpById(w, id, sort)
}

func (apiCfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, id int, sort string) {
	chirp, err := apiCfg.DB.GetChirp(id)
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

func (apiCfg *apiConfig) handlerGetChirpsByAuthorId(
	w http.ResponseWriter,
	authorIdString string,
) {
	authorId, err := strconv.Atoi(authorIdString)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	chirps, err := apiCfg.DB.GetChirpsByAuthorId(authorId)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(chirps)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func (apiCfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter) {
	chirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(chirps)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func possiblySortChirps(chirps []database.Chirp, sort string) []database.Chirp {
	if sort == "asc" {
		return sortChirpsAsc(chirps)
	}
	if sort == "desc" {
		return sortChirpsDesc(chirps)
	}
	return chirps
}

func sortChirpsAsc(chirps []database.Chirp) []database.Chirp {
	for i := 0; i < len(chirps); i++ {
		for j := 0; j < len(chirps); j++ {
			if chirps[i].Id < chirps[j].Id {
				chirps[i], chirps[j] = chirps[j], chirps[i]
			}
		}
	}
	return chirps
}

func sortChirpsDesc(chirps []database.Chirp) []database.Chirp {
	for i := 0; i < len(chirps); i++ {
		for j := 0; j < len(chirps); j++ {
			if chirps[i].Id > chirps[j].Id {
				chirps[i], chirps[j] = chirps[j], chirps[i]
			}
		}
	}
	return chirps
}
