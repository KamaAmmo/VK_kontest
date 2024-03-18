package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"vk_app/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// @Summary getActor
// @Description get actor by id
// @ID get-actor-id
// @Tags actors
// @Param id path integer true "actor id"
// @Security ApiKeyAuth
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Router /actors/{id} [get]
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

// @Summary create actor
// @Description get actor by id
// @ID create-actor-id
// @Tags actors
// @Param input body storage.Person true "Actor personal data"
// @Security ApiKeyAuth
// @Accept json
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Router /actors/ [post]
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

// @Summary delete actor
// @Description get actor by id
// @ID delete-actor-id
// @Tags actors
// @Param id path integer true "actor id"
// @Security ApiKeyAuth
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Router /actors/{id} [delete]
func (app *application) deleteActor(w http.ResponseWriter, r *http.Request) {

	id := app.getID(w, r, getActorRe)

	err := app.filmsActors.DeleteByActor(id)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.people.Delete(id)
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

// @Summary edit actor
// @Description edit actor
// @ID edit-actor-id
// @Tags actors
// @Param data body storage.Person true "new data about actor"
// @Security ApiKeyAuth
// @Accept json
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Router /actors/ [put]
func (app *application) editActor(w http.ResponseWriter, r *http.Request) {
	p := storage.Person{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		app.serverError(w, err)
		return
	}

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

// @Summary list actors and films
// @ID list-actor-film-id
// @Security ApiKeyAuth
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Router /actors/ [get]
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

// @Summary getFilm
// @Description get actor by id
// @ID get-film-id
// @Tags films
// @Param id path integer true "film id"
// @Security ApiKeyAuth
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Router /films/{id} [get]
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

// @Summary create film
// @ID create-film-id
// @Tags films
// @Param input body storage.Film true "Data about film"
// @Security ApiKeyAuth
// @Accept json
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Router /films/ [post]
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

	for _, name := range f.Cast {
		idActor, err := app.people.InsertOnlyName(name)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.filmsActors.Insert(idActor, lid)
		if err != nil {
			app.serverError(w, err)
			return
		}
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

// @Summary delete film
// @ID delete-film-id
// @Tags films
// @Param id path integer true "film id"
// @Security ApiKeyAuth
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Router /films/{id} [delete]
func (app *application) deleteFilm(w http.ResponseWriter, r *http.Request) {

	id := app.getID(w, r, getFilmRe)

	err := app.filmsActors.DeleteByFilm(id)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.films.Delete(id)

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

// @Summary edit film
// @ID edit-film-id
// @Tags films
// @Param input body storage.Film true "Data about film"
// @Security ApiKeyAuth
// @Accept json
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Router /films/ [put]
func (app *application) editFilm(w http.ResponseWriter, r *http.Request) {
	f := storage.Film{}
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		app.serverError(w, err)
		return
	}

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

// @Summary listFilms
// @Description get movies with actors
// @ID list-films-id
// @Tags films
// @Security ApiKeyAuth
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Router /films/ [get]
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

// @Summary searchFilmByPattern
// @Description get a list of films by a fragment of the film title or actor's name
// @ID search-film-id
// @Tags films
// @Security ApiKeyAuth
// @Param title query string false "by title movie"
// @Param name query string false "by actor name"
// @Procuce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Router /films [get]
func (app *application) searchFilmByPattern(w http.ResponseWriter, r *http.Request) {
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

// @Summary registerUser
// @Description register new user
// @Tags registration and logging
// @Param input body RegisterRequest true "New user info"
// @Accept json
// @Produce json
// @success 201 {integer} integer "New user registered"
// @Failure 405,409 {object} error
// @Router /register [post]
func (app application) registerUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	u := RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if u.Role == "" {
		u.Role = "user"
	} else if u.Role != "user" && u.Role != "admin" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.users.Insert(u.Name, u.Password, u.Role)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateUserName) {
			app.clientError(w, http.StatusBadRequest)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary loginUser
// @Description login existing user
// @Tags registration and logging
// @Param input body LoginRequest true "New user info"
// @Accept json
// @Produce json
// @success 200 {integer} integer "Succesfuly logged in
// @Failure 405,409 {object} error
// @Router /login [post]
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	logReq := LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&logReq)
	if err != nil {
		app.serverError(w, err)
	}

	u, err := app.users.Get(logReq.Name)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			app.clientError(w, http.StatusUnauthorized)
		} else {
			app.serverError(w, err)
		}
		return
	}

	if err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(logReq.Password)); err != nil {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	payload := jwt.MapClaims{
		"sub":  u.Name,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
		"role": u.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(secretKey)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data, err := json.MarshalIndent(LoginToken{AccessToken: t}, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
