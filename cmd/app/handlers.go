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

	lid, err := app.people.Insert(p.Name, p.Gender, p.BirthDate)
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

func (app *application) listActorsFilms(w http.ResponseWriter, r *http.Request) {
	actorFilms, err := app.people.List()
	if err != nil{
		app.serverError(w, err)
		return
	}
	data, err := json.MarshalIndent(actorFilms, "", "	")
	if err != nil{
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
