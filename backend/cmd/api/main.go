package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"backend/internal/app"
)

func main() {
	env := getEnv("APP_ENV", "development")
	addr := getEnv("APP_HTTP_ADDR", ":8080")
	dsn := strings.TrimSpace(os.Getenv("APP_DB_DSN"))
	if dsn == "" {
		log.Fatal("APP_DB_DSN is required")
	}

	params := app.Params{
		Env:      strings.ToLower(env),
		HTTPAddr: addr,
		DBDSN:    dsn,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, params)
	if err != nil {
		log.Fatalf("init error: %v", err)
	}
	defer application.Close()

	if err := application.Run(ctx); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
