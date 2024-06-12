package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var profaneWords = [3]string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnErr struct {
		Error string `json:"error"`
	}

	type returValid struct {
		Valid bool `json:"valid"`
	}

	type returnCleanedBody struct {
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(returnErr{Error: "Chirp is too long"})
	}

	replaceProfanity(&params.Body)

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(returnCleanedBody{Body: params.Body})
}

func replaceProfanity(body *string) {
	for _, word := range profaneWords {
		index := strings.Index(strings.ToLower(*body), word)
		if index != -1 {
			*body = (*body)[:index] + "****" + (*body)[index+len(word):]
		}
	}
}
