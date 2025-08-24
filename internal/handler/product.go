package handler

import (
	"context"

	"github.com/aldngrha/ecommerce-be/internal/service"
	"github.com/aldngrha/ecommerce-be/internal/utils"
	"github.com/aldngrha/ecommerce-be/pb/product"
)

type productHandler struct {
	product.UnimplementedProductServiceServer
	productServive service.IProductService
}

func (ph *productHandler) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	validationErrors, err := utils.CheckValidations(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.CreateProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Process the registration logic here
	res, err := ph.productServive.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewProductHandler(productService service.IProductService) *productHandler {
	return &productHandler{
		productServive: productService,
	}
}
