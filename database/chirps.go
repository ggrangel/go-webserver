package database

import (
	"fmt"
	"strings"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
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

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	if len(body) > 140 {
		return Chirp{}, fmt.Errorf("Chirp is too long")
	}

	replaceProfanity(&body)

	dbStructure, err := db.loadDb()
	if err != nil {
		return Chirp{}, err
	}

	nextKey := len(dbStructure.Chirps) + 1

	chirp := Chirp{Id: nextKey, Body: body, AuthorId: authorId}
	dbStructure.Chirps[nextKey] = chirp

	err = db.writeDb(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return dbStructure.Chirps[nextKey], nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDb()
	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, id)

	err = db.writeDb(dbStructure)
	if err != nil {
		return err
	}

	return nil
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
