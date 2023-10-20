package service

import (
	"context"
	"github.com/leta/order-management-system/orders/pkg/utils"
)

type OrderItem struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Order struct {
	Id          string            `json:"id"`
	CustomerId  string            `json:"customer_id"`
	Items       []*OrderItem      `json:"items"`
	OrderStatus utils.OrderStatus `json:"order_status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

type CheckoutService interface {
	ProcessCheckout(ctx context.Context, orderID string) (*Order, error)
}
