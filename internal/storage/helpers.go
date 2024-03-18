package storage

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

type JsonDate time.Time

func (j *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	*j = JsonDate(t)
	return nil
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j).Format(time.DateOnly))
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
