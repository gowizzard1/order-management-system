package payments

import (
	"context"
	"github.com/leta/order-management-system/payments/generated"
	"github.com/leta/order-management-system/payments/internal/interfaces/api/payment"
)

type PaymentService struct {
	paymentSvc  payment.PaymentsService
	callBackURL string
}

func (s *PaymentService) ProcessPayment(ctx context.Context, in *generated.MpesaPaymentRequest) (*generated.MpesaPaymentResponse, error) {

	//callBackURL := "https://order-management-system.herokuapp.com/payments/callback"

	p, err := s.paymentSvc.ProcessPayment(ctx, &payment.Payment{
		Amount:      uint(in.GetAmount()),
		PhoneNumber: uint(in.GetPhoneNumber()),
		CallbackURL: s.callBackURL,
		Reference:   in.GetReference(),
		Description: in.GetDescription(),
	})
	if err != nil {
		return nil, err
	}

	return &generated.MpesaPaymentResponse{
		CheckoutRequestId: p.CheckoutRequestID,
		CustomerMessage:   p.CustomerMessage,
		ResponseCode:      p.ResponseCode,
		MerchantRequestId: p.MerchantRequestID,
	}, nil
}
