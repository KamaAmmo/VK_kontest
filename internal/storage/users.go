package storage

import (
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	DB *sql.DB
}

type User struct {
	Name           string
	Role           string
	HashedPassword []byte
}

func (s *UserStorage) Insert(name, pass, role string) error {
	stmt := `INSERT INTO users (name, password, role) VALUES($1, $2, $3)`

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		return err
	}

	name = strings.TrimSpace(name)
	role = strings.ToLower(strings.TrimSpace(role))
	_, err = s.DB.Exec(stmt, name, string(hashedPass), role)
	if err != nil {
		if IsErrorCode(err, UniqueViolationErr) {
			return ErrDuplicateUserName
		}
		return err
	}
	return nil
}

func (s *UserStorage) Get(name string) (User, error) {
	stmt := `SELECT name, password, role from users where name = $1 `

	u := User{}
	err := s.DB.QueryRow(stmt, name).Scan(&u.Name, &u.HashedPassword, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		} else {
			return User{}, err
		}
	}

	return u, nil
}

func (s *UserStorage) Delete(id int) error {
	stmt := `DELETE FROM users WHERE id = $1`

	_, err := s.DB.Exec(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}

	return nil
}
