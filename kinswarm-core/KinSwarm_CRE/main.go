package main

/*
#cgo LDFLAGS: -L./lib -lkin_kernel
#include "include/kernel.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	root := [32]byte{0xbd, 0xae, 0xe7, 0xb2, 0x05, 0xce, 0x72, 0x6b}
	totalAmount := uint64(350000000000000000)
	workers := uint32(1000)

	outcome := C.execute_settlement_batch(
		(*C.uint8_t)(unsafe.Pointer(&root[0])),
		C.uint64_t(totalAmount),
		C.uint32(workers),
	)

	if outcome.success {
		fmt.Printf("[EVMAdapter] Pushing Root: %x\n", outcome.root_output)
		fmt.Printf("[SolanaAdapter] Pushing Root: %x\n", outcome.root_output)
		fmt.Printf("[CosmosAdapter] Pushing Root: %x\n", outcome.root_output)
	}
}
