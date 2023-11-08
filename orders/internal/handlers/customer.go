package handlers

import (
	"context"
	"fmt"
	"github.com/leta/order-management-system/orders/generated"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/customers"
	"github.com/leta/order-management-system/orders/pkg/utils"
	"log"
)

func (s *GRPCServer) CreateCustomer(
	ctx context.Context, in *generated.CreateCustomerRequest) (*generated.CreateCustomerResponse, error) {

	log.Printf("Received: %v", in.GetFirstName())

	p, err := s.CustomerRepository.CreateCustomer(ctx, &customers.Customer{
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

func (s *GRPCServer) GetCustomer(
	ctx context.Context, in *generated.GetCustomerRequest) (*generated.GetCustomerResponse, error) {

	log.Printf("Received: %v", in.GetId())

	p, err := s.CustomerRepository.GetCustomer(ctx, in.GetId())
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

func (s *GRPCServer) ListCustomers(
	ctx context.Context, in *generated.ListCustomersRequest) (*generated.ListCustomersResponse, error) {

	log.Printf("Received: %v", in)

	customers, err := s.CustomerRepository.ListCustomers(ctx)
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

func (s *GRPCServer) UpdateCustomer(
	ctx context.Context, in *generated.UpdateCustomerRequest) (*generated.UpdateCustomerResponse, error) {

	log.Printf("Received: %v", in)

	p, err := s.CustomerRepository.UpdateCustomer(ctx, in.GetId(), &customers.CustomerUpdate{
		FirstName: utils.StringPtr(in.GetUpdate().GetFirstName()),
		LastName:  utils.StringPtr(in.GetUpdate().GetLastName()),
		Phone:     utils.StringPtr(in.GetUpdate().GetPhone()),
		Email:     utils.StringPtr(in.GetUpdate().GetEmail()),
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

func (s *GRPCServer) DeleteCustomer(
	ctx context.Context, in *generated.DeleteCustomerRequest) (*generated.DeleteCustomerResponse, error) {

	log.Printf("Received: %v", in.GetId())

	err := s.CustomerRepository.DeleteCustomer(ctx, in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete customers: %w", err)
	}

	return &generated.DeleteCustomerResponse{}, nil
}
