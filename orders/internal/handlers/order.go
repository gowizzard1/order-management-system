package handlers

import (
	"context"
	"github.com/leta/order-management-system/orders/generated"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/orders"
	"github.com/leta/order-management-system/orders/pkg/utils"
)

func (s *GRPCServer) CreateOrder(ctx context.Context, in *generated.CreateOrderRequest) (*generated.CreateOrderResponse, error) {

	var items []*orders.OrderItem
	for _, item := range in.GetOrderItems() {
		items = append(items, &orders.OrderItem{
			ProductId: item.GetProductId(),
			Quantity:  uint(item.GetQuantity()),
		})
	}

	p, err := s.OrderRepository.CreateOrder(ctx, &orders.Order{
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

func (s *GRPCServer) GetOrder(ctx context.Context, in *generated.GetOrderRequest) (*generated.GetOrderResponse, error) {

	order, err := s.OrderRepository.GetOrder(ctx, in.GetId())
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

func (s *GRPCServer) ListOrders(ctx context.Context, in *generated.ListOrdersRequest) (*generated.ListOrdersResponse, error) {

	orders, err := s.OrderRepository.ListOrders(ctx)
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

func (s *GRPCServer) UpdateOrderStatus(
	ctx context.Context, in *generated.UpdateOrderStatusRequest) (*generated.UpdateOrderStatusResponse, error) {

	order, err := s.OrderRepository.UpdateOrderStatus(ctx, in.GetId(), getOrderStatus(in.GetStatus()))
	if err != nil {
		return nil, err
	}

	return &generated.UpdateOrderStatusResponse{
		Id:         order.Id,
		CustomerId: order.CustomerId,
		Status:     getGRPCOrderStatus(order.OrderStatus),
	}, nil
}

func (s *GRPCServer) DeleteOrder(ctx context.Context, in *generated.DeleteOrderRequest) (*generated.DeleteOrderResponse, error) {

	err := s.OrderRepository.DeleteOrder(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &generated.DeleteOrderResponse{}, nil
}

func (s *GRPCServer) CreateOrderItem(
	ctx context.Context, in *generated.CreateOrderItemRequest) (*generated.CreateOrderItemResponse, error) {

	orderItem, err := s.OrderRepository.CreateOrderItem(ctx, in.GetOrderId(), &orders.OrderItem{
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

func (s *GRPCServer) GetOrderItem(ctx context.Context, in *generated.GetOrderItemRequest) (*generated.GetOrderItemResponse, error) {

	orderItem, err := s.OrderRepository.GetOrderItem(ctx, in.GetOrderId(), in.GetId())
	if err != nil {
		return nil, err
	}

	return &generated.GetOrderItemResponse{
		Id:        orderItem.Id,
		ProductId: orderItem.ProductId,
		Quantity:  uint32(orderItem.Quantity),
	}, nil
}

func (s *GRPCServer) ListOrderItems(
	ctx context.Context, in *generated.ListOrderItemsRequest) (*generated.ListOrderItemsResponse, error) {

	orderItems, err := s.OrderRepository.ListOrderItems(ctx, in.GetOrderId())
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

func (s *GRPCServer) UpdateOrderItem(
	ctx context.Context, in *generated.UpdateOrderItemRequest) (*generated.UpdateOrderItemResponse, error) {

	var quantity *uint
	if in.GetUpdate().GetQuantity() > 0 {
		quantity = new(uint)
		*quantity = uint(in.GetUpdate().GetQuantity())
	}
	orderItem, err := s.OrderRepository.UpdateOrderItem(ctx, in.GetOrderId(), in.GetId(), &orders.OrderItemUpdate{
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

func (s *GRPCServer) DeleteOrderItem(
	ctx context.Context, in *generated.DeleteOrderItemRequest) (*generated.DeleteOrderItemResponse, error) {

	err := s.OrderRepository.DeleteOrderItem(ctx, in.GetOrderId(), in.GetId())
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
