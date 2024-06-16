package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DbStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	dbStructure := DbStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}

	return db.writeDb(dbStructure)
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDB()
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	return db, err
}

func (db *DB) loadDb() (DbStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	file, err := os.ReadFile(db.path)
	if err != nil {
		return DbStructure{}, fmt.Errorf("error reading file: %v", err)
	}

	var dbStructure DbStructure
	if err := json.Unmarshal(file, &dbStructure); err != nil {
		return DbStructure{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return dbStructure, nil
}

func (db *DB) writeDb(dbStructure DbStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	jsonData, err := json.MarshalIndent(dbStructure, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
