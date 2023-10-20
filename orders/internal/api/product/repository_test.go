package product_test

import (
	"context"
	db "github.com/leta/order-management-system/orders/db/firebase"
	"github.com/leta/order-management-system/orders/internal/repository"
	"reflect"
	"testing"
	"time"
)

func deleteTestProduct(t *testing.T, ctx context.Context, productRepository repository.ProductRepository, id string) {
	err := productRepository.DeleteProduct(ctx, id)
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}
}

func TestProductRepository_CheckPreconditions(t *testing.T) {
	type fields struct {
		db *db.FirestoreService
	}
	tests := []struct {
		name      string
		fields    fields
		wantPanic bool
	}{
		{
			name: "Check Preconditions Failed - nil DB",
			fields: fields{
				db: nil,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := db.NewProductService(tt.fields.db)
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ProductRepository.CheckPreconditions() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			s.CheckPreconditions()
		})
	}
}

func TestProductRepository_CreateProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productRepository := db.NewProductService(firestoreService)

	type args struct {
		ctx     context.Context
		product *repository.Product
	}
	tests := []struct {
		name    string
		args    args
		want    *repository.Product
		wantErr bool
	}{
		{
			name: "Create Product Success",
			args: args{
				ctx: context.Background(),
				product: &repository.Product{
					Name:        "Test Product",
					Description: "Test Description",
					Price:       100,
				},
			},
			want: &repository.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
			wantErr: false,
		},
		{
			name: "Create Product Failed - Invalid Product",
			args: args{
				ctx: context.Background(),
				product: &repository.Product{
					Name:        "",
					Description: "Test Description",
					Price:       100,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := productRepository.CreateProduct(tt.args.ctx, tt.args.product)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductRepository.CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {

				// Clean up the created product
				defer deleteTestProduct(t, ctx, productRepository, got.Id)

				// Ignore the ID in the comparison since it's unpredictable
				got.Id = ""
			}

			if tt.want != nil {
				tt.want.Id = ""
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductRepository.CreateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductRepository_GetProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productRepository := db.NewProductService(firestoreService)

	p, err := productRepository.CreateProduct(ctx, &repository.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}
	defer deleteTestProduct(t, ctx, productRepository, p.Id)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *repository.Product
		wantErr bool
	}{
		{
			name: "Get Product Success",
			args: args{
				ctx: context.Background(),
				id:  p.Id,
			},
			want:    p,
			wantErr: false,
		},
		{
			name: "Get Product Failed - Invalid ID",
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

			got, err := productRepository.GetProduct(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductRepository.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductRepository.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductRepository_ListProducts(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productRepository := db.NewProductService(firestoreService)

	p, err := productRepository.CreateProduct(ctx, &repository.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	defer deleteTestProduct(t, ctx, productRepository, p.Id)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*repository.Product
		wantErr bool
	}{
		{
			name: "List Products Success",
			args: args{
				ctx: context.Background(),
			},
			want: []*repository.Product{
				p,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := productRepository.ListProducts(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductRepository.ListProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductRepository.ListProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductRepository_UpdateProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productRepository := db.NewProductService(firestoreService)

	p, err := productRepository.CreateProduct(ctx, &repository.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	defer deleteTestProduct(t, ctx, productRepository, p.Id)

	type args struct {
		ctx    context.Context
		id     string
		update *repository.ProductUpdate
	}
	tests := []struct {
		name    string
		args    args
		want    *repository.Product
		wantErr bool
	}{
		{
			name: "Update Product Success",
			args: args{
				ctx: context.Background(),
				id:  p.Id,
				update: &repository.ProductUpdate{
					Name:        pkg.StringPtr("Updated Test Product"),
					Description: pkg.StringPtr("Updated Test Description"),
					Price:       pkg.UintPtr(200),
				},
			},
			want: &repository.Product{
				Id:          p.Id,
				Name:        "Updated Test Product",
				Description: "Updated Test Description",
				Price:       200,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
			wantErr: false,
		},
		{
			name: "Update Product Failed - Invalid Product",
			args: args{
				ctx: context.Background(),
				id:  p.Id,
				update: &repository.ProductUpdate{
					Name:        pkg.StringPtr(""),
					Description: pkg.StringPtr("Updated Test Description"),
					Price:       pkg.UintPtr(200),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := productRepository.UpdateProduct(tt.args.ctx, tt.args.id, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductRepository.UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductRepository.UpdateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductRepository_DeleteProduct(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	productRepository := db.NewProductService(firestoreService)

	p, err := productRepository.CreateProduct(ctx, &repository.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}
	defer deleteTestProduct(t, ctx, productRepository, p.Id)

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
			name: "Delete Product Success",
			args: args{
				ctx: context.Background(),
				id:  p.Id,
			},
			wantErr: false,
		},
		{
			name: "Delete Product Failed - Invalid ID",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := productRepository.DeleteProduct(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ProductRepository.DeleteProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
