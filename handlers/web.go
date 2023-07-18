package handlers

import (
	"io"
	"net/http"
	"os"
	"serverFordownDrive/database"
	"serverFordownDrive/jwtAuth"
	"serverFordownDrive/model"
)

//todo
//redirect to some informative page if file open fails rather than just returning empty response

func checkAuth(r *http.Request) bool {
	attemptUserId, err := jwtAuth.IsAuthorized2(r)
	if err != nil {
		println(err.Error())
		return false
	}
	if attemptUserId == "" {
		return false
	}

	userDb, err := database.NewTokenDb()
	if err != nil {
		println(err.Error())
		return false
	}

	var temUser model.User
	err = userDb.Where("user_id = ?", attemptUserId).First(&temUser).Error
	if err != nil {
		println(err.Error())
		return false
	}
	return true
}

func tableWeb(w http.ResponseWriter, r *http.Request) {

	if checkAuth(r) {
		file, err := os.Open("web/nonPublicAssets/table.html")

		if err != nil {
			println(err.Error())
			return
		}
		io.Copy(w, file)
		return
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func dashboardWeb(w http.ResponseWriter, r *http.Request) {

	if checkAuth(r) {
		file, err := os.Open("web/nonPublicAssets/dashboard.html")

		if err != nil {
			println(err.Error())
			return
		}
		io.Copy(w, file)
		return
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func profileWeb(w http.ResponseWriter, r *http.Request) {

	if checkAuth(r) {
		println("auth check passes")
		file, err := os.Open("web/nonPublicAssets/profile.html")
		if err != nil {
			println(err.Error())
			return
		}
		io.Copy(w, file)
		return
	} else {
		println("Auth failed in profile")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}

}

func registerWeb(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) {
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}
	file, err := os.Open("web/nonPublicAssets/register.html")
	if err != nil {
		println(err.Error())
		return
	}
	io.Copy(w, file)

}

func loginWeb(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) {
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}
	file, err := os.Open("web/nonPublicAssets/login.html")
	if err != nil {
		println(err.Error())
		return
	}
	io.Copy(w, file)

}
