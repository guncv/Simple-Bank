package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/guncv/Simple-Bank/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	fmt.Println("metadata: ", metadata)

	values := metadata.Get(authorizationHeaderKey)
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header is not provided")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header format must be Bearer {token}")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header format must be Bearer {token}")
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return payload, nil
}
