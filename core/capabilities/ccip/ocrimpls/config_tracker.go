package ocrimpls

import (
	"context"
	"fmt"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type configTracker struct {
	cfg            cctypes.OCR3ConfigWithMeta
	addressCodec   ccipcommon.AddressCodec
	contractConfig types.ContractConfig
}

func NewConfigTracker(cfg cctypes.OCR3ConfigWithMeta, addressCodec ccipcommon.AddressCodec) (*configTracker, error) {
	contractConfig, err := contractConfigFromOCRConfig(cfg, addressCodec)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract config from ocr config: %w", err)
	}
	return &configTracker{
		cfg:            cfg,
		addressCodec:   addressCodec,
		contractConfig: contractConfig,
	}, nil
}

// LatestBlockHeight implements types.ContractConfigTracker.
func (c *configTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, nil
}

// LatestConfig implements types.ContractConfigTracker.
func (c *configTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	return c.contractConfig, nil
}

// LatestConfigDetails implements types.ContractConfigTracker.
func (c *configTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	return 0, c.cfg.ConfigDigest, nil
}

// Notify implements types.ContractConfigTracker.
func (c *configTracker) Notify() <-chan struct{} {
	return nil
}

func contractConfigFromOCRConfig(cfg cctypes.OCR3ConfigWithMeta, addressCodec ccipcommon.AddressCodec) (types.ContractConfig, error) {
	var signers [][]byte
	var transmitters [][]byte
	var err error
	for oracleID, node := range cfg.Config.Nodes {
		signers = append(signers, node.SignerKey)

		// nil transmitters in the OCR config are valid, it just means that this oracle does not support the destination chain.
		// we generate a canonical address with the oracle ID for the transmitter here so that we don't get an error when calling ocr3confighelper.PublicConfigFromContractConfig.
		// the transmitters will never be used as part of the transmission protocol because the custom schedule should exclude nodes
		// that cannot transmit to the destination chain.
		// this canonical address is defined like so to make it clear that this particular oracle is not able to transmit to the destination chain.
		transmitter := node.TransmitterKey
		if len(transmitter) == 0 {
			// #nosec G115 - Overflow is not a concern in this test scenario
			transmitter, err = addressCodec.OracleIDAsAddressBytes(uint8(oracleID), cfg.Config.ChainSelector)
			if err != nil {
				return types.ContractConfig{}, fmt.Errorf("failed to get transmitter from oracle ID: %w", err)
			}
		}
		transmitters = append(transmitters, transmitter)
	}

	transmitterAccounts, err := toOCRAccounts(transmitters, addressCodec, cfg.Config.ChainSelector)
	if err != nil {
		return types.ContractConfig{}, fmt.Errorf("failed to get transmitter accounts: %w", err)
	}

	return types.ContractConfig{
		ConfigDigest:          cfg.ConfigDigest,
		ConfigCount:           uint64(cfg.Version),
		Signers:               toOnchainPublicKeys(signers),
		Transmitters:          transmitterAccounts,
		F:                     cfg.Config.FRoleDON,
		OnchainConfig:         []byte{},
		OffchainConfigVersion: cfg.Config.OffchainConfigVersion,
		OffchainConfig:        cfg.Config.OffchainConfig,
	}, nil
}

// PublicConfig returns the OCR configuration as a PublicConfig so that we can
// access ReportingPluginConfig and other fields prior to launching the plugins.
func (c *configTracker) PublicConfig() (ocr3confighelper.PublicConfig, error) {
	return ocr3confighelper.PublicConfigFromContractConfig(false, c.contractConfig)
}

func toOnchainPublicKeys(signers [][]byte) []types.OnchainPublicKey {
	keys := make([]types.OnchainPublicKey, len(signers))
	for i, signer := range signers {
		keys[i] = types.OnchainPublicKey(signer)
	}
	return keys
}

func toOCRAccounts(transmitters [][]byte, addressCodec ccipcommon.AddressCodec, chainSelector ccipocr3.ChainSelector) ([]types.Account, error) {
	accounts := make([]types.Account, len(transmitters))
	for i, transmitter := range transmitters {
		address, err := addressCodec.TransmitterBytesToString(transmitter, chainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to get address from transmitter bytes: %w", err)
		}
		accounts[i] = types.Account(address)
	}
	return accounts, nil
}

var _ types.ContractConfigTracker = (*configTracker)(nil)
