package database

import (
	"fmt"
	"time"

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
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDb()
	if err != nil {
		return User{}, err
	}

	nextKey := getNextUserKey(&dbStructure.Users)

	dbStructure.addUserToStorage(nextKey, email, password)

	err = db.writeDb(dbStructure)
	if err != nil {
		return User{}, err
	}

	return dbStructure.Users[nextKey], nil
}

func (db *DB) GetUser(email, password string) (User, error) {
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

func (db *DB) SaveUserToken(userId int, token string) error {
	dbStructure, err := db.loadDb()
	if err != nil {
		return err
	}

	for _, user := range dbStructure.Users {
		if user.Id == userId {
			user.Token = token
			dbStructure.Users[userId] = user
			err = db.writeDb(dbStructure)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("User not found")
}

func (db *DB) SaveUserRefreshToken(userId int, refreshToken string, expiryAt int) error {
	dbStructure, err := db.loadDb()
	if err != nil {
		return err
	}

	for _, user := range dbStructure.Users {
		if user.Id == userId {
			user.RefreshToken.Token = refreshToken
			user.RefreshToken.Expiry = expiryAt
			dbStructure.Users[userId] = user
			err = db.writeDb(dbStructure)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("User not found")
}

func (db *DB) GetUserByRefreshToken(refreshToken string) (User, error) {
	dbStructure, err := db.loadDb()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.RefreshToken.Token == refreshToken {
			if user.RefreshToken.Expiry > int(time.Now().Unix()) {
				return user, nil
			}
		}
	}

	return User{}, fmt.Errorf("User not found")
}

func (db *DB) RevokeUserRefreshToken(userId int) error {
	dbStructure, err := db.loadDb()
	if err != nil {
		return err
	}

	for _, user := range dbStructure.Users {
		if user.Id == userId {
			user.RefreshToken.Token = ""
			user.RefreshToken.Expiry = 0
			dbStructure.Users[userId] = user
			err = db.writeDb(dbStructure)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("User not found")
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

func (dbStructure *DbStructure) addUserToStorage(id int, email, password string) {
	user := User{Id: id, Email: email, Password: password}
	dbStructure.Users[id] = user
}
