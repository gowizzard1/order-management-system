package handlers

import (
	"context"
	"fmt"
)

func (s *GRPCServer) ProcessCheckout(
	ctx context.Context, in *orders_service.ProcessCheckoutRequest) (*orders_service.ProcessCheckoutResponse, error) {

	o, err := s.CheckoutService.ProcessCheckout(ctx, in.GetOrderId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to process checkout: %w", err))
	}

	return &orders_service.ProcessCheckoutResponse{
		OrderId:    o.Id,
		CustomerId: o.CustomerId,
		Status:     getGRPCOrderStatus(o.OrderStatus),
	}, nil
}
