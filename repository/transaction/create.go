package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"mock-payment-provider/repository"
)

func (r *Repository) Create(ctx context.Context, params repository.CreateTransactionParam) error {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("returning connection back to pool: %s", err.Error())
		}
	}()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO
			transaction_log
			(
				 order_id,
				 amount,
				 payment_type,
				 status,
				 expired_at,
				 created_at,
				 updated_at
			)
		VALUES 
			(?, ?, ?, ?, ?, ?, ?)`,
		params.OrderID,
		params.Amount,
		params.PaymentType,
		params.Status,
		params.ExpiredAt,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("rolling back transaction: %w", e)
		}

		return fmt.Errorf("executing insert statement: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
