package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"mock-payment-provider/business/payment_service"
	"mock-payment-provider/business/transaction_service"
	"mock-payment-provider/presentation"
	"mock-payment-provider/repository/emoney"
	"mock-payment-provider/repository/transaction"
	"mock-payment-provider/repository/virtual_account"
	"mock-payment-provider/repository/webhook"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := parseConfig()
	database, err := sql.Open("sqlite3", cfg.databasePath)
	if err != nil {
		log.Fatalf("opening sql connection: %s", err.Error())
	}
	defer func() {
		err := database.Close()
		if err != nil {
			log.Printf("closing database connection: %s", err.Error())
		}
	}()

	transactionRepository, err := transaction.NewTransactionRepository(database)
	if err != nil {
		log.Fatalf("creating transaction repository: %s", err.Error())
	}

	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(database)
	if err != nil {
		log.Fatalf("creating virtual account repository: %s", err.Error())
	}

	emoneyRepository, err := emoney.NewEmoneyRepository(database)
	if err != nil {
		log.Fatalf("creating emoney repository: %s", err.Error())
	}

	webhookClient, err := webhook.NewWebhookClient(cfg.webhookTargetURL)
	if err != nil {
		log.Fatalf("creating webhook client: %s", err.Error())
	}

	transactionService, err := transaction_service.NewTransactionService(transaction_service.Config{
		ServerKey:                cfg.serverKey,
		TransactionRepository:    transactionRepository,
		WebhookClient:            webhookClient,
		VirtualAccountRepository: virtualAccountRepository,
		EMoneyRepository:         emoneyRepository,
	})
	if err != nil {
		log.Fatalf("creating transaction service: %s", err.Error())
	}

	paymentService, err := payment_service.NewPaymentService(payment_service.Config{
		ServerKey:                cfg.serverKey,
		TransactionRepository:    transactionRepository,
		WebhookClient:            webhookClient,
		EMoneyRepository:         emoneyRepository,
		VirtualAccountRepository: virtualAccountRepository,
	})
	if err != nil {
		log.Fatalf("creating payment service: %s", err.Error())
	}

	httpServer, err := presentation.NewPresenter(presentation.PresenterConfig{
		Hostname: cfg.httpHostname,
		Port:     cfg.httpPort,
		Dependency: &presentation.Dependency{
			TransactionService: transactionService,
			PaymentService:     paymentService,
		},
	})
	if err != nil {
		log.Fatalf("creating new presenter: %s", err.Error())
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		<-exitSignal

		log.Printf("Interrupt signal received, exiting...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
		defer shutdownCancel()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("shutting down HTTP server: %s", err.Error())
		}
	}()

	log.Printf("HTTP server listening on %s", httpServer.Addr)

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serving HTTP server: %s", err.Error())
	}
}
