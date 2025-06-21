package handler

import (
	"context"
	"fmt"
	"github.com/aldngrha/ecommerce-be/internal/utils"
	"github.com/aldngrha/ecommerce-be/pb/service"
)

type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (h *serviceHandler) SayHello(ctx context.Context, req *service.HelloRequest) (*service.HelloResponse, error) {
	validationErrors, err := utils.CheckValidations(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &service.HelloResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	return &service.HelloResponse{
		Message: fmt.Sprintf("Hello %s", req.Name),
		Base:    utils.SuccessResponse("Successfully processed request"),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}
