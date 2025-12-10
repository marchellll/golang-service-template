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

// SendSMSNotification sends an SMS notification (mock implementation)
// This is a mock activity that simulates sending an SMS.
// To adapt for real SMS sending, replace the mock logic with your SMS service.
func SendSMSNotification(ctx context.Context, taskID string, notificationType string) error {
	logger := activity.GetLogger(ctx)

	// Simulate network delay
	delay := time.Duration(rand.Intn(500)+100) * time.Millisecond
	time.Sleep(delay)

	// Simulate occasional failures for retry demonstration (10% failure rate)
	if rand.Float32() < 0.1 {
		logger.Error(fmt.Sprintf("mock SMS send failed: task_id=%s type=%s", taskID, notificationType))
		return temporal.NewApplicationError("sms_send_failed", "Failed to send SMS", nil)
	}

	logger.Info(fmt.Sprintf("mock SMS notification sent successfully: task_id=%s type=%s", taskID, notificationType))

	log.Info().
		Str("task_id", taskID).
		Str("type", notificationType).
		Msg("[MOCK] SMS notification sent")

	return nil
}


