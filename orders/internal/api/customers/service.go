package customers

import (
	"context"
	"fmt"
	"github.com/leta/order-management-system/orders/generated"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/customers"
	"log"
)

type CustomerService struct {
	customersRepo customers.CustomerRepositoryInterface
}

func NewCustomerService(customersRepo customers.CustomerRepository) customers.CustomerServiceInterface {
	return &CustomerService{
		customersRepo: customersRepo,
	}
}

func (s *CustomerService) CreateCustomer(
	ctx context.Context, in *generated.CreateCustomerRequest) (*generated.CreateCustomerResponse, error) {

	log.Printf("Received: %v", in.GetFirstName())

	p, err := s.customersRepo.CreateCustomer(ctx, &customers.Customer{
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Email:     in.GetEmail(),
		Phone:     in.GetPhone(),
	})
	if err != nil {
		return nil, err
	}
	return &generated.CreateCustomerResponse{
		Id: p.Id,
	}, nil
}

func (s *CustomerService) GetCustomer(
	ctx context.Context, in *generated.GetCustomerRequest) (*generated.GetCustomerResponse, error) {

	log.Printf("Received: %v", in.GetId())

	p, err := s.customersRepo.GetCustomer(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return &generated.GetCustomerResponse{
		Id:        p.Id,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Phone:     p.Phone,
		Email:     p.Email,
	}, nil
}

func (s *CustomerService) ListCustomers(
	ctx context.Context, in *generated.ListCustomersRequest) (*generated.ListCustomersResponse, error) {

	log.Printf("Received: %v", in)

	customers, err := s.customersRepo.ListCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	var responseCustomers []*generated.Customer
	for _, p := range customers {
		responseCustomers = append(responseCustomers, &generated.Customer{
			Id:        p.Id,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Phone:     p.Phone,
			Email:     p.Email,
		})
	}

	return &generated.ListCustomersResponse{
		Customers: responseCustomers,
	}, nil
}

func (s *CustomerService) UpdateCustomer(
	ctx context.Context, in *generated.UpdateCustomerRequest) (*generated.UpdateCustomerResponse, error) {

	log.Printf("Received: %v", in)

	p, err := s.customersRepo.UpdateCustomer(ctx, in.GetId(), &customers.CustomerUpdate{
		FirstName: utils.StringPtr(in.GetUpdate().GetFirstName()),
		LastName:  pkg.StringPtr(in.GetUpdate().GetLastName()),
		Phone:     pkg.StringPtr(in.GetUpdate().GetPhone()),
		Email:     pkg.StringPtr(in.GetUpdate().GetEmail()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update customers: %w", err)
	}

	return &generated.UpdateCustomerResponse{
		Id:        p.Id,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
	}, nil
}

func (s *CustomerService) DeleteCustomer(
	ctx context.Context, in *generated.DeleteCustomerRequest) (*generated.DeleteCustomerResponse, error) {

	log.Printf("Received: %v", in.GetId())

	err := s.customersRepo.DeleteCustomer(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete customers: %w", err)
	}

	return &generated.DeleteCustomerResponse{}, nil
}
