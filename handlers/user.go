package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"serverFordownDrive/database"
	"serverFordownDrive/model"
)

type retrunUser struct {
	Name        string `json:"username"`
	Email       string `json:"email"`
	Image       string `json:"userimage"`
	DataRemaing uint64 `json:"transferused"`
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	userDb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
		return
	}
	temUser := model.User{}
	r.ParseForm()
	userIdCookie, err := r.Cookie("user")
	if err != nil {
		println(err.Error())
		return
	}
	userId := userIdCookie.Value
	err = userDb.Where("user_id = ?", userId).First(&temUser).Error
	if err == nil {
		fmt.Fprintf(w, "user already regitered")
		return
	}
	userDb.AutoMigrate(&model.User{})
	regUser := model.User{
		UserId:               userId,
		FisrtName:            r.Form["firstname"][0],
		LastName:             r.Form["lastname"][0],
		AccountType:          r.Form["accounttype"][0],
		ConsumedDataTransfer: 0,
	}
	switch regUser.AccountType {
	case "1":
		{
			regUser.AllowedThreads = 1
			regUser.AllowedSpeed = 5
			regUser.AllowedDataTransfer = 1 * 1024 * 1024 * 1024
		}
	case "2":
		{
			regUser.AllowedThreads = 1
			regUser.AllowedSpeed = 10
			regUser.AllowedDataTransfer = 2 * 1024 * 1024 * 1024
		}
	case "3":
		{
			regUser.AllowedThreads = 1
			regUser.AllowedSpeed = 20
			regUser.AllowedDataTransfer = 3 * 1024 * 1024 * 1024
		}
	}

	userDb.Create(&regUser)
	http.Redirect(w, r, "/profile", http.StatusPermanentRedirect)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userDb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
		return
	}

	userIdCookie, err := r.Cookie("user")
	if err != nil {
		println(err.Error())
		return
	}
	userId := userIdCookie.Value
	var retUser model.User

	userDb.Where("user_id = ?", userId).First(&retUser)

	outUser := retrunUser{
		Name:        retUser.FisrtName + " " + retUser.LastName,
		Email:       "",
		Image:       "",
		DataRemaing: retUser.AllowedDataTransfer - retUser.ConsumedDataTransfer,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(outUser)

}
