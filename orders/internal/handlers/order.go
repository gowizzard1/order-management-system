package handlers

import (
	"context"
	"fmt"
	"github.com/leta/order-management-system/orders/generated"
	"github.com/leta/order-management-system/shared"

	"github.com/leta/order-management-system/orders/internal/repository"
)

func (s *GRPCServer) CreateOrder(ctx context.Context, in *generated.CreateOrderRequest) (*generated.CreateOrderResponse, error) {

	var items []*repository.OrderItem
	for _, item := range in.GetOrderItems() {
		items = append(items, &repository.OrderItem{
			ProductId: item.GetProductId(),
			Quantity:  uint(item.GetQuantity()),
		})
	}

	p, err := s.OrderRepository.CreateOrder(ctx, &repository.Order{
		CustomerId: in.GetCustomerId(),
		Items:      items,
	})
	if err != nil {
		return nil, Error(fmt.Errorf("failed to create orders: %w", err))
	}

	return &orders_service.CreateOrderResponse{
		Id: p.Id,
	}, nil
}

func (s *GRPCServer) GetOrder(ctx context.Context, in *orders_service.GetOrderRequest) (*orders_service.GetOrderResponse, error) {

	order, err := s.OrderRepository.GetOrder(ctx, in.GetId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to get orders: %w", err))
	}

	var orderItems []*orders_service.OrderItem
	for _, item := range order.Items {
		orderItems = append(orderItems, &orders_service.OrderItem{
			ProductId: item.ProductId,
			Quantity:  uint32(item.Quantity),
		})
	}

	return &orders_service.GetOrderResponse{
		Id:         order.Id,
		CustomerId: order.CustomerId,
		OrderItems: orderItems,
		Status:     getGRPCOrderStatus(order.OrderStatus),
	}, nil
}

func (s *GRPCServer) ListOrders(ctx context.Context, in *orders_service.ListOrdersRequest) (*orders_service.ListOrdersResponse, error) {

	orders, err := s.OrderRepository.ListOrders(ctx)
	if err != nil {
		return nil, Error(fmt.Errorf("failed to list orders: %w", err))
	}

	var responseOrders []*orders_service.Order
	for _, p := range orders {
		var orderItems []*orders_service.OrderItem
		for _, item := range p.Items {
			orderItems = append(orderItems, &orders_service.OrderItem{
				ProductId: item.ProductId,
				Quantity:  uint32(item.Quantity),
			})
		}

		responseOrders = append(responseOrders, &orders_service.Order{
			Id:         p.Id,
			CustomerId: p.CustomerId,
			OrderItems: orderItems,
			Status:     getGRPCOrderStatus(p.OrderStatus),
		})
	}

	return &orders_service.ListOrdersResponse{
		Orders: responseOrders,
	}, nil
}

func (s *GRPCServer) UpdateOrderStatus(
	ctx context.Context, in *orders_service.UpdateOrderStatusRequest) (*orders_service.UpdateOrderStatusResponse, error) {

	order, err := s.OrderRepository.UpdateOrderStatus(ctx, in.GetId(), getOrderStatus(in.GetStatus()))
	if err != nil {
		return nil, Error(fmt.Errorf("failed to update orders status: %w", err))
	}

	return &orders_service.UpdateOrderStatusResponse{
		Id:         order.Id,
		CustomerId: order.CustomerId,
		Status:     getGRPCOrderStatus(order.OrderStatus),
	}, nil
}

func (s *GRPCServer) DeleteOrder(ctx context.Context, in *orders_service.DeleteOrderRequest) (*orders_service.DeleteOrderResponse, error) {

	err := s.OrderRepository.DeleteOrder(ctx, in.GetId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to delete orders: %w", err))
	}

	return &orders_service.DeleteOrderResponse{}, nil
}

