package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"strconv"
)

type FilmStorage struct {
	DB *sql.DB
}

type Film struct {
	ID          int      `json:"id"`
	Title       string   `json:"title" example:"Fight Club"`
	Description string   `json:"description" example:"The first rule of Fight Club is"`
	ReleaseDate JsonDate `json:"release_date" example:"1999-09-10"`
	Rating      float32  `json:"rating" example:"8.7"`
	Cast        []string `json:"cast" example:"Alex Lesly"`
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

	row := s.DB.QueryRow(stmt, f.Title, f.Description, time.Time(f.ReleaseDate), f.Rating)

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

	row := s.DB.QueryRow(stmt, f.ID, f.Title, f.Description, time.Time(f.ReleaseDate), f.Rating)
	err := row.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}
	return nil
}

func (s *FilmStorage) List(sortColumn, sortOrder string) ([]string, error) {

	stmt := fmt.Sprintf("SELECT id, title FROM films ORDER BY %s %s", sortColumn, sortOrder)

	rows, err := s.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	result := []string{}
	for rows.Next() {
		var title string
		var id int
		err = rows.Scan(&id, &title)
		if err != nil {
			return nil, err
		}
		title = strconv.Itoa(id) + " : " + title
		result = append(result, title)
	}
	return result, nil
}

func (s *FilmStorage) GetByTitle(title string) ([]string, error) {
	pattern := "%" + title + "%"
	pattern = strings.ToLower(pattern)
	stmt := fmt.Sprintf(`SELECT DISTINCT f.title FROM films f WHERE LOWER(f.title) LIKE '%s'`, pattern)

	rows, err := s.DB.Query(stmt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	result := []string{}
	for rows.Next() {

		var title string
		err := rows.Scan(&title)
		if err != nil {
			return nil, err
		}
		result = append(result, title)
	}

	return result, nil
}

func (s *FilmStorage) GetByActorName(name string) ([]string, error) {
	pattern := "%" + name + "%"
	pattern = strings.ToLower(pattern)
	stmt := fmt.Sprintf(`SELECT DISTINCT f.title FROM films f LEFT JOIN films_actors fa ON fa.film_id = f.id 
		LEFT JOIN people p ON p.id = fa.actor_id WHERE LOWER(p.name) LIKE '%s'`, pattern)
	fmt.Println(stmt)

	rows, err := s.DB.Query(stmt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	result := []string{}
	for rows.Next() {

		var title string
		err := rows.Scan(&title)
		if err != nil {
			return nil, err
		}
		result = append(result, title)
	}

	return result, nil
}
