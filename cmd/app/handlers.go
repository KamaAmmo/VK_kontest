package main

import (
	// "fmt"
	"encoding/json"
	"errors"
	"net/http"
	"vk_app/internal/storage"
	// "strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All is good"))
}

func (app *application) getActor(w http.ResponseWriter, r *http.Request) {

	id := app.getID(w, r, getActorRe)
	person, err := app.people.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data, err := json.MarshalIndent(person, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *application) createActor(w http.ResponseWriter, r *http.Request) {
	p := storage.Person{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		app.serverError(w, err)
		return
	}

	lid, err := app.people.Insert(p)
	if err != nil {
		app.serverError(w, err)
		return
	}

	p.ID = lid
	data, err := json.MarshalIndent(p, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (app *application) deleteActor(w http.ResponseWriter, r *http.Request) {

	id := app.getID(w, r, getActorRe)
	err := app.people.Delete(id)

	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) editActor(w http.ResponseWriter, r *http.Request) {
	id := app.getID(w, r, getActorRe)
	p := storage.Person{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		app.serverError(w, err)
		return
	}

	p.ID = id
	err = app.people.Update(p)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) listActorsFilms(w http.ResponseWriter) {
	actorFilms, err := app.people.List()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data, err := json.MarshalIndent(actorFilms, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *application) getFilm(w http.ResponseWriter, r *http.Request) {

	id := app.getID(w, r, getFilmRe)
	film, err := app.films.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data, err := json.MarshalIndent(film, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *application) createFilm(w http.ResponseWriter, r *http.Request) {

	f := storage.Film{}
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = f.Validate()
	if err != nil {
		app.serverError(w, err)
	}

	lid, err := app.films.Insert(f)
	if err != nil {
		app.serverError(w, err)
		return
	}

	f.ID = lid
	data, err := json.MarshalIndent(f, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (app *application) deleteFilm(w http.ResponseWriter, r *http.Request) {

	id := app.getID(w, r, getFilmRe)
	err := app.films.Delete(id)

	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) editFilm(w http.ResponseWriter, r *http.Request) {
	id := app.getID(w, r, getFilmRe)
	f := storage.Film{}
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		app.serverError(w, err)
		return
	}

	f.ID = id
	err = app.films.Update(f)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) listFilms(w http.ResponseWriter, r *http.Request) {
	sortColumn, sortOrder := `rating`, `DESC`

	v := r.URL.Query()
	if value := v.Get("sortcolumn"); value != "" {
		sortColumn = value
	}
	if value := v.Get("sortorder"); value != "" {
		sortOrder = value
	}

	Films, err := app.films.List(sortColumn, sortOrder)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data, err := json.MarshalIndent(Films, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *application) searchFilm(w http.ResponseWriter, r *http.Request) {
	var pattern string

	var findFilms func(pattern string) ([]string, error)

	v := r.URL.Query()
	if title := v.Get("title"); title != "" {
		pattern = title
		findFilms = app.films.GetByTitle
	} else if name := v.Get("name"); name != "" {
		pattern = name
		findFilms = app.films.GetByActorName
	} else {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	titles, err := findFilms(pattern)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data, err := json.MarshalIndent(titles, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
