// +build sgx_enclave

package adapters

/*
#cgo LDFLAGS: -L ../sgx/target/ -ladapters
#include "../sgx/libadapters/adapters.h"
*/
import "C"

import (
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
	_, err := C.http_get(C.CString(hga.URL.String()))
	if err != nil {
		return input.WithError(err)
	}
	return input.WithValue("HTTP GET request performed")
}

// HTTPPost requires a URL which is used for a POST request when the adapter is called.
type HTTPPost struct {
	URL models.WebURL `json:"url"`
}

// Perform ensures that the adapter's URL responds to a POST request without
// errors and returns the response body as the "value" field of the result.
func (hpa *HTTPPost) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	_, err := C.http_post(C.CString(hpa.URL.String()), C.CString(input.Data.String()))
	if err != nil {
		return input.WithError(err)
	}
	return input.WithValue("HTTP POST request performed")
}
