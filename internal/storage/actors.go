package storage

import (
	"database/sql"
	"errors"
	"time"
)

type PersonStorage struct {
	DB *sql.DB
}

type Person struct {
	ID        int
	Name      string
	Gender    string
	BirthDate time.Time
}

func OpenDB(params string) (*sql.DB, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db, nil
}

func (s *PersonStorage) Get(id int) (*Person, error) {
	stmt := `SELECT id, name, gender, birth_date FROM people WHERE id = $1`

	p := Person{}

	row := s.DB.QueryRow(stmt, id)

	err := row.Scan(&p.ID, &p.Name, &p.Gender, &p.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &p, nil

}
