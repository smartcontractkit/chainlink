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

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	txm "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc"
	types "github.com/smartcontractkit/chainlink/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ relaytypes.Relayer = &Relayer{}

type RelayerConfig interface {
	MercuryCredentials(url string) (username, password string, err error)
}

type Relayer struct {
	db       *sqlx.DB
	chainSet evm.ChainSet
	lggr     logger.Logger
	cfg      RelayerConfig
	ks       keystore.Master
}

func NewRelayer(db *sqlx.DB, chainSet evm.ChainSet, lggr logger.Logger, cfg RelayerConfig, ks keystore.Master) *Relayer {
	return &Relayer{
		db:       db,
		chainSet: chainSet,
		lggr:     lggr.Named("Relayer"),
		cfg:      cfg,
		ks:       ks,
	}
}

func (r *Relayer) Name() string {
	return r.lggr.Name()
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

func (r *Relayer) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

// This is a stub, to be added in smartcontractkit/chainlink#8340
func (r *Relayer) NewMercuryProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.MercuryProvider, error) {
	return nil, errors.New("mercury is not supported")
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
	offchainDigester ocrtypes.OffchainConfigDigester
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
	offchainDigester ocrtypes.OffchainConfigDigester,
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

func (c *configWatcher) Name() string {
	return c.lggr.Name()
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

func (c *configWatcher) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
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

func newContractTransmitter(lggr logger.Logger, rargs relaytypes.RelayArgs, transmitterID string, configWatcher *configWatcher, ethKeystore keystore.Eth) (*contractTransmitter, error) {
	var relayConfig types.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	var fromAddresses []common.Address
	sendingKeys := relayConfig.SendingKeys
	if !relayConfig.EffectiveTransmitterAddress.Valid {
		return nil, errors.New("EffectiveTransmitterAddress must be specified")
	}
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterAddress.String)

	sendingKeysLength := len(sendingKeys)
	if sendingKeysLength == 0 {
		return nil, errors.New("no sending keys provided")
	}

	// If we are using multiple sending keys, then a forwarder is needed to rotate transmissions.
	// Ensure that this forwarder is not set to a local sending key.
	for _, s := range sendingKeys {
		if sendingKeysLength > 1 && s == effectiveTransmitterAddress.String() {
			return nil, errors.New("the transmitter is a local sending key with transaction forwarding enabled")
		}
		fromAddresses = append(fromAddresses, common.HexToAddress(s))
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

	transmitter, err := ocrcommon.NewTransmitter(
		configWatcher.chain.TxManager(),
		fromAddresses,
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txm.TransmitCheckerSpec{},
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
	)
}

func newPipelineContractTransmitter(lggr logger.Logger, rargs relaytypes.RelayArgs, transmitterID string, pluginGasLimit *uint32, configWatcher *configWatcher, spec job.Job, pr pipeline.Runner) (*contractTransmitter, error) {
	var relayConfig types.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}

	if !relayConfig.EffectiveTransmitterAddress.Valid {
		return nil, errors.New("EffectiveTransmitterAddress must be specified")
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

	var relayConfig types.RelayConfig
	if err = json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	var contractTransmitter ContractTransmitter
	var reportCodec median.ReportCodec
	if relayConfig.MercuryConfig != nil {
		r.lggr.Debugf("Mercury mode enabled for job %d", rargs.JobID)
		contractTransmitter, reportCodec, err = r.NewMercuryMedianProvider(relayConfig)
	} else {
		r.lggr.Debugf("On-chain mode enabled for job %d", rargs.JobID)
		reportCodec = evmreportcodec.ReportCodec{}
		contractTransmitter, err = newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher, r.ks.Eth())
	}
	if err != nil {
		return nil, err
	}

	medianContract, err := newMedianContract(configWatcher.ContractConfigTracker(), configWatcher.contractAddress, configWatcher.chain, rargs.JobID, r.db, r.lggr, relayConfig.MercuryConfig != nil)
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

func (r *Relayer) NewMercuryMedianProvider(relayConfig types.RelayConfig) (contractTransmitter ContractTransmitter, reportCodec median.ReportCodec, err error) {
	// Override on-chain transmitter with Mercury if the relevant config is set
	reportURL := relayConfig.MercuryConfig.URL
	if reportURL == nil {
		return contractTransmitter, reportCodec, errors.New("Mercury URL must be specified")
	}
	if !relayConfig.EffectiveTransmitterAddress.Valid {
		return contractTransmitter, reportCodec, errors.New("EffectiveTransmitterAddress must be specified")
	}
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterAddress.String)
	var username, password string
	username, password, err = r.cfg.MercuryCredentials(reportURL.String())
	if err != nil {
		return contractTransmitter, reportCodec, errors.Wrapf(err, "failed to get mercury credentials for URL: %s", reportURL.String())
	}
	privKey, err := r.ks.CSA().Get(relayConfig.MercuryConfig.ClientPrivKeyID)
	if err != nil {
		return contractTransmitter, reportCodec, errors.Wrap(err, "failed to get CSA key for mercury connection")
	}
	serverPubKey := relayConfig.MercuryConfig.ServerPubKey
	contractTransmitter = mercury.NewTransmitter(r.lggr, wsrpc.NewClient(privKey, serverPubKey, reportURL.URL()), effectiveTransmitterAddress, reportURL.String(), username, password)
	if relayConfig.MercuryConfig.FeedID == (common.Hash{}) {
		return contractTransmitter, reportCodec, errors.New("FeedID must be specified")
	}
	reportCodec = mercury.ReportCodec{FeedID: relayConfig.MercuryConfig.FeedID}
	return
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

func (p *medianProvider) Healthy() error {
	return multierr.Combine(p.configWatcher.Healthy(), p.contractTransmitter.Healthy())
}

func (p *medianProvider) HealthReport() map[string]error {
	return map[string]error{p.Name(): p.Healthy()}
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
