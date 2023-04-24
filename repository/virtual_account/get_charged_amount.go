package virtual_account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"mock-payment-provider/repository"
)

func (r *Repository) GetChargedAmount(ctx context.Context, virtualAccountNumber string) (int64, error) {
	if virtualAccountNumber == "" {
		return 0, fmt.Errorf("virtualAccountNumber is empty")
	}

	conn, err := r.db.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil && !errors.Is(err, sql.ErrConnDone) {
			log.Printf("returning connection back to pool: %s", err.Error())
		}
	}()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return 0, fmt.Errorf("creating transaction: %w", err)
	}

	var currentOrderId sql.NullString
	err = tx.QueryRowContext(
		ctx,
		`SELECT current_order_id FROM virtual_accounts WHERE virtual_account_number = ?`,
		virtualAccountNumber,
	).Scan(&currentOrderId)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return 0, fmt.Errorf("rolling back transaction: %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return 0, repository.ErrNotFound
		}

		return 0, fmt.Errorf("executing query: %w", err)
	}

	var chargedAmount int64
	err = tx.QueryRowContext(
		ctx,
		`SELECT amount FROM virtual_account_entries WHERE order_id = ?`,
		currentOrderId.String,
	).Scan(
		&chargedAmount,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return 0, fmt.Errorf("rolling back transaction: %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return 0, repository.ErrNotFound
		}

		return 0, fmt.Errorf("executing query: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return 0, fmt.Errorf("rolling back transaction: %w", err)
		}

		return 0, fmt.Errorf("commiting transaction: %w", err)
	}

	return chargedAmount, nil
}
