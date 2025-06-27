package handler

import (
	"context"
	"github.com/aldngrha/ecommerce-be/internal/service"
	"github.com/aldngrha/ecommerce-be/internal/utils"
	"github.com/aldngrha/ecommerce-be/pb/auth"
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer
	authServive service.IAuthService
}

func (sh *authHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	validationErrors, err := utils.CheckValidations(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Process the registration logic here
	res, err := sh.authServive.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	validationErrors, err := utils.CheckValidations(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &auth.LoginResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Process the registration logic here
	res, err := sh.authServive.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authServive: authService,
	}
}
