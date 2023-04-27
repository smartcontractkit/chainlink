package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	txm "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	mercuryconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ relaytypes.Relayer = &Relayer{}

type RelayerConfig interface {
}

type Relayer struct {
	db          *sqlx.DB
	chainSet    evm.ChainSet
	lggr        logger.Logger
	cfg         RelayerConfig
	ks          keystore.Master
	mercuryPool wsrpc.Pool
}

func NewRelayer(db *sqlx.DB, chainSet evm.ChainSet, lggr logger.Logger, cfg RelayerConfig, ks keystore.Master) *Relayer {
	return &Relayer{
		db:          db,
		chainSet:    chainSet,
		lggr:        lggr.Named("Relayer"),
		cfg:         cfg,
		ks:          ks,
		mercuryPool: wsrpc.NewPool(lggr.Named("Mercury.WSRPCPool")),
	}
}

func (r *Relayer) Name() string {
	return r.lggr.Name()
}

// Start does noop: no subservices started on relay start, but when the first job is started
func (r *Relayer) Start(context.Context) error {
	return nil
}

func (r *Relayer) Close() error {
	return r.mercuryPool.Close()
}

// Ready does noop: always ready
func (r *Relayer) Ready() error {
	return r.mercuryPool.Ready()
}

func (r *Relayer) HealthReport() (report map[string]error) {
	report = make(map[string]error)
	maps.Copy(report, r.chainSet.HealthReport())
	maps.Copy(report, r.mercuryPool.HealthReport())
	return
}

func (r *Relayer) NewMercuryProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.MercuryProvider, error) {
	var relayConfig types.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, errors.WithStack(err)
	}

	var mercuryConfig mercuryconfig.PluginConfig
	if err := json.Unmarshal(pargs.PluginConfig, &mercuryConfig); err != nil {
		return nil, errors.WithStack(err)
	}

	if relayConfig.FeedID == nil {
		return nil, errors.New("FeedID must be specified")
	}

	configWatcher, err := newConfigProvider(r.lggr, r.chainSet, rargs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	reportCodec := reportcodec.NewEVMReportCodec(*relayConfig.FeedID, r.lggr.Named("ReportCodec"))

	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, errors.New("EffectiveTransmitterID must be specified")
	}
	privKey, err := r.ks.CSA().Get(relayConfig.EffectiveTransmitterID.String)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CSA key for mercury connection")
	}

	client, err := r.mercuryPool.Checkout(context.Background(), privKey, mercuryConfig.ServerPubKey, mercuryConfig.ServerURL())
	if err != nil {
		return nil, err
	}
	transmitter := mercury.NewTransmitter(r.lggr, configWatcher.ContractConfigTracker(), client, privKey.PublicKey, *relayConfig.FeedID)

	return NewMercuryProvider(configWatcher, transmitter, reportCodec, r.lggr), nil
}

func (r *Relayer) NewConfigProvider(args relaytypes.RelayArgs) (relaytypes.ConfigProvider, error) {
	configProvider, err := newConfigProvider(r.lggr, r.chainSet, args)
	if err != nil {
		// Never return (*configProvider)(nil)
		return nil, err
	}
	return configProvider, err
}

func FilterNamesFromRelayArgs(args relaytypes.RelayArgs) (filterNames []string, err error) {
	var addr ethkey.EIP55Address
	if addr, err = ethkey.NewEIP55Address(args.ContractID); err != nil {
		return nil, err
	}
	var relayConfig types.RelayConfig
	if err = json.Unmarshal(args.RelayConfig, &relayConfig); err != nil {
		return nil, errors.WithStack(err)
	}

	if relayConfig.FeedID != nil {
		filterNames = []string{mercury.FilterName(addr.Address())}
	} else {
		filterNames = []string{configPollerFilterName(addr.Address()), transmitterFilterName(addr.Address())}
	}
	return filterNames, err
}

type ConfigPoller interface {
	ocrtypes.ContractConfigTracker

	Replay(ctx context.Context, fromBlock int64) error
}

