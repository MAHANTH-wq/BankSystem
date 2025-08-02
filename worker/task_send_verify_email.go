package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/util"
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

	verifyEmail, err := rtp.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email for user %s: %w", user.Username, err)
	}

	subject := "Welcome to Mahanth's Bank System"
	verifyEmailLink := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf("Hello %s,\n\n"+
		"Please verify your email by clicking the link below:\n"+
		"%s\n\n"+
		"Thank you for joining us!\n\n", user.FullName, verifyEmailLink)

	to := []string{user.Email}

	err = rtp.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email to %s: %w", user.Email, err)
	}

	fmt.Println("process task for user ", user.Username)
	return nil
}
