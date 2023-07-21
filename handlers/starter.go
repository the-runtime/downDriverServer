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
	println("For debug 'starter satred'")

	id := r.Header.Get("user")

	downUrl := r.FormValue("url")

	if downUrl == "" {
		println("Empty Url")
		return
	}

	userDb, err := database.NewUserDb()
	if err != nil {
		println("UserDb error")
		println(err.Error())
		return
	}

	var temUser model.User
	//userDb.AutoMigrate(&model.User{})
	userDb.Where("user_id = ?", id).First(&temUser)
	if (int64(temUser.AllowedDataTransfer) - int64(temUser.ConsumedDataTransfer)) <= 0 {
		println("user surpassed transfer limit")
		fmt.Fprintf(w, "user surpassed transfer limit")
		return
	}

	println("job pushed to queue")
	job := workers.NewJob(downUrl, id, GoogleOauthConfig, &temUser)
	workers.JobQueue <- job

	//fmt.Fprintf(w, "Work in progrss check your drive after some time")

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)

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

	userId := r.Header.Get("user")
	reqListProcess := controller.GetProgressList(userId)

	//fmt.Printf("\n", reqListProcess)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// no process is on going
	if len(reqListProcess) == 0 {
		return
	}
	//testing for multiple progress bars
	//err := json.NewEncoder(w).Encode(reqListProcess[len(reqListProcess)-1])
	//fmt.Printf("process list is %+v", reqListProcess)
	err := json.NewEncoder(w).Encode(reqListProcess)
	if err != nil {
		return
	}

	//fmt.Fprintf(w, "Information rageding your process \n"+"filename: "+reqProcess.Filename+"\n"+"File size: %d \n"+"Downloaded: %d MBs", reqProcess.Total, reqProcess.Transferred)
}

// should not be used from now on
// soon it will be removed permanently
func frontAuth(w http.ResponseWriter, r *http.Request) {
	attemptUserId := r.FormValue("user")
	println("user from frontent is ", attemptUserId)
	if attemptUserId == "" {
		fmt.Fprintf(w, "0")
		return
	}

	userDb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
		return
	}

	var temUser model.User
	err = userDb.Where("user_id = ?", attemptUserId).First(&temUser).Error
	if err != nil {
		fmt.Fprintf(w, "0")
		return
	}

	fmt.Fprintf(w, "1")

}

func keepAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Keep it alive")
}

func clearHistory(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewHistoryDb()
	if err != nil {
		println(err.Error())
		return
	}
	err = db.Exec("DROP TABLE IF EXISTS single_histories;").Error
	if err != nil {
		fmt.Fprintf(w, "Problem deleting history")
	}

}
