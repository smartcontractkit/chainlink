package evm

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	txm "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	relaytypes "github.com/smartcontractkit/chainlink/core/services/relay/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ relaytypes.Relayer = &Relayer{}

type Relayer struct {
	db       *sqlx.DB
	chainSet evm.ChainSet
	lggr     logger.Logger
}

func NewRelayer(db *sqlx.DB, chainSet evm.ChainSet, lggr logger.Logger) *Relayer {
	return &Relayer{
		db:       db,
		chainSet: chainSet,
		lggr:     lggr.Named("Relayer"),
	}
}

// Start does noop: no subservices started on relay start, but when the first job is started
func (r *Relayer) Start(context.Context) error {
	return nil
}

// Close does noop: no persistent subservices to close on relay close
func (r *Relayer) Close() error {
	return nil
}

// Ready does noop: always ready
func (r *Relayer) Ready() error {
	return nil
}

// Healthy does noop: always healthy
func (r *Relayer) Healthy() error {
	return nil
}

func (r *Relayer) NewConfigWatcher(args relaytypes.ConfigWatcherArgs) (relaytypes.ConfigWatcher, error) {
	return newConfigWatcher(r.lggr, r.chainSet, args)
}

type configWatcher struct {
	contractAddress  common.Address
	contractABI      abi.ABI
	offchainDigester types.OffchainConfigDigester
	configTracker    *ConfigTracker
	chain            evm.Chain
}

func (c configWatcher) Start(ctx context.Context) error {
	return c.configTracker.Start()
}

func (c configWatcher) Close() error {
	return c.configTracker.Close()
}

func (c configWatcher) Ready() error {
	return nil
}

func (c configWatcher) Healthy() error {
	return nil
}

func (c configWatcher) OffchainConfigDigester() types.OffchainConfigDigester {
	return c.offchainDigester
}

func (c configWatcher) ContractConfigTracker() types.ContractConfigTracker {
	return c.configTracker
}

func newConfigWatcher(lggr logger.Logger, chainSet evm.ChainSet, args relaytypes.ConfigWatcherArgs) (*configWatcher, error) {
	var relayConfig RelayConfig
	err := json.Unmarshal(args.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	chain, err := chainSet.Get(relayConfig.ChainID.ToInt())
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(args.ContractID) {
		return nil, errors.Errorf("invalid contractID, expected hex address")
	}
	contractAddress := common.HexToAddress(args.ContractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get contract ABI JSON")
	}
	configTracker := NewConfigTracker(lggr, contractABI,
		chain.Client(),
		contractAddress,
		chain.Config().ChainType(),
		chain.HeadBroadcaster())

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().ChainID().Uint64(),
		ContractAddress: contractAddress,
	}
	return &configWatcher{
		contractAddress:  contractAddress,
		contractABI:      contractABI,
		configTracker:    configTracker,
		offchainDigester: offchainConfigDigester,
		chain:            chain,
	}, nil
}

func (r *Relayer) NewMedianProvider(args relaytypes.PluginArgs) (relaytypes.MedianProvider, error) {
	configWatcher, err := newConfigWatcher(r.lggr, r.chainSet, args.ConfigWatcherArgs)
	if err != nil {
		return nil, err
	}
	transmitterAddress := common.HexToAddress(args.TransmitterID)
	strategy := txm.NewQueueingTxStrategy(args.ExternalJobID, configWatcher.chain.Config().OCRDefaultTransactionQueueDepth())

	contractTransmitter := NewOCRContractTransmitter(
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		configWatcher.contractABI,
		ocrcommon.NewTransmitter(configWatcher.chain.TxManager(), transmitterAddress, configWatcher.chain.Config().EvmGasLimitDefault(), strategy, txm.TransmitCheckerSpec{}),
		r.lggr,
	)
	medianContract, err := newMedianContract(configWatcher.contractAddress, configWatcher.chain, args.JobID, r.db, r.lggr)
	if err != nil {
		return nil, errors.Wrap(err, "error during median contract setup")
	}
	return &medianProvider{
		configWatcher:       configWatcher,
		reportCodec:         evmreportcodec.ReportCodec{},
		contractTransmitter: contractTransmitter,
		medianContract:      medianContract,
	}, nil
}

type RelayConfig struct {
	ChainID *utils.Big `json:"chainID"`
}

var _ relaytypes.MedianProvider = (*medianProvider)(nil)

type medianProvider struct {
	*configWatcher
	contractTransmitter *ContractTransmitter
	reportCodec         median.ReportCodec
	medianContract      *medianContract
}

func (p medianProvider) ContractTransmitter() types.ContractTransmitter {
	return p.contractTransmitter
}

func (p medianProvider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p medianProvider) MedianContract() median.MedianContract {
	return p.medianContract
}
