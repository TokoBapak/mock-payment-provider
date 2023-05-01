package virtual_account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

func (r *Repository) CreateCharge(ctx context.Context, virtualAccountNumber string, orderId string, amount int64, expiresAt time.Time) (account string, err error) {
	if orderId == "" {
		return "", fmt.Errorf("orderId is empty")
	}

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
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		return "", fmt.Errorf("creating transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO virtual_account_entries
			(
			 	order_id,
			 	virtual_account_number,
			 	amount,
			 	expired_at,
			 	created_at,
			 	updated_at
			)
			VALUES 
				(?, ?, ?, ?, ?, ?)`,
		orderId,
		virtualAccountNumber,
		amount,
		expiresAt,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return "", fmt.Errorf("rolling back transaction: %w", err)
		}

		return "", fmt.Errorf("executing query: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE OR FAIL
    		virtual_accounts
		SET
		    current_order_id = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE
		    virtual_account_number = ?`,
		orderId,
		virtualAccountNumber,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return "", fmt.Errorf("rolling back transaction: %w", err)
		}

		return "", fmt.Errorf("executing query: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return "", fmt.Errorf("rolling back transaction: %w", err)
		}

		return "", fmt.Errorf("commiting transaction: %w", err)
	}

	return virtualAccountNumber, nil
}
