package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrMissingDSN = errors.New("db: missing dsn")

func NewConnection(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, ErrMissingDSN
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %v", err)
	}

	return db, nil
}
