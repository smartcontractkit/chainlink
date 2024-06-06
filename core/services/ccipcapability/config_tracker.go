package ccipcapability

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	_ ocrtypes.ContractConfigTracker = (*configTracker)(nil)
)

type configTracker struct {
	configUpdates chan OCRConfig
	latestConfig  ocrtypes.ContractConfig
}

func newConfigTracker(configUpdates chan OCRConfig) *configTracker {
	return &configTracker{configUpdates: configUpdates}
}

// Notify may optionally emit notification events when the contract's
// configuration changes. This is purely used as an optimization reducing
// the delay between a configuration change and its enactment. Implementors
// who don't care about this may simply return a nil channel.
//
// The returned channel should never be closed.
func (c *configTracker) Notify() <-chan struct{} {
	return nil
}

// LatestConfigDetails returns information about the latest configuration,
// but not the configuration itself.
func (c *configTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	// TODO: implement
	return 0, [32]byte{}, nil
}

// LatestConfig returns the latest configuration.
func (c *configTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	select {
	case newConfig := <-c.configUpdates:
		c.latestConfig = c.toContractConfig(newConfig)
	default:
	}
	return c.latestConfig, nil
}

// LatestBlockHeight returns the height of the most recent block in the chain.
func (c *configTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	// TODO: implement
	return 0, nil
}

func (c *configTracker) toContractConfig(config OCRConfig) ocrtypes.ContractConfig {
	// TODO: implement
	return ocrtypes.ContractConfig{}
}
