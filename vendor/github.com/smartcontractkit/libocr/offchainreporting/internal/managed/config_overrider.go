package managed

import "github.com/smartcontractkit/libocr/offchainreporting/types"

var _ types.ConfigOverrider = ConfigOverriderWrapper{}

// A wrapper around a types.ConfigOverrider that gracefully handles nil ConfigOverriders
type ConfigOverriderWrapper struct {
	wrapped types.ConfigOverrider
}

func (cow ConfigOverriderWrapper) ConfigOverride() *types.ConfigOverride {
	if cow.wrapped == nil {
		return nil
	}
	return cow.wrapped.ConfigOverride()
}
