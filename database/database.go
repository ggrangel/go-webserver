package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type DbStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func NewDB(path string) (*DB, error) {
	if path == "" {
		path = "database.json"
	}

	if _, err := os.ReadFile(path); err == nil {
		return &DB{path: path, mux: &sync.RWMutex{}}, nil
	}

	jsonData, err := json.MarshalIndent(DbStructure{}, "", "  ")

	if err != nil {
		return nil, err
	}

	err = os.WriteFile(path, jsonData, 0644)

	if err != nil {
		return nil, err
	}

	return &DB{path: path, mux: &sync.RWMutex{}}, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	if len(body) > 140 {
		return Chirp{}, fmt.Errorf("Chirp is too long")
	}

	replaceProfanity(&body)

	db.mux.Lock()
	defer db.mux.Unlock()

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

func (db *DB) GetChirp(id int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

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

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

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

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDb()
	if err != nil {
		return User{}, err
	}

	nextKey := getNextUserKey(&dbStructure.Users)

	dbStructure.addUserToStorage(nextKey, email)

	err = db.writeDb(dbStructure)
	if err != nil {
		return User{}, err
	}

	return dbStructure.Users[nextKey], nil
}

func getNextUserKey(users *map[int]User) int {
	lastKey := 0
	for key := range *users {
		if key > lastKey {
			lastKey = key
		}
	}
	return lastKey + 1
}

// Find the highest key to determine the new chirp ID
func getNextChirpKey(chirps *map[int]Chirp) int {
	lastKey := 0
	for key := range *chirps {
		if key > lastKey {
			lastKey = key
		}
	}
	return lastKey + 1
}

func (db *DB) loadDb() (DbStructure, error) {
	file, err := os.ReadFile(db.path)
	if err != nil {
		return DbStructure{}, fmt.Errorf("error reading file: %v", err)
	}

	var dbStructure DbStructure
	if err := json.Unmarshal(file, &dbStructure); err != nil {
		return DbStructure{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	if dbStructure.Chirps == nil {
		dbStructure.Chirps = make(map[int]Chirp)
	}
	if dbStructure.Users == nil {
		dbStructure.Users = make(map[int]User)
	}

	return dbStructure, nil
}

func (db *DB) writeDb(dbStructure DbStructure) error {
	jsonData, err := json.MarshalIndent(dbStructure, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}
	err = os.WriteFile(db.path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func (dbStructure *DbStructure) addUserToStorage(id int, email string) {
	user := User{Id: id, Email: email}
	dbStructure.Users[id] = user
}

func (dbStructure *DbStructure) addChirpToStorage(id int, body string) {
	chirp := Chirp{Id: id, Body: body}
	dbStructure.Chirps[id] = chirp
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
