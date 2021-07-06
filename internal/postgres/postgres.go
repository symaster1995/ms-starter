package postgres

import (
	"context"
	"database/sql"
	"time"
)

type DB struct {
	db     *sql.DB
	ctx    context.Context
	cancel func()

	DSN string

	Now func() time.Time
}

func NewDB(dsn string) *DB {
	db := &DB{
		Now: time.Now,
		DSN: dsn,
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())

	return db
}

func (db *DB) Open() (err error) {
	if db.db, err = sql.Open("postgres", db.DSN); err != nil {
		return err
	}
	defer db.Close()

	//ping to see if database connection is working
	if err = db.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() error {
	db.cancel()

	if db.db != nil {
		return db.db.Close()
	}

	return nil
}
