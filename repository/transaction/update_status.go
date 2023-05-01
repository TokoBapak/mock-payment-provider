package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"mock-payment-provider/primitive"
)

func (r *Repository) UpdateStatus(ctx context.Context, orderId string, status primitive.TransactionStatus) error {
	if orderId == "" {
		return fmt.Errorf("empty order id")
	}

	conn, err := r.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection from pool: %w", err)
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
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE transaction_log SET status = ? WHERE order_id = ?`,
		status,
		orderId,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("rolling back transaction: %w", e)
		}

		return fmt.Errorf("executing update statement: %w", err)
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
