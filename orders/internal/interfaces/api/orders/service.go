package orders

import (
	"context"
	"github.com/leta/order-management-system/orders/generated"
)

// OrderServiceInterface defines the methods that should be implemented by a service that handles order-related operations.
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, in *generated.CreateOrderRequest) (*generated.CreateOrderResponse, error)
	GetOrder(ctx context.Context, in *generated.GetOrderRequest) (*generated.GetOrderResponse, error)
	ListOrders(ctx context.Context, in *generated.ListOrdersRequest) (*generated.ListOrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, in *generated.UpdateOrderStatusRequest) (*generated.UpdateOrderStatusResponse, error)
	DeleteOrder(ctx context.Context, in *generated.DeleteOrderRequest) (*generated.DeleteOrderResponse, error)
	CreateOrderItem(ctx context.Context, in *generated.CreateOrderItemRequest) (*generated.CreateOrderItemResponse, error)
	GetOrderItem(ctx context.Context, in *generated.GetOrderItemRequest) (*generated.GetOrderItemResponse, error)
	ListOrderItems(ctx context.Context, in *generated.ListOrderItemsRequest) (*generated.ListOrderItemsResponse, error)
	UpdateOrderItem(ctx context.Context, in *generated.UpdateOrderItemRequest) (*generated.UpdateOrderItemResponse, error)
	DeleteOrderItem(ctx context.Context, in *generated.DeleteOrderItemRequest) (*generated.DeleteOrderItemResponse, error)
}
