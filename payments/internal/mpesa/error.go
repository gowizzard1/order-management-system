package mpesa

import (
	"regexp"

	"github.com/leta/order-management-system/payments/internal/service"
)

const (
	MPESA_ERR_INVALID_REQUEST = "400.002.05"
)

func MpesaErrorToInternalError(err error) error {
	if err == nil {
		return nil
	}

	code, ok := findCodeInError(err, MPESA_ERR_INVALID_REQUEST)
	if ok {
		return service.Errorf(service.INVALID_ERROR, "invalid request: %s", code)
	}

	return service.Errorf(service.INTERNAL_ERROR, "internal error: %v", err)

}

func findCodeInError(err error, code string) (string, bool) {
	if err == nil {
		return "", false
	}

	quotedPattern := regexp.QuoteMeta(code)

	re := regexp.MustCompile(`(?m)` + quotedPattern + `: (.*)`)
	match := re.FindStringSubmatch(err.Error())

	return match[0], len(match) > 0
}
