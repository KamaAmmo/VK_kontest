package main 

import (
	// "fmt"
	"net/http"
)


func (app *application) route() http.Handler{
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	
	return mux
}