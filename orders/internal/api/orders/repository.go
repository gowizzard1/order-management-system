package orders

import (
	"github.com/leta/order-management-system/orders/pkg/models"
	"github.com/leta/order-management-system/orders/pkg/utils"

	"context"
	"github.com/leta/order-management-system/orders/db/firebase"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/leta/order-management-system/orders/internal/repository"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ repository.OrderRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	db *firebase.FirestoreService
}

func NewOrderRepository(db *firebase.FirestoreService) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (s *OrderRepository) CheckPreconditions() {
	if s.db == nil {
		panic("no DB service provided")
	}
}

func (r *OrderRepository) orderCollection() *firestore.CollectionRef {
	r.CheckPreconditions()

	return r.db.Client.Collection("orders")
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *repository.Order) (*repository.Order, error) {
	r.CheckPreconditions()

	// Set CreatedAt and UpdatedAt to the current time
	currentTime := time.Now()

	order.CreatedAt = currentTime.Format(time.RFC3339)
	order.UpdatedAt = currentTime.Format(time.RFC3339)

	order.OrderStatus = utils.OrderStatusNew

	err := order.Validate()
	if err != nil {
		return nil, utils.Errorf(utils.INVALID_ERROR, "invalid orders details provided: %v", err)
	}

	orderModel := r.marshallOrder(order)

	docRef, _, err := r.orderCollection().Add(ctx, orderModel)
	if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to create orders: %v", err)
	}

	order.Id = docRef.ID

	_, err = r.CreateOrderItems(ctx, order.Id, order.Items)
	if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to create orders items: %v", err)
	}

	return order, nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, id string) (*repository.Order, error) {
	r.CheckPreconditions()

	if id == "" {
		return nil, utils.Errorf(utils.INVALID_ERROR, "id is required")
	}

	docRef, err := r.orderCollection().Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, utils.Errorf(utils.NOT_FOUND_ERROR, "orders not found")
	} else if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to get orders: %v", err)
	}

	orderModel := &models.OrderModel{}
	if err := docRef.DataTo(orderModel); err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to unmarshall orders: %v", err)
	}

	order := r.unmarshallOrder(orderModel)

	order.Id = id

	orderItems, err := r.ListOrderItems(ctx, id)
	if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to get orders items: %v", err)
	}

	order.Items = orderItems

	return order, nil
}

func (r *OrderRepository) ListOrders(ctx context.Context) ([]*repository.Order, error) {
	r.CheckPreconditions()

	iter := r.orderCollection().Documents(ctx)

	orders := make([]*repository.Order, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to iterate orders: %v", err)
		}

		orderModel := &models.OrderModel{}
		if err := doc.DataTo(orderModel); err != nil {
			return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to unmarshall orders: %v", err)
		}

		order := r.unmarshallOrder(orderModel)

		order.Id = doc.Ref.ID

		orderItems, err := r.ListOrderItems(ctx, order.Id)
		if err != nil {
			return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to get orders items: %v", err)
		}

		order.Items = orderItems

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(
	ctx context.Context, orderId string, status utils.OrderStatus) (*repository.Order, error) {
	r.CheckPreconditions()

	order, err := r.GetOrder(ctx, orderId)
	if err != nil {
		return nil, err
	}

	order.OrderStatus = status
	order.UpdatedAt = time.Now().Format(time.RFC3339)

	orderModel := r.marshallOrder(order)
	_, err = r.orderCollection().Doc(orderId).Set(ctx, orderModel)
	if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to update orders status: %v", err)
	}

	return r.GetOrder(ctx, orderId)
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, id string) error {
	r.CheckPreconditions()

	_, err := r.orderCollection().Doc(id).Delete(ctx)
	if err != nil {
		return utils.Errorf(utils.INTERNAL_ERROR, "failed to delete orders: %v", err)
	}

	return nil
}

func (r *OrderRepository) orderItemCollection(orderId string) *firestore.CollectionRef {
	r.CheckPreconditions()

	return r.orderCollection().Doc(orderId).Collection("items")
}

func (r *OrderRepository) CreateOrderItem(
	ctx context.Context, orderId string, orderItem *repository.OrderItem) (*repository.OrderItem, error) {

	r.CheckPreconditions()

	// Set CreatedAt and UpdatedAt to the current time
	currentTime := time.Now()

	orderItem.CreatedAt = currentTime.Format(time.RFC3339)
	orderItem.UpdatedAt = currentTime.Format(time.RFC3339)

	err := orderItem.Validate()
	if err != nil {
		return nil, utils.Errorf(utils.INVALID_ERROR, "invalid orders item details provided: %v", err)
	}

	orderItemModel := r.marshallOrderItem(orderItem)

	docRef, _, err := r.orderItemCollection(orderId).Add(ctx, orderItemModel)
	if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to create orders item: %v", err)
	}

	orderItem.Id = docRef.ID

	return orderItem, nil
}

func (r *OrderRepository) CreateOrderItems(
	ctx context.Context, orderId string, orderItems []*repository.OrderItem) ([]*repository.OrderItem, error) {

	r.CheckPreconditions()

	if orderId == "" {
		return nil, utils.Errorf(utils.INVALID_ERROR, "orders id is required")
	}

	bulkWriter := r.db.Client.BulkWriter(ctx)

	currentTime := time.Now().Format(time.RFC3339)
	var createdOrderItems []*repository.OrderItem

	for _, orderItem := range orderItems {
		// Set CreatedAt and UpdatedAt to the current time
		orderItem.CreatedAt = currentTime
		orderItem.UpdatedAt = currentTime

		err := orderItem.Validate()
		if err != nil {
			return nil, utils.Errorf(utils.INVALID_ERROR, "invalid orders item details provided: %v", err)
		}

		orderItemModel := r.marshallOrderItem(orderItem)
		docRef := r.orderItemCollection(orderId).NewDoc() // Create a new document reference.

		orderItem.Id = docRef.ID
		createdOrderItems = append(createdOrderItems, orderItem)

		_, err = bulkWriter.Create(docRef, orderItemModel)
		if err != nil {
			return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to create orders item: %v", err)
		}
	}

	bulkWriter.Flush()

	return createdOrderItems, nil
}

