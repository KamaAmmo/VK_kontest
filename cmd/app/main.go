package main

import (
	"flag"
	"fmt"
	// "fmt"
	"log"
	"net/http"
	"os"
	"vk_app/cmd/config"
	// "vk_app/internal/storage"

	_ "github.com/lib/pq"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

// const (
// 	host     = os.Getenv("HOST")
// 	port     = os.Getenv("PORT")
// 	user     = os.Getenv("DB_USER")
// 	password = os.Getenv("DB_PASSWORD")
// 	dbname   = os.Getenv("DB_NAME")
// )

func main() {
	db_c := config.GetConfig()
	addr := flag.String("addr", ":5000", "HTTP network address")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", db_c.Host, db_c.Port, db_c.User, db_c.Password, db_c.DB_name)
	// params := flag.String("params", )
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	// _, err := storage.OpenDB(psqlInfo)
	fmt.Println(psqlInfo)

	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.route(),
	}
	fmt.Println("Getting started")
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
