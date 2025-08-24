package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	grpcmiddleware2 "github.com/aldngrha/ecommerce-be/internal/grpcmiddleware"
	"github.com/aldngrha/ecommerce-be/internal/handler"
	"github.com/aldngrha/ecommerce-be/internal/repository"
	"github.com/aldngrha/ecommerce-be/internal/service"
	"github.com/aldngrha/ecommerce-be/pb/auth"
	"github.com/aldngrha/ecommerce-be/pb/product"
	"github.com/aldngrha/ecommerce-be/pkg/database"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	godotenv.Load()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}

	db := database.ConnectionDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to database successfully")

	cacheService := gocache.New(time.Hour*24, time.Hour)

	authMiddleware := grpcmiddleware2.NewAuthMiddleware(cacheService)

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, cacheService)
	authHandler := handler.NewAuthHandler(authService)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware2.ErrorMiddleware,
			authMiddleware.Middleware,
		),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)
	product.RegisterProductServiceServer(serv, productHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection registered")
	}

	log.Println("Server is starting on port :50052...")

	if err := serv.Serve(lis); err != nil {
		log.Panicf("Error serving: %v", err)
	}
}
