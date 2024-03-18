package main

import (
	// "fmt"
	"net/http"
	"regexp"
)

var (
	// listActorRe = regexp.MustCompile(`^\/actors[\/]*$`)
	getActorRe    = regexp.MustCompile(`^\/actors\/(\d+)$`)
	createActorRe = regexp.MustCompile(`^\/actors[\/]*$`)
)

type actorsHandler struct {
	app *application
}

func (h *actorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	switch {
	case r.Method == http.MethodGet && getActorRe.MatchString(path):
		h.app.getActor(w, r)
		return
	case r.Method == http.MethodDelete && getActorRe.MatchString(path):
		h.app.deleteActor(w, r)
		return
	case r.Method == http.MethodPost && createActorRe.MatchString(path):
		h.app.createActor(w, r)
		return
	case r.Method == http.MethodPut && getActorRe.MatchString(path):
		h.app.editActor(w, r)
	case r.Method == http.MethodGet && createActorRe.MatchString(path):
		h.app.listActorsFilms(w, r)
	default:
		h.app.notFound(w)
	}
}

func (app *application) route() http.Handler {
	mux := http.NewServeMux()
	actorsHandler := &actorsHandler{app}

	mux.HandleFunc("/", app.home)

	//user
	mux.Handle("/actors/", actorsHandler)

	//admin

	return app.logRequest(mux)
}
