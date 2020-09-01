// +build !sgx_enclave

package cmd

// InitEnclave is a stub in non SGX enabled builds.
func InitEnclave() error {
	logger.Infow("SGX enclave *NOT* loaded")
	logger.Infow("This version of chainlink was not built with support for SGX tasks")
	return nil
}
