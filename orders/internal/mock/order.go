package mock

import (
	"context"

	"github.com/leta/order-management-system/orders/internal/repository"
)

var _ repository.OrderRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	CreateOrderFunc       func(ctx context.Context, order *repository.Order) (*repository.Order, error)
	GetOrderFunc          func(ctx context.Context, id string) (*repository.Order, error)
	ListOrdersFunc        func(ctx context.Context) ([]*repository.Order, error)
	UpdateOrderStatusFunc func(ctx context.Context, id string, status shared.OrderStatus) (*repository.Order, error)
	DeleteOrderFunc       func(ctx context.Context, id string) error

	CreateOrderItemFunc func(ctx context.Context, orderId string, item *repository.OrderItem) (*repository.OrderItem, error)
	GetOrderItemFunc    func(ctx context.Context, orderId, itemId string) (*repository.OrderItem, error)
	ListOrderItemsFunc  func(ctx context.Context, orderId string) ([]*repository.OrderItem, error)
	UpdateOrderItemFunc func(
		ctx context.Context, orderId string, itemId string, update *repository.OrderItemUpdate) (*repository.OrderItem, error)
	DeleteOrderItemFunc func(ctx context.Context, orderId string, itemId string) error
}

func (m *OrderRepository) CreateOrder(ctx context.Context, order *repository.Order) (*repository.Order, error) {
	return m.CreateOrderFunc(ctx, order)
}

func (m *OrderRepository) GetOrder(ctx context.Context, id string) (*repository.Order, error) {
	return m.GetOrderFunc(ctx, id)
}

func (m *OrderRepository) ListOrders(ctx context.Context) ([]*repository.Order, error) {
	return m.ListOrdersFunc(ctx)
}

func (m *OrderRepository) UpdateOrderStatus(
	ctx context.Context, id string, status shared.OrderStatus) (*repository.Order, error) {
	return m.UpdateOrderStatusFunc(ctx, id, status)
}

func (m *OrderRepository) DeleteOrder(ctx context.Context, id string) error {
	return m.DeleteOrderFunc(ctx, id)
}

func (m *OrderRepository) CreateOrderItem(
	ctx context.Context, orderId string, item *repository.OrderItem) (*repository.OrderItem, error) {
	return m.CreateOrderItemFunc(ctx, orderId, item)
}

func (m *OrderRepository) GetOrderItem(ctx context.Context, orderId string, itemId string) (*repository.OrderItem, error) {
	return m.GetOrderItemFunc(ctx, orderId, itemId)
}

func (m *OrderRepository) ListOrderItems(ctx context.Context, orderId string) ([]*repository.OrderItem, error) {
	return m.ListOrderItemsFunc(ctx, orderId)
}

func (m *OrderRepository) UpdateOrderItem(
	ctx context.Context, orderId string, itemId string, update *repository.OrderItemUpdate) (*repository.OrderItem, error) {
	return m.UpdateOrderItemFunc(ctx, orderId, itemId, update)
}

func (m *OrderRepository) DeleteOrderItem(ctx context.Context, orderId string, itemId string) error {
	return m.DeleteOrderItemFunc(ctx, orderId, itemId)
}
