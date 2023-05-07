package virtual_account_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"mock-payment-provider/repository/virtual_account"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error = nil
	db, err = sql.Open("sqlite3", ":memory:?_txlock=exclusive&_foreign_keys=1&")
	if err != nil {
		log.Fatalf("Opening sql database: %s", err.Error())
	}

	db.SetMaxOpenConns(1)

	setupCtx, setupCancel := context.WithTimeout(context.Background(), time.Minute)
	defer setupCancel()

	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		log.Fatalf("Creating virtual account repository: %s", err.Error())
	}

	err = virtualAccountRepository.Migrate(setupCtx)
	if err != nil {
		log.Fatalf("migrating database: %s", err.Error())
	}

	exitCode := m.Run()

	err = db.Close()
	if err != nil {
		log.Printf("Closing database: %s", err.Error())
	}

	os.Exit(exitCode)
}

func TestNewVirtualAccountRepository(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		repository, err := virtual_account.NewVirtualAccountRepository(&sql.DB{})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if repository == nil {
			t.Errorf("expecting repository to be not nil, got nil instead")
		}
	})

	t.Run("NilDatabase", func(t *testing.T) {
		_, err := virtual_account.NewVirtualAccountRepository(nil)
		if err == nil {
			t.Errorf("expecting an error, got nil instead")
		}

		if err.Error() != "db is nil" {
			t.Errorf("expecting an error of 'db is nil', instead got %s", err.Error())
		}
	})
}
