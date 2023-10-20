package orders

import (
	"github.com/leta/order-management-system/orders/pkg/utils"

	"context"
)

// OrderRepositoryInterface is an interface for interacting with order data.
type OrderRepository interface {

	// Order CRUD
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	ListOrders(ctx context.Context) ([]*Order, error)
	UpdateOrderStatus(ctx context.Context, orderId string, status utils.OrderStatus) (*Order, error)
	DeleteOrder(ctx context.Context, id string) error

	// OrderItem CRUD
	CreateOrderItem(ctx context.Context, orderId string, orderItem *OrderItem) (*OrderItem, error)
	GetOrderItem(ctx context.Context, orderId string, orderItemId string) (*OrderItem, error)
	ListOrderItems(ctx context.Context, orderId string) ([]*OrderItem, error)
	UpdateOrderItem(ctx context.Context, orderId string, orderItemId string, update *OrderItemUpdate) (*OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderId string, orderItemId string) error
}
