package main

import (
	"context"
	"fmt"
	"golang-service-template/internal"
	"os"
)

func run(
	ctx    context.Context,
) error {
	container := internal.NewContainer()
	srv := internal.NewServer(container)

	err := srv.Start(":"+container.Config.Port)

	return err
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}