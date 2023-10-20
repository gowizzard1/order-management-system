package customers

import (
	"context"
	"github.com/leta/order-management-system/orders/generated"
)

type CustomerServiceInterface interface {
	CreateCustomer(ctx context.Context, in *generated.CreateCustomerRequest) (*generated.CreateCustomerResponse, error)
	GetCustomer(ctx context.Context, in *generated.GetCustomerRequest) (*generated.GetCustomerResponse, error)
	ListCustomers(ctx context.Context, in *generated.ListCustomersRequest) (*generated.ListCustomersResponse, error)
	UpdateCustomer(ctx context.Context, in *generated.UpdateCustomerRequest) (*generated.UpdateCustomerResponse, error)
	DeleteCustomer(ctx context.Context, in *generated.DeleteCustomerRequest) (*generated.DeleteCustomerResponse, error)
}
