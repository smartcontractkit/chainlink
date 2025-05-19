package dummy

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type configTracker struct {
	lggr logger.Logger
	cfg  ocrtypes.ContractConfig

	changedInBlock uint64
	blockHeight    uint64
}

func NewContractConfigTracker(lggr logger.Logger, cfg ConfigTrackerCfg) (ocrtypes.ContractConfigTracker, error) {
	contractConfig, err := cfg.ToContractConfig()
	if err != nil {
		return nil, err
	}
	return &configTracker{lggr.Named("DummyConfigTracker"), contractConfig, cfg.ChangedInBlock, cfg.BlockHeight}, nil
}

// Notify may optionally emit notification events when the contract's
// configuration changes. This is purely used as an optimization reducing
// the delay between a configuration change and its enactment. Implementors
// who don't care about this may simply return a nil channel.
//
// The returned channel should never be closed.
func (ct *configTracker) Notify() <-chan struct{} {
	return nil
}

// LatestConfigDetails returns information about the latest configuration,
// but not the configuration itself.
func (ct *configTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	return ct.changedInBlock, ct.cfg.ConfigDigest, nil
}

// LatestConfig returns the latest configuration.
func (ct *configTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	return ct.cfg, nil
}

// LatestBlockHeight returns the height of the most recent block in the chain.
func (ct *configTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return ct.blockHeight, nil // set LatestBlockHeight to a high number so that OCR considers it to be confirmed
}
