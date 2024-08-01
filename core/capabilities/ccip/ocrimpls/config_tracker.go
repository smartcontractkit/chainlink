package ocrimpls

import (
	"context"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type configTracker struct {
	cfg cctypes.OCR3ConfigWithMeta
}

func NewConfigTracker(cfg cctypes.OCR3ConfigWithMeta) *configTracker {
	return &configTracker{cfg: cfg}
}

// LatestBlockHeight implements types.ContractConfigTracker.
func (c *configTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, nil
}

// LatestConfig implements types.ContractConfigTracker.
func (c *configTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	return types.ContractConfig{
		ConfigDigest:          c.cfg.ConfigDigest,
		ConfigCount:           c.cfg.ConfigCount,
		Signers:               toOnchainPublicKeys(c.cfg.Config.Signers),
		Transmitters:          toOCRAccounts(c.cfg.Config.Transmitters),
		F:                     c.cfg.Config.F,
		OnchainConfig:         []byte{},
		OffchainConfigVersion: c.cfg.Config.OffchainConfigVersion,
		OffchainConfig:        c.cfg.Config.OffchainConfig,
	}, nil
}

// LatestConfigDetails implements types.ContractConfigTracker.
func (c *configTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	return 0, c.cfg.ConfigDigest, nil
}

// Notify implements types.ContractConfigTracker.
func (c *configTracker) Notify() <-chan struct{} {
	return nil
}

func toOnchainPublicKeys(signers [][]byte) []types.OnchainPublicKey {
	keys := make([]types.OnchainPublicKey, len(signers))
	for i, signer := range signers {
		keys[i] = types.OnchainPublicKey(signer)
	}
	return keys
}

func toOCRAccounts(transmitters [][]byte) []types.Account {
	accounts := make([]types.Account, len(transmitters))
	for i, transmitter := range transmitters {
		// TODO: string-encode the transmitter appropriately to the dest chain family.
		accounts[i] = types.Account(gethcommon.BytesToAddress(transmitter).Hex())
	}
	return accounts
}

var _ types.ContractConfigTracker = (*configTracker)(nil)
