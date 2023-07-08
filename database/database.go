package database

import (
	//"github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"serverFordownDrive/config"
	"serverFordownDrive/model"
)

var databaseConn *gorm.DB

func init() {
	initDb, err := gorm.Open(postgres.Open(config.GetPostgresUrl()), &gorm.Config{})
	if err != nil {
		println(err.Error())
		return
	}
	databaseConn = initDb
}

func NewUserDb() (*gorm.DB, error) {

	//userdb, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	//if err != nil {
	//	println(err.Error())
	//	return nil, err
	//}
	//userdb.AutoMigrate(&model.User{})
	//return userdb, nil

	//change from sqlite to Postgres
	userDb := databaseConn
	userDb.AutoMigrate(&model.User{})
	return userDb, nil

}

func NewTokenDb() (*gorm.DB, error) {

	tokenDb := databaseConn
	tokenDb.AutoMigrate(&model.UserToken{})
	return tokenDb, nil
}

func NewHistoryDb() (*gorm.DB, error) {
	//historyDb, err := gorm.Open(sqlite.Open("history.db"), &gorm.Config{})
	//if err != nil {
	//	println(err.Error())
	//	return nil, err
	//}
	//historyDb.AutoMigrate(&model.SingleHistory{})
	//return historyDb, nil

	historyDb := databaseConn
	historyDb.AutoMigrate(&model.SingleHistory{})
	return historyDb, nil
}
