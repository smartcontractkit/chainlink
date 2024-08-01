package rpc

import "context"

// SpecVersion returns the version of the Starknet JSON-RPC specification being used
// Parameters: None
// Returns: String of the Starknet JSON-RPC specification
func (provider *Provider) SpecVersion(ctx context.Context) (string, error) {
	var result string
	err := do(ctx, provider.c, "starknet_specVersion", &result)
	return result, Err(InternalError, err)
}
