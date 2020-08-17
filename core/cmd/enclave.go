// +build !sgx_enclave

package cmd

import (
	"github.com/smartcontractkit/chainlink/core/logger"
)

// InitEnclave is a stub in non SGX enabled builds.
func InitEnclave() error {
	logger.Infow("SGX enclave *NOT* loaded")
	logger.Infow("This version of chainlink was not built with support for SGX tasks")
	return nil
}
