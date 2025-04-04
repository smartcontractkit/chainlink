package ocrimpls

import (
	"context"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type configTracker struct {
	cfg          cctypes.OCR3ConfigWithMeta
	addressCodec ccipocr3.AddressCodec
}

func NewConfigTracker(cfg cctypes.OCR3ConfigWithMeta, addressCodec ccipocr3.AddressCodec) *configTracker {
	return &configTracker{cfg: cfg, addressCodec: addressCodec}
}

// LatestBlockHeight implements types.ContractConfigTracker.
func (c *configTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, nil
}

// LatestConfig implements types.ContractConfigTracker.
func (c *configTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	return c.contractConfig(), nil
}

// LatestConfigDetails implements types.ContractConfigTracker.
func (c *configTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	return 0, c.cfg.ConfigDigest, nil
}

// Notify implements types.ContractConfigTracker.
func (c *configTracker) Notify() <-chan struct{} {
	return nil
}

func (c *configTracker) contractConfig() types.ContractConfig {
	var signers [][]byte
	var transmitters [][]byte
	for _, node := range c.cfg.Config.Nodes {
		signers = append(signers, node.SignerKey)
		transmitters = append(transmitters, node.TransmitterKey)
	}

	return types.ContractConfig{
		ConfigDigest:          c.cfg.ConfigDigest,
		ConfigCount:           uint64(c.cfg.Version),
		Signers:               toOnchainPublicKeys(signers),
		Transmitters:          toOCRAccounts(transmitters, c.addressCodec, c.cfg.Config.ChainSelector),
		F:                     c.cfg.Config.FRoleDON,
		OnchainConfig:         []byte{},
		OffchainConfigVersion: c.cfg.Config.OffchainConfigVersion,
		OffchainConfig:        c.cfg.Config.OffchainConfig,
	}
}

// PublicConfig returns the OCR configuration as a PublicConfig so that we can
// access ReportingPluginConfig and other fields prior to launching the plugins.
func (c *configTracker) PublicConfig() (ocr3confighelper.PublicConfig, error) {
	return ocr3confighelper.PublicConfigFromContractConfig(false, c.contractConfig())
}

func toOnchainPublicKeys(signers [][]byte) []types.OnchainPublicKey {
	keys := make([]types.OnchainPublicKey, len(signers))
	for i, signer := range signers {
		keys[i] = types.OnchainPublicKey(signer)
	}
	return keys
}

func toOCRAccounts(transmitters [][]byte, addressCodec ccipocr3.AddressCodec, chainSelector ccipocr3.ChainSelector) []types.Account {
	accounts := make([]types.Account, len(transmitters))
	for i, transmitter := range transmitters {
		address, _ := addressCodec.AddressBytesToString(transmitter, chainSelector)
		accounts[i] = types.Account(address)
	}
	return accounts
}

var _ types.ContractConfigTracker = (*configTracker)(nil)
