package transaction

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
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
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
		`CREATE TABLE IF NOT EXISTS transaction_log (
    		order_id TEXT PRIMARY KEY,
    		amount INT NOT NULL,
    		payment_type INT NOT NULL,
    		status INT NOT NULL,
    		expired_at TEXT NOT NULL,
    		created_at TEXT NOT NULL,
    		updated_at TEXT NOT NULL
		)`)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("rolling back transaction: %w", err)
		}

		return fmt.Errorf("executing create table transaction log: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return fmt.Errorf("rolling back transaction: %w", e)
		}

		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
