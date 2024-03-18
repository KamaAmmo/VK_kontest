package main

import (
	// "fmt"

	"net/http"
	"regexp"
	_ "vk_app/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	var protected func(w http.ResponseWriter, r *http.Request)
	switch {
	case r.Method == http.MethodGet && getFilmRe.MatchString(path):
		h.app.getFilm(w, r)
		return
	case r.Method == http.MethodGet && listFilmRe.MatchString(path):
		if r.URL.Query().Get("name") != "" || r.URL.Query().Get("title") != "" {
			h.app.searchFilmByPattern(w, r)
		} else {
			h.app.listFilms(w, r)
		}
		return
	case r.Method == http.MethodDelete && getFilmRe.MatchString(path):
		protected = h.app.deleteFilm
	case r.Method == http.MethodPost && createFilmRe.MatchString(path):
		protected = h.app.createFilm
	case r.Method == http.MethodPut && createFilmRe.MatchString(path):
		protected = h.app.editFilm
	default:
		h.app.notFound(w)
		return
	}

	if r.Context().Value(roleKey) == "admin" {
		protected(w, r)
	} else {
		h.app.clientError(w, http.StatusForbidden)
	}
}

func (h *actorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	var protected func(w http.ResponseWriter, r *http.Request)
	switch {
	case r.Method == http.MethodGet && getActorRe.MatchString(path):
		h.app.getActor(w, r)
		return
	case r.Method == http.MethodGet && createActorRe.MatchString(path):
		h.app.listActorsFilms(w)
		return
	case r.Method == http.MethodDelete && getActorRe.MatchString(path):
		protected = h.app.deleteActor
	case r.Method == http.MethodPost && createActorRe.MatchString(path):
		protected = h.app.createActor
	case r.Method == http.MethodPut && createActorRe.MatchString(path):
		protected = h.app.editActor
	default:
		h.app.notFound(w)
		return
	}
	if r.Context().Value(roleKey) == "admin" {
		protected(w, r)
	} else {
		h.app.clientError(w, http.StatusForbidden)
	}

}

func (app *application) route() http.Handler {
	mux := http.NewServeMux()
	actorsHandler := &actorsHandler{app}
	filmsHandler := &filmsHandler{app}

	mux.Handle("/actors/", app.authenticate(actorsHandler))
	mux.Handle("/films/", app.authenticate(filmsHandler))

	mux.HandleFunc("/register", app.registerUser)
	mux.HandleFunc("/login", app.loginUser)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return app.logRequest((app.recoverPanic(mux)))
}
