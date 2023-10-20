package service

import (
	"context"

	"github.com/jwambugu/mpesa-golang-sdk"
)

type Payment struct {
	Id          string `json:"id"`
	OrderId     string `json:"orderId"`
	PhoneNumber uint   `json:"phoneNumber"`
	Amount      uint   `json:"amount"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
	CallbackURL string `json:"callbackUrl"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type PaymentResponse struct {
	CheckoutRequestID string `json:"CheckoutRequestID"`
	CustomerMessage   string `json:"CustomerMessage"`
	MerchantRequestID string `json:"MerchantRequestID"`
	ResponseCode      string `json:"ResponseCode"`
	ResponseMessage   string `json:"ResponseMessage"`
}

// I know I know, this is a bit of a hack but it works for now
type PaymentCallback = mpesa.STKCallback

type PaymentsService interface {
	ProcessPayment(ctx context.Context, payment *Payment) (*PaymentResponse, error)
	HandleMpesaCallback(ctx context.Context, callback *PaymentCallback) error
}
