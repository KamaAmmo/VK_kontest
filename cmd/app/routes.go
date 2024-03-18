package main

import (
	// "fmt"
	"net/http"
)

func (app *application) route() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)

	//user
	mux.HandleFunc("/actors", app.getActor)
	// mux.HandleFunc("/actors", app.)

	//admin

	return app.logRequest(mux)
}
