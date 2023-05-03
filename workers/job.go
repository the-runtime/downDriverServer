package workers

import (
	"context"
	"golang.org/x/oauth2"
	"serverFordownDrive/database"
	"serverFordownDrive/fileController"
	//"serverFordownDrive/handlers"
	"serverFordownDrive/model"
)

type JobHandler func(ctx context.Context, args []interface{}) error

var JobQueue chan Job

// job for downloading the file

type Job struct {
	url              string
	userid           string
	googleAuthConfig *oauth2.Config
	CurrentUser      *model.User
}

func NewJob(url, id string, googleAuthConfig *oauth2.Config, temUser *model.User) Job {
	return Job{
		url:              url,
		userid:           id,
		googleAuthConfig: googleAuthConfig,
		CurrentUser:      temUser,
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

	filename, num := fileController.StartDown(j.url, j.CurrentUser)
	if num != 1 {
		println("Problem while downloading file")
		return nil

	}

	fileController.UploadFile(token, j.googleAuthConfig, filename)
	return nil
}
