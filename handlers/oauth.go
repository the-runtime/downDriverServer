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
	"net/url"
	"serverFordownDrive/config"
	"serverFordownDrive/database"
	"serverFordownDrive/jwtAuth"
	"serverFordownDrive/model"
	"strings"
	"time"
)

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  config.GetRedirectUrl(),
	ClientID:     config.GetClientId(),
	ClientSecret: config.GetClientSecret(),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/drive"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

//const oauthGoogleRevokeURi = "https://oauth2.googleapis.com/revoke?token="

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
	//tokenDb.AutoMigrate(&model.UserToken{})

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

	var temToken model.UserToken

	//check if token struct contains refresh token
	if token.RefreshToken != "" {
		tokenDb.Where(model.UserToken{UserId: structData.Id}).Assign(model.UserToken{AccessToken: token.AccessToken,
			TokenType:    token.TokenType,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry,
		}).FirstOrCreate(&temToken)
	} else {
		tokenDb.Where(model.UserToken{UserId: structData.Id}).Assign(model.UserToken{AccessToken: token.AccessToken,
			TokenType: token.TokenType,
			Expiry:    token.Expiry,
		}).FirstOrCreate(&temToken)
	}

	//tokenDb.Where(model.UserToken{UserId: structData.Id}).Assign(model.UserToken{AuthCode: r.FormValue("code")}).FirstOrCreate(&temToken)

	//tokenDb.FirstOrCreate(&model.UserToken{
	//	UserId: structData.Id,
	//	Token:  *token,
	//})

	//creating a new user if not found in database already
	//var temUser model.User
	//userDb.Where(model.User{UserId: structData.Id}).Attrs(model.User{AllowedSpeed: 5, AllowedDataTransfer: 500 * 1024 * 1024, ConsumedDataTransfer: 0, AllowedThreads: 2}).FirstOrCreate(&temUser)
	////userDb.FirstOrCreate(&model.User{
	//	UserId:            structData.Id,
	//	AllowedBandwidth:  10, // data allowed to be consumed in MB
	//	AllowedSpeed:      1,
	//	ConsumedBandwidth: 0,
	//})

	//set userid as a cookie to the cleint for verification

	tokenString, err := jwtAuth.GeenerateJWT(structData.Id)
	if err != nil {
		println(err.Error())
		return
	}

	setUserCookie(w, tokenString)
	//fmt.Fprintf(w, "UserInfo: %s\n", data)
	http.Redirect(w, r, "/register", http.StatusTemporaryRedirect)
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
	var expiration = time.Now().Add(30 * 24 * time.Hour)
	//b := []byte(str)
	println("code is %s", str)
	//state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:    "token",
		Value:   str,
		Path:    "/",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)

}

func revokeToken(refreshToken string) error {
	revocationURL := "https://oauth2.googleapis.com/revoke"
	data := url.Values{
		"token":           {refreshToken},
		"token_type_hint": {"refresh_token"},
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, revocationURL, strings.NewReader(data.Encode()))
	if err != nil {
		println("Error with new request", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(req)
	if err != nil {
		println("error with client do", err)
		return err
	}

	return nil
}
