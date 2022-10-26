package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"
	"gopkg.in/guregu/null.v4"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	txm "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
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
	lggr             logger.Logger
	contractAddress  common.Address
	contractABI      abi.ABI
	offchainDigester types.OffchainConfigDigester
	configPoller     *ConfigPoller
	chain            evm.Chain
	runReplay        bool
	fromBlock        uint64
	replayCtx        context.Context
	replayCancel     context.CancelFunc
	wg               sync.WaitGroup
}

func newConfigWatcher(lggr logger.Logger,
	contractAddress common.Address,
	contractABI abi.ABI,
	offchainDigester types.OffchainConfigDigester,
	configPoller *ConfigPoller,
	chain evm.Chain,
	fromBlock uint64,
	runReplay bool,
) *configWatcher {
	replayCtx, replayCancel := context.WithCancel(context.Background())
	return &configWatcher{
		StartStopOnce:    utils.StartStopOnce{},
		lggr:             lggr,
		contractAddress:  contractAddress,
		contractABI:      contractABI,
		offchainDigester: offchainDigester,
		configPoller:     configPoller,
		chain:            chain,
		runReplay:        runReplay,
		fromBlock:        fromBlock,
		replayCtx:        replayCtx,
		replayCancel:     replayCancel,
		wg:               sync.WaitGroup{},
	}

}

func (c *configWatcher) Start(ctx context.Context) error {
	return c.StartOnce(fmt.Sprintf("configWatcher %x", c.contractAddress), func() error {
		if c.runReplay && c.fromBlock != 0 {
			// Only replay if its a brand runReplay job.
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.lggr.Infow("starting replay for config", "fromBlock", c.fromBlock)
				if err := c.configPoller.destChainLogPoller.Replay(c.replayCtx, int64(c.fromBlock)); err != nil {
					c.lggr.Errorf("error replaying for config", "err", err)
				} else {
					c.lggr.Infow("completed replaying for config", "fromBlock", c.fromBlock)
				}
			}()
		}
		return nil
	})
}

func (c *configWatcher) Close() error {
	return c.StopOnce(fmt.Sprintf("configWatcher %x", c.contractAddress), func() error {
		c.replayCancel()
		c.wg.Wait()
		return nil
	})
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
	configPoller, err := NewConfigPoller(lggr,
		chain.LogPoller(),
		contractAddress,
	)
	if err != nil {
		return nil, err
	}

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().ChainID().Uint64(),
		ContractAddress: contractAddress,
	}
	return newConfigWatcher(lggr, contractAddress, contractABI, offchainConfigDigester, configPoller, chain, relayConfig.FromBlock, args.New), nil
}

func newContractTransmitter(lggr logger.Logger, rargs relaytypes.RelayArgs, transmitterID string, configWatcher *configWatcher) (*ContractTransmitter, error) {
	var relayConfig RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	var fromAddresses []common.Address
	sendingKeys := relayConfig.SendingKeys
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterAddress.String)
	useForwarders := configWatcher.chain.Config().EvmUseForwarders()

	if useForwarders {
		// If using the forwarder, ensure sending keys are provided.
		if len(sendingKeys) == 0 {
			return nil, errors.New("no sending keys found in job spec with forwarder enabled")
		}

		// The sending keys provided are used as the from addresses.
		for _, s := range sendingKeys {
			// Ensure the transmitter is not contained in the sending keys slice.
			if s == effectiveTransmitterAddress.String() {
				return nil, errors.New("the transmitter is a local sending key with transaction forwarding enabled")
			}
			fromAddresses = append(fromAddresses, common.HexToAddress(s))
		}
	} else {
		// Ensure the transmitter is contained in the sending keys slice.
		var transmitterFoundLocally bool
		for _, s := range sendingKeys {
			if s == effectiveTransmitterAddress.String() {
				transmitterFoundLocally = true
				break
			}
		}
		if !transmitterFoundLocally {
			return nil, errors.New("the transmitter was not found in the list of sending keys, perhaps EvmUseForwarders needs to be enabled")
		}

		// If not using the forwarder, the effectiveTransmitterAddress (TransmitterID) is used as the from address.
		fromAddresses = append(fromAddresses, effectiveTransmitterAddress)
	}

	scoped := configWatcher.chain.Config()
	strategy := txm.NewQueueingTxStrategy(rargs.ExternalJobID, scoped.OCRDefaultTransactionQueueDepth(), scoped.DatabaseDefaultQueryTimeout())

	var checker txm.TransmitCheckerSpec
	if configWatcher.chain.Config().OCRSimulateTransactions() {
		checker.CheckerType = txm.TransmitCheckerTypeSimulate
	}

	gasLimit := configWatcher.chain.Config().EvmGasLimitDefault()
	if configWatcher.chain.Config().EvmGasLimitOCRJobType() != nil {
		gasLimit = *configWatcher.chain.Config().EvmGasLimitOCRJobType()
	}

	return NewOCRContractTransmitter(
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		configWatcher.contractABI,
		ocrcommon.NewTransmitter(configWatcher.chain.TxManager(), fromAddresses, gasLimit, effectiveTransmitterAddress, strategy, txm.TransmitCheckerSpec{}),
		configWatcher.chain.LogPoller(),
		lggr,
	)
}

func newPipelineContractTransmitter(lggr logger.Logger, rargs relaytypes.RelayArgs, transmitterID string, pluginGasLimit *uint32, configWatcher *configWatcher, spec job.Job, pr pipeline.Runner) (*ContractTransmitter, error) {
	var relayConfig RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}

	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterAddress.String)
	transmitterAddress := common.HexToAddress(transmitterID)
	scoped := configWatcher.chain.Config()
	strategy := txm.NewQueueingTxStrategy(rargs.ExternalJobID, scoped.OCRDefaultTransactionQueueDepth(), scoped.DatabaseDefaultQueryTimeout())

	var checker txm.TransmitCheckerSpec
	if configWatcher.chain.Config().OCRSimulateTransactions() {
		checker.CheckerType = txm.TransmitCheckerTypeSimulate
	}

	gasLimit := configWatcher.chain.Config().EvmGasLimitDefault()
	if configWatcher.chain.Config().EvmGasLimitOCRJobType() != nil {
		gasLimit = *configWatcher.chain.Config().EvmGasLimitOCRJobType()
	}
	if pluginGasLimit != nil {
		gasLimit = *pluginGasLimit
	}

	return NewOCRContractTransmitter(
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		configWatcher.contractABI,
		ocrcommon.NewPipelineTransmitter(
			lggr,
			transmitterAddress,
			gasLimit,
			effectiveTransmitterAddress,
			strategy,
			checker,
			pr,
			spec,
			configWatcher.chain.ID().String(),
		),
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
	ChainID                     *utils.Big     `json:"chainID"`
	FromBlock                   uint64         `json:"fromBlock"`
	EffectiveTransmitterAddress null.String    `json:"effectiveTransmitterAddress"`
	SendingKeys                 pq.StringArray `json:"sendingKeys"`
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

func (p *medianProvider) OnchainConfigCodec() median.OnchainConfigCodec {
	return median.StandardOnchainConfigCodec{}
}
