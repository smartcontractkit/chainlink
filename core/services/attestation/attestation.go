// +build sgx_enclave

package attestation

/*
#cgo LDFLAGS: -L../../../sgx/target/ -ladapters
#include <stdlib.h>
#include "../../../sgx/libadapters/adapters.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Report retrieves an enclave attestation report from the attached enclave
func Report() (string, error) {
	buffer := make([]byte, 8192)
	output := (*C.char)(unsafe.Pointer(&buffer[0]))
	bufferCapacity := C.int(len(buffer))
	outputLen := C.int(0)
	outputLenPtr := (*C.int)(unsafe.Pointer(&outputLen))

	if _, err := C.report(output, bufferCapacity, outputLenPtr); err != nil {
		return "", fmt.Errorf("SGX report: %v", err)
	}

	return C.GoStringN(output, outputLen), nil
}
