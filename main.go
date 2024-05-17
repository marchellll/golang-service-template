package main

import (
	"context"
	"fmt"
	"golang-service-template/internal"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func run(
	ctx    context.Context,
	args   []string,
	stdin  io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	container := internal.NewContainer()
	srv := internal.NewServer(container)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(container.Config.Host, container.Config.Port),
		Handler: srv,
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}