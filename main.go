package main

import (
	"context"
	"fmt"
	"golang-service-template/internal"
	"io"
	"os"
)

func run(
	ctx    context.Context,
	args   []string,
	stdin  io.Reader,
	stdout, stderr io.Writer,
) error {
	container := internal.NewContainer()
	srv := internal.NewServer(container)

	err := srv.Start(":1323")

	return err
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}