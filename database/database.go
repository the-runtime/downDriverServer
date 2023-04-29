package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewUserDb() (*gorm.DB, error) {
	userdb, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		println(err.Error())
		return nil, err
	}
	return userdb, nil
}

func NewTokenDb() (*gorm.DB, error) {
	tokenDb, err := gorm.Open(sqlite.Open("token.db"), &gorm.Config{})
	if err != nil {
		println(err.Error())
		return nil, err
	}
	return tokenDb, nil
}
