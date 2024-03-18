package main

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
)

func (app *application) clientError(w http.ResponseWriter, Status int) {
	http.Error(w, http.StatusText(Status), Status)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// некрасиво, но мне не хотелось дублировать код
func (app *application) getID(w http.ResponseWriter, r *http.Request, re *regexp.Regexp) int {
	match := re.FindStringSubmatch(r.URL.Path)
	if len(match) < 2 {
		app.notFound(w)
		return 0
	}
	id, err := strconv.Atoi(match[1])

	if err != nil || id < 1 {
		app.notFound(w)
		return 0
	}

	return id
}




