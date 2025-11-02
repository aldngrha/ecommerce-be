package service

import (
	"context"
	"fmt"
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
	DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error)
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

func (ps *productService) DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	// query to db with data id
	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// if null, return not found
	if productEntity == nil {
		return &product.DetailProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	// send response
	return &product.DetailProductResponse{
		Base:         utils.SuccessResponse("Get product detail successfully"),
		Id:           productEntity.Id,
		Name:         productEntity.Name,
		Description:  productEntity.Description,
		Price:        productEntity.Price,
		ImageFileUrl: fmt.Sprintf("%s/images/products/%s", os.Getenv("STORAGE_SERVICE_URL"), productEntity.ImageFileName),
	}, nil
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}
