package evm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/sqlx"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
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
	pr    pipeline.Runner
	spec  job.Job
	lggr  logger.Logger
}

// NewOCR2KeeperRelayer is the constructor of ocr2keeperRelayer
func NewOCR2KeeperRelayer(db *sqlx.DB, chain evm.Chain, pr pipeline.Runner, spec job.Job, lggr logger.Logger) OCR2KeeperRelayer {
	return &ocr2keeperRelayer{
		db:    db,
		chain: chain,
		pr:    pr,
		spec:  spec,
		lggr:  lggr,
	}
}

func (r *ocr2keeperRelayer) NewOCR2KeeperProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2KeeperProvider, error) {
	cfgWatcher, err := newOCR2KeeperConfigProvider(r.lggr, r.chain, rargs)
	if err != nil {
		return nil, err
	}

	gasLimit := cfgWatcher.chain.Config().OCR2AutomationGasLimit()
	contractTransmitter, err := newPipelineContractTransmitter(r.lggr, rargs, pargs.TransmitterID, &gasLimit, cfgWatcher, r.spec, r.pr)
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
	contractTransmitter ContractTransmitter
}

func (c *ocr2keeperProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}

func newOCR2KeeperConfigProvider(lggr logger.Logger, chain evm.Chain, rargs relaytypes.RelayArgs) (*configWatcher, error) {
	var relayConfig types.RelayConfig
	err := json.Unmarshal(rargs.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(rargs.ContractID) {
		return nil, fmt.Errorf("invalid contract address '%s'", rargs.ContractID)
	}

	contractAddress := common.HexToAddress(rargs.ContractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get OCR2Aggregator ABI JSON")
	}

	configPoller, err := NewConfigPoller(
		lggr.With("contractID", rargs.ContractID),
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

	return newConfigWatcher(
		lggr,
		contractAddress,
		contractABI,
		offchainConfigDigester,
		configPoller,
		chain,
		relayConfig.FromBlock,
		rargs.New,
	), nil
}
