package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type config struct {
	port int
	env  string
}

type application struct {
	cfg *config
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

	app := application{cfg: &cfg}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Starting %s server on %d", app.cfg.env, app.cfg.port)
	return srv.ListenAndServe()
}
