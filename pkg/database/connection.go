package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
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
	db := &DB{
		Now: time.Now,
		log: logger,
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())

	pgxConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pgxConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// Ping the connection to see if it is still valid. Ping returns an error if
		// it fails.
		return conn.Ping(ctx) == nil
	}

	db.Pool, err = pgxpool.ConnectConfig(db.ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	db.cancel()
	if db.Pool != nil {
		return nil
	}
	return nil
}