func (s *GRPCServer) CreateOrderItem(
	ctx context.Context, in *orders_service.CreateOrderItemRequest) (*orders_service.CreateOrderItemResponse, error) {

	orderItem, err := s.OrderRepository.CreateOrderItem(ctx, in.GetOrderId(), &repository.OrderItem{
		ProductId: in.GetProductId(),
		Quantity:  uint(in.GetQuantity()),
	})
	if err != nil {
		return nil, Error(fmt.Errorf("failed to create orders item: %w", err))
	}

	return &orders_service.CreateOrderItemResponse{
		Id: orderItem.Id,
	}, nil
}

func (s *GRPCServer) GetOrderItem(ctx context.Context, in *orders_service.GetOrderItemRequest) (*orders_service.GetOrderItemResponse, error) {

	orderItem, err := s.OrderRepository.GetOrderItem(ctx, in.GetOrderId(), in.GetId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to get orders item: %w", err))
	}

	return &orders_service.GetOrderItemResponse{
		Id:        orderItem.Id,
		ProductId: orderItem.ProductId,
		Quantity:  uint32(orderItem.Quantity),
	}, nil
}

func (s *GRPCServer) ListOrderItems(
	ctx context.Context, in *orders_service.ListOrderItemsRequest) (*orders_service.ListOrderItemsResponse, error) {

	orderItems, err := s.OrderRepository.ListOrderItems(ctx, in.GetOrderId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to list orders items: %w", err))
	}

	var responseOrderItems []*orders_service.OrderItem
	for _, item := range orderItems {
		responseOrderItems = append(responseOrderItems, &orders_service.OrderItem{
			Id:        item.Id,
			ProductId: item.ProductId,
			Quantity:  uint32(item.Quantity),
		})
	}

	return &orders_service.ListOrderItemsResponse{
		OrderItems: responseOrderItems,
	}, nil
}

func (s *GRPCServer) UpdateOrderItem(
	ctx context.Context, in *orders_service.UpdateOrderItemRequest) (*orders_service.UpdateOrderItemResponse, error) {

	var quantity *uint
	if in.GetUpdate().GetQuantity() > 0 {
		quantity = new(uint)
		*quantity = uint(in.GetUpdate().GetQuantity())
	}
	orderItem, err := s.OrderRepository.UpdateOrderItem(ctx, in.GetOrderId(), in.GetId(), &repository.OrderItemUpdate{
		Quantity: quantity,
	})
	if err != nil {
		return nil, Error(fmt.Errorf("failed to update orders item: %w", err))
	}

	return &orders_service.UpdateOrderItemResponse{
		Id:        orderItem.Id,
		ProductId: orderItem.ProductId,
		Quantity:  uint32(orderItem.Quantity),
	}, nil
}

func (s *GRPCServer) DeleteOrderItem(
	ctx context.Context, in *orders_service.DeleteOrderItemRequest) (*orders_service.DeleteOrderItemResponse, error) {

	err := s.OrderRepository.DeleteOrderItem(ctx, in.GetOrderId(), in.GetId())
	if err != nil {
		return nil, Error(fmt.Errorf("failed to delete orders item: %w", err))
	}

	return &orders_service.DeleteOrderItemResponse{}, nil
}

func getGRPCOrderStatus(status shared.OrderStatus) orders_service.OrderStatus {
	switch status {
	case shared.OrderStatusNew:
		return orders_service.OrderStatus_NEW
	case shared.OrderStatusPending:
		return orders_service.OrderStatus_PENDING
	case shared.OrderStatusProcessing:
		return orders_service.OrderStatus_PROCESSING
	case shared.OrderStatusPaid:
		return orders_service.OrderStatus_PAID
	default:
		return orders_service.OrderStatus_UNKNOWN
	}
}

func getOrderStatus(status orders_service.OrderStatus) shared.OrderStatus {
	switch status {
	case orders_service.OrderStatus_NEW:
		return shared.OrderStatusNew
	case orders_service.OrderStatus_PENDING:
		return shared.OrderStatusPending
	case orders_service.OrderStatus_PROCESSING:
		return shared.OrderStatusProcessing
	case orders_service.OrderStatus_PAID:
		return shared.OrderStatusPaid
	default:
		return shared.OrderStatusNew
	}
}
