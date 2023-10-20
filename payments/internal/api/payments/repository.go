package payments

import (
	"cloud.google.com/go/firestore"
	"context"
	db "github.com/leta/order-management-system/payments/db/firebase"
	"github.com/leta/order-management-system/payments/internal/repository"
	"github.com/leta/order-management-system/payments/internal/service"
	"github.com/leta/order-management-system/payments/pkg/models"
	"time"
)

var _ repository.PaymentsRepository = (*PaymentsRepository)(nil)

type PaymentsRepository struct {
	db *db.FirestoreService
}

func NewPaymentsRepository(db *db.FirestoreService) *PaymentsRepository {
	return &PaymentsRepository{
		db: db,
	}
}

func (r *PaymentsRepository) CheckPreconditions() {
	if r.db == nil {
		panic("no DB service provided")
	}
}

func (r *PaymentsRepository) paymentsCollection() *firestore.CollectionRef {
	r.CheckPreconditions()

	return r.db.Client.Collection("payments")
}

func (r *PaymentsRepository) CreatePayment(ctx context.Context, payment *repository.Payment) (string, error) {
	r.CheckPreconditions()

	currentTime := time.Now()
	payment.CreatedAt = currentTime.Format(time.RFC3339)
	payment.UpdatedAt = currentTime.Format(time.RFC3339)

	err := payment.Validate()
	if err != nil {
		return "", service.Errorf(service.INVALID_ERROR, "invalid payment details provided: %v", err)
	}
	paymentModel := r.marshallPayment(payment)

	docRef, _, err := r.paymentsCollection().Add(ctx, paymentModel)
	if err != nil {
		return "", service.Errorf(service.INTERNAL_ERROR, "failed to create payment: %v", err)
	}

	payment.Id = docRef.ID

	return payment.Id, nil
}

func (r *PaymentsRepository) GetPaymentByID(ctx context.Context, paymentID string) (*repository.Payment, error) {
	r.CheckPreconditions()

	return r.getPaymentByID(ctx, paymentID)
}

func (r *PaymentsRepository) getPaymentByID(ctx context.Context, paymentID string) (*repository.Payment, error) {
	r.CheckPreconditions()

	if paymentID == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid payment ID provided")
	}

	docRef := r.paymentsCollection().Doc(paymentID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get payment: %v", err)
	}

	var paymentModel models.PaymentModel
	err = doc.DataTo(&paymentModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to decode payment: %v", err)
	}

	payment := r.unmarshallPayment(&paymentModel)

	return payment, nil
}

func (r *PaymentsRepository) UpdatePaymentStatus(
	ctx context.Context, paymentID string, status repository.PaymentStatus) error {
	r.CheckPreconditions()

	if paymentID == "" {
		return service.Errorf(service.INVALID_ERROR, "invalid payment ID provided")
	}

	payment, err := r.getPaymentByID(ctx, paymentID)
	if err != nil {
		return err
	}

	payment.Status = status
	payment.UpdatedAt = time.Now().Format(time.RFC3339)

	opaymentModel := r.marshallPayment(payment)
	_, err = r.paymentsCollection().Doc(paymentID).Set(ctx, opaymentModel)
	if err != nil {
		return service.Errorf(service.INTERNAL_ERROR, "failed to update orders status: %v", err)
	}
	return nil
}

func (r *PaymentsRepository) GetPaymentByMerchantRequestID(
	ctx context.Context, merchantRequestID string) (*repository.Payment, error) {
	r.CheckPreconditions()

	if merchantRequestID == "" {
		return nil, service.Errorf(service.INVALID_ERROR, "invalid merchant request ID provided")
	}

	query := r.paymentsCollection().Where("reference", "==", merchantRequestID).Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to get payment: %v", err)
	}

	if len(docs) == 0 {
		return nil, service.Errorf(service.NOT_FOUND_ERROR, "payment not found")
	}

	var paymentModel models.PaymentModel
	err = docs[0].DataTo(&paymentModel)
	if err != nil {
		return nil, service.Errorf(service.INTERNAL_ERROR, "failed to decode payment: %v", err)
	}

	payment := r.unmarshallPayment(&paymentModel)

	return payment, nil
}

func (r *PaymentsRepository) marshallPayment(payment *repository.Payment) *models.PaymentModel {
	return &models.PaymentModel{
		Amount:      payment.Amount,
		Status:      string(payment.Status),
		OrderID:     payment.OrderID,
		Phone:       payment.Phone,
		Reference:   payment.Reference,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}
}

func (r *PaymentsRepository) unmarshallPayment(paymentModel *models.PaymentModel) *repository.Payment {
	return &repository.Payment{
		Id:          paymentModel.ID,
		Amount:      paymentModel.Amount,
		Status:      repository.PaymentStatus(paymentModel.Status),
		OrderID:     paymentModel.OrderID,
		Phone:       paymentModel.Phone,
		Reference:   paymentModel.Reference,
		Description: paymentModel.Description,
		CreatedAt:   paymentModel.CreatedAt,
		UpdatedAt:   paymentModel.UpdatedAt,
	}
}
