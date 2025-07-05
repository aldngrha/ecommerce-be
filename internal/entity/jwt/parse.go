package jwt

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func ParseTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "no metadata found in context")
	}

	bearerToken, ok := md["authorization"]
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "no authorization token found in metadata")
	}

	if len(bearerToken) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is empty")
	}

	tokenSplit := strings.Split(bearerToken[0], " ")

	if len(tokenSplit) != 2 {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization token format")
	}

	if tokenSplit[0] != "Bearer" {
		return "", status.Errorf(codes.Unauthenticated, "authorization token must start with Bearer")

	}
	return tokenSplit[1], nil
}
