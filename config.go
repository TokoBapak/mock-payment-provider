package main

import "os"

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
