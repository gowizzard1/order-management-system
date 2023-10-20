package models

type ProductModel struct {
	Name        string `firestore:"name"`
	Description string `firestore:"description"`
	Price       int    `firestore:"price"`
	CreatedAt   string `firestore:"created_at"`
	UpdatedAt   string `firestore:"updated_at"`
}

type CustomerModel struct {
	FirstName string `firestore:"first_name"`
	LastName  string `firestore:"last_name"`
	Email     string `firestore:"email"`
	Phone     string `firestore:"phone"`
	CreatedAt string `firestore:"created_at"`
	UpdatedAt string `firestore:"updated_at"`
}

type OrderModel struct {
	CustomerId  string            `firestore:"customer_id"`
	Items       []*OrderItemModel `firestore:"items"`
	OrderStatus string            `firestore:"order_status"`
	CreatedAt   string            `firestore:"created_at"`
	UpdatedAt   string            `firestore:"updated_at"`
}

type OrderItemModel struct {
	ProductId string `firestore:"product_id"`
	Quantity  int    `firestore:"quantity"`
	CreatedAt string `firestore:"created_at"`
	UpdatedAt string `firestore:"updated_at"`
}
