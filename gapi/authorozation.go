package gapi

import (
	"context"
	"fmt"
	"simple-bank/token"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizationUser(ctx context.Context, accessableRole []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	authHeader := values[0]

	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type %s", authType)

	}
	accesToken := fields[len(fields)-1]

	payload, err := server.tokenMaker.VerifyToken(accesToken, token.TokenTypeAccessToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorization")
	}

	if !hashPermision(payload.Role, accessableRole) {
		return nil, fmt.Errorf("permission denied")
	}
	return payload, nil
}

func hashPermision(userRole string, accesableRole []string) bool {
	for _, role := range accesableRole {
		if userRole == role {
			return true
		}

	}
	return false
}
