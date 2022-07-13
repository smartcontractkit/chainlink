package evm

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
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

func (r *Relayer) NewConfigProvider(args relaytypes.RelayArgs) (relaytypes.ConfigProvider, error) {
	configProvider, err := newConfigProvider(r.lggr, r.chainSet, args)
	if err != nil {
		// Never return (*configProvider)(nil)
		return nil, err
	}
	return configProvider, err
}

type configWatcher struct {
	utils.StartStopOnce
	contractAddress  common.Address
	contractABI      abi.ABI
	offchainDigester types.OffchainConfigDigester
	configPoller     *ConfigPoller
	chain            evm.Chain
}

func (c *configWatcher) Start(ctx context.Context) error {
	return nil
}

func (c *configWatcher) Close() error {
	return nil
}

func (c *configWatcher) OffchainConfigDigester() types.OffchainConfigDigester {
	return c.offchainDigester
}

func (c *configWatcher) ContractConfigTracker() types.ContractConfigTracker {
	return c.configPoller
}

func newConfigProvider(lggr logger.Logger, chainSet evm.ChainSet, args relaytypes.RelayArgs) (*configWatcher, error) {
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
	configPoller := NewConfigPoller(lggr,
		chain.LogPoller(),
		contractAddress,
	)

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

func newContractTransmitter(lggr logger.Logger, rargs relaytypes.RelayArgs, transmitterID string, configWatcher *configWatcher) (*ContractTransmitter, error) {
	transmitterAddress := common.HexToAddress(transmitterID)
	strategy := txm.NewQueueingTxStrategy(rargs.ExternalJobID, configWatcher.chain.Config().OCRDefaultTransactionQueueDepth())
	var checker txm.TransmitCheckerSpec
	if configWatcher.chain.Config().OCRSimulateTransactions() {
		checker.CheckerType = txm.TransmitCheckerTypeSimulate
	}
	return NewOCRContractTransmitter(
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		configWatcher.contractABI,
		ocrcommon.NewTransmitter(configWatcher.chain.TxManager(), transmitterAddress, configWatcher.chain.Config().EvmGasLimitDefault(), strategy, txm.TransmitCheckerSpec{}),
		configWatcher.chain.LogPoller(),
		lggr,
	)
}

func (r *Relayer) NewMedianProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.MedianProvider, error) {
	configWatcher, err := newConfigProvider(r.lggr, r.chainSet, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher)
	if err != nil {
		return nil, err
	}
	medianContract, err := newMedianContract(configWatcher.contractAddress, configWatcher.chain, rargs.JobID, r.db, r.lggr)
	if err != nil {
		return nil, err
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

func (p *medianProvider) ContractTransmitter() types.ContractTransmitter {
	return p.contractTransmitter
}

func (p *medianProvider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p *medianProvider) MedianContract() median.MedianContract {
	return p.medianContract
}
