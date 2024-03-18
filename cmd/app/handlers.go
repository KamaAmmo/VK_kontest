package main

import (
	// "fmt"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"vk_app/internal/storage"
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
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

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
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
