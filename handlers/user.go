package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"serverFordownDrive/database"
	"serverFordownDrive/model"
)

type returnUser struct {
	Name         string `json:"username"`
	Email        string `json:"email"`
	Image        string `json:"user_image"`
	DataRemains  int64  `json:"data_remains"`
	DataAllotted int64  `json:"data_allotted"`
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	return
	//}
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
	//userDb.AutoMigrate(&model.User{})
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

	tokenDb, err := database.NewTokenDb()
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
	var userToken model.UserToken

	tokenDb.Where("user_id = ?", userId).First(&userToken)
	userDb.Where("user_id = ?", userId).First(&retUser)

	AccessToken := userToken.AccessToken
	RefreshToken := userToken.RefreshToken
	TokenType := userToken.TokenType
	Expiry := userToken.Expiry

	token := &oauth2.Token{AccessToken: AccessToken,
		TokenType:    TokenType,
		RefreshToken: RefreshToken,
		Expiry:       Expiry,
	}

	googleUserData, err := getUserDataFromGoogle(token)
	if err != nil {
		println(err.Error())
		return
	}
	var structData model.GoogleUserData
	err = json.Unmarshal(googleUserData, &structData)
	if err != nil {
		println("error in Unmarshalling GoogleUserData")
		return
	}

	outUser := returnUser{
		Name:  retUser.FisrtName + " " + retUser.LastName,
		Email: structData.Email,
		Image: structData.Picture,
		DataRemains: func() int64 {
			if int64(retUser.AllowedDataTransfer)-int64(retUser.ConsumedDataTransfer) < 0 {
				return int64(0)
			}
			return int64((retUser.AllowedDataTransfer - retUser.ConsumedDataTransfer) / (1024 * 1024))
		}(),
		DataAllotted: int64(retUser.AllowedDataTransfer / (1024 * 1024)),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(outUser)

}

func getTable(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("user")
	if err != nil {
		println(err.Error())
		return
	}
	id := idCookie.Value

	historyDb, err := database.NewHistoryDb()
	if err != nil {
		println(err.Error())
		return
	}

	historyList := []model.SingleHistory{}
	historyDb.Order("finishedat asc").Where("user_id = ?", id, &historyList)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(historyList)

}

func resetLimit(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("user")
	if err != nil {
		println(err.Error())
		return
	}
	id := idCookie.Value

	userDb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
		return
	}

	userDb.Where("user_id = ?", id).Update("consumed_data_transfer", 0)
}
