package handlers

// to implement a getgoogleoauthconfig func so this varibale is prvate not public
import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"serverFordownDrive/database"
	"serverFordownDrive/model"
	"time"
)

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "https://theruntime.software/auth/google/callback",
	ClientID:     "882134345746-3fo1qd40q4p0m0fbdm31f453frjhu60e.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/drive"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	//oauthState := generateStateOauthCookie(w)

	//u := GoogleOauthConfig.AuthCodeURL(oauthState)
	u := GoogleOauthConfig.AuthCodeURL("state-oauth", oauth2.AccessTypeOffline)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {

	//oauthState, _ := r.Cookie("oauthstate")

	//if r.FormValue("state") != oauthState.Value {
	//	fmt.Println("invalid oauth google state")
	//	http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
	//	return
	//}

	token, err := GoogleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	//GoogleOauthConfig.AuthCodeURL("state-token",oauth2.AccessTypeOffline)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// set user database

	userDb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
		return
	}

	err = userDb.AutoMigrate(&model.User{})
	if err != nil {
		return
	}

	//set token to database
	tokenDb, err := database.NewTokenDb()
	if err != nil {
		println(err.Error())
		return
	}
	tokenDb.AutoMigrate(&model.UserToken{})

	// user data from google
	data, err := getUserDataFromGoogle(token)
	if err != nil {
		fmt.Println("get user data from google")
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//converting []byte of user data from google to a struct
	var structData model.GoogleUserData
	err = json.Unmarshal(data, &structData)
	if err != nil {
		println("error in Unmarshalling GoogleUserData")
		return
	}

	// assigning a token to tokenDB
	var temToken model.UserToken
	//tokenDb.Where(model.UserToken{UserId: structData.Id}).Assign(model.UserToken{AuthCode: r.FormValue("code")}).FirstOrCreate(&temToken)
	tokenDb.Where(model.UserToken{UserId: structData.Id}).Assign(model.UserToken{AccessToken: token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}).FirstOrCreate(&temToken)
	//tokenDb.FirstOrCreate(&model.UserToken{
	//	UserId: structData.Id,
	//	Token:  *token,
	//})

	//creating a new user if not found in database already
	var temUser model.User
	userDb.Where(model.User{UserId: structData.Id}).Attrs(model.User{AllowedSpeed: 5, AllowedDataTransfer: 500 * 1024 * 1024, ConsumedDataTransfer: 0}).FirstOrCreate(&temUser)
	//userDb.FirstOrCreate(&model.User{
	//	UserId:            structData.Id,
	//	AllowedBandwidth:  10, // data allowed to be consumed in MB
	//	AllowedSpeed:      1,
	//	ConsumedBandwidth: 0,
	//})

	//set userid as a cookie to the cleint for verification

	setUserCookie(w, structData.Id)
	fmt.Fprintf(w, "UserInfo: %s\n", data)
}

func getUserDataFromGoogle(token *oauth2.Token) ([]byte, error) {
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: #{err.Error()}")
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	return contents, nil
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func setUserCookie(w http.ResponseWriter, str string) {
	var expiration = time.Now().Add(20 * time.Minute)
	//b := []byte(str)
	println("code is %s", str)
	//state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:    "user",
		Value:   str,
		Path:    "/process",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)

}
