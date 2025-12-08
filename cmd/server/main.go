package main

import (
	"context"
	"fmt"
	"golang-service-template/internal/app"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/samber/do"
)

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdout, stderr io.Writer,
) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	injector := app.NewInjector(
		ctx,
		getenv,
		stdout,
		stderr,
	)
	logger := do.MustInvoke[zerolog.Logger](injector)

	shutdownFn := app.RunNewServer(
		injector,
	)

	<-ctx.Done()
	// Create shutdown context with timeout, but derive from a fresh background context
	// since the original ctx is already cancelled (ctx.Done() was triggered)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := shutdownFn(shutdownCtx); err != nil {
		logger.Err(err).Msg("failed to shutdown server gracefully")
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	ctx := context.Background()

	err = run(
		ctx,
		os.Args,
		os.Getenv,
		os.Stdout,
		os.Stderr,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
