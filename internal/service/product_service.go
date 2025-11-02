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
	EditProduct(ctx context.Context, request *product.EditProductRequest) (*product.EditProductResponse, error)
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

func (ps *productService) EditProduct(ctx context.Context, request *product.EditProductRequest) (*product.EditProductResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)

	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return &product.EditProductResponse{
			Base: utils.BadRequestResponse("only admin can update product"),
		}, nil
	}

	// validate if id available on db
	productEntity, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if productEntity == nil {
		return &product.EditProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	// if image change, remove old image
	if productEntity.ImageFileName != request.ImageFileName {
		newImagePath := filepath.Join("storage", "images", "products", request.ImageFileName)
		_, err := os.Stat(newImagePath)
		if err != nil {
			if os.IsNotExist(err) {
				return &product.EditProductResponse{
					Base: utils.NotFoundResponse("Image not found"),
				}, nil
			}
			return nil, err
		}

		oldImagePath := filepath.Join("storage", "images", "products", productEntity.ImageFileName)
		err = os.Remove(oldImagePath)
		if err != nil {
			return nil, err
		}
	}

	// update db
	newProduct := entity.Product{
		Id:            request.Id,
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		ImageFileName: request.ImageFileName,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &claims.FullName,
	}

	err = ps.productRepository.UpdateProduct(ctx, &newProduct)

	if err != nil {
		return nil, err
	}

	// send response
	return &product.EditProductResponse{
		Base: utils.SuccessResponse("Edit product successfully"),
		Id:   request.Id,
	}, nil
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}
