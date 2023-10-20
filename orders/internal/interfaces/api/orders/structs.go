package orders

import (
	"github.com/leta/order-management-system/orders/pkg/utils"
)

type OrderItem struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (o *OrderItem) Validate() error {
	if o.ProductId == "" {
		return utils.Errorf(utils.INVALID_ERROR, "product_id is required")
	}

	if o.Quantity == 0 {
		return utils.Errorf(utils.INVALID_ERROR, "quantity is required")
	}

	return nil
}

type OrderItemUpdate struct {
	Quantity *uint `json:"quantity"`
}

type Order struct {
	Id          string            `json:"id"`
	CustomerId  string            `json:"customer_id"`
	Items       []*OrderItem      `json:"items"`
	OrderStatus utils.OrderStatus `json:"order_status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

func (o *Order) Validate() error {
	if o.CustomerId == "" {
		return utils.Errorf(utils.INVALID_ERROR, "customer_id is required")
	}

	if len(o.Items) == 0 {
		return utils.Errorf(utils.INVALID_ERROR, "items are required")
	}

	for _, item := range o.Items {
		err := item.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
