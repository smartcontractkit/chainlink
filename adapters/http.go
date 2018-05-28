package adapters

/*
#cgo LDFLAGS: -L./http/target/release/ -lhttp
#include "./http.h"
*/
import "C"

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// HTTPGet requires a URL which is used for a GET request when the adapter is called.
type HTTPGet struct {
	URL models.WebURL `json:"url"`
}

// Perform ensures that the adapter's URL responds to a GET request without
// errors and returns the response body as the "value" field of the result.
func (hga *HTTPGet) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	body, err := C.perform_http_get(C.CString(hga.URL.String()))
	if err != nil {
		return input.WithError(fmt.Errorf(C.GoString(body)))
	}
	return input.WithValue(C.GoString(body))
}

// HTTPPost requires a URL which is used for a POST request when the adapter is called.
type HTTPPost struct {
	URL models.WebURL `json:"url"`
}

// Perform ensures that the adapter's URL responds to a POST request without
// errors and returns the response body as the "value" field of the result.
func (hpa *HTTPPost) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	body, err := C.perform_http_post(C.CString(hpa.URL.String()), C.CString(input.Data.String()))
	if err != nil {
		return input.WithError(fmt.Errorf(C.GoString(body)))
	}
	return input.WithValue(C.GoString(body))
}
