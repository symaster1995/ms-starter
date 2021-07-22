package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type Tx struct {
}

func (db *DB) BeginTransaction(ctx context.Context, isoLevel pgx.TxIsoLevel) error {

	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel})

	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(db.ctx)

	_, err = tx.Exec(db.ctx,"insert into foo(id) values (1)")
	if err != nil {
		return err
	}

	err = tx.Commit(db.ctx)
	if err != nil {
		return err
	}

	return nil
}
