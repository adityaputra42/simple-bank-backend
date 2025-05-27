package gapi

import (
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"simple-bank/worker"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymetricKey:    util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServerGrpc(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}
