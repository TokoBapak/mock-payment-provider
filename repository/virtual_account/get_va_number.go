package virtual_account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"mock-payment-provider/repository"
)

func (r *Repository) GetByVirtualAccountNumber(ctx context.Context, virtualAccountNumber string) (repository.Entry, error) {
	if virtualAccountNumber == "" {
		return repository.Entry{}, fmt.Errorf("virtualAccountNumber is empty")
	}

	conn, err := r.db.Conn(ctx)
	if err != nil {
		return repository.Entry{}, fmt.Errorf("acquiring connection from pool: %w", err)
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
		return repository.Entry{}, fmt.Errorf("creating transaction: %w", err)
	}

	var currentOrderId string
	err = tx.QueryRowContext(
		ctx,
		`SELECT current_order_id FROM virtual_accounts WHERE virtual_account_number = ?`,
		virtualAccountNumber,
	).Scan(&currentOrderId)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return repository.Entry{}, fmt.Errorf("rolling back transaction: %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return repository.Entry{}, repository.ErrNotFound
		}

		return repository.Entry{}, fmt.Errorf("executing query: %w", err)
	}

	var entry repository.Entry
	err = tx.QueryRowContext(
		ctx,
		`SELECT order_id, amount, expired_at FROM virtual_account_entries WHERE order_id = ?`,
		currentOrderId,
	).Scan(
		&entry.OrderId,
		&entry.ChargedAmount,
		&entry.ExpiresAt,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return repository.Entry{}, fmt.Errorf("rolling back transaction: %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return repository.Entry{}, repository.ErrNotFound
		}

		return repository.Entry{}, fmt.Errorf("executing query: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return repository.Entry{}, fmt.Errorf("rolling back transaction: %w", err)
		}

		return repository.Entry{}, fmt.Errorf("commiting transaction: %w", err)
	}

	return entry, nil
}
