package handlers

import (
	"net/http"
	"serverFordownDrive/workers"
)

func New() http.Handler {
	mux := http.NewServeMux()

	// start workers with dispatcher

	dispatch := workers.NewDispatcher(5)
	dispatch.Run()
	workers.InitJobQueue()

	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, "/index.html", http.StatusPermanentRedirect)
	//})
	mux.Handle("/", http.FileServer(http.Dir("template")))

	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)
	mux.HandleFunc("/process/", startGdrive)
	mux.HandleFunc("/progressbar", progressBar)

	return mux
}
