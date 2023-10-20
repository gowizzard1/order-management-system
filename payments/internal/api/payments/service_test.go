package payments_test

import (
	"context"
	"github.com/leta/order-management-system/payments/generated"
	"github.com/leta/order-management-system/payments/internal/handlers/grpc"
	"github.com/leta/order-management-system/payments/internal/service"
	"reflect"
	"testing"
)

const (
	ERROR_PAYMENT_TRIGGER = "Error Payment"
)

func mockProcessPaymentFunc(ctx context.Context, p *service.Payment) (*service.PaymentResponse, error) {

	if p.Reference == ERROR_PAYMENT_TRIGGER {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid request: %s", p.Reference)
	}

	return &service.PaymentResponse{
		CheckoutRequestID: "checkoutRequestID",
		CustomerMessage:   "customerMessage",
		ResponseCode:      "responseCode",
		MerchantRequestID: "merchantRequestID",
	}, nil
}

func TestGRPCServer_ProcessPayment(t *testing.T) {

	s := grpc.NewTestGRPCServer(t)

	s.PaymentsService.ProcessPaymentFunc = mockProcessPaymentFunc

	type args struct {
		ctx context.Context
		in  *generated.MpesaPaymentRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *generated.MpesaPaymentResponse
		wantErr bool
	}{
		{
			name: "Process Payment Success",
			args: args{
				ctx: context.Background(),
				in: &generated.MpesaPaymentRequest{
					Amount:      100,
					PhoneNumber: 254700000000,
					Reference:   "reference",
					Description: "description",
				},
			},
			want: &generated.MpesaPaymentResponse{
				CheckoutRequestId: "checkoutRequestID",
				CustomerMessage:   "customerMessage",
				ResponseCode:      "responseCode",
				MerchantRequestId: "merchantRequestID",
			},
			wantErr: false,
		},
		{
			name: "Process Payment Error",
			args: args{
				ctx: context.Background(),
				in: &generated.MpesaPaymentRequest{
					Amount:      100,
					PhoneNumber: 254700000000,
					Reference:   ERROR_PAYMENT_TRIGGER,
					Description: "description",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.ProcessPayment(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.ProcessPayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.ProcessPayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
