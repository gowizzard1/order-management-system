package orders_test

import (
	"context"
	"encoding/json"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/orders"
	"github.com/leta/order-management-system/shared"
	"reflect"
	"testing"
	"time"

	db "github.com/leta/order-management-system/orders/db/firebase"
)

func deleteTestOrder(t *testing.T, ctx context.Context, orderRepository orders.OrderRepository, id string) {
	err := orderRepository.DeleteOrder(ctx, id)
	if err != nil {
		t.Fatalf("failed to delete orders: %v", err)
	}
}

func TestOrderService_CheckPreconditions(t *testing.T) {
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
			s := db.NewOrderRepository(tt.fields.db)
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("OrderRepository.CheckPreconditions() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			s.CheckPreconditions()
		})
	}
}

func TestOrderService_CreateOrder(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	type args struct {
		ctx   context.Context
		order *orders.Order
	}
	tests := []struct {
		name    string
		args    args
		want    *orders.Order
		wantErr bool
	}{
		{
			name: "Create Order Success",
			args: args{
				ctx: ctx,
				order: &orders.Order{
					CustomerId: "customers-1",
					Items: []*orders.OrderItem{
						{
							ProductId: "product-1",
							Quantity:  1,
						},
					},
				},
			},
			want: &orders.Order{
				CustomerId:  "customers-1",
				OrderStatus: shared.OrderStatusNew,
				Items: []*orders.OrderItem{
					{
						ProductId: "product-1",
						Quantity:  1,
						UpdatedAt: time.Now().Format(time.RFC3339),
						CreatedAt: time.Now().Format(time.RFC3339),
					},
				},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
		},
		{
			name: "Create Order Success - Multiple Items",
			args: args{
				ctx: ctx,
				order: &orders.Order{
					CustomerId: "customers-1",
					Items: []*orders.OrderItem{
						{
							ProductId: "product-1",
							Quantity:  1,
						},
						{
							ProductId: "product-2",
							Quantity:  1,
						},
					},
				},
			},
			want: &orders.Order{
				CustomerId:  "customers-1",
				OrderStatus: shared.OrderStatusNew,
				Items: []*orders.OrderItem{
					{
						ProductId: "product-1",
						Quantity:  1,
						UpdatedAt: time.Now().Format(time.RFC3339),
						CreatedAt: time.Now().Format(time.RFC3339),
					},
					{
						ProductId: "product-2",
						Quantity:  1,
						UpdatedAt: time.Now().Format(time.RFC3339),
						CreatedAt: time.Now().Format(time.RFC3339),
					},
				},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
		},
		{
			name: "Create Order Error - Missing Customer ID",
			args: args{
				ctx: ctx,
				order: &orders.Order{
					Items: []*orders.OrderItem{
						{
							ProductId: "product-1",
							Quantity:  1,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create Order Error - Missing Items",
			args: args{
				ctx: ctx,
				order: &orders.Order{
					CustomerId: "customers-1",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderRepository.CreateOrder(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				// Clear out the fields that are set by the DB
				defer deleteTestOrder(t, ctx, orderRepository, got.Id)

				got.Id = ""
				got.CreatedAt = ""
				got.UpdatedAt = ""
				for _, item := range got.Items {
					item.Id = ""
					item.CreatedAt = ""
					item.UpdatedAt = ""
				}
			}

			if tt.want != nil {
				tt.want.Id = ""
				tt.want.CreatedAt = ""
				tt.want.UpdatedAt = ""
				for _, item := range tt.want.Items {
					item.Id = ""
					item.CreatedAt = ""
					item.UpdatedAt = ""
				}
			}

			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderRepository.CreateOrder() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderRepository.CreateOrder() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderRepository.CreateOrder() = %s, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_GetOrder(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *orders.Order
		wantErr bool
	}{
		{
			name: "Get Order Success",
			args: args{
				ctx: ctx,
				id:  testOrder.Id,
			},
			want:    testOrder,
			wantErr: false,
		},
		{
			name: "Get Order Error - Missing ID",
			args: args{
				ctx: ctx,
				id:  "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderRepository.GetOrder(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.GetOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderRepository.GetOrder() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderRepository.GetOrder() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderRepository.GetOrder() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_ListOrders(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*orders.Order
		wantErr bool
	}{
		{
			name: "List Orders Success",
			args: args{
				ctx: ctx,
			},
			want: []*orders.Order{
				testOrder,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderRepository.ListOrders(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderRepository.ListOrders() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderRepository.ListOrders() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderRepository.ListOrders() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_DeleteOrder(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

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
			name:    "Delete Order Success",
			args:    args{ctx: ctx, id: testOrder.Id},
			wantErr: false,
		},
		{
			name:    "Delete Order Error - Missing ID",
			args:    args{ctx: ctx, id: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := orderRepository.DeleteOrder(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.DeleteOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderService_CreateOrderItem(t *testing.T) {

	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx       context.Context
		orderId   string
		orderItem *orders.OrderItem
	}
	tests := []struct {
		name    string
		args    args
		want    *orders.OrderItem
		wantErr bool
	}{
		{
			name: "Create Order Item Success",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				orderItem: &orders.OrderItem{
					ProductId: "product-2",
					Quantity:  1,
				},
			},
			want: &orders.OrderItem{
				ProductId: "product-2",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
		{
			name: "Create Order Item Error - Missing Order ID",
			args: args{
				ctx:       ctx,
				orderId:   "",
				orderItem: &orders.OrderItem{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create Order Item Error - Missing Product ID",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				orderItem: &orders.OrderItem{
					Quantity: 1,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create Order Item Error - Missing Quantity",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				orderItem: &orders.OrderItem{
					ProductId: "product-2",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// no need to delete the orders item as it will be deleted when the orders is deleted
			got, err := orderRepository.CreateOrderItem(tt.args.ctx, tt.args.orderId, tt.args.orderItem)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.CreateOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				got.Id = ""
			}

			if tt.want != nil {
				tt.want.Id = ""
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.CreateOrderItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_CreateOrderItems(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx        context.Context
		orderId    string
		orderItems []*orders.OrderItem
	}
	tests := []struct {
		name    string
		args    args
		want    []*orders.OrderItem
		wantErr bool
	}{
		{
			name: "Create Order Items Success",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				orderItems: []*orders.OrderItem{
					{
						ProductId: "product-2",
						Quantity:  1,
					},
					{
						ProductId: "product-3",
						Quantity:  1,
					},
				},
			},
			want: []*orders.OrderItem{
				{
					ProductId: "product-2",
					Quantity:  1,
					UpdatedAt: time.Now().Format(time.RFC3339),
					CreatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ProductId: "product-3",
					Quantity:  1,
					UpdatedAt: time.Now().Format(time.RFC3339),
					CreatedAt: time.Now().Format(time.RFC3339),
				},
			},
		},
		{
			name: "Create Order Items Error - Missing Order ID",
			args: args{
				ctx:        ctx,
				orderId:    "",
				orderItems: []*orders.OrderItem{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create Order Items Error - Missing Product ID",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				orderItems: []*orders.OrderItem{
					{
						Quantity: 1,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create Order Items Error - Missing Quantity",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				orderItems: []*orders.OrderItem{
					{
						ProductId: "product-2",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderRepository.CreateOrderItems(tt.args.ctx, tt.args.orderId, tt.args.orderItems)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.CreateOrderItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, item := range got {
				if item != nil {
					item.Id = ""
				}
			}

			for _, item := range tt.want {
				if item != nil {
					item.Id = ""
				}
			}

			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderRepository.CreateOrderItems() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderRepository.CreateOrderItems() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderRepository.CreateOrderItems() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestOrderService_GetOrderItem(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx         context.Context
		orderId     string
		orderItemId string
	}
	tests := []struct {
		name string

		args    args
		want    *orders.OrderItem
		wantErr bool
	}{
		{
			name: "Get Order Item Success",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: testOrder.Items[0].Id,
			},
			want: testOrder.Items[0],
		},
		{
			name: "Get Order Item Error - Missing Order ID",
			args: args{
				ctx:         ctx,
				orderId:     "",
				orderItemId: testOrder.Items[0].Id,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get Order Item Error - Missing Order Item ID",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderRepository.GetOrderItem(tt.args.ctx, tt.args.orderId, tt.args.orderItemId)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.GetOrderItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.GetOrderItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_ListOrderItems(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx     context.Context
		orderId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*orders.OrderItem
		wantErr bool
	}{
		{
			name: "List Order Items Success",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
			},
			want: testOrder.Items,
		},
		{
			name: "List Order Items Error - Missing Order ID",
			args: args{
				ctx:     ctx,
				orderId: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderRepository.ListOrderItems(tt.args.ctx, tt.args.orderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.ListOrderItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.ListOrderItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_UpdateOrderItem(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx         context.Context
		orderId     string
		orderItemId string
		update      *orders.OrderItemUpdate
	}
	tests := []struct {
		name    string
		args    args
		want    *orders.OrderItem
		wantErr bool
	}{
		{
			name: "Update Order Item Success",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: testOrder.Items[0].Id,
				update: &orders.OrderItemUpdate{
					Quantity: func(i uint) *uint { return &i }(2),
				},
			},
			want: &orders.OrderItem{
				Id:        testOrder.Items[0].Id,
				ProductId: "product-1",
				Quantity:  2,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
		{
			name: "Update Order Item Error - Missing Order ID",
			args: args{
				ctx:         ctx,
				orderId:     "",
				orderItemId: testOrder.Items[0].Id,
				update: &orders.OrderItemUpdate{
					Quantity: func(i uint) *uint { return &i }(2),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Update Order Item Error - Missing Order Item ID",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: "",
				update: &orders.OrderItemUpdate{
					Quantity: func(i uint) *uint { return &i }(2),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orderRepository.UpdateOrderItem(tt.args.ctx, tt.args.orderId, tt.args.orderItemId, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.UpdateOrderItem() error = %v, wantErr %v", err, tt.wantErr)
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

func TestOrderService_DeleteOrderItem(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx         context.Context
		orderId     string
		orderItemId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete Order Item Success",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: testOrder.Items[0].Id,
			},
		},
		{
			name: "Delete Order Item Error - Missing Order ID",
			args: args{
				ctx:         ctx,
				orderId:     "",
				orderItemId: testOrder.Items[0].Id,
			},
			wantErr: true,
		},
		{
			name: "Delete Order Item Error - Missing Order Item ID",
			args: args{
				ctx:         ctx,
				orderId:     testOrder.Id,
				orderItemId: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := orderRepository.DeleteOrderItem(tt.args.ctx, tt.args.orderId, tt.args.orderItemId); (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.DeleteOrderItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
	ctx := context.Background()

	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	orderRepository := db.NewOrderRepository(firestoreService)

	testOrder, err := orderRepository.CreateOrder(ctx, &orders.Order{
		CustomerId: "customers-1",
		Items: []*orders.OrderItem{
			{
				ProductId: "product-1",
				Quantity:  1,
				UpdatedAt: time.Now().Format(time.RFC3339),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test orders: %v", err)
	}
	defer deleteTestOrder(t, ctx, orderRepository, testOrder.Id)

	type args struct {
		ctx     context.Context
		orderId string
		status  shared.OrderStatus
	}
	tests := []struct {
		name string

		args    args
		want    *orders.Order
		wantErr bool
	}{
		{
			name: "Update Order Status Success",
			args: args{
				ctx:     ctx,
				orderId: testOrder.Id,
				status:  shared.OrderStatusPaid,
			},
			want: &orders.Order{
				Id:          testOrder.Id,
				CustomerId:  "customers-1",
				OrderStatus: shared.OrderStatusPaid,
				Items: []*orders.OrderItem{
					{
						Id:        testOrder.Items[0].Id,
						ProductId: "product-1",
						Quantity:  1,
						UpdatedAt: time.Now().Format(time.RFC3339),
						CreatedAt: time.Now().Format(time.RFC3339),
					},
				},
				UpdatedAt: testOrder.UpdatedAt,
				CreatedAt: testOrder.CreatedAt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := orderRepository.UpdateOrderStatus(tt.args.ctx, tt.args.orderId, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.UpdateOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				// Clear out the fields that are set by the DB
				defer deleteTestOrder(t, ctx, orderRepository, got.Id)

				got.Id = ""
				got.CreatedAt = ""
				got.UpdatedAt = ""
				for _, item := range got.Items {
					item.Id = ""
					item.CreatedAt = ""
					item.UpdatedAt = ""
				}
			}

			if tt.want != nil {
				tt.want.Id = ""
				tt.want.CreatedAt = ""
				tt.want.UpdatedAt = ""
				for _, item := range tt.want.Items {
					item.Id = ""
					item.CreatedAt = ""
					item.UpdatedAt = ""
				}
			}

			// Convert structs to JSON for easier comparison
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("OrderRepository.UpdateOrderStatus() error = %v", err)
			}
			wantJSON, _ := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("OrderRepository.UpdateOrderStatus() error = %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("OrderRepository.UpdateOrderStatus() = %s, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}
