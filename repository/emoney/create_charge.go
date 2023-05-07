package emoney

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (r *Repository) CreateCharge(ctx context.Context, orderId string, amount int64, expiresAt time.Time) (id string, err error) {
	id = uuid.NewString()

	conn, err := r.db.Conn(ctx)
	if err != nil {
		return "", fmt.Errorf("acquiring connection from pool: %w", err)
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
		ReadOnly:  false,
	})
	if err != nil {
		return "", fmt.Errorf("creating transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO
			emoney_entries
			(
			 	order_id,
			 	id,
			 	amount,
			 	expired_at,
			 	created_at,
			 	updated_at
			)
		VALUES 
			(?, ?, ?, ?, ?, ?)`,
		orderId,
		id,
		amount,
		expiresAt,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return "", fmt.Errorf("rolling back transaction: %w", e)
		}

		return "", fmt.Errorf("executing query: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return "", fmt.Errorf("rolling back transaction: %w", e)
		}

		return "", fmt.Errorf("commiting transaction: %w", err)
	}

	return id, nil
}
