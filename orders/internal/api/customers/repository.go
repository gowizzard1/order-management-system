package customers

import (
	"github.com/leta/order-management-system/orders/internal/interfaces/api/customers"
	"github.com/leta/order-management-system/orders/pkg/models"
	"github.com/leta/order-management-system/orders/pkg/utils"

	"cloud.google.com/go/firestore"
	"context"
	"github.com/leta/order-management-system/orders/db/firebase"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

//var _ repository.CustomerRepository = (*CustomerRepository)(nil)

type customerRepository struct {
	db *firebase.FirestoreService
}

func NewCustomerRepository(db *firebase.FirestoreService) customers.CustomerRepositoryInterface {
	return &customerRepository{
		db: db,
	}
}

func (s *customerRepository) CheckPreconditions() {
	if s.db == nil {
		panic("no DB service provided")
	}
}

func (s *customerRepository) customerCollection() *firestore.CollectionRef {
	s.CheckPreconditions()

	return s.db.Client.Collection("customers")
}

func (s *customerRepository) CreateCustomer(ctx context.Context, customer *customers.Customer) (*customers.Customer, error) {
	s.CheckPreconditions()

	currentTime := time.Now()
	customer.CreatedAt = currentTime.Format(time.RFC3339)
	customer.UpdatedAt = currentTime.Format(time.RFC3339)

	//err := customer.Validate()
	//if err != nil {
	//	return nil, utils.Errorf(utils.INVALID_ERROR, "invalid customers provided: %v", err)
	//}
	customerModel := s.marshallCustomer(customer)

	docRef, _, err := s.customerCollection().Add(ctx, customerModel)
	if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to create customers: %v", err)
	}

	customer.Id = docRef.ID

	return customer, nil
}

func (s *customerRepository) GetCustomer(ctx context.Context, id string) (*customers.Customer, error) {

	s.CheckPreconditions()

	if id == "" {
		return nil, utils.Errorf(utils.INVALID_ERROR, "id is required")
	}

	docRef, err := s.customerCollection().Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, utils.Errorf(utils.NOT_FOUND_ERROR, "customers not found")
	} else if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to get customers: %v", err)
	}

	customerModel := &models.CustomerModel{}
	if err := docRef.DataTo(customerModel); err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to unmarshall customers: %v", err)
	}

	customer := s.unmarshallCustomer(customerModel)

	customer.Id = docRef.Ref.ID

	return customer, nil
}

func (s *customerRepository) ListCustomers(ctx context.Context) ([]*customers.Customer, error) {
	s.CheckPreconditions()

	iter := s.customerCollection().Documents(ctx)

	var customers []*customers.Customer

	for {
		docRef, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to iterate customers: %v", err)
		}

		customerModel := &models.CustomerModel{}
		if err := docRef.DataTo(customerModel); err != nil {
			return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to unmarshall customers: %v", err)
		}

		customer := s.unmarshallCustomer(customerModel)
		customer.Id = docRef.Ref.ID

		customers = append(customers, customer)
	}

	return customers, nil
}

func (s *customerRepository) UpdateCustomer(ctx context.Context, id string, update *customers.CustomerUpdate) (*customers.Customer, error) {

	s.CheckPreconditions()

	customer, err := s.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	if c := update.FirstName; c != nil {
		customer.FirstName = *c
	}

	if c := update.LastName; c != nil {
		customer.LastName = *c
	}

	if c := update.Email; c != nil {
		customer.Email = *c
	}

	if c := update.Phone; c != nil {
		customer.Phone = *c
	}

	timeNow := time.Now()
	customer.UpdatedAt = timeNow.Format(time.RFC3339)

	//err = customer.Validate()
	//if err != nil {
	//	return nil, utils.Errorf(utils.INVALID_ERROR, "invalid customers details provided: %v", err)
	//}

	customerModel := s.marshallCustomer(customer)

	_, err = s.customerCollection().Doc(id).Set(ctx, customerModel)
	if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to update customers: %v", err)
	}

	return customer, nil

}

func (s *customerRepository) DeleteCustomer(ctx context.Context, id string) error {
	s.CheckPreconditions()

	if id == "" {
		return utils.Errorf(utils.INVALID_ERROR, "id is required")
	}

	_, err := s.customerCollection().Doc(id).Delete(ctx)

	return err
}

func (s *customerRepository) marshallCustomer(customer *customers.Customer) *models.CustomerModel {

	return &models.CustomerModel{
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

func (s *customerRepository) unmarshallCustomer(customerModel *models.CustomerModel) *customers.Customer {

	return &customers.Customer{
		FirstName: customerModel.FirstName,
		LastName:  customerModel.LastName,
		Email:     customerModel.Email,
		Phone:     customerModel.Phone,
		CreatedAt: customerModel.CreatedAt,
		UpdatedAt: customerModel.UpdatedAt,
	}
}
