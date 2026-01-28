package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	httpserver "backend/internal/http/server"
	"backend/internal/repo/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

// Params collects runtime settings for the application.
type Params struct {
	Env      string
	HTTPAddr string
	DBDSN    string
}

// App wires configuration, database and HTTP server together.
type App struct {
	params Params
	pool   *pgxpool.Pool
	server *http.Server
}

// New constructs application dependencies.
func New(ctx context.Context, params Params) (*App, error) {
	pool, err := connectDB(ctx, params.DBDSN)
	if err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)
	engine := httpserver.NewRouter(params.Env, pool, queries)
	srv := &http.Server{
		Addr:              params.HTTPAddr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			fmt.Sprintf("%s/api/v1/auth/google/callback",
				os.Getenv("PUBLIC_URL"),
			),
			"email",
			"profile",
		),
	)

	return &App{params: params, pool: pool, server: srv}, nil
}

// Run starts the HTTP server and blocks until context cancellation.
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func connectDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 1
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute
	poolCfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

// Close releases resources.
func (a *App) Close() {
	if a.pool != nil {
		a.pool.Close()
	}
}
