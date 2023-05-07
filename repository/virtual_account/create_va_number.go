package virtual_account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

// CreateOrGetVirtualAccountNumber will accept the incoming customerUniqueField (can be anything ranging from
// customer's unique ID to customer's email that's supposed to be unique). If the customerUniqueField haven't
// been registered or submitted before, we will create a new virtual account number. Otherwise, we will
// retrieve the virtual account number directly that's supposed to be exists.
func (r *Repository) CreateOrGetVirtualAccountNumber(ctx context.Context, customerUniqueField string) (string, error) {
	if customerUniqueField == "" {
		return "", fmt.Errorf("customerUniqueField is empty")
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
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return "", fmt.Errorf("creating transaction: %w", err)
	}

	var virtualAccountNumber string
	err = tx.QueryRowContext(
		ctx,
		`SELECT virtual_account_number FROM virtual_accounts WHERE unique_identifier = ?`,
		customerUniqueField,
	).Scan(&virtualAccountNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create one
			virtualAccountNumber = generateVirtualAccountNumber()

			_, err := tx.ExecContext(
				ctx,
				`INSERT INTO
					virtual_accounts
					(
					 	unique_identifier,
					 	virtual_account_number,
					 	current_order_id,
					 	created_at,
					 	updated_at
					)
				VALUES 
					(?, ?, NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
				customerUniqueField,
				virtualAccountNumber,
			)
			if err != nil {
				if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
					return "", fmt.Errorf("rolling back transaction: %w", err)
				}

				return "", fmt.Errorf("executing query: %w", err)
			}
		} else {
			if e := tx.Rollback(); e != nil && !errors.Is(err, sql.ErrTxDone) {
				return "", fmt.Errorf("rolling back transaction: %w", err)
			}

			return "", fmt.Errorf("executing query: %w", err)
		}
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
