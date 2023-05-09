package emoney_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"mock-payment-provider/repository/emoney"
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

	emoneyRepository, err := emoney.NewEmoneyRepository(db)
	if err != nil {
		log.Fatalf("Creating emoney repository: %s", err.Error())
	}

	err = emoneyRepository.Migrate(setupCtx)
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

func TestNewEmoneyRepository(t *testing.T) {
	t.Run("Nil Database", func(t *testing.T) {
		_, err := emoney.NewEmoneyRepository(nil)
		if err.Error() != "db is nil" {
			t.Errorf("expecting an error of 'db is nil', instead got %s", err.Error())
		}
	})

	t.Run("Happy Case", func(t *testing.T) {
		repository, err := emoney.NewEmoneyRepository(&sql.DB{})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if repository == nil {
			t.Errorf("nil repository")
		}
	})
}
