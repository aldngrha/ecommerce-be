package grpcmiddleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"runtime/debug"
)

func ErrorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Printf("Recovered from panic: %v", recovered)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()
	res, err := handler(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, err
}
