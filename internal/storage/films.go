package storage

import (
	"database/sql"
	"errors"
	"time"
)

type FilmStorage struct {
	DB *sql.DB
}

type Film struct {
	ID          int
	Title       string
	Description string
	ReleaseDate time.Time
	Rating      float32
}

func (s *FilmStorage) Get(id int) (*Film, error) {
	stmt := `SELECT id, title, description, release_date, rating FROM films WHERE id = $1`

	f := Film{}

	row := s.DB.QueryRow(stmt, id)

	err := row.Scan(&f.ID, &f.Title, &f.Description, &f.ReleaseDate, &f.Rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &f, nil
}

func (s *FilmStorage) Insert(f Film) (int, error) {
	stmt := `INSERT INTO films (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id`

	row := s.DB.QueryRow(stmt, f.Title, f.Description, f.ReleaseDate, f.Rating)

	var lastInsertedId int
	err := row.Scan(&lastInsertedId)
	if err != nil {
		return 0, err
	}

	return int(lastInsertedId), nil
}

func (s *FilmStorage) Delete(id int) error {
	stmt := `DELETE FROM films where id = $1`

	res, err := s.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNoRecord
	}

	return nil
}

func (s *FilmStorage) Update(f Film) error {
	stmt := `UPDATE films SET (title, description, release_date, rating) = ($2, $3, $4, $5)
			WHERE id = $1 RETURNING $1`

	row := s.DB.QueryRow(stmt, f.ID, f.Title, f.Description, f.ReleaseDate, f.Rating)
	err := row.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}
	return nil
}
