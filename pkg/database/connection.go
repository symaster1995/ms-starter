package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type DB struct {
	ctx    context.Context
	cancel func()
	log    *zerolog.Logger
	Pool   *pgxpool.Pool
	DSN    string
	Now    func() time.Time
}

func NewDB(dbUrl string, logger *zerolog.Logger) (*DB, error) {

	pgxConfig, err := pgxpool.ParseConfig(dbUrl)

	db := &DB{
		Now: time.Now,
		log: logger,
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())

	return db, nil
}

func (db *DB) Open() error {

	//pgxConfig, err := pgxpool.ParseConfig()

	dbPool, err := pgxpool.Connect(db.ctx, db.DSN)
	//dbPool, err := pgxpool.ConnectConfig(db.ctx, db.DSN)
	if err != nil {
		db.log.Error().Err(err).Msg("unable to connect to database")
		os.Exit(1)
	}

	defer dbPool.Close()

	if err := dbPool.Ping(db.ctx); err != nil {
		return err
	}
	db.Pool = dbPool

	/*conn, err := dbPool.Acquire(db.ctx)*/

	return nil
}

func (db *DB) Close() error {
	db.cancel()
	if db.Pool != nil {
		return nil
	}
	return nil
}
