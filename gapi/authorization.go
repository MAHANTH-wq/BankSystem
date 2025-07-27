package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/mahanth/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	mtdt, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &token.Payload{}, fmt.Errorf("metadata is missing")
	}

	values := mtdt.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	//<authorization-type> <authorization token>

	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])

	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)

	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	return payload, nil

}
