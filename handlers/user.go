package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"serverFordownDrive/database"
	"serverFordownDrive/model"
	"time"
)

type returnUser struct {
	Name         string `json:"username"`
	Email        string `json:"email"`
	Image        string `json:"user_image"`
	DataRemains  int64  `json:"data_remains"`
	DataAllotted int64  `json:"data_allotted"`
}

type returnTable struct {
	Id          string                `json:"userid"`
	HistoryList []model.SingleHistory `json:"history_list"`
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

	userId := r.Header.Get("user")
	err = userDb.Where("user_id = ?", userId).First(&temUser).Error
	if err == nil {
		fmt.Fprintf(w, "user already regitered")
		return
	}
	//userDb.AutoMigrate(&model.User{})
	regUser := model.User{
		UserId:               userId,
		FirstName:            r.Form["firstname"][0],
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
	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
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

	userId := r.Header.Get("user")

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

	//push possibly new access token to database
	var tempToken model.UserToken
	tokenDb.Where("user_id = ? ", userId).Assign(model.UserToken{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		Expiry:      token.Expiry,
	}).FirstOrCreate(&tempToken)

	var structData model.GoogleUserData
	err = json.Unmarshal(googleUserData, &structData)
	fmt.Printf("structData is ", structData)
	if err != nil {
		println("error in Unmarshalling GoogleUserData")
		return
	}

	outUser := returnUser{
		Name:  retUser.FirstName + " " + retUser.LastName,
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

	id := r.Header.Get("user")

	historyDb, err := database.NewHistoryDb()
	if err != nil {
		println(err.Error())
		return
	}

	historyList := []model.SingleHistory{}
	historyDb.Order("finishedat asc").Where("user_id = ?", id).Find(&historyList)

	retTable := returnTable{id, historyList}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(retTable)

}

func resetLimit(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("user")

	userDb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
		return
	}

	userDb.Model(&model.User{}).Where("user_id = ?", id).Update("consumed_data_transfer", 0)

	fmt.Fprintf(w, "User limit is reset")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("user")
	userDb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
		return
	}

	delToken := model.UserToken{}
	userDb.Model(&delToken)
	userDb.Where("user_id = ?", id).Find(&delToken)

	err = revokeToken(delToken.RefreshToken)
	if err != nil {
		println("revoke token failed")
		println(err.Error())
		fmt.Fprint(w, err)
		return
	}
	println("access token revoked")
	userDb.Unscoped().Delete(&delToken) //permanently delete from the table
	delUser := model.User{}
	userDb.Model(&delUser).Where("user_id = ?", id).Find(&delUser)
	userDb.Unscoped().Delete(&delUser)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func signOut(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
