package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

type JwtEntityContextKey string

var JwtEntityContextKeyValue JwtEntityContextKey = "JwtEntity"

type JwtClaims struct {
	jwt.RegisteredClaims
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func (jc *JwtClaims) SendToContext(ctx context.Context) context.Context {

	return context.WithValue(ctx, JwtEntityContextKeyValue, jc)
}

func GetClaimsFromToken(jwtToken string) (*JwtClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(jwtToken, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		// return secret key
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	if !tokenClaims.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "token is not valid")
	}

	if claims, ok := tokenClaims.Claims.(*JwtClaims); ok {
		return claims, nil
	}

	return nil, status.Errorf(codes.Unauthenticated, "token is not valid")
}

func GetClaimsFromContext(ctx context.Context) (*JwtClaims, error) {
	claims, ok := ctx.Value(JwtEntityContextKeyValue).(*JwtClaims)

	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: no JWT claims found in context")
	}

	return claims, nil
}
