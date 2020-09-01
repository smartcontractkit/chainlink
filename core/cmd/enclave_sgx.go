// +build sgx_enclave

package cmd

// InitEnclave initialized the SGX enclave for use by adapters
func InitEnclave() error {
	logger.Infow("SGX Enclave Loaded")
	return nil
}
