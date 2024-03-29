package emoney

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

func (r *Repository) Migrate(ctx context.Context) error {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil && !errors.Is(err, sql.ErrConnDone) {
			log := zerolog.Ctx(ctx)
			log.Err(err).Msg("returning connection back to pool")
		}
	}()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS emoney_entries (
			order_id TEXT PRIMARY KEY,
    		id TEXT NOT NULL,
			amount INT NOT NULL,
			expired_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return fmt.Errorf("rolling back transaction: %w", err)
		}

		return fmt.Errorf("executing query: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE INDEX IF NOT EXISTS idx_emoney_entries_id ON emoney_entries (id)`,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return fmt.Errorf("rolling back transaction: %w", err)
		}

		return fmt.Errorf("executing query: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return fmt.Errorf("rolling back transaction: %w", err)
		}

		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
