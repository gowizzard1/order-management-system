package payment

import "context"

type PaymentsService interface {
	ProcessPayment(ctx context.Context, payment *Payment) (*PaymentResponse, error)
	HandleMpesaCallback(ctx context.Context, callback *PaymentCallback) error
}
