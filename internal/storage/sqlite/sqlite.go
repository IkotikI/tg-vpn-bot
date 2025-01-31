package sqlite

import (
	"database/sql"
	"os"
	"vpn-tg-bot/pkg/e"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLStorage struct {
	db *sqlx.DB
}

func New(source string) (*SQLStorage, error) {
	db, err := sqlx.Open("sqlite3", source)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err = db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &SQLStorage{db: db}, nil
}

func (s *SQLStorage) Init() error {

	if err := s.createUsersTable(); err != nil {
		return e.Wrap("can't create users table", err)
	}

	if err := s.createServersTable(); err != nil {
		return e.Wrap("can't create servers table", err)
	}

	if err := s.createUsersServersTable(); err != nil {
		return e.Wrap("can't create users-servers table", err)
	}

	if err := s.createCountriesTable(); err != nil {
		return e.Wrap("can't create countries table", err)
	}

	if err := s.createServersAuthorizationsTable(); err != nil {
		return e.Wrap("can't create servers authorizations table", err)
	}

	return nil
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}

func (s *SQLStorage) SQLStorageInstance() (*sql.DB, error) {
	return s.db.DB, nil
}

func (s *SQLStorage) Drop(source string) error {
	return os.Remove(source)
}
