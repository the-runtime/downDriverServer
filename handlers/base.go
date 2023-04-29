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

	mux.Handle("/", http.FileServer(http.Dir("templates")))

	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)
	mux.HandleFunc("/process/", startGdrive)

	return mux
}
