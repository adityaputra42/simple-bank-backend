package worker

import (
	"context"
	db "simple-bank/db/sqlc"

	"github.com/hibiken/asynq"
)

type TaskProsessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProsesor struct {
	server *asynq.Server
	store  db.Store
}

// Start implements TaskProsessor.
func (processor *RedisTaskProsesor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, db db.Store) TaskProsessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{},
	)

	return &RedisTaskProsesor{server: server, store: db}
}
