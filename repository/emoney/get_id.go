package emoney

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"mock-payment-provider/repository"
)

func (r *Repository) GetByID(ctx context.Context, id string) (repository.Entry, error) {
	if id == "" {
		return repository.Entry{}, fmt.Errorf("id is empty")
	}

	conn, err := r.db.Conn(ctx)
	if err != nil {
		return repository.Entry{}, fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil && !errors.Is(err, sql.ErrConnDone) {
			log := zerolog.Ctx(ctx)
			log.Err(err).Msg("returning connection back to pool")
		}
	}()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return repository.Entry{}, fmt.Errorf("creating transaction: %w", err)
	}

	var entry repository.Entry
	err = tx.QueryRowContext(
		ctx,
		`SELECT
    		order_id,
    		id,
    		amount,
    		expired_at
		FROM
		    emoney_entries
		WHERE
			id = ?`,
		id,
	).Scan(
		&entry.OrderId,
		&entry.EMoneyID,
		&entry.ChargedAmount,
		&entry.ExpiresAt,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return repository.Entry{}, fmt.Errorf("rolling back transaction: %w", e)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return repository.Entry{}, repository.ErrNotFound
		}

		return repository.Entry{}, fmt.Errorf("executing query: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return repository.Entry{}, fmt.Errorf("rolling back transaction: %w", e)
		}

		return repository.Entry{}, fmt.Errorf("commiting transaction: %w", err)
	}

	return entry, nil
}
