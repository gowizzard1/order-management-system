package product

import (
	"context"
	"github.com/leta/order-management-system/orders/generated"
)

// ProductServiceInterface defines the methods that should be implemented
// by a service that handles product-related operations.
type ProductServiceInterface interface {
	CreateProduct(ctx context.Context, in *generated.CreateProductRequest) (*generated.CreateProductResponse, error)
	GetProduct(ctx context.Context, in *generated.GetProductRequest) (*generated.GetProductResponse, error)
	ListProducts(ctx context.Context, in *generated.ListProductsRequest) (*generated.ListProductsResponse, error)
	UpdateProduct(ctx context.Context, in *generated.UpdateProductRequest) (*generated.UpdateProductResponse, error)
	DeleteProduct(ctx context.Context, in *generated.DeleteProductRequest) (*generated.DeleteProductResponse, error)
}
