package database

import (
	"fmt"

	"github.com/ggrangel/go-webserver/auth"
)

type RefreshToken struct {
	Token  string `json:"refresh_token"`
	Expiry int    `json:"expiry"`
}

type User struct {
	Id           int          `json:"id"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	Token        string       `json:"token"`
	RefreshToken RefreshToken `json:"refresh_token"`
	IsChirpyRed  bool         `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDb()
	if err != nil {
		return User{}, err
	}

	nextKey := getNextUserKey(&dbStructure.Users)

	dbStructure.addUserToStorage(nextKey, email, password, false)

	err = db.writeDb(dbStructure)
	if err != nil {
		return User{}, err
	}

	return dbStructure.Users[nextKey], nil
}

func (db *DB) GetUserById(id int) (User, error) {
	dbStructure, err := db.loadDb()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Id == id {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("User not found")
}

func (db *DB) GetUserByEmail(email, password string) (User, error) {
	dbStructure, err := db.loadDb()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			err := auth.CheckPasswordHash(password, user.Password)
			if err != nil {
				return User{}, err
			} else {
				return user, nil
			}
		}
	}

	return User{}, fmt.Errorf("User not found")
}

func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	dbStructure, err := db.loadDb()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Id == id {
			user.Email = email
			user.Password = password
			dbStructure.Users[id] = user
			err = db.writeDb(dbStructure)
			if err != nil {
				return User{}, err
			}
			return user, nil
		}
	}

	return User{}, fmt.Errorf("User not found")
}

func (db *DB) SetUserToRedMember(id int) error {
	dbStructure, err := db.loadDb()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return fmt.Errorf("User not found")
	}

	user.IsChirpyRed = true

	dbStructure.Users[id] = user
	err = db.writeDb(dbStructure)
	if err != nil {
		return err
	}

	return nil
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

func (dbStructure *DbStructure) addUserToStorage(id int, email, password string, isChirpyRed bool) {
	user := User{Id: id, Email: email, Password: password, IsChirpyRed: isChirpyRed}
	dbStructure.Users[id] = user
}
