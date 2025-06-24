package service

import (
	"context"
	"github.com/aldngrha/ecommerce-be/internal/entity"
	"github.com/aldngrha/ecommerce-be/internal/repository"
	"github.com/aldngrha/ecommerce-be/internal/utils"
	"github.com/aldngrha/ecommerce-be/pb/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
}

func (s *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.Password != req.ConfirmPassword {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password and confirm password do not match"),
		}, nil
	}

	// Check email from db
	user, err := s.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	// if email already exists, return error
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("User with this email already exists"),
		}, nil
	}

	// if email does not exist, proceed with registration logic insert to db
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	// insert user into db
	newUser := entity.User{
		Id:        uuid.NewString(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		RoleCode:  entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: &req.FullName,
	}
	err = s.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("User registered successfully"),
	}, nil
}

func NewAuthService(authRepository repository.IAuthRepository) IAuthService {
	return &authService{
		authRepository: authRepository,
	}
}
