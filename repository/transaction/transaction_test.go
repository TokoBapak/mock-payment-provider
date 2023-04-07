package transaction_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"mock-payment-provider/repository/transaction"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error = nil
	db, err = sql.Open("sqlite3", ":memory:?_txlock=exclusive&_foreign_keys=1&")
	if err != nil {
		log.Fatalf("Opening sql database: %s", err.Error())
	}

	setupCtx, setupCancel := context.WithTimeout(context.Background(), time.Minute)
	defer setupCancel()

	repository, err := transaction.NewTransactionRepository(db)
	if err != nil {
		log.Fatalf("Creating transaction repository: %s", err.Error())
	}

	err = repository.Migrate(setupCtx)
	if err != nil {
		log.Fatalf("Migrating database: %s", err.Error())
	}

	exitCode := m.Run()

	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), time.Minute)
	defer cleanupCancel()

	err = destroy(cleanupCtx, db)
	if err != nil {
		log.Printf("Destroying tables: %s", err.Error())
	}

	err = db.Close()
	if err != nil {
		log.Printf("Closing database: %s", err.Error())
	}

	os.Exit(exitCode)
}

func destroy(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
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
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DROP TABLE IF EXISTS transaction_log`)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("rolling back transaction: %w", err)
		}

		return fmt.Errorf("executing create table transaction log: %w", err)
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
