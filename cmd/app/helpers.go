package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) clientError(w http.ResponseWriter, Status int) {
	http.Error(w, http.StatusText(Status), Status)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) notFound(w http.ResponseWriter){
	app.clientError(w, http.StatusNotFound)
}




