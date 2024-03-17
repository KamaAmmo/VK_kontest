package storage

import (
	"database/sql"
)

func OpenDB(params string) (*sql.DB, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db, nil
}
