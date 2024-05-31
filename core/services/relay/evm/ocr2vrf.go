package evm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// DKGProvider provides all components needed for a DKG plugin.
type DKGProvider interface {
	commontypes.Plugin
}

// OCR2VRFProvider provides all components needed for a OCR2VRF plugin.
type OCR2VRFProvider interface {
	commontypes.Plugin
}

// OCR2VRFRelayer contains the relayer and instantiating functions for OCR2VRF providers.
type OCR2VRFRelayer interface {
	NewDKGProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (DKGProvider, error)
	NewOCR2VRFProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (OCR2VRFProvider, error)
}

var (
	_ OCR2VRFRelayer  = (*ocr2vrfRelayer)(nil)
	_ DKGProvider     = (*dkgProvider)(nil)
	_ OCR2VRFProvider = (*ocr2vrfProvider)(nil)
)

// Relayer with added DKG and OCR2VRF provider functions.
type ocr2vrfRelayer struct {
	chain       legacyevm.Chain
	lggr        logger.Logger
	ethKeystore keystore.Eth
}

func NewOCR2VRFRelayer(chain legacyevm.Chain, lggr logger.Logger, ethKeystore keystore.Eth) OCR2VRFRelayer {
	return &ocr2vrfRelayer{
		chain:       chain,
		lggr:        lggr,
		ethKeystore: ethKeystore,
	}
}

func (r *ocr2vrfRelayer) NewDKGProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (DKGProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()

	configWatcher, err := newOCR2VRFConfigProvider(ctx, r.lggr, r.chain, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, pargs.TransmitterID, r.ethKeystore, configWatcher, configTransmitterOpts{}, OCR2AggregatorTransmissionContractABI, 0)
	if err != nil {
		return nil, err
	}

	var pluginConfig config.PluginConfig
	err = json.Unmarshal(pargs.PluginConfig, &pluginConfig)
	if err != nil {
		return nil, err
	}

	return &dkgProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
		pluginConfig:        pluginConfig,
	}, nil
}

func (r *ocr2vrfRelayer) NewOCR2VRFProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (OCR2VRFProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()

	configWatcher, err := newOCR2VRFConfigProvider(ctx, r.lggr, r.chain, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, pargs.TransmitterID, r.ethKeystore, configWatcher, configTransmitterOpts{}, OCR2AggregatorTransmissionContractABI, 0)
	if err != nil {
		return nil, err
	}
	return &ocr2vrfProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

type dkgProvider struct {
	*configWatcher
	contractTransmitter ContractTransmitter
	pluginConfig        config.PluginConfig
}

func (c *dkgProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}

func (c *dkgProvider) ChainReader() commontypes.ContractReader {
	return nil
}

func (c *dkgProvider) Codec() commontypes.Codec {
	return nil
}

type ocr2vrfProvider struct {
	*configWatcher
	contractTransmitter ContractTransmitter
}

func (c *ocr2vrfProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}

func (c *ocr2vrfProvider) ChainReader() commontypes.ContractReader {
	return nil
}

func (c *ocr2vrfProvider) Codec() commontypes.Codec {
	return nil
}

func newOCR2VRFConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, rargs commontypes.RelayArgs) (*configWatcher, error) {
	var relayConfig types.RelayConfig
	err := json.Unmarshal(rargs.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(rargs.ContractID) {
		return nil, fmt.Errorf("invalid contract address '%s'", rargs.ContractID)
	}

	contractAddress := common.HexToAddress(rargs.ContractID)
	configPoller, err := NewConfigPoller(
		ctx,
		lggr.With("contractID", rargs.ContractID),
		CPConfig{
			chain.Client(),
			chain.LogPoller(),
			contractAddress,
			// TODO: Does ocr2vrf need to support config contract? DF-19182
			nil,
			OCR2AggregatorLogDecoder,
		},
	)
	if err != nil {
		return nil, err
	}

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().EVM().ChainID().Uint64(),
		ContractAddress: contractAddress,
	}

	return newConfigWatcher(
		lggr,
		contractAddress,
		offchainConfigDigester,
		configPoller,
		chain,
		relayConfig.FromBlock,
		rargs.New,
	), nil
}
