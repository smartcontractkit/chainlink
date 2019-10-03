// +build sgx_enclave

package cmd

import "chainlink/core/logger"

// InitEnclave initialized the SGX enclave for use by adapters
func InitEnclave() error {
	logger.Infow("SGX Enclave Loaded")
	return nil
}
