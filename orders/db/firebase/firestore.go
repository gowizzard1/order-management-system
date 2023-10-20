package firebase

import (
	"cloud.google.com/go/firestore"
)

type FirestoreService struct {
	Client *firestore.Client
}

func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		Client: client,
	}
}
