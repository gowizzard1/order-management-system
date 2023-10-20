package mpesa

import (
	"context"
	"fmt"
	"github.com/jwambugu/mpesa-golang-sdk"
	"github.com/leta/order-management-system/payments/internal/repository"
	"github.com/leta/order-management-system/payments/internal/service"
	"github.com/leta/order-management-system/payments/pkg/utils"

	orders "github.com/leta/order-management-system/orders/pkg/client"
)

const (
	MPESA_BUSINESS_SHORT_CODE = "MPESA_BUSINESS_SHORT_CODE" // #nosec G101 - This is an env variable name
	MPESA_PASSKEY             = "MPESA_PASSKEY"             // #nosec G101 - This is an env variable name
)

var _ service.PaymentsService = (*PaymentsService)(nil)

type PaymentsService struct {
	mpesa        *Mpesa
	db           repository.PaymentsRepository
	ordersClient orders.OrdersClient
}

func NewPaymentsService(
	mpesa *Mpesa, orderClient orders.OrdersClient, db repository.PaymentsRepository) *PaymentsService {
	return &PaymentsService{
		mpesa:        mpesa,
		ordersClient: orderClient,
	}
}

func (s *PaymentsService) CheckPreconditions() {
	if s.mpesa == nil {
		panic("no Mpesa service provided")
	}
}

func (s *PaymentsService) ProcessPayment(ctx context.Context, payment *service.Payment) (*service.PaymentResponse, error) {
	s.CheckPreconditions()

	// stored in an environemnt variable for now:- assumption is that the system handles orders for a single business
	businessShortCode, err := utils.StringToUint(utils.MustGetEnv(MPESA_BUSINESS_SHORT_CODE))
	if err != nil {
		return nil, fmt.Errorf("failed to convert business short code to uint: %v", err)
	}

	passKey := utils.MustGetEnv(MPESA_PASSKEY)

	_, err = s.ordersClient.UpdateOrderStatus(ctx, &orders.UpdateOrderStatusRequest{
		Id:     "1",
		Status: orders.OrderStatusPending,
	})
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to update orders status: %v", err)
	}

	stkPushRes, err := s.mpesa.app.STKPush(ctx, passKey, mpesa.STKPushRequest{
		BusinessShortCode: businessShortCode,
		TransactionType:   "CustomerBuyGoodsOnline",
		Amount:            payment.Amount,
		PartyA:            payment.PhoneNumber,
		PartyB:            businessShortCode,
		PhoneNumber:       uint64(payment.PhoneNumber),
		CallBackURL:       payment.CallbackURL,
		AccountReference:  payment.Reference,
		TransactionDesc:   payment.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %v", MpesaErrorToInternalError(err))
	}

	_, err = s.db.CreatePayment(ctx, &repository.Payment{
		Amount:            payment.Amount,
		Phone:             fmt.Sprint(payment.PhoneNumber),
		Reference:         payment.Reference,
		Description:       payment.Description,
		MerchantRequestID: stkPushRes.MerchantRequestID,
		Status:            repository.PaymentStatusPending,
	})
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to store payment record: %v", err)
	}

	return &service.PaymentResponse{
		MerchantRequestID: stkPushRes.MerchantRequestID,
		CheckoutRequestID: stkPushRes.CheckoutRequestID,
		ResponseCode:      stkPushRes.ResponseCode,
		CustomerMessage:   stkPushRes.CustomerMessage,
	}, nil

}

func (s *PaymentsService) HandleMpesaCallback(ctx context.Context, callback *service.PaymentCallback) error {
	s.CheckPreconditions()

	payment, err := s.db.GetPaymentByMerchantRequestID(ctx, callback.MerchantRequestID)
	if err != nil {
		return service.Errorf(service.INTERNAL_ERROR, "failed to get payment: %v", err)
	}

	if callback.ResultCode != 0 {
		_, err := s.ordersClient.UpdateOrderStatus(ctx, &orders.UpdateOrderStatusRequest{
			Id:     payment.OrderID,
			Status: orders.OrderStatusFailed,
		})
		if err != nil {
			return service.Errorf(service.INTERNAL_ERROR, "failed to update orders status(Failed): %v", err)
		}

		err = s.db.UpdatePaymentStatus(ctx, payment.Id, repository.PaymentStatusFailed)
		if err != nil {
			return service.Errorf(service.INTERNAL_ERROR, "failed to update payment status(Failed): %v", err)
		}

		return service.Errorf(service.NOT_IMPLEMENTED_ERROR, "invalid request: %s", callback.ResultDesc)
	}

	_, err = s.ordersClient.UpdateOrderStatus(ctx, &orders.UpdateOrderStatusRequest{
		Id:     payment.OrderID,
		Status: orders.OrderStatusPaid,
	})
	if err != nil {
		return service.Errorf(service.INTERNAL_ERROR, "failed to update orders status(Paid): %v", err)
	}

	return nil
}
