package product

import "context"

// RepositoryInterface is an interface for product repo
type RepositoryInterface interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context) ([]*Product, error)
	UpdateProduct(ctx context.Context, id string, update *ProductUpdate) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
}
