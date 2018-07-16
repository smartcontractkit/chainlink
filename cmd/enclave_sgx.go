// +build sgx_enclave

package cmd

/*
#cgo LDFLAGS: -L ../sgx/target/ -ladapters
#include "../sgx/libadapters/adapters.h"
*/
import "C"
import (
	"fmt"

	"github.com/smartcontractkit/chainlink/logger"
)

// InitEnclave initialized the SGX enclave for use by adapters
func InitEnclave() error {
	_, err := C.init_enclave()
	if err != nil {
		return fmt.Errorf("error initializing SGX enclave: %+v", err)
	}
	logger.Infow("SGX Enclave Loaded")
	return nil
}
