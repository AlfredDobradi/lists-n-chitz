package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // driver for postgres
)

func New(uri string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Try(err error, tx *sqlx.Tx) error {
	if err == nil {
		err = tx.Commit()
	} else {
		err = tx.Rollback()
	}

	return err
}
