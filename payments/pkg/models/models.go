package models

type PaymentModel struct {
	ID          string `firestore:"id"`
	Amount      uint   `firestore:"amount"`
	Status      string `firestore:"status"`
	OrderID     string `firestore:"orderId"`
	Phone       string `firestore:"phone"`
	Reference   string `firestore:"reference"`
	Description string `firestore:"description"`
	CreatedAt   string `firestore:"createdAt"`
	UpdatedAt   string `firestore:"updatedAt"`
}
