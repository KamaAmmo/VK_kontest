package main

import (
	"flag"
	"fmt"
	"log"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":5000", "HTTP network address")
	//db

	flag.Parse()


	srv := 
}
