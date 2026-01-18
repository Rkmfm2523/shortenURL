package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("err in storage ni new func: %s", err)
	}

	stmt, err := db.Prepare(`
	create table if not exists url(
	id integer,
	alias text not null ubique,
	url text not null);
	create index if not exists idx_alias on (url(alias));
	`)

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("error stmt in new func")
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("insert into url(url, alias) values (?, ?)")
	if err != nil {
		return fmt.Errorf("error in save url: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		return fmt.Errorf("error in save url: %s", err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("select url from url where alias = ? ")
	if err != nil {
		return "", fmt.Errorf("error get response this alias: %s", err)
	}
	defer stmt.Close()

	var resString string

	err = stmt.QueryRow(alias).Scan(&resString)
	if err != nil {
		return "", fmt.Errorf("error exec this url or alias: %s", err)
	}

	return resString, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	v, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("affected err: %s", err)
	}
	if v == 0 {
		return errors.New("No url deleted")
	}

	return nil
}
