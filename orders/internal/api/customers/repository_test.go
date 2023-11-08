package customers_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	c "github.com/leta/order-management-system/orders/internal/api/customers"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/customers"

	db "github.com/leta/order-management-system/orders/db/firebase"
)

func deleteTestCustomer(t *testing.T, ctx context.Context, cs customers.CustomerRepositoryInterface, id string) {
	err := cs.DeleteCustomer(ctx, id)
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}
}

func TestCustomerService_CheckPreconditions(t *testing.T) {
	type fields struct {
		repo customers.CustomerRepositoryInterface
	}
	tests := []struct {
		name      string
		fields    fields
		wantPanic bool
	}{
		{
			name: "Check Preconditions Failed - nil DB",
			fields: fields{
				repo: nil,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := c.NewCustomerService(tt.fields.repo)
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("CustomerRepository.CheckPreconditions() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			s.CheckPreconditions()
		})
	}
}

func TestCustomerService_CreateCustomer(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	customerRepository := db.NewCustomerService(firestoreService)

	type args struct {
		ctx      context.Context
		customer *customers.Customer
	}
	tests := []struct {
		name    string
		args    args
		want    *customers.Customer
		wantErr bool
	}{
		{
			name: "Create Customer Success",
			args: args{
				ctx: context.Background(),
				customer: &customers.Customer{
					FirstName: "Test",
					LastName:  "Customer",
					Email:     "test@test.com",
					Phone:     "254722000000",
				},
			},
			want: &customers.Customer{
				Id:        "",
				FirstName: "Test",
				LastName:  "Customer",
				Email:     "test@test.com",
				Phone:     "254722000000",
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			wantErr: false,
		},
		{
			name: "Create Customer Failure - Invalid Customer",
			args: args{
				ctx: context.Background(),
				customer: &customers.Customer{
					FirstName: "",
					LastName:  "Customer",
					Email:     "test@test.com",
					Phone:     "1234567890",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customerRepository.CreateCustomer(tt.args.ctx, tt.args.customer)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomerRepository.CreateCustomers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				// Clean up the created customers
				defer deleteTestCustomer(t, ctx, customerRepository, got.Id)
				// Ignore the ID in the comparison since it's unpredictable
				got.Id = ""
			}

			if tt.want != nil {
				tt.want.Id = ""
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomerRepository.CreateCustomers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomerService_GetCustomer(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	customerRepository := db.NewCustomerService(firestoreService)

	c, err := customerRepository.CreateCustomer(ctx, &customers.Customer{
		FirstName: "Test",
		LastName:  "Customer",
		Email:     "test@email.com",
		Phone:     "254722000000",
	})
	if err != nil {
		t.Fatalf("failed to create customers: %v", err)
	}
	defer deleteTestCustomer(t, ctx, customerRepository, c.Id)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *customers.Customer
		wantErr bool
	}{
		{
			name: "Get Customer Success",
			args: args{
				ctx: context.Background(),
				id:  c.Id,
			},
			want:    c,
			wantErr: false,
		},
		{
			name: "Get Customer Failure - Invalid ID",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := customerRepository.GetCustomer(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomerRepository.GetCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomerRepository.GetCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomerService_ListCustomers(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	customerRepository := db.NewCustomerService(firestoreService)

	c, err := customerRepository.CreateCustomer(ctx, &customers.Customer{
		FirstName: "Test",
		LastName:  "Customer",
		Email:     "test@test.com",
		Phone:     "254722000000",
	})
	if err != nil {
		t.Fatalf("failed to create customers: %v", err)
	}

	defer deleteTestCustomer(t, ctx, customerRepository, c.Id)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*customers.Customer
		wantErr bool
	}{
		{
			name: "List Products Success",
			args: args{
				ctx: context.Background(),
			},
			want: []*customers.Customer{
				c,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customerRepository.ListCustomers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomerRepository.ListCustomers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomerRepository.ListCustomers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomerService_UpdateCustomer(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	customerRepository := db.NewCustomerService(firestoreService)

	c, err := customerRepository.CreateCustomer(ctx, &customers.Customer{
		FirstName: "Test",
		LastName:  "Customer",
		Email:     "test@test.com",
		Phone:     "254722000000",
	})
	if err != nil {
		t.Fatalf("failed to create customers: %v", err)
	}

	defer deleteTestCustomer(t, ctx, customerRepository, c.Id)

	type args struct {
		ctx    context.Context
		id     string
		update *customers.CustomerUpdate
	}
	tests := []struct {
		name    string
		args    args
		want    *customers.Customer
		wantErr bool
	}{
		{
			name: "Update Product Success",
			args: args{
				ctx: context.Background(),
				id:  c.Id,
				update: &customers.CustomerUpdate{
					FirstName: func(s string) *string { return &s }("Updated Test"),
					LastName:  func(s string) *string { return &s }("Customer"),
					Phone:     func(s string) *string { return &s }("254722000000"),
					Email:     func(s string) *string { return &s }("test@test.com"),
				},
			},
			want: &customers.Customer{
				Id:        c.Id,
				FirstName: "Updated Test",
				LastName:  "Customer",
				Email:     "test@test.com",
				Phone:     "254722000000",
				CreatedAt: c.CreatedAt,
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			wantErr: false,
		},
		{
			name: "Update Product Failed - Invalid Product",
			args: args{
				ctx: context.Background(),
				id:  c.Id,
				update: &customers.CustomerUpdate{
					FirstName: func(s string) *string { return &s }(""),
					LastName:  func(s string) *string { return &s }("Customer"),
					Email:     func(s string) *string { return &s }("test@test.com"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customerRepository.UpdateCustomer(tt.args.ctx, tt.args.id, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("customerRepository.UpdateCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderRepository.UpdateOrderItem() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderRepository.UpdateOrderItem() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderRepository.UpdateOrderItem() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestCustomerService_DeleteProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	customerRepository := db.NewCustomerService(firestoreService)

	c, err := customerRepository.CreateCustomer(ctx, &customers.Customer{
		FirstName: "Test",
		LastName:  "Customer",
		Phone:     "254722000000",
		Email:     "test@test.com",
	})
	if err != nil {
		t.Fatalf("failed to create customers: %v", err)
	}
	defer deleteTestCustomer(t, ctx, customerRepository, c.Id)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete Customer Success",
			args: args{
				ctx: context.Background(),
				id:  c.Id,
			},
			wantErr: false,
		},
		{
			name: "Delete Customer Failure - Invalid ID",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := customerRepository.DeleteCustomer(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("customerRepository.DeleteCustomer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
