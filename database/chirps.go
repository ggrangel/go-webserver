package database

import (
	"fmt"
	"strings"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDb()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	fmt.Println("here")
	dbStructure, err := db.loadDb()
	fmt.Println(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, fmt.Errorf("Chirp not found")
	}

	return chirp, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	if len(body) > 140 {
		return Chirp{}, fmt.Errorf("Chirp is too long")
	}

	replaceProfanity(&body)

	dbStructure, err := db.loadDb()
	if err != nil {
		return Chirp{}, err
	}

	nextKey := getNextChirpKey(&dbStructure.Chirps)

	dbStructure.addChirpToStorage(nextKey, body)

	err = db.writeDb(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return dbStructure.Chirps[nextKey], nil
}

func (dbStructure *DbStructure) addChirpToStorage(id int, body string) {
	chirp := Chirp{Id: id, Body: body}
	dbStructure.Chirps[id] = chirp
}

func getNextChirpKey(chirps *map[int]Chirp) int {
	lastKey := 0
	for key := range *chirps {
		if key > lastKey {
			lastKey = key
		}
	}
	return lastKey + 1
}

func replaceProfanity(body *string) {
	var profaneWords = [3]string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	for _, word := range profaneWords {
		index := strings.Index(strings.ToLower(*body), word)
		if index != -1 {
			*body = (*body)[:index] + "****" + (*body)[index+len(word):]
		}
	}
}
