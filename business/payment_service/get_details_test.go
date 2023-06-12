package payment_service_test

import (
	"context"
	"database/sql"
	"errors"
	"mock-payment-provider/business"
	"mock-payment-provider/business/payment_service"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
	"mock-payment-provider/repository/emoney"
	"mock-payment-provider/repository/transaction"
	"mock-payment-provider/repository/virtual_account"
	"mock-payment-provider/repository/webhook"
	"os"
	"testing"
	"time"

	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db  *sql.DB
	cfg config
)

func TestMain(m *testing.M) {
	var err error = nil
	// db, err = sql.Open("sqlite3", ":memory:?_txlock=exclusive&_foreign_keys=1&")
	// if err != nil {
	// 	log.Fatalf("Opening sql database: %s", err.Error())
	// }
	cfg = parseConfig()
	db, err = sql.Open("sqlite3", cfg.databasePath)
	if err != nil {
		log.Fatalf("Opening sql database: %s", err.Error())
	}

	db.SetMaxOpenConns(1)

	setupCtx, setupCancel := context.WithTimeout(context.Background(), time.Minute)
	defer setupCancel()

	transactionRepository, err := transaction.NewTransactionRepository(db)
	if err != nil {
		log.Fatalf("Creating transaction repository: %s", err.Error())
	}
	err = transactionRepository.Migrate(setupCtx)
	if err != nil {
		log.Fatalf("migrating transaction repository: %s", err.Error())
	}

	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		log.Fatalf("Creating virtual account repository: %s", err.Error())
	}
	err = virtualAccountRepository.Migrate(setupCtx)
	if err != nil {
		log.Fatalf("migrating virtual account repository: %s", err.Error())
	}

	emoneyRepository, err := emoney.NewEmoneyRepository(db)
	if err != nil {
		log.Fatalf("Creating emoney repository: %s", err.Error())
	}
	err = emoneyRepository.Migrate(setupCtx)
	if err != nil {
		log.Fatalf("migrating emoney repository: %s", err.Error())
	}

	exitCode := m.Run()

	err = db.Close()
	if err != nil {
		log.Printf("Closing database: %s", err.Error())
	}

	os.Exit(exitCode)
}

type config struct {
	httpHostname     string
	httpPort         string
	databasePath     string
	webhookTargetURL string
	serverKey        string
}

func defaultConfig() config {
	return config{
		httpHostname: "localhost",
		httpPort:     "3000",
		databasePath: "payment.db",
	}
}

func parseConfig() config {
	result := defaultConfig()

	if v, ok := os.LookupEnv("HTTP_HOSTNAME"); ok {
		result.httpHostname = v
	}

	if v, ok := os.LookupEnv("HTTP_PORT"); ok {
		result.httpPort = v
	}

	if v, ok := os.LookupEnv("DATABASE_PATH"); ok {
		result.databasePath = v
	}

	if v, ok := os.LookupEnv("WEBHOOK_TARGET_URL"); ok {
		result.webhookTargetURL = v
	}

	if v, ok := os.LookupEnv("SERVER_KEY"); ok {
		result.serverKey = v
	}

	return result
}

