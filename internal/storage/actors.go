package storage

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

type ActorStorage struct {
	DB *sql.DB
}

type Person struct {
	ID        int      `json:"id"`
	Name      string   `json:"name" example:"John Doe"`
	Gender    string   `json:"gender" example:"Male"`
	BirthDate JsonDate `json:"birthDate" example:"1970-01-01"`
}

func (s *ActorStorage) Get(id int) (Person, error) {
	stmt := `SELECT id, name, gender, birth_date FROM people WHERE id = $1`

	row := s.DB.QueryRow(stmt, id)

	p := Person{}
	var NullGender sql.NullString
	var NullBDate sql.NullTime
	err := row.Scan(&p.ID, &p.Name, &NullGender, &NullBDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Person{}, ErrNoRecord
		}
		return Person{}, err
	}

	if NullGender.Valid {
		p.Gender = NullGender.String
	}
	if NullBDate.Valid {
		p.BirthDate = JsonDate(NullBDate.Time)
	}

	return p, nil
}

func (s *ActorStorage) Insert(p Person) (int, error) {
	stmt := `INSERT INTO people (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id`

	row := s.DB.QueryRow(stmt, p.Name, p.Gender, time.Time(p.BirthDate))

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

	row := s.DB.QueryRow(stmt, p.ID, p.Name, p.Gender, time.Time(p.BirthDate))
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
	stmt := `SELECT p.id, p.name, f.title FROM people AS p LEFT JOIN films_actors fa ON p.id = fa.actor_id
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
		var id int
		var titleNull sql.NullString

		err = rows.Scan(&id, &name, &titleNull)
		if err != nil {
			return nil, err
		}
		if titleNull.Valid {
			title = titleNull.String
		}
		name = strconv.Itoa(id) + " - " + name
		res[name] = append(res[name], title)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ActorStorage) InsertOnlyName(name string) (int, error) {
	stmt := `
		WITH inserted_person AS (
		INSERT INTO people (name)
		SELECT $1  
		WHERE NOT EXISTS (
			SELECT id
			FROM people
			WHERE name = $1
		)
		RETURNING id
		)
		SELECT id FROM inserted_person
		UNION ALL
		SELECT id FROM people WHERE name = $1
		`
	row := s.DB.QueryRow(stmt, name)

	var lastInsertedId int
	err := row.Scan(&lastInsertedId)
	if err != nil {
		return 0, err
	}

	return int(lastInsertedId), nil
}