func (r *OrderRepository) GetOrderItem(
	ctx context.Context, orderId string, orderItemId string) (*repository.OrderItem, error) {

	r.CheckPreconditions()

	if orderId == "" {
		return nil, utils.Errorf(utils.INVALID_ERROR, "orders id is required")
	} else if orderItemId == "" {
		return nil, utils.Errorf(utils.INVALID_ERROR, "orders item id is required")
	}

	docRef, err := r.orderItemCollection(orderId).Doc(orderItemId).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, utils.Errorf(utils.NOT_FOUND_ERROR, "orders item not found")
	} else if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to get orders item: %v", err)
	}

	orderItemModel := &models.OrderItemModel{}
	if err := docRef.DataTo(orderItemModel); err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to unmarshall orders item: %v", err)
	}

	orderItem := r.unmarshallOrderItem(orderItemModel)

	orderItem.Id = docRef.Ref.ID

	return orderItem, nil
}

func (r *OrderRepository) ListOrderItems(ctx context.Context, orderId string) ([]*repository.OrderItem, error) {
	r.CheckPreconditions()

	if orderId == "" {
		return nil, utils.Errorf(utils.INVALID_ERROR, "orders id is required")
	}

	iter := r.orderItemCollection(orderId).Documents(ctx)

	orderItems := make([]*repository.OrderItem, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		orderItemModel := &models.OrderItemModel{}
		if err := doc.DataTo(orderItemModel); err != nil {
			return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to unmarshall orders item: %v", err)
		}

		orderItem := r.unmarshallOrderItem(orderItemModel)

		orderItem.Id = doc.Ref.ID

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (r *OrderRepository) UpdateOrderItem(
	ctx context.Context, orderId string, orderItemId string, update *repository.OrderItemUpdate) (*repository.OrderItem, error) {

	r.CheckPreconditions()

	if orderId == "" {
		return nil, utils.Errorf(utils.INVALID_ERROR, "orders id is required")
	} else if orderItemId == "" {
		return nil, utils.Errorf(utils.INVALID_ERROR, "orders item id is required")
	}

	orderItem, err := r.GetOrderItem(ctx, orderId, orderItemId)
	if err != nil {
		return nil, err
	}

	if v := update.Quantity; v != nil {
		orderItem.Quantity = *v
	}

	// Set UpdatedAt to the current time
	orderItem.UpdatedAt = time.Now().Format(time.RFC3339)

	orderItemModel := r.marshallOrderItem(orderItem)
	_, err = r.orderItemCollection(orderId).Doc(orderItemId).Set(ctx, orderItemModel)
	if err != nil {
		return nil, utils.Errorf(utils.INTERNAL_ERROR, "failed to update orders item: %v", err)
	}

	return r.GetOrderItem(ctx, orderId, orderItemId)
}

func (r *OrderRepository) DeleteOrderItem(ctx context.Context, orderId string, orderItemId string) error {
	r.CheckPreconditions()

	if orderId == "" {
		return utils.Errorf(utils.INVALID_ERROR, "orders id is required")
	} else if orderItemId == "" {
		return utils.Errorf(utils.INVALID_ERROR, "orders item id is required")
	}

	_, err := r.orderItemCollection(orderId).Doc(orderItemId).Delete(ctx)
	if err != nil {
		return utils.Errorf(utils.INTERNAL_ERROR, "failed to delete orders item: %v", err)
	}

	return nil
}

func (r *OrderRepository) marshallOrder(order *repository.Order) *models.OrderModel {
	return &models.OrderModel{
		CustomerId:  order.CustomerId,
		Items:       r.marshallOrderItems(order.Items),
		OrderStatus: string(order.OrderStatus),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

func (r *OrderRepository) unmarshallOrder(order *models.OrderModel) *repository.Order {
	return &repository.Order{
		CustomerId:  order.CustomerId,
		Items:       r.unmarshallOrderItems(order.Items),
		OrderStatus: utils.OrderStatus(order.OrderStatus),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

func (r *OrderRepository) marshallOrderItems(items []*repository.OrderItem) []*models.OrderItemModel {
	orderItems := make([]*models.OrderItemModel, 0)

	for _, item := range items {
		orderItems = append(orderItems, r.marshallOrderItem(item))
	}

	return orderItems
}

func (r *OrderRepository) unmarshallOrderItems(items []*models.OrderItemModel) []*repository.OrderItem {
	orderItems := make([]*repository.OrderItem, 0)

	for _, item := range items {
		orderItems = append(orderItems, r.unmarshallOrderItem(item))
	}

	return orderItems
}

func (r *OrderRepository) marshallOrderItem(item *repository.OrderItem) *models.OrderItemModel {
	return &models.OrderItemModel{
		ProductId: item.ProductId,
		Quantity:  int(item.Quantity),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func (r *OrderRepository) unmarshallOrderItem(item *models.OrderItemModel) *repository.OrderItem {
	return &repository.OrderItem{
		ProductId: item.ProductId,
		Quantity:  uint(item.Quantity),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
