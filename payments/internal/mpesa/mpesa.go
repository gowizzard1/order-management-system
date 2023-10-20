package mpesa

import (
	"github.com/leta/order-management-system/payments/pkg/utils"
	"net/http"

	"github.com/jwambugu/mpesa-golang-sdk"
)

const (
	MPESA_CONSUMER_KEY    = "MPESA_CONSUMER_KEY"    // #nosec G101 - This is an env variable name
	MPESA_CONSUMER_SECRET = "MPESA_CONSUMER_SECRET" // #nosec G101 - This is an env variable name
	ENVIRONMENT           = "ENVIRONMENT"
)

type Mpesa struct {
	app *mpesa.Mpesa
}

func NewMpesaService() *Mpesa {

	consumerKey := utils.MustGetEnv(MPESA_CONSUMER_KEY)
	consumerSecret := utils.MustGetEnv(MPESA_CONSUMER_SECRET)

	env := utils.MustGetEnv(ENVIRONMENT)

	var mpesaEnv mpesa.Environment
	if env == "prod" {
		mpesaEnv = mpesa.Production
	} else {
		mpesaEnv = mpesa.Sandbox
	}

	mpesaApp := mpesa.NewApp(http.DefaultClient, consumerKey, consumerSecret, mpesaEnv)

	return &Mpesa{
		app: mpesaApp,
	}
}
