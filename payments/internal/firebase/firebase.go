package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
)

const (
	//  These are the environment variable names that need to be set in orders for the
	//  Firebase SDK to work.
	GOOGLE_APPLICATION_CREDENTIALS = "GOOGLE_APPLICATION_CREDENTIALS" // #nosec G101 - This is an env variable name
	GOOGLE_CLOUD_PROJECT           = "GOOGLE_CLOUD_PROJECT"           // #nosec G101 - This is an env variable name
	FIRESTORE_EMULATOR_HOST        = "FIRESTORE_EMULATOR_HOST"
	ENVIRONMENT                    = "ENVIRONMENT"
)

type FirebaseService struct {
	app *firebase.App
}

func NewFirebaseService() *FirebaseService {
	// Check preconditions
	// https://firebase.google.com/docs/admin/setup/#initialize_the_sdk_in_non-google_environments
	c := shared.MustGetEnv(GOOGLE_APPLICATION_CREDENTIALS)
	log.Println(c)

	env := shared.MustGetEnv(ENVIRONMENT)
	if env == "dev" || env == "test" {
		shared.MustGetEnv(FIRESTORE_EMULATOR_HOST)
	}

	projectID := shared.MustGetEnv(GOOGLE_CLOUD_PROJECT)

	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return &FirebaseService{
		app: app,
	}
}

func (s *FirebaseService) GetApp() *firebase.App {
	return s.app
}
