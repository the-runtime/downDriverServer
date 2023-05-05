package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"serverFordownDrive/controller"
	"serverFordownDrive/database"
	"serverFordownDrive/model"
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

	userDb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
		return
	}

	var temUser model.User
	userDb.Where("user_id = ?", id).First(&temUser)
	if (temUser.AllowedDataTransfer - temUser.ConsumedDataTransfer) <= 0 {
		println("user surpassed transfer limit")
		fmt.Fprintf(w, "user surpassed transfer limit")
		return
	}

	println("job pushed to queue")
	job := workers.NewJob(downUrl, id, GoogleOauthConfig, &temUser)
	workers.JobQueue <- job

	fmt.Fprintf(w, "Work in progrss check your drive after some time")

	//http.Redirect(w, r, "/progressbar", http.StatusPermanentRedirect)

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

func progressBar(w http.ResponseWriter, r *http.Request) {
	cokkieUserId, err := r.Cookie("user")
	if err != nil {
		print("error in progresserror")
		println(err.Error())
	}
	userId := cokkieUserId.Value
	reqListProcess := controller.GetProgressList(userId)

	fmt.Printf("\n", reqListProcess)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(reqListProcess)
	if err != nil {
		return
	}

	//fmt.Fprintf(w, "Information rageding your process \n"+"filename: "+reqProcess.Filename+"\n"+"File size: %d \n"+"Downloaded: %d MBs", reqProcess.Total, reqProcess.Transferred)
}
