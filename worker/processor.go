package worker

import (
	"context"
	db "simple-bank/db/sqlc"
	"simple-bank/mail"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProsessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProsesor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmailSender
}

// Start implements TaskProsessor.
func (processor *RedisTaskProsesor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, db db.Store, mailer mail.EmailSender) TaskProsessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)

	return &RedisTaskProsesor{server: server, store: db, mailer: mailer}
}
