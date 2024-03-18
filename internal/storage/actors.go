package storage

import (
	"database/sql"
	"errors"
	"time"
)

type ActorStorage struct {
	DB *sql.DB
}

type Person struct {
	ID        int
	Name      string
	Gender    string
	BirthDate time.Time
}


// type ActorFilms map[string][]string

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

func (s *ActorStorage) Get(id int) (*Person, error) {
	stmt := `SELECT id, name, gender, birth_date FROM people WHERE id = $1`

	p := Person{}

	row := s.DB.QueryRow(stmt, id)

	err := row.Scan(&p.ID, &p.Name, &p.Gender, &p.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &p, nil
}

func (s *ActorStorage) Insert(p Person) (int, error) {
	stmt := `INSERT INTO people (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id`

	row := s.DB.QueryRow(stmt, p.Name, p.Gender, p.BirthDate)

	var lastInsertedId int
	err := row.Scan(&lastInsertedId)
	if err != nil {
		return 0, err
	}

	return int(lastInsertedId), nil
}

func (s *ActorStorage) Delete(id int) error {
	stmt := `DELETE FROM people where id = $1`

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

func (s *ActorStorage) Update(p Person) error {
	stmt := `UPDATE people SET (name, gender, birth_date) = ($2, $3, $4)
			WHERE id = $1 RETURNING $1`

	row := s.DB.QueryRow(stmt, p.ID, p.Name, p.Gender, p.BirthDate)
	err := row.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}
	return nil
}

func (s *ActorStorage) List() (map[string][]string, error) {
	stmt := `SELECT p.name, f.title FROM people AS p LEFT JOIN films_actors fa ON p.id = fa.actor_id
			LEFT JOIN films f ON f.id = fa.film_id
			ORDER BY p.id`

	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[string][]string)

	for rows.Next() {
		var name, title string
		var titleNull sql.NullString

		err = rows.Scan(&name, &titleNull)
		if err != nil {
			return nil, err
		}
		if titleNull.Valid{
			title = titleNull.String
		}
		res[name] = append(res[name], title)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	

	return res, nil
}

// package storage 

// import (
// 	"encoding/json"
// )

// func (s TitleType) MarhsalJSON() ([]byte, error) {
// 	if s.Valid {
// 		return json.Marshal(s.String)
// 	}
// 	return []byte{}, nil
// }