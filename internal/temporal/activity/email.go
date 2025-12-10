package activity

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

// SendEmailNotification sends an email notification (mock implementation)
// This is a mock activity that simulates sending an email.
// To adapt for real email sending, replace the mock logic with your email service.
func SendEmailNotification(ctx context.Context, taskID string, notificationType string) error {
	logger := activity.GetLogger(ctx)

	// Simulate network delay
	delay := time.Duration(rand.Intn(500)+100) * time.Millisecond
	time.Sleep(delay)

	// Simulate occasional failures for retry demonstration (10% failure rate)
	if rand.Float32() < 0.1 {
		logger.Error(fmt.Sprintf("mock email send failed: task_id=%s type=%s", taskID, notificationType))
		return temporal.NewApplicationError("email_send_failed", "Failed to send email", nil)
	}

	logger.Info(fmt.Sprintf("mock email notification sent successfully: task_id=%s type=%s", taskID, notificationType))

	log.Info().
		Str("task_id", taskID).
		Str("type", notificationType).
		Msg("[MOCK] Email notification sent")

	return nil
}


