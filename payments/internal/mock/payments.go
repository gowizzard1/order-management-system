package mock

import (
	"context"

	"github.com/leta/order-management-system/payments/internal/service"
)

var _ service.PaymentsService = (*PaymentsService)(nil)

type PaymentsService struct {
	ProcessPaymentFunc      func(ctx context.Context, p *service.Payment) (*service.PaymentResponse, error)
	HandleMpesaCallbackFunc func(ctx context.Context, p *service.PaymentCallback) error
}

func (m *PaymentsService) ProcessPayment(ctx context.Context, p *service.Payment) (*service.PaymentResponse, error) {
	return m.ProcessPaymentFunc(ctx, p)
}

func (m *PaymentsService) HandleMpesaCallback(ctx context.Context, p *service.PaymentCallback) error {
	return m.HandleMpesaCallbackFunc(ctx, p)
}
