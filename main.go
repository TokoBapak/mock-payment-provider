package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
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

	log := zerolog.New(os.Stdout)

	database, err := sql.Open("sqlite3", cfg.databasePath)
	if err != nil {
		log.Fatal().Msgf("opening sql connection: %s", err.Error())
	}
	defer func() {
		err := database.Close()
		if err != nil {
			log.Err(err).Msg("closing database connection")
		}
	}()

	transactionRepository, err := transaction.NewTransactionRepository(database)
	if err != nil {
		log.Fatal().Msgf("creating transaction repository: %s", err.Error())
	}

	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(database)
	if err != nil {
		log.Fatal().Msgf("creating virtual account repository: %s", err.Error())
	}

	emoneyRepository, err := emoney.NewEmoneyRepository(database)
	if err != nil {
		log.Fatal().Msgf("creating emoney repository: %s", err.Error())
	}

	webhookClient, err := webhook.NewWebhookClient(cfg.webhookTargetURL)
	if err != nil {
		log.Fatal().Msgf("creating webhook client: %s", err.Error())
	}

	transactionService, err := transaction_service.NewTransactionService(transaction_service.Config{
		ServerKey:                cfg.serverKey,
		TransactionRepository:    transactionRepository,
		WebhookClient:            webhookClient,
		VirtualAccountRepository: virtualAccountRepository,
		EMoneyRepository:         emoneyRepository,
	})
	if err != nil {
		log.Fatal().Msgf("creating transaction service: %s", err.Error())
	}

	paymentService, err := payment_service.NewPaymentService(payment_service.Config{
		ServerKey:                cfg.serverKey,
		TransactionRepository:    transactionRepository,
		WebhookClient:            webhookClient,
		EMoneyRepository:         emoneyRepository,
		VirtualAccountRepository: virtualAccountRepository,
	})
	if err != nil {
		log.Fatal().Msgf("creating payment service: %s", err.Error())
	}

	httpServer, err := presentation.NewPresenter(presentation.PresenterConfig{
		Hostname:  cfg.httpHostname,
		Port:      cfg.httpPort,
		ServerKey: cfg.serverKey,
		Dependency: &presentation.Dependency{
			TransactionService: transactionService,
			PaymentService:     paymentService,
		},
	})
	if err != nil {
		log.Fatal().Msgf("creating new presenter: %s", err.Error())
	}

	// Migrate repositories on startup
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err = transactionRepository.Migrate(ctx)
	if err != nil {
		log.Fatal().Msgf("migrating transaction repository: %s", err.Error())
	}

	err = virtualAccountRepository.Migrate(ctx)
	if err != nil {
		log.Fatal().Msgf("migrating virtual account repository: %s", err.Error())
	}

	err = emoneyRepository.Migrate(ctx)
	if err != nil {
		log.Fatal().Msgf("migrating emoney repository: %s", err.Error())
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		<-exitSignal

		log.Info().Msg("Interrupt signal received, exiting...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
		defer shutdownCancel()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Err(err).Msg("shutting down HTTP server")
		}
	}()

	log.Printf("HTTP server listening on %s", httpServer.Addr)

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msgf("serving HTTP server: %s", err.Error())
	}
}
