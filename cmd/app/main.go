package main

import (
	"flag"
	"fmt"

	"log"
	"net/http"
	"os"
	"vk_app/cmd/config"
	"vk_app/internal/storage"

	_ "github.com/lib/pq"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	people   *storage.ActorStorage
	films    *storage.FilmStorage
}

func main() {
	db_c := config.GetConfig()

	addr := flag.String("addr", ":5000", "HTTP network address")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", db_c.Host, db_c.Port, db_c.User, db_c.Password, db_c.DB_name)

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	db, err := storage.OpenDB(psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
		people:   &storage.ActorStorage{DB: db},
		films:    &storage.FilmStorage{DB: db},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.route(),
	}
	infoLog.Println("Starting server on", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
