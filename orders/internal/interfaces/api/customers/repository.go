package customers

import "context"

type CustomerRepositoryInterface interface {
	CreateCustomer(ctx context.Context, customer *Customer) (*Customer, error)
	GetCustomer(ctx context.Context, id string) (*Customer, error)
	ListCustomers(ctx context.Context) ([]*Customer, error)
	UpdateCustomer(ctx context.Context, id string, update *CustomerUpdate) (*Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}
