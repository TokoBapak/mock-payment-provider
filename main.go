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

	"mock-payment-provider/business/transaction_service"
	"mock-payment-provider/presentation"
	"mock-payment-provider/repository/transaction"
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

	webhookClient, err := webhook.NewWebhookClient(cfg.webhookTargetURL)
	if err != nil {
		log.Fatalf("creating webhook client: %s", err.Error())
	}

	transactionService, err := transaction_service.NewTransactionService(transaction_service.Dependency{
		TransactionRepository: transactionRepository,
		WebhookClient:         webhookClient,
	})
	if err != nil {
		log.Fatalf("creating transaction service: %s", err.Error())
	}

	httpServer, err := presentation.NewPresenter(presentation.PresenterConfig{
		Hostname:   cfg.httpHostname,
		Port:       cfg.httpPort,
		Dependency: &presentation.Dependency{TransactionService: transactionService},
	})
	if err != nil {
		log.Fatalf("creating new presenter: %s", err.Error())
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		log.Printf("HTTP server listening on %s", httpServer.Addr)

		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("serving HTTP server: %s", err.Error())
		}
	}()

	<-exitSignal

	log.Printf("Interrupt signal received, exiting...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
	defer shutdownCancel()

	err = httpServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("shutting down HTTP server: %s", err.Error())
	}
}
