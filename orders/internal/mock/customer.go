package mock

import (
	"context"

	"github.com/leta/order-management-system/orders/internal/repository"
)

var _ repository.CustomerRepository = (*CustomerRepository)(nil)

type CustomerRepository struct {
	CreateCustomerFunc func(ctx context.Context, p *repository.Customer) (*repository.Customer, error)
	GetCustomerFunc    func(ctx context.Context, id string) (*repository.Customer, error)
	ListCustomersFunc  func(ctx context.Context) ([]*repository.Customer, error)
	UpdateCustomerFunc func(ctx context.Context, id string, update *repository.CustomerUpdate) (*repository.Customer, error)
	DeleteCustomerFunc func(ctx context.Context, id string) error
}

func (m *CustomerRepository) CreateCustomer(ctx context.Context, p *repository.Customer) (*repository.Customer, error) {
	return m.CreateCustomerFunc(ctx, p)
}

func (m *CustomerRepository) GetCustomer(ctx context.Context, id string) (*repository.Customer, error) {
	return m.GetCustomerFunc(ctx, id)
}

func (m *CustomerRepository) ListCustomers(ctx context.Context) ([]*repository.Customer, error) {
	return m.ListCustomersFunc(ctx)
}

func (m *CustomerRepository) UpdateCustomer(ctx context.Context, id string, update *repository.CustomerUpdate,
) (*repository.Customer, error) {
	return m.UpdateCustomerFunc(ctx, id, update)
}

func (m *CustomerRepository) DeleteCustomer(ctx context.Context, id string) error {
	return m.DeleteCustomerFunc(ctx, id)
}
