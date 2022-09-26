package evm

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var (
	_ OCR2KeeperRelayer  = (*ocr2keeperRelayer)(nil)
	_ OCR2KeeperProvider = (*ocr2keeperProvider)(nil)
)

// OCR2KeeperProviderOpts is the custom options to create a keeper provider
type OCR2KeeperProviderOpts struct {
	RArgs      relaytypes.RelayArgs
	PArgs      relaytypes.PluginArgs
	InstanceID int
}

// OCR2KeeperProvider provides all components needed for a OCR2Keeper plugin.
type OCR2KeeperProvider interface {
	relaytypes.Plugin
}

// OCR2KeeperRelayer contains the relayer and instantiating functions for OCR2Keeper providers.
type OCR2KeeperRelayer interface {
	NewOCR2KeeperProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2KeeperProvider, error)
}

// ocr2keeperRelayer is the relayer with added DKG and OCR2Keeper provider functions.
type ocr2keeperRelayer struct {
	db    *sqlx.DB
	chain evm.Chain
	lggr  logger.Logger
}

// NewOCR2KeeperRelayer is the constructor of ocr2keeperRelayer
func NewOCR2KeeperRelayer(db *sqlx.DB, chain evm.Chain, lggr logger.Logger) OCR2KeeperRelayer {
	return &ocr2keeperRelayer{
		db:    db,
		chain: chain,
		lggr:  lggr,
	}
}

func (r *ocr2keeperRelayer) NewOCR2KeeperProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2KeeperProvider, error) {
	cfgWatcher, err := newOCR2KeeperConfigProvider(r.lggr, r.chain, rargs.ContractID)
	if err != nil {
		return nil, err
	}

	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, cfgWatcher)
	if err != nil {
		return nil, err
	}

	return &ocr2keeperProvider{
		configWatcher:       cfgWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

type ocr2keeperProvider struct {
	*configWatcher
	contractTransmitter *ContractTransmitter
}

func (c *ocr2keeperProvider) ContractTransmitter() types.ContractTransmitter {
	return c.contractTransmitter
}

func newOCR2KeeperConfigProvider(lggr logger.Logger, chain evm.Chain, contractID string) (*configWatcher, error) {
	if !common.IsHexAddress(contractID) {
		return nil, fmt.Errorf("invalid contract address '%s'", contractID)
	}

	contractAddress := common.HexToAddress(contractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get OCR2Aggregator ABI JSON")
	}

	configPoller, err := NewConfigPoller(
		lggr.With("contractID", contractID),
		chain.LogPoller(),
		contractAddress,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create config poller")
	}

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().ChainID().Uint64(),
		ContractAddress: contractAddress,
	}

	return &configWatcher{
		contractAddress:  contractAddress,
		contractABI:      contractABI,
		configPoller:     configPoller,
		offchainDigester: offchainConfigDigester,
		chain:            chain,
	}, nil
}
