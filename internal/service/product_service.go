package service

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/aldngrha/ecommerce-be/internal/entity"
	jwtentity "github.com/aldngrha/ecommerce-be/internal/entity/jwt"
	"github.com/aldngrha/ecommerce-be/internal/repository"
	"github.com/aldngrha/ecommerce-be/internal/utils"
	"github.com/aldngrha/ecommerce-be/pb/product"
	"github.com/google/uuid"
)

type IProductService interface {
	CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
}

func (ps *productService) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {

	//check if role user is admin?
	claims, err := jwtentity.GetClaimsFromContext(ctx)

	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return &product.CreateProductResponse{
			Base: utils.BadRequestResponse("only admin can create product"),
		}, nil
	}

	// check if image exists
	imagePath := filepath.Join("storage", "images", "products", req.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &product.CreateProductResponse{
				Base: utils.BadRequestResponse("image file not found"),
			}, nil
		}
	}

	// insert to db

	// success

	productEntity := entity.Product{
		Id:            uuid.NewString(),
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		ImageFileName: req.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.FullName,
	}

	err = ps.productRepository.CreateNewProduct(ctx, &productEntity)

	if err != nil {
		return nil, err
	}

	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("Product created successfully"),
		Id:   productEntity.Id,
	}, nil
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}
