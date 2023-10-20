package customers

import (
	"github.com/leta/order-management-system/orders/pkg/utils"
)

type Customer struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	CreatedAt string
	UpdatedAt string
}

func (c *Customer) Validate() error {
	if c.FirstName == "" {
		return utils.Errorf(utils.INVALID_ERROR, "first_name is required")
	}

	if c.LastName == "" {
		return utils.Errorf(utils.INVALID_ERROR, "last_name is required")
	}

	if c.Email == "" {
		return utils.Errorf(utils.INVALID_ERROR, "email is required")
	}

	if c.Phone == "" {
		return utils.Errorf(utils.INVALID_ERROR, "phone is required")
	}

	return nil
}

type CustomerUpdate struct {
	FirstName *string
	LastName  *string
	Phone     *string
	Email     *string
}
