package mock

import (
	"context"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/product"
)

type ProductRepository struct {
	CreateProductFunc func(ctx context.Context, p *product.Product) (*product.Product, error)
	GetProductFunc    func(ctx context.Context, id string) (*product.Product, error)
	ListProductsFunc  func(ctx context.Context) ([]*product.Product, error)
	UpdateProductFunc func(ctx context.Context, id string, update *product.ProductUpdate) (*product.Product, error)
	DeleteProductFunc func(ctx context.Context, id string) error
}

func (m *ProductRepository) CreateProduct(ctx context.Context, p *product.Product) (*product.Product, error) {
	return m.CreateProductFunc(ctx, p)
}

func (m *ProductRepository) GetProduct(ctx context.Context, id string) (*product.Product, error) {
	return m.GetProductFunc(ctx, id)
}

func (m *ProductRepository) ListProducts(ctx context.Context) ([]*product.Product, error) {
	return m.ListProductsFunc(ctx)
}

func (m *ProductRepository) UpdateProduct(ctx context.Context, id string, update *product.ProductUpdate,
) (*product.Product, error) {
	return m.UpdateProductFunc(ctx, id, update)
}

func (m *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	return m.DeleteProductFunc(ctx, id)
}
