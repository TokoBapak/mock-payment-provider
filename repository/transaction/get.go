package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
)

func (r *Repository) GetByOrderId(ctx context.Context, orderId string) (primitive.Transaction, error) {
	if orderId == "" {
		return primitive.Transaction{}, fmt.Errorf("empty order id")
	}

	conn, err := r.db.Conn(ctx)
	if err != nil {
		return primitive.Transaction{}, fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log := zerolog.Ctx(ctx)
			log.Err(err).Msg("closing connection")
		}
	}()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return primitive.Transaction{}, fmt.Errorf("creating transaction: %w", err)
	}

	var transaction primitive.Transaction
	err = tx.QueryRowContext(
		ctx,
		`SELECT
    		order_id,
    		amount,
    		payment_type,
    		status,
    		expired_at,
    		created_at
		FROM
			transaction_log
		WHERE
			order_id = ?`,
		orderId,
	).Scan(
		&transaction.OrderId,
		&transaction.TransactionAmount,
		&transaction.PaymentType,
		&transaction.TransactionStatus,
		&transaction.ExpiresAt,
		&transaction.TransactionTime,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return primitive.Transaction{}, fmt.Errorf("rolling back transaction: %w", e)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return primitive.Transaction{}, repository.ErrNotFound
		}

		return primitive.Transaction{}, fmt.Errorf("querying row: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return primitive.Transaction{}, fmt.Errorf("rolling back transaction: %w", e)
		}

		return primitive.Transaction{}, fmt.Errorf("commiting transaction: %w", err)
	}

	return transaction, nil
}
