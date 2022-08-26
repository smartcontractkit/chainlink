package configtracker

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type configTracker struct {
	types.ContractConfigTracker
	instance uint8
}

// New is the constructor of configTracker
func New(instance uint8) types.ContractConfigTracker {
	return &configTracker{
		instance: instance,
	}
}

// Notify may optionally emit notification events when the contract's
// configuration changes. This is purely used as an optimization reducing
// the delay between a configuration change and its enactment. Implementors
// who don't care about this may simply return a nil channel.
//
// The returned channel should never be closed.
func (c *configTracker) Notify() <-chan struct{} {
	return c.ContractConfigTracker.Notify()
}

// LatestConfigDetails returns information about the latest configuration,
// but not the configuration itself.
func (c *configTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	changedInBlock, configDigest, err = c.ContractConfigTracker.LatestConfigDetails(ctx)
	if err != nil {
		return changedInBlock, configDigest, err
	}

	return changedInBlock, c.thresholdDigitalDigest(configDigest), nil
}

// LatestConfig returns the latest configuration.
func (c *configTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (cc types.ContractConfig, err error) {
	cc, err = c.ContractConfigTracker.LatestConfig(ctx, changedInBlock)
	if err != nil {
		return cc, err
	}

	cc.ConfigDigest = c.thresholdDigitalDigest(cc.ConfigDigest)

	return cc, nil
}

// LatestBlockHeight returns the height of the most recent block in the chain.
func (c *configTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return c.ContractConfigTracker.LatestBlockHeight(ctx)
}

func (c *configTracker) thresholdDigitalDigest(root types.ConfigDigest) types.ConfigDigest {
	var thresholdBytes types.ConfigDigest
	for i, b := range root {
		thresholdBytes[i] = b ^ c.instance
	}
	return thresholdBytes
}
