package database

import (
	"fmt"
	"time"
)

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
