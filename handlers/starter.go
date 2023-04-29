package handlers

import (
	"fmt"
	"net/http"
	"serverFordownDrive/workers"
)

func startGdrive(w http.ResponseWriter, r *http.Request) {
	//tokenDb, err := database.NewTokenDb()
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	////tokenDb.AutoMigrate(&model.UserToken{})
	//var Token oauth2.Token
	//var temTokenUser model.UserToken

	// to do implement a page if the user  used all bandwidth
	temUserCookie, _ := r.Cookie("user")
	id := temUserCookie.Value

	downUrl := r.FormValue("url")

	if downUrl == "" {
		println("Empty Url")
		return
	}

	println("job pushed to queue")
	job := workers.NewJob(downUrl, id)
	workers.JobQueue <- job

	fmt.Fprintf(w, "Work in progrss check your drive after some time")

	//tokenDb.Where("UserId = ?", id).First(&temTokenUser)
	//Token = temTokenUser.Token
	//token := &Token

	//filename, num := fileController.StartDown(downUrl)
	//if num == 0 {
	//	return
	//}
	//
	//fileController.UploadFile(token, googleOauthConfig, filename)

}
