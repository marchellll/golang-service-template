package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"

	"golang-service-template/internal/app"
	"golang-service-template/internal/common"
	"golang-service-template/internal/temporal/activity"
	"golang-service-template/internal/temporal/workflow"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/samber/do"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func run(
	ctx context.Context,
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
	config := do.MustInvoke[common.Config](injector)

	// Get Temporal client
	temporalClient, err := do.Invoke[client.Client](injector)
	if err != nil || temporalClient == nil {
		logger.Warn().Msg("Temporal client not available, worker will not start")
		return nil
	}

	// Create worker
	taskQueue := config.TemporalConfig.TaskQueue
	if taskQueue == "" {
		taskQueue = "task-notifications"
	}

	w := worker.New(temporalClient, taskQueue, worker.Options{})

	// Register workflow
	w.RegisterWorkflow(workflow.TaskNotificationWorkflow)

	// Register activities
	w.RegisterActivity(activity.SendEmailNotification)
	w.RegisterActivity(activity.SendSMSNotification)
	w.RegisterActivity(activity.SendPushNotification)

	// Start worker
	logger.Info().
		Str("task_queue", taskQueue).
		Str("namespace", config.TemporalConfig.Namespace).
		Msg("Starting Temporal worker")

	if err := w.Start(); err != nil {
		logger.Fatal().Err(err).Msg("Unable to start worker")
		return err
	}

	// Wait for interrupt signal
	<-ctx.Done()
	logger.Info().Msg("Shutting down Temporal worker")

	w.Stop()
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
		os.Getenv,
		os.Stdout,
		os.Stderr,
	)

	if err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}
