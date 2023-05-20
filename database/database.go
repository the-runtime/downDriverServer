package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"serverFordownDrive/model"
)

func NewUserDb() (*gorm.DB, error) {
	userdb, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		println(err.Error())
		return nil, err
	}
	userdb.AutoMigrate(&model.User{})
	return userdb, nil
}

func NewTokenDb() (*gorm.DB, error) {
	tokenDb, err := gorm.Open(sqlite.Open("token.db"), &gorm.Config{})
	if err != nil {
		println(err.Error())
		return nil, err
	}
	tokenDb.AutoMigrate(&model.UserToken{})
	return tokenDb, nil
}

func NewHistoryDb() (*gorm.DB, error) {
	historyDb, err := gorm.Open(sqlite.Open("history.db"), &gorm.Config{})
	if err != nil {
		println(err.Error())
		return nil, err
	}
	historyDb.AutoMigrate(&model.SingleHistory{})
	return historyDb, nil
}
