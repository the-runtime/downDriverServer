package handlers

import (
	"io"
	"net/http"
	"os"
	"serverFordownDrive/database"
	"serverFordownDrive/model"
)

//todo
//redirect to some informative page if file open fails rather than just returning empty response

func checkAuth(r *http.Request) bool {
	attemptUserId := r.FormValue("user")
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
		return false
	}
	return true
}

func rootWeb(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open("web/nonPublicAssets/index.html")
	if err != nil {
		println(err.Error())
		return
	}
	io.Copy(w, file)

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

	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
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

	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}

func profileWeb(w http.ResponseWriter, r *http.Request) {

	if checkAuth(r) {
		file, err := os.Open("web/nonPublicAssets/profile.html")

		if err != nil {
			println(err.Error())
			return
		}
		io.Copy(w, file)
		return
	}

	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}

func registerWeb(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) {
		http.Redirect(w, r, "/dashboard", http.StatusPermanentRedirect)
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
		http.Redirect(w, r, "/dashboard", http.StatusPermanentRedirect)
		return
	}
	file, err := os.Open("web/nonPublicAssets/login.html")
	if err != nil {
		println(err.Error())
		return
	}
	io.Copy(w, file)

}
