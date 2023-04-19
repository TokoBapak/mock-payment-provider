package virtual_account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

func (r *Repository) CreateCharge(ctx context.Context, orderId string, amount int64, expiresAt time.Time) (account string, err error) {
	if orderId == "" {
		return "", fmt.Errorf("orderId is empty")
	}

	// Generate a virtual account number
	var virtualAccountNumber string
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return "", fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil && !errors.Is(err, sql.ErrConnDone) {
			log.Printf("closing connection back to pool: %s", err.Error())
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
		`INSERT INTO virtual_accounts
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

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return "", fmt.Errorf("rolling back transaction: %w", err)
		}

		return "", fmt.Errorf("commiting transaction: %w", err)
	}

	return "", nil
}
