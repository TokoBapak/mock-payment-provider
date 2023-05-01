package virtual_account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

func (r *Repository) DeductCharge(ctx context.Context, virtualAccountNumber string) error {
	if virtualAccountNumber == "" {
		return fmt.Errorf("empty virtual account number")
	}

	conn, err := r.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection from pool: %w", err)
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
		return fmt.Errorf("creating transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE virtual_accounts SET current_order_id = NULL WHERE virtual_account_number = ?`,
		virtualAccountNumber,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return fmt.Errorf("rolling back transaction: %w", err)
		}

		return fmt.Errorf("executing query: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
			return fmt.Errorf("rolling back transaction: %w", err)
		}

		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
