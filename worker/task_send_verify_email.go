package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (rtd *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, options ...asynq.Option) error {

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, options...)

	taskInfo, err := rtd.client.EnqueueContext(ctx, task)

	if err != nil {
		fmt.Errorf("failed to enqueue task %w", err)
	}

	fmt.Println("enqeued task ", taskInfo.ID)
	return nil
}

func (rtp *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := rtp.store.GetUser(ctx, payload.Username)

	if err != nil {
		return fmt.Errorf("failed to get user with username %s: %w", payload.Username, err)
	}

	//todo : send email to user

	fmt.Println("process task for user ", user.Username)
	return nil
}
