package main

import (
	"context"
	"github.com/aldngrha/ecommerce-be/internal/handler"
	"github.com/aldngrha/ecommerce-be/package/database"
	"github.com/aldngrha/ecommerce-be/pb/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	ctx := context.Background()
	godotenv.Load()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}

	database.ConnectionDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to database successfully")

	serviceHandler := handler.NewServiceHandler()

	serv := grpc.NewServer()

	service.RegisterHelloWorldServiceServer(serv, serviceHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection registered for HelloWorldService")
	}

	log.Println("Server is starting on port :50051...")

	if err := serv.Serve(lis); err != nil {
		log.Panicf("Error serving: %v", err)
	}
}
