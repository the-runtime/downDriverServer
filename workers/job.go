package workers

import (
	"golang.org/x/oauth2"
	"serverFordownDrive/controller"
	"serverFordownDrive/database"
	"serverFordownDrive/fileController"
	//"serverFordownDrive/handlers"
	"serverFordownDrive/model"
)

//type JobHandler func(ctx context.Context, args []interface{}) error
// if custom JobHandler is required

var JobQueue chan Job

// job for downloading the file

type Job struct {
	url              string
	userid           string
	googleAuthConfig *oauth2.Config
	CurrentUser      *model.User
	ProgressId       int
}

func NewJob(url, userid string, googleAuthConfig *oauth2.Config, temUser *model.User) Job {
	progressId := controller.NewProgress("", userid, 0)
	//dataProgress := *controller.GetDataProgress()
	//temp2Progress := dataProgress[id][progressId]

	return Job{
		url:              url,
		userid:           userid,
		googleAuthConfig: googleAuthConfig,
		CurrentUser:      temUser,
		ProgressId:       progressId,
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

	//to handle progress info

	filename, num := fileController.StartDown(j.url, j.CurrentUser, j.ProgressId)
	if num != 1 {
		println("Problem while downloading file")
		return nil

	}

	fileController.UploadFile(token, j.googleAuthConfig, filename, j.CurrentUser, j.ProgressId)

	//print("trying to update the token after refresh")
	var tempToken model.UserToken
	tokenDb.Where("user_id = ? ", j.userid).Assign(model.UserToken{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		Expiry:      token.Expiry,
	}).FirstOrCreate(&tempToken)

	return nil
}
