package handlers

import (
	"context"
	"github.com/leta/order-management-system/orders/generated"
)

func (s *GRPCServer) ProcessCheckout(
	ctx context.Context, in *generated.ProcessCheckoutRequest) (*generated.ProcessCheckoutResponse, error) {

	o, err := s.CheckoutService.ProcessCheckout(ctx, in.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &generated.ProcessCheckoutResponse{
		OrderId:    o.Id,
		CustomerId: o.CustomerId,
		Status:     getGRPCOrderStatus(o.OrderStatus),
	}, nil
}
