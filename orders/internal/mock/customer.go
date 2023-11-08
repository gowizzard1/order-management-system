package mock

import (
	"context"

	"github.com/leta/order-management-system/orders/internal/api/interfaces/customers"
)

var _ customers.CustomerRepository = (*CustomerRepository)(nil)

type CustomerRepository struct {
	CreateCustomerFunc func(ctx context.Context, p *customers.Customer) (*customers.Customer, error)
	GetCustomerFunc    func(ctx context.Context, id string) (*customers.Customer, error)
	ListCustomersFunc  func(ctx context.Context) ([]*customers.Customer, error)
	UpdateCustomerFunc func(ctx context.Context, id string, update *customers.CustomerUpdate) (*customers.Customer, error)
	DeleteCustomerFunc func(ctx context.Context, id string) error
}

func (m *CustomerRepository) CreateCustomer(ctx context.Context, p *customers.Customer) (*customers.Customer, error) {
	return m.CreateCustomerFunc(ctx, p)
}

func (m *CustomerRepository) GetCustomer(ctx context.Context, id string) (*customers.Customer, error) {
	return m.GetCustomerFunc(ctx, id)
}

func (m *CustomerRepository) ListCustomers(ctx context.Context) ([]*customers.Customer, error) {
	return m.ListCustomersFunc(ctx)
}

func (m *CustomerRepository) UpdateCustomer(ctx context.Context, id string, update *customers.CustomerUpdate,
) (*customers.Customer, error) {
	return m.UpdateCustomerFunc(ctx, id, update)
}

func (m *CustomerRepository) DeleteCustomer(ctx context.Context, id string) error {
	return m.DeleteCustomerFunc(ctx, id)
}
