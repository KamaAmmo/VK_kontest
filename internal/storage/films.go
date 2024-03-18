package storage 

import (
	"database/sql"
)


type FilmStorage struct{
	DB *sql.DB
}

