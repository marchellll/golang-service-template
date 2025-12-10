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

// SendPushNotification sends a push notification (mock implementation)
// This is a mock activity that simulates sending a push notification.
// To adapt for real push notifications, replace the mock logic with your push service.
func SendPushNotification(ctx context.Context, taskID string, notificationType string) error {
	logger := activity.GetLogger(ctx)

	// Simulate network delay
	delay := time.Duration(rand.Intn(500)+100) * time.Millisecond
	time.Sleep(delay)

	// Simulate occasional failures for retry demonstration (10% failure rate)
	if rand.Float32() < 0.1 {
		logger.Error(fmt.Sprintf("mock push send failed: task_id=%s type=%s", taskID, notificationType))
		return temporal.NewApplicationError("push_send_failed", "Failed to send push notification", nil)
	}

	logger.Info(fmt.Sprintf("mock push notification sent successfully: task_id=%s type=%s", taskID, notificationType))

	log.Info().
		Str("task_id", taskID).
		Str("type", notificationType).
		Msg("[MOCK] Push notification sent")

	return nil
}


