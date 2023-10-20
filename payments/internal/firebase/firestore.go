package firebase

import (
	"cloud.google.com/go/firestore"
)

type FirestoreService struct {
	client *firestore.Client
}

func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		client: client,
	}
}
