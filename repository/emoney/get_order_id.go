package emoney

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"mock-payment-provider/repository"
)

func (r *Repository) GetByOrderId(ctx context.Context, orderId string) (repository.Entry, error) {
	if orderId == "" {
		return repository.Entry{}, fmt.Errorf("orderId is empty")
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

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return repository.Entry{}, fmt.Errorf("rolling back transaction: %w", e)
		}

		return repository.Entry{}, fmt.Errorf("commiting transaction: %w", err)
	}

	return repository.Entry{}, nil
}
