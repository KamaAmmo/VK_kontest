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

	getFilmRe    = regexp.MustCompile(`^\/films\/(\d+)$`)
	createFilmRe = regexp.MustCompile(`^\/films[\/]*$`)
	listFilmRe   = regexp.MustCompile(`^\/films`)
)

type actorsHandler struct {
	app *application
}

type filmsHandler struct {
	app *application
}

func (h *filmsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	switch {
	case r.Method == http.MethodGet && getFilmRe.MatchString(path):
		h.app.getFilm(w, r)
		return
	case r.Method == http.MethodDelete && getFilmRe.MatchString(path):
		h.app.deleteFilm(w, r)
		return
	case r.Method == http.MethodPost && createFilmRe.MatchString(path):
		h.app.createFilm(w, r)
		return
	case r.Method == http.MethodPut && getFilmRe.MatchString(path):
		h.app.editFilm(w, r)
		return 
	case r.Method == http.MethodGet && listFilmRe.MatchString(path):
		v := r.URL.Query()
		if v.Get("name") != "" || v.Get("title") != ""{
			h.app.searchFilm(w, r)
		} else {
			h.app.listFilms(w, r)
		}
		return
	default:
		h.app.notFound(w)
	}
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
		h.app.listActorsFilms(w)
	default:
		h.app.notFound(w)
	}
}

func (app *application) route() http.Handler {
	mux := http.NewServeMux()
	actorsHandler := &actorsHandler{app}
	filmsHandler := &filmsHandler{app}
	mux.HandleFunc("/", app.home)

	//user
	mux.Handle("/actors/", actorsHandler)
	mux.Handle("/films/", filmsHandler)
	//admin

	return app.logRequest(app.recoverPanic(mux))
}