type configWatcher struct {
	utils.StartStopOnce
	lggr             logger.Logger
	contractAddress  common.Address
	contractABI      abi.ABI
	offchainDigester ocrtypes.OffchainConfigDigester
	configPoller     ConfigPoller
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
	offchainDigester ocrtypes.OffchainConfigDigester,
	configPoller ConfigPoller,
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

func (c *configWatcher) Name() string {
	return c.lggr.Name()
}

func (c *configWatcher) Start(ctx context.Context) error {
	return c.StartOnce(fmt.Sprintf("configWatcher %x", c.contractAddress), func() error {
		if c.runReplay && c.fromBlock != 0 {
			// Only replay if it's a brand runReplay job.
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.lggr.Infow("starting replay for config", "fromBlock", c.fromBlock)
				if err := c.configPoller.Replay(c.replayCtx, int64(c.fromBlock)); err != nil {
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

func (c *configWatcher) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.StartStopOnce.Healthy()}
}

func (c *configWatcher) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return c.offchainDigester
}

func (c *configWatcher) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return c.configPoller
}

func newConfigProvider(lggr logger.Logger, chainSet evm.ChainSet, args relaytypes.RelayArgs) (*configWatcher, error) {
	var relayConfig types.RelayConfig
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
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get contract ABI JSON")
	}
	var cp ConfigPoller

	if relayConfig.FeedID != nil {
		cp, err = mercury.NewConfigPoller(lggr,
			chain.LogPoller(),
			contractAddress,
			*relayConfig.FeedID,
		)
	} else {
		cp, err = NewConfigPoller(lggr,
			chain.LogPoller(),
			contractAddress,
		)
	}
	if err != nil {
		return nil, err
	}

	var offchainConfigDigester ocrtypes.OffchainConfigDigester
	if relayConfig.FeedID != nil {
		// Mercury
		offchainConfigDigester = mercury.NewOffchainConfigDigester(*relayConfig.FeedID, chain.Config().ChainID().Uint64(), contractAddress)
	} else {
		// Non-mercury
		offchainConfigDigester = evmutil.EVMOffchainConfigDigester{
			ChainID:         chain.Config().ChainID().Uint64(),
			ContractAddress: contractAddress,
		}
	}
	return newConfigWatcher(lggr, contractAddress, contractABI, offchainConfigDigester, cp, chain, relayConfig.FromBlock, args.New), nil
}

func newContractTransmitter(lggr logger.Logger, rargs relaytypes.RelayArgs, transmitterID string, configWatcher *configWatcher, ethKeystore keystore.Eth) (*contractTransmitter, error) {
	var relayConfig types.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	var fromAddresses []common.Address
	sendingKeys := relayConfig.SendingKeys
	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, errors.New("EffectiveTransmitterID must be specified")
	}
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterID.String)

	sendingKeysLength := len(sendingKeys)
	if sendingKeysLength == 0 {
		return nil, errors.New("no sending keys provided")
	}

	// If we are using multiple sending keys, then a forwarder is needed to rotate transmissions.
	// Ensure that this forwarder is not set to a local sending key, and ensure our sending keys are enabled.
	for _, s := range sendingKeys {
		if sendingKeysLength > 1 && s == effectiveTransmitterAddress.String() {
			return nil, errors.New("the transmitter is a local sending key with transaction forwarding enabled")
		}
		if err := ethKeystore.CheckEnabled(common.HexToAddress(s), configWatcher.chain.Config().ChainID()); err != nil {
			return nil, errors.Wrap(err, "one of the sending keys given is not enabled")
		}
		fromAddresses = append(fromAddresses, common.HexToAddress(s))
	}

	scoped := configWatcher.chain.Config()
	strategy := txm.NewQueueingTxStrategy(rargs.ExternalJobID, scoped.OCRDefaultTransactionQueueDepth(), scoped.DatabaseDefaultQueryTimeout())

	var checker txm.EvmTransmitCheckerSpec
	if configWatcher.chain.Config().OCRSimulateTransactions() {
		checker.CheckerType = txm.TransmitCheckerTypeSimulate
	}

	gasLimit := configWatcher.chain.Config().EvmGasLimitDefault()
	if configWatcher.chain.Config().EvmGasLimitOCRJobType() != nil {
		gasLimit = *configWatcher.chain.Config().EvmGasLimitOCRJobType()
	}

	transmitter, err := ocrcommon.NewTransmitter(
		configWatcher.chain.TxManager(),
		fromAddresses,
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txm.EvmTransmitCheckerSpec{},
		configWatcher.chain.ID(),
		ethKeystore,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create transmitter")
	}

	return NewOCRContractTransmitter(
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		configWatcher.contractABI,
		transmitter,
		configWatcher.chain.LogPoller(),
		lggr,
		nil,
	)
}

func newPipelineContractTransmitter(lggr logger.Logger, rargs relaytypes.RelayArgs, transmitterID string, pluginGasLimit *uint32, configWatcher *configWatcher, spec job.Job, pr pipeline.Runner) (*contractTransmitter, error) {
	var relayConfig types.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}

	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, errors.New("EffectiveTransmitterID must be specified")
	}
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterID.String)
	transmitterAddress := common.HexToAddress(transmitterID)
	scoped := configWatcher.chain.Config()
	strategy := txm.NewQueueingTxStrategy(rargs.ExternalJobID, scoped.OCRDefaultTransactionQueueDepth(), scoped.DatabaseDefaultQueryTimeout())

	var checker txm.EvmTransmitCheckerSpec
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
		nil,
	)
}

func (r *Relayer) NewMedianProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.MedianProvider, error) {
	configWatcher, err := newConfigProvider(r.lggr, r.chainSet, rargs)
	if err != nil {
		return nil, err
	}

	var relayConfig types.RelayConfig
	if err = json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	var contractTransmitter ContractTransmitter
	var reportCodec median.ReportCodec

	reportCodec = evmreportcodec.ReportCodec{}
	contractTransmitter, err = newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher, r.ks.Eth())
	if err != nil {
		return nil, err
	}

	medianContract, err := newMedianContract(configWatcher.ContractConfigTracker(), configWatcher.contractAddress, configWatcher.chain, rargs.JobID, r.db, r.lggr)
	if err != nil {
		return nil, err
	}
	return &medianProvider{
		configWatcher:       configWatcher,
		reportCodec:         reportCodec,
		contractTransmitter: contractTransmitter,
		medianContract:      medianContract,
	}, nil
}

var _ relaytypes.MedianProvider = (*medianProvider)(nil)

type medianProvider struct {
	configWatcher       *configWatcher
	contractTransmitter ContractTransmitter
	reportCodec         median.ReportCodec
	medianContract      *medianContract

	ms services.MultiStart
}

func (p *medianProvider) Name() string {
	return "EVM.MedianProvider"
}

func (p *medianProvider) Start(ctx context.Context) error {
	return p.ms.Start(ctx, p.configWatcher, p.contractTransmitter)
}

func (p *medianProvider) Close() error {
	return p.ms.Close()
}

func (p *medianProvider) Ready() error {
	return multierr.Combine(p.configWatcher.Ready(), p.contractTransmitter.Ready())
}

func (p *medianProvider) HealthReport() map[string]error {
	report := p.configWatcher.HealthReport()
	maps.Copy(report, p.contractTransmitter.HealthReport())
	return report
}

func (p *medianProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
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

func (p *medianProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.configWatcher.OffchainConfigDigester()
}

func (p *medianProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.configWatcher.ContractConfigTracker()
}
