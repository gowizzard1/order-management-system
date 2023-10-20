package product

import (
	"github.com/leta/order-management-system/orders/pkg/utils"
)

type Product struct {
	Id          string
	Name        string
	Description string
	Price       uint
	CreatedAt   string
	UpdatedAt   string
}

type ProductUpdate struct {
	Name        *string
	Description *string
	Price       *uint
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return utils.Errorf(utils.INVALID_ERROR, "name is required")
	}

	return nil
}
