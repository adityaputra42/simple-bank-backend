package worker

import (
	"context"
	"encoding/json"
	"fmt"
	db "simple-bank/db/sqlc"
	"simple-bank/util"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

// DistributeTaskSendVerifyEmail implements TaskDistributor.
func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload : %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task : %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

// ProcessTaskSendVerifyEmail implements TaskProsessor.
func (processor *RedisTaskProsesor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload : %w", asynq.SkipRetry)
	}
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf("user not found : %w", asynq.SkipRetry)
		// }
		return fmt.Errorf("failed to get user : %w", err)
	}
	arg := db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create verify enail : %w", err)
	}

	subject := "welcome to simple bank"
	verifyUrl := fmt.Sprintf("http://simple-bank.org?id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us! <br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>`, user.FullName, verifyUrl)

	err = processor.mailer.SendEmail(subject, content, []string{verifyEmail.Email}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify enail : %w", err)
	}

	// Send email to user
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}
