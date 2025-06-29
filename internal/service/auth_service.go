package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/aldngrha/ecommerce-be/internal/entity"
	"github.com/aldngrha/ecommerce-be/internal/repository"
	"github.com/aldngrha/ecommerce-be/internal/utils"
	"github.com/aldngrha/ecommerce-be/pb/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"strings"
	"time"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
	cacheService   *gocache.Cache
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
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
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

func (s *authService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// get token from metadata grpc
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata found in context")
	}

	bearerToken, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no authorization token found in metadata")
	}

	if len(bearerToken) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is empty")
	}

	tokenSplit := strings.Split(bearerToken[0], " ")

	if len(tokenSplit) != 2 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token format")
	}

	if tokenSplit[0] != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token must start with Bearer")

	}

	jwtToken := tokenSplit[1]

	// return token
	tokenClaims, err := jwt.ParseWithClaims(jwtToken, &entity.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	var claims *entity.JwtClaims
	if claims, ok = tokenClaims.Claims.(*entity.JwtClaims); !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token claims")
	}

	// insert token to memory cache or database to invalidate the token
	s.cacheService.Set(jwtToken, "", time.Duration(claims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)

	// send response

	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout successful"),
	}, nil
}

func NewAuthService(authRepository repository.IAuthRepository, cacheService *gocache.Cache) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:   cacheService,
	}
}
