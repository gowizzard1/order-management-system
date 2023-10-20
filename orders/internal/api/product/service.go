package product

import (
	"context"
	"fmt"
	"github.com/leta/order-management-system/orders/generated"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/product"
	"github.com/leta/order-management-system/orders/pkg/utils"
	"log"
)

type productService struct {
	productRepo product.ProductRepositoryInterface
}

func NewProductService(productRepo product.ProductRepositoryInterface) product.ProductServiceInterface {
	return &productService{
		productRepo: productRepo,
	}
}
func (s *productService) CreateProduct(ctx context.Context, in *generated.CreateProductRequest) (*generated.CreateProductResponse, error) {

	log.Printf("Received: %v", in.GetName())

	p, err := s.productRepo.CreateProduct(ctx, &product.Product{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Price:       uint(in.GetPrice()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &generated.CreateProductResponse{
		Id: p.Id,
	}, nil
}

func (s *productService) GetProduct(ctx context.Context, in *generated.GetProductRequest) (*generated.GetProductResponse, error) {

	log.Printf("Received: %v", in.GetId())

	p, err := s.productRepo.GetProduct(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &generated.GetProductResponse{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       uint32(p.Price),
	}, nil
}

func (s *productService) ListProducts(ctx context.Context, in *generated.ListProductsRequest) (*generated.ListProductsResponse, error) {

	log.Printf("Received: %v", in)

	products, err := s.productRepo.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	var responseProducts []*generated.Product
	for _, p := range products {
		responseProducts = append(responseProducts, &generated.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       uint32(p.Price),
		})
	}

	return &generated.ListProductsResponse{
		Products: responseProducts,
	}, nil
}

func (s *productService) UpdateProduct(ctx context.Context, in *generated.UpdateProductRequest) (*generated.UpdateProductResponse, error) {

	log.Printf("Received: %v", in)

	p, err := s.productRepo.UpdateProduct(ctx, in.GetId(), &product.ProductUpdate{
		Name:        utils.StringPtr(in.GetUpdate().GetName()),
		Description: utils.StringPtr(in.GetUpdate().GetDescription()),
		Price:       utils.UintPtr(uint(in.GetUpdate().GetPrice())),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &generated.UpdateProductResponse{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       uint32(p.Price),
	}, nil
}

func (s *productService) DeleteProduct(ctx context.Context, in *generated.DeleteProductRequest) (*generated.DeleteProductResponse, error) {

	log.Printf("Received: %v", in)

	err := s.productRepo.DeleteProduct(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}

	return &generated.DeleteProductResponse{}, nil
}