func TestBusinessGetDetails(t *testing.T) {

	ctx := context.Background()

	transactionRepository, err := transaction.NewTransactionRepository(db)
	if err != nil {
		t.Errorf("creating transaction repository: %s", err.Error())
	}
	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		t.Errorf("creating virtual account repository: %s", err.Error())
	}

	emoneyRepository, err := emoney.NewEmoneyRepository(db)
	if err != nil {
		t.Errorf("creating emoney repository: %s", err.Error())
	}
	webhookClient, err := webhook.NewWebhookClient(cfg.webhookTargetURL)
	if err != nil {
		t.Errorf("creating webhook client: %s", err.Error())
	}

	paymentService, err := payment_service.NewPaymentService(payment_service.Config{
		ServerKey:                cfg.serverKey,
		TransactionRepository:    transactionRepository,
		EMoneyRepository:         emoneyRepository,
		WebhookClient:            webhookClient,
		VirtualAccountRepository: virtualAccountRepository,
	})
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	t.Run("GetDetails should return not found if the order id is empty", func(t *testing.T) {
		_, err = paymentService.GetDetail(ctx, "")
		if err == nil {
			t.Errorf("expecting error to be not nil, but got nil")
		}
	})
	t.Run("GetDetails should return err 'not found' if the order id is not found", func(t *testing.T) {
		_, err = paymentService.GetDetail(ctx, "not-exist")
		if err == nil {
			t.Errorf("expecting error to be not nil, but got nil")
		}
		if !errors.Is(err, business.ErrTransactionNotFound) {
			t.Errorf("expecting error %s, instead got %v", business.ErrTransactionNotFound, err)
		}
	})

	t.Run("GetDetails should return the correct details", func(t *testing.T) {
		vaNumber, err := virtualAccountRepository.CreateOrGetVirtualAccountNumber(ctx, "annedoe@example.com")
		if err != nil {
			_, err = db.Exec("DELETE FROM virtual_accounts")
			if err != nil {
				log.Printf("Deleting virtual accounts: %s", err.Error())
			}
			t.Errorf("unexpected error: %s", err.Error())
		}

		orderId := "order-id"
		_, err = virtualAccountRepository.CreateCharge(ctx, vaNumber, orderId, 50000, time.Now().Add(time.Hour))
		if err != nil {
			_, err = db.Exec("DELETE FROM virtual_accounts")
			if err != nil {
				log.Printf("Deleting virtual accounts: %s", err.Error())
			}
			_, err = db.Exec("DELETE FROM virtual_account_entries")
			if err != nil {
				log.Printf("Deleting virtual account entries: %s", err.Error())
			}
			t.Errorf("unexpected error: %s", err.Error())
		}

		_, err = emoneyRepository.CreateCharge(ctx, orderId, 50000, time.Now().Add(time.Hour))
		if err != nil {
			_, err = db.Exec("DELETE FROM virtual_accounts")
			if err != nil {
				log.Printf("Deleting virtual accounts: %s", err.Error())
			}
			_, err = db.Exec("DELETE FROM virtual_account_entries")
			if err != nil {
				log.Printf("Deleting virtual account entries: %s", err.Error())
			}
			_, err = db.Exec("DELETE FROM emoney_entries")
			if err != nil {
				log.Printf("Deleting emoney entries: %s", err.Error())
			}
			t.Errorf("unexpected error: %s", err.Error())
		}

		_, err = paymentService.GetDetail(ctx, vaNumber)
		if err == nil {
			_, err = db.Exec("DELETE FROM virtual_accounts")
			if err != nil {
				log.Printf("Deleting virtual accounts: %s", err.Error())
			}
			_, err = db.Exec("DELETE FROM virtual_account_entries")
			if err != nil {
				log.Printf("Deleting virtual account entries: %s", err.Error())
			}
			_, err = db.Exec("DELETE FROM emoney_entries")
			if err != nil {
				log.Printf("Deleting emoney entries: %s", err.Error())
			}
			t.Errorf("expecting error to be not nil, but got nil")
		}

		err = transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     orderId,
			Amount:      50000,
			PaymentType: primitive.PaymentTypeEMoneyQRIS,
			Status:      primitive.TransactionStatusPending,
			ExpiredAt:   time.Now().Add(time.Hour),
		})
		if err != nil {
			_, err = db.Exec("DELETE FROM transaction_log")
			if err != nil {
				log.Printf("Deleting transaction log: %s", err.Error())
			}
			_, err = db.Exec("DELETE FROM virtual_accounts")
			if err != nil {
				log.Printf("Deleting virtual accounts: %s", err.Error())
			}
			_, err = db.Exec("DELETE FROM virtual_account_entries")
			if err != nil {
				log.Printf("Deleting virtual account entries: %s", err.Error())
			}
			_, err = db.Exec("DELETE FROM emoney_entries")
			if err != nil {
				log.Printf("Deleting emoney entries: %s", err.Error())
			}
			t.Errorf("expecting not error when inserting transaction, instead got %s", err.Error())
		}

		_, err = paymentService.GetDetail(ctx, vaNumber)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		_, err = db.Exec("DELETE FROM transaction_log")
		if err != nil {
			log.Printf("Deleting transaction log: %s", err.Error())
		}
		_, err = db.Exec("DELETE FROM virtual_accounts")
		if err != nil {
			log.Printf("Deleting virtual accounts: %s", err.Error())
		}
		_, err = db.Exec("DELETE FROM virtual_account_entries")
		if err != nil {
			log.Printf("Deleting virtual account entries: %s", err.Error())
		}
		_, err = db.Exec("DELETE FROM emoney_entries")
		if err != nil {
			log.Printf("Deleting emoney entries: %s", err.Error())
		}
	})
}
