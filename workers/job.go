package workers

import (
	"context"
	"golang.org/x/oauth2"
	"serverFordownDrive/database"
	"serverFordownDrive/fileController"
	//"serverFordownDrive/handlers"
	"serverFordownDrive/model"
)

//var GoogleOauthConfig = &oauth2.Config{
//	RedirectURL:  "http://127.0.0.1:8000/auth/google/callback",
//	ClientID:     "882134345746-3fo1qd40q4p0m0fbdm31f453frjhu60e.apps.googleusercontent.com",
//	ClientSecret: "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",
//	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/drive"},
//	Endpoint:     google.Endpoint,
//}

type JobHandler func(ctx context.Context, args []interface{}) error

var JobQueue chan Job

// job for downloading the file

type Job struct {
	url              string
	userid           string
	googleAuthConfig *oauth2.Config
}

func NewJob(url, id string, googleAuthConfig *oauth2.Config) Job {
	return Job{
		url:              url,
		userid:           id,
		googleAuthConfig: googleAuthConfig,
	}
}

func (j Job) DoJob() error {

	tokenDb, err := database.NewTokenDb()
	if err != nil {
		println(err.Error())
		return err
	}

	//tokenDb.AutoMigrate(&model.UserToken{})
	//var Token oauth2.Token
	var temTokenUser model.UserToken

	tokenDb.Where("user_id = ?", j.userid).First(&temTokenUser)
	AccessToken := temTokenUser.AccessToken
	RefreshToken := temTokenUser.RefreshToken
	TokenType := temTokenUser.TokenType
	Expiry := temTokenUser.Expiry

	token := &oauth2.Token{AccessToken: AccessToken,
		TokenType:    TokenType,
		RefreshToken: RefreshToken,
		Expiry:       Expiry,
	}

	//token, err := j.googleAuthConfig.Exchange(context.Background(), temTokenUser.AuthCode)
	//if err != nil {
	//	println("error is in token exchange")
	//	println(err.Error())
	//	return err
	//}

	filename, num := fileController.StartDown(j.url)
	if num != 1 {
		println("Problem while downloading file")
		return nil

	}

	fileController.UploadFile(token, j.googleAuthConfig, filename)
	return nil
}
