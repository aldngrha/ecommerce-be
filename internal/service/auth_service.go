package service

import (
	"context"
	"errors"
	"github.com/aldngrha/ecommerce-be/internal/entity"
	"github.com/aldngrha/ecommerce-be/internal/repository"
	"github.com/aldngrha/ecommerce-be/internal/utils"
	"github.com/aldngrha/ecommerce-be/pb/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
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

func (s *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// Check email from db
	user, err := s.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("User with this email does not exist"),
		}, nil
	}
	// check if password is correct

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid password")
		}
		return nil, err // return error if there is an issue with comparing passwords
	}

	// generate JWT token
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			IssuedAt:  jwt.NewNumericDate(now.Add(time.Hour * 24)),
			ExpiresAt: jwt.NewNumericDate(now),
		},
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.RoleCode,
	})

	secretKey := os.Getenv("JWT_SECRET")

	accessToken, err := token.SignedString([]byte(secretKey)) //
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login successful"),
		AccessToken: accessToken,
	}, nil

}

func NewAuthService(authRepository repository.IAuthRepository) IAuthService {
	return &authService{
		authRepository: authRepository,
	}
}
