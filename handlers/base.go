package handlers

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"
	"serverFordownDrive/config"
	"serverFordownDrive/jwtAuth"
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
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/process/", jwtAuth.IsAuthorized(startGdrive)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/progress", jwtAuth.IsAuthorized(progressBar)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/frontauth/", jwtAuth.IsAuthorized(frontAuth))) //we don't use it. will be removed soon
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/register", jwtAuth.IsAuthorized(registerUser)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/getuser", jwtAuth.IsAuthorized(getUser)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/password/reset", jwtAuth.IsAuthorized(resetLimit))) //for testing only  to be removed if in production
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/table", jwtAuth.IsAuthorized(getTable)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/delete", jwtAuth.IsAuthorized(deleteUser)))
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/api/account/signout", jwtAuth.IsAuthorized(signOut)))

	//for maintenance
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, "/maintain/keepalive", keepAlive))

	//for delete history every 24 hours (using cron job)
	secretPath := fmt.Sprintf("/maintain/delete/%s", config.GetMaintainSecret())
	mux.HandleFunc(newrelic.WrapHandleFunc(relicApp, secretPath, clearHistory))

	return mux
}
