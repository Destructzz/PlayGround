// @title           PlayGround API
// @version         1.0
// @description     Backend API for PlayGround
// @host            localhost
// @BasePath        /
package main

import (
	_ "backend/docs"
	"backend/internal/app"
	"backend/internal/observability"
	"backend/pkg"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	env := pkg.GetEnv("APP_ENV", "development")
	addr := pkg.GetEnv("APP_HTTP_ADDR", ":8080")
	dsn := pkg.GetEnv("APP_DB_DSN", "")
	params := app.Params{
		Env:      strings.ToLower(env),
		HTTPAddr: addr,
		DBDSN:    dsn,
	}
	logger, err := observability.NewLogger(params.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
	defer func() {
		_ = logger.Sync()
	}()

	if params.DBDSN == "" {
		zap.L().Fatal("APP_DB_DSN is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, params)
	if err != nil {
		zap.L().Fatal("init error", zap.Error(err))
	}
	defer application.Close()

	if err := application.Run(ctx); err != nil {
		zap.L().Fatal("server error", zap.Error(err))
	}
}
