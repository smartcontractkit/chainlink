package main
/*
#cgo LDFLAGS: -L${SRCDIR}/lib -llogic_core
#cgo CFLAGS: -I${SRCDIR}
#include "kernel.h"
*/
import "C"
import (
	"encoding/hex"
	"fmt"
	"time"
	"unsafe"
)
func main() {
	fmt.Println("====================================================")
	fmt.Println("   KINSWARM: CHAINLINK FORK INTEGRATION BENCHMARK")
	fmt.Println("====================================================")
	baseRoot := "4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b"
	b, _ := hex.DecodeString(baseRoot)
	start := time.Now()

	// Explicitly cast Go uint32 to C.uint32_t
	nonce := C.uint32_t(time.Now().UnixNano())

	outcome := C.execute_settlement_anchor(
		(*C.uint8_t)(unsafe.Pointer(&b[0])), 
		C.uint64_t(1000), 
		nonce,
	)
	if bool(outcome.success) {
		res := C.GoBytes(unsafe.Pointer(&outcome.root_output[0]), 32)
		fmt.Printf("Status         : [COMMITTED]\n")
		fmt.Printf("EVM Anchor     : 0x%s\n", hex.EncodeToString(res))
		fmt.Printf("Engine Latency : %v\n", time.Since(start))
		fmt.Println("====================================================")
	}
}
