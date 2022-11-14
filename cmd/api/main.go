package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/jsonlog"
)

const (
	version = "0.0.1"
)

type config struct {
	port int
	env  string
}

type application struct {
	cfg    *config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {

	// TODO: Move this a separate file that passes config from file
	apiPort, err := strconv.Atoi(os.Getenv("apiPort"))
	if err != nil {
		return err
	}

	var cfg config

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	cfg.port = *flags.Int("port", apiPort, "Port to listen on")
	cfg.env = *flags.String("env", "development", "Environment ([development]|production)")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	app := application{
		cfg:    &cfg,
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		wg:     sync.WaitGroup{},
	}

	return app.serve()
}
