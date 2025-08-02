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

func (server *Server) authorizeUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
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

	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("user does not have permission to access this resource")
	}
	return payload, nil

}

func hasPermission(role string, accessibleRoles []string) bool {
	for _, accessibleRole := range accessibleRoles {
		if role == accessibleRole {
			return true
		}
	}
	return false
}
