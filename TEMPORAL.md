# Temporal Workflow Integration

This template includes a Temporal workflow integration demonstrating asynchronous task notifications. When tasks are created or updated, a workflow is triggered that sends notifications via multiple channels (email, SMS, push) with automatic retry logic.

## Overview

The Temporal integration consists of:

- **Workflow**: `TaskNotificationWorkflow` - Orchestrates notification activities
- **Activities**: Mock implementations for email, SMS, and push notifications
- **Worker**: Separate process that executes workflows and activities
- **Integration**: Automatic workflow triggering from TaskService

## Setup

### 1. Start Temporal Server

```sh
# Start Temporal server and UI
docker compose --profile temporal up temporal temporal-ui

# Or start all services including Temporal
docker compose --profile temporal --profile dev up
```

The Temporal UI will be available at http://localhost:8088

### 2. Configure Environment Variables

Add these to your `.env` file:

```env
TEMPORAL_ADDRESS=localhost:7233
TEMPORAL_NAMESPACE=default
TEMPORAL_TASK_QUEUE=task-notifications
```

### 3. Start the Worker

The worker executes workflows and activities. Run it alongside your main server:

```sh
# In a separate terminal
go run ./cmd/temporal_worker/main.go
```

Or run it in Docker:

```sh
docker compose --profile temporal up temporal-worker
```

### 4. Start the Main Server

```sh
go run ./cmd/server/main.go
```

## How It Works

1. **Task Creation/Update**: When a task is created or updated via the API, the `TaskService` triggers a Temporal workflow asynchronously (fire-and-forget).

2. **Workflow Execution**: The `TaskNotificationWorkflow` orchestrates three activities in parallel:
   - Send email notification
   - Send SMS notification
   - Send push notification

3. **Retry Logic**: Each activity has automatic retry logic configured:
   - Initial retry after 1 second
   - Exponential backoff (coefficient: 2.0)
   - Maximum 3 attempts
   - Maximum interval: 1 minute

4. **Worker Processing**: The Temporal worker picks up workflow and activity executions and processes them.

## Mock Activities

All notification activities are currently mock implementations that:
- Simulate network delay (100-600ms)
- Randomly fail 10% of the time (to demonstrate retries)
- Log their execution

To replace with real implementations:
1. Edit the activity files in `internal/temporal/activity/`
2. Replace mock logic with your actual notification service calls
3. Update activity registration in `cmd/temporal_worker/main.go` if needed

## Adapting for Other Use Cases

This workflow pattern can be easily adapted for other use cases:

### Example: Order Processing Workflow

1. **Copy the workflow structure**:
   ```sh
   cp internal/temporal/workflow/task_notification.go internal/temporal/workflow/order_processing.go
   ```

2. **Modify the workflow**:
   - Change input/output types
   - Add/remove activities as needed
   - Adjust retry policies

3. **Create new activities**:
   ```sh
   cp internal/temporal/activity/email.go internal/temporal/activity/payment.go
   ```

4. **Register in worker**:
   - Add workflow registration
   - Add activity registrations

5. **Trigger from your service**:
   - Add workflow execution call similar to TaskService

### Example: Bulk Data Processing

For processing large batches:
- Use `workflow.ExecuteActivity` with arrays
- Implement pagination within the workflow
- Use signals for progress updates

## Monitoring

### Temporal UI

Access the Temporal UI at http://localhost:8088 to:
- View workflow executions
- See activity results
- Monitor retries and failures
- Debug workflow issues

### Logs

Workflows and activities log to stdout. Check logs for:
- Workflow start/completion
- Activity execution
- Retry attempts
- Errors

## Removing Temporal

To completely remove Temporal from the project:

1. **Delete directories**:
   ```sh
   rm -rf internal/temporal/
   rm internal/app/temporal.go
   rm cmd/temporal_worker/
   ```

2. **Remove from `internal/common/configs.go`**:
   - Delete `TemporalConfig` struct
   - Remove `TemporalConfig` field from `Config` struct

3. **Remove from `internal/app/configs.go`**:
   - Delete Temporal config parsing in `NewConfig()`

4. **Remove from `internal/app/di.go`**:
   - Delete `do.Provide(injector, NewTemporalClient)` line

5. **Remove from `internal/service/task.go`**:
   - Remove Temporal client import
   - Remove `temporalClient` field from `taskService` struct
   - Remove Temporal client initialization in `NewTaskService()`
   - Remove workflow trigger calls in `Create()` and `Update()` methods (2-3 lines each)

6. **Remove from `compose.yml`**:
   - Delete `temporal` and `temporal-ui` service definitions

7. **Remove dependency**:
   ```sh
   go mod tidy
   ```

8. **Delete this file**:
   ```sh
   rm TEMPORAL.md
   ```

## Troubleshooting

### Worker not starting

- Check that Temporal server is running: `docker compose ps`
- Verify `TEMPORAL_ADDRESS` in `.env` matches the server address
- Check worker logs for connection errors

### Workflows not executing

- Ensure worker is running
- Verify task queue name matches in both service and worker
- Check Temporal UI for workflow status

### Activities failing

- Check activity logs in worker output
- Verify activity registration in worker
- Review retry policy configuration

## References

- [Temporal Documentation](https://docs.temporal.io/)
- [Temporal Go SDK](https://github.com/temporalio/sdk-go)
- [Temporal UI](https://docs.temporal.io/dev-guide/ui)


