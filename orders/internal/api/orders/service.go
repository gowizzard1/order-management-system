package orders

import (
	"context"
	"github.com/leta/order-management-system/orders/generated"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/orders"
	"github.com/leta/order-management-system/orders/pkg/utils"
)

type orderService struct {
	orderRepo OrderRepository
}

func NewOrdersService(orderRepo OrderRepository) orders.OrderServiceInterface {
	return &orderService{
		orderRepo: orderRepo,
	}
}
func (s *orderService) CreateOrder(ctx context.Context, in *generated.CreateOrderRequest) (*generated.CreateOrderResponse, error) {

	var items []*orders.OrderItem
	for _, item := range in.GetOrderItems() {
		items = append(items, &orders.OrderItem{
			ProductId: item.GetProductId(),
			Quantity:  uint(item.GetQuantity()),
		})
	}

	p, err := s.orderRepo.CreateOrder(ctx, &orders.Order{
		CustomerId: in.GetCustomerId(),
		Items:      items,
	})
	if err != nil {
		return nil, err
	}

	return &generated.CreateOrderResponse{
		Id: p.Id,
	}, nil
}

func (s *orderService) GetOrder(ctx context.Context, in *generated.GetOrderRequest) (*generated.GetOrderResponse, error) {

	order, err := s.orderRepo.GetOrder(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	var orderItems []*generated.OrderItem
	for _, item := range order.Items {
		orderItems = append(orderItems, &generated.OrderItem{
			ProductId: item.ProductId,
			Quantity:  uint32(item.Quantity),
		})
	}

	return &generated.GetOrderResponse{
		Id:         order.Id,
		CustomerId: order.CustomerId,
		OrderItems: orderItems,
		Status:     getGRPCOrderStatus(order.OrderStatus),
	}, nil
}

func (s *orderService) ListOrders(ctx context.Context, in *generated.ListOrdersRequest) (*generated.ListOrdersResponse, error) {

	orders, err := s.orderRepo.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	var responseOrders []*generated.Order
	for _, p := range orders {
		var orderItems []*generated.OrderItem
		for _, item := range p.Items {
			orderItems = append(orderItems, &generated.OrderItem{
				ProductId: item.ProductId,
				Quantity:  uint32(item.Quantity),
			})
		}

		responseOrders = append(responseOrders, &generated.Order{
			Id:         p.Id,
			CustomerId: p.CustomerId,
			OrderItems: orderItems,
			Status:     getGRPCOrderStatus(p.OrderStatus),
		})
	}

	return &generated.ListOrdersResponse{
		Orders: responseOrders,
	}, nil
}

func (s *orderService) UpdateOrderStatus(
	ctx context.Context, in *generated.UpdateOrderStatusRequest) (*generated.UpdateOrderStatusResponse, error) {

	order, err := s.orderRepo.UpdateOrderStatus(ctx, in.GetId(), getOrderStatus(in.GetStatus()))
	if err != nil {
		return nil, err
	}

	return &generated.UpdateOrderStatusResponse{
		Id:         order.Id,
		CustomerId: order.CustomerId,
		Status:     getGRPCOrderStatus(order.OrderStatus),
	}, nil
}

func (s *orderService) DeleteOrder(ctx context.Context, in *generated.DeleteOrderRequest) (*generated.DeleteOrderResponse, error) {

	err := s.orderRepo.DeleteOrder(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &generated.DeleteOrderResponse{}, nil
}

func (s *orderService) CreateOrderItem(
	ctx context.Context, in *generated.CreateOrderItemRequest) (*generated.CreateOrderItemResponse, error) {

	orderItem, err := s.orderRepo.CreateOrderItem(ctx, in.GetOrderId(), &orders.OrderItem{
		ProductId: in.GetProductId(),
		Quantity:  uint(in.GetQuantity()),
	})
	if err != nil {
		return nil, err
	}

	return &generated.CreateOrderItemResponse{
		Id: orderItem.Id,
	}, nil
}

func (s *orderService) GetOrderItem(ctx context.Context, in *generated.GetOrderItemRequest) (*generated.GetOrderItemResponse, error) {

	orderItem, err := s.orderRepo.GetOrderItem(ctx, in.GetOrderId(), in.GetId())
	if err != nil {
		return nil, err
	}

	return &generated.GetOrderItemResponse{
		Id:        orderItem.Id,
		ProductId: orderItem.ProductId,
		Quantity:  uint32(orderItem.Quantity),
	}, nil
}

func (s *orderService) ListOrderItems(
	ctx context.Context, in *generated.ListOrderItemsRequest) (*generated.ListOrderItemsResponse, error) {

	orderItems, err := s.orderRepo.ListOrderItems(ctx, in.GetOrderId())
	if err != nil {
		return nil, err
	}

	var responseOrderItems []*generated.OrderItem
	for _, item := range orderItems {
		responseOrderItems = append(responseOrderItems, &generated.OrderItem{
			Id:        item.Id,
			ProductId: item.ProductId,
			Quantity:  uint32(item.Quantity),
		})
	}

	return &generated.ListOrderItemsResponse{
		OrderItems: responseOrderItems,
	}, nil
}

func (s *orderService) UpdateOrderItem(
	ctx context.Context, in *generated.UpdateOrderItemRequest) (*generated.UpdateOrderItemResponse, error) {

	var quantity *uint
	if in.GetUpdate().GetQuantity() > 0 {
		quantity = new(uint)
		*quantity = uint(in.GetUpdate().GetQuantity())
	}
	orderItem, err := s.orderRepo.UpdateOrderItem(ctx, in.GetOrderId(), in.GetId(), &orders.OrderItemUpdate{
		Quantity: quantity,
	})
	if err != nil {
		return nil, err
	}

	return &generated.UpdateOrderItemResponse{
		Id:        orderItem.Id,
		ProductId: orderItem.ProductId,
		Quantity:  uint32(orderItem.Quantity),
	}, nil
}

func (s *orderService) DeleteOrderItem(
	ctx context.Context, in *generated.DeleteOrderItemRequest) (*generated.DeleteOrderItemResponse, error) {

	err := s.orderRepo.DeleteOrderItem(ctx, in.GetOrderId(), in.GetId())
	if err != nil {
		return nil, err
	}

	return &generated.DeleteOrderItemResponse{}, nil
}

func getGRPCOrderStatus(status utils.OrderStatus) generated.OrderStatus {
	switch status {
	case utils.OrderStatusNew:
		return generated.OrderStatus_NEW
	case utils.OrderStatusPending:
		return generated.OrderStatus_PENDING
	case utils.OrderStatusProcessing:
		return generated.OrderStatus_PROCESSING
	case utils.OrderStatusPaid:
		return generated.OrderStatus_PAID
	default:
		return generated.OrderStatus_UNKNOWN
	}
}

func getOrderStatus(status generated.OrderStatus) utils.OrderStatus {
	switch status {
	case generated.OrderStatus_NEW:
		return utils.OrderStatusNew
	case generated.OrderStatus_PENDING:
		return utils.OrderStatusPending
	case generated.OrderStatus_PROCESSING:
		return utils.OrderStatusProcessing
	case generated.OrderStatus_PAID:
		return utils.OrderStatusPaid
	default:
		return utils.OrderStatusNew
	}
}
