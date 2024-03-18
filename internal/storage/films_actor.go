package storage

import (
	"database/sql"
)

type FilmActor struct {
	FilmId  int
	ActorId int
}

type FilmActorStorage struct{
	DB *sql.DB
}

func (s *FilmActorStorage) Insert(idActor, idFilm int) error {
	stmt := `INSERT INTO films_actors (film_id, actor_id) VALUES ($1, $2)`

	_, err := s.DB.Exec(stmt, idFilm, idActor)
	if err != nil {
		return err
	}

	return nil
}

func (s *FilmActorStorage) DeleteByActor(id int) error {
	stmt := `DELETE FROM films_actors where actor_id = $1`

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

func (s *FilmActorStorage) DeleteByFilm(id int) error {
	stmt := `DELETE FROM films_actors where film_id = $1`

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