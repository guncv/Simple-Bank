package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/token"
	"github.com/guncv/Simple-Bank/util"
	"github.com/guncv/Simple-Bank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(taskDistributor, config, store)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {
	ctx := context.Background()
	newToken, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationTypeBearer, newToken)
	md := metadata.MD{
		authorizationHeaderKey: []string{bearerToken},
	}
	return metadata.NewIncomingContext(ctx, md)
}
