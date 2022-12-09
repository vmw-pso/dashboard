package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/data"
	"github.com/vmw-pso/delivery-dashboard/back-end/internal/jsonlog"

	_ "github.com/lib/pq"
)

const (
	version = "0.0.1"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	cfg    *config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	if err := run(os.Args, logger); err != nil {
		logger.PrintFatal(err, nil)
		os.Exit(1)
	}
}

func run(args []string, logger *jsonlog.Logger) error {

	var cfg config

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	cfg.port = *flags.Int("port", 8086, "Port to listen on")
	cfg.env = *flags.String("env", "development", "Environment ([development]|production)")
	cfg.db.dsn = *flags.String("dsn", "postgres://postgres:password@localhost/dashboard?sslmode=disable", "PostgreSQL DSN")
	cfg.db.maxOpenConns = *flags.Int("db-max-open-conns", 25, "Database maximum open connections")
	cfg.db.maxIdleConns = *flags.Int("db-max-idle-conns", 25, "Database maximum idle connections")
	cfg.db.maxIdleTime = *flags.String("db-max-idle-time", "15m", "Database maximum idle time")

	flags.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	db, err := openDB(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		return err
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := application{
		cfg:    &cfg,
		logger: logger,
		models: *data.NewModels(db),
		wg:     sync.WaitGroup{},
	}

	return app.serve()
}

func openDB(dsn string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
