package handlers

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"
	"serverFordownDrive/config"
	"serverFordownDrive/jwtauth"
	"serverFordownDrive/workers"
)

func New() http.Handler {
	mux := http.NewServeMux()

	// start workers with dispatcher

	dispatch := workers.NewDispatcher(config.GetNumWorkers())
	dispatch.Run()
	workers.InitJobQueue()

	//start newRelic agent for monitoring of the service

	relicApp, err := newrelic.NewApplication(newrelic.ConfigAppName("Downdive"),
		newrelic.ConfigLicense(config.GetNewRelic()),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		println("Problem with newRelic init")
		println(err.Error())
	}

	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, "/index.html", http.StatusPermanentRedirect)
	//})

	//For public assests like css, js and images
	mux.Handle("/", http.FileServer(http.Dir("web/public")))

	//Handle interaction with web Browser
	//mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/", rootWeb))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/login", loginWeb))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/register", registerWeb))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/dashboard", dashboardWeb))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/table", tableWeb))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/profile", profileWeb))

	//Handle interaction with website
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/auth/google/login", oauthGoogleLogin))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/auth/google/callback", oauthGoogleCallback))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/process/", jwtauth.IsAuthorized(startGdrive)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/progress", jwtauth.IsAuthorized(progressBar)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/frontauth/", jwtauth.IsAuthorized(frontAuth)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/register", jwtauth.IsAuthorized(registerUser)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/getuser", jwtauth.IsAuthorized(getUser)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/password/reset", jwtauth.IsAuthorized(resetLimit))) //for testing only  to be removed if in production
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/table", jwtauth.IsAuthorized(getTable)))

	return mux
}
