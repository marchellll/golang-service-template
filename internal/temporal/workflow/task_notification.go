package workflow

import (
	"fmt"
	"time"

	"golang-service-template/internal/temporal/activity"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// TaskNotificationInput represents the input for the task notification workflow
type TaskNotificationInput struct {
	TaskID           string
	NotificationType string // "create" or "update"
}

// TaskNotificationResult represents the result of the notification workflow
type TaskNotificationResult struct {
	EmailSent bool
	SMSSent   bool
	PushSent  bool
	Errors    []string
}

// TaskNotificationWorkflow orchestrates sending notifications via multiple channels
// This workflow demonstrates:
// - Parallel activity execution
// - Error handling and retries (configured at activity level)
// - Workflow result aggregation
//
// To adapt for other use cases:
// 1. Modify TaskNotificationInput to include your workflow parameters
// 2. Add/remove activities as needed
// 3. Adjust the result structure
// 4. Update activity registration in the worker
func TaskNotificationWorkflow(ctx workflow.Context, input TaskNotificationInput) (TaskNotificationResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info(fmt.Sprintf("Starting task notification workflow: task_id=%s type=%s", input.TaskID, input.NotificationType))

	// Configure activity options with retry policy
	// Activities will retry automatically on failure
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute notification activities in parallel
	// This demonstrates Temporal's ability to handle concurrent operations
	emailFuture := workflow.ExecuteActivity(ctx, activity.SendEmailNotification, input.TaskID, input.NotificationType)
	smsFuture := workflow.ExecuteActivity(ctx, activity.SendSMSNotification, input.TaskID, input.NotificationType)
	pushFuture := workflow.ExecuteActivity(ctx, activity.SendPushNotification, input.TaskID, input.NotificationType)

	result := TaskNotificationResult{
		Errors: []string{},
	}

	// Wait for email notification
	if err := emailFuture.Get(ctx, nil); err != nil {
		logger.Error(fmt.Sprintf("Email notification failed: %v", err))
		result.Errors = append(result.Errors, "email: "+err.Error())
	} else {
		result.EmailSent = true
	}

	// Wait for SMS notification
	if err := smsFuture.Get(ctx, nil); err != nil {
		logger.Error(fmt.Sprintf("SMS notification failed: %v", err))
		result.Errors = append(result.Errors, "sms: "+err.Error())
	} else {
		result.SMSSent = true
	}

	// Wait for push notification
	if err := pushFuture.Get(ctx, nil); err != nil {
		logger.Error(fmt.Sprintf("Push notification failed: %v", err))
		result.Errors = append(result.Errors, "push: "+err.Error())
	} else {
		result.PushSent = true
	}

	logger.Info(fmt.Sprintf("Task notification workflow completed: task_id=%s email_sent=%v sms_sent=%v push_sent=%v error_count=%d",
		input.TaskID, result.EmailSent, result.SMSSent, result.PushSent, len(result.Errors)))

	return result, nil
}


