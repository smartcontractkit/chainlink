package evm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txm "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	mercuryconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	reportcodecv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/reportcodec"
	reportcodecv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2/reportcodec"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var _ commontypes.Relayer = &Relayer{} //nolint:staticcheck

type Relayer struct {
	db          *sqlx.DB
	chain       legacyevm.Chain
	lggr        logger.Logger
	ks          CSAETHKeystore
	mercuryPool wsrpc.Pool
	pgCfg       pg.QConfig
	chainReader commontypes.ChainReader
	codec       commontypes.Codec
}

type CSAETHKeystore interface {
	CSA() keystore.CSA
	Eth() keystore.Eth
}

type RelayerOpts struct {
	*sqlx.DB
	pg.QConfig
	CSAETHKeystore
	MercuryPool wsrpc.Pool
}

func (c RelayerOpts) Validate() error {
	var err error
	if c.DB == nil {
		err = errors.Join(err, errors.New("nil DB"))
	}
	if c.QConfig == nil {
		err = errors.Join(err, errors.New("nil QConfig"))
	}
	if c.CSAETHKeystore == nil {
		err = errors.Join(err, errors.New("nil Keystore"))
	}

	if err != nil {
		err = fmt.Errorf("invalid RelayerOpts: %w", err)
	}
	return err
}

func NewRelayer(lggr logger.Logger, chain legacyevm.Chain, opts RelayerOpts) (*Relayer, error) {
	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("cannot create evm relayer: %w", err)
	}
	lggr = lggr.Named("Relayer")
	return &Relayer{
		db:          opts.DB,
		chain:       chain,
		lggr:        lggr,
		ks:          opts.CSAETHKeystore,
		mercuryPool: opts.MercuryPool,
		pgCfg:       opts.QConfig,
	}, nil
}

func (r *Relayer) Name() string {
	return r.lggr.Name()
}

// Start does noop: no subservices started on relay start, but when the first job is started
func (r *Relayer) Start(context.Context) error {
	return nil
}

func (r *Relayer) Close() error {
	return nil
}

// Ready does noop: always ready
func (r *Relayer) Ready() error {
	return r.chain.Ready()
}

func (r *Relayer) HealthReport() (report map[string]error) {
	report = make(map[string]error)
	maps.Copy(report, r.chain.HealthReport())
	return
}

func (r *Relayer) NewMercuryProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.MercuryProvider, error) {
	lggr := r.lggr.Named("MercuryProvider").Named(rargs.ExternalJobID.String())
	relayOpts := types.NewRelayOpts(rargs)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}

	var mercuryConfig mercuryconfig.PluginConfig
	if err = json.Unmarshal(pargs.PluginConfig, &mercuryConfig); err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if relayConfig.FeedID == nil {
		return nil, pkgerrors.New("FeedID must be specified")
	}
	feedID := mercuryutils.FeedID(*relayConfig.FeedID)

	if relayConfig.ChainID.String() != r.chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), r.chain.ID().String())
	}
	cw, err := newConfigProvider(lggr, r.chain, relayOpts)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, pkgerrors.New("EffectiveTransmitterID must be specified")
	}
	privKey, err := r.ks.CSA().Get(relayConfig.EffectiveTransmitterID.String)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get CSA key for mercury connection")
	}

	client, err := r.mercuryPool.Checkout(context.Background(), privKey, mercuryConfig.ServerPubKey, mercuryConfig.ServerURL())
	if err != nil {
		return nil, err
	}

	// FIXME: We actually know the version here since it's in the feed ID, can
	// we use generics to avoid passing three of this?
	// https://smartcontract-it.atlassian.net/browse/MERC-1414
	reportCodecV1 := reportcodecv1.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV1"))
	reportCodecV2 := reportcodecv2.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV2"))
	reportCodecV3 := reportcodecv3.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV3"))

	var transmitterCodec mercury.TransmitterReportDecoder
	switch feedID.Version() {
	case 1:
		transmitterCodec = reportCodecV1
	case 2:
		transmitterCodec = reportCodecV2
	case 3:
		transmitterCodec = reportCodecV3
	default:
		return nil, fmt.Errorf("invalid feed version %d", feedID.Version())
	}
	transmitter := mercury.NewTransmitter(lggr, cw.ContractConfigTracker(), client, privKey.PublicKey, rargs.JobID, *relayConfig.FeedID, r.db, r.pgCfg, transmitterCodec)

	return NewMercuryProvider(cw, r.chainReader, r.codec, NewMercuryChainReader(r.chain.HeadTracker()), transmitter, reportCodecV1, reportCodecV2, reportCodecV3, lggr), nil
}

func (r *Relayer) NewFunctionsProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.FunctionsProvider, error) {
	lggr := r.lggr.Named("FunctionsProvider").Named(rargs.ExternalJobID.String())
	// TODO(FUN-668): Not ready yet (doesn't implement FunctionsEvents() properly)
	return NewFunctionsProvider(r.chain, rargs, pargs, lggr, r.ks.Eth(), functions.FunctionsPlugin)
}

func (r *Relayer) NewConfigProvider(args commontypes.RelayArgs) (commontypes.ConfigProvider, error) {
	lggr := r.lggr.Named("ConfigProvider").Named(args.ExternalJobID.String())
	relayOpts := types.NewRelayOpts(args)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	expectedChainID := relayConfig.ChainID.String()
	if expectedChainID != r.chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), r.chain.ID().String())
	}

	configProvider, err := newConfigProvider(lggr, r.chain, relayOpts)
	if err != nil {
		// Never return (*configProvider)(nil)
		return nil, err
	}
	return configProvider, err
}

func FilterNamesFromRelayArgs(args commontypes.RelayArgs) (filterNames []string, err error) {
	var addr ethkey.EIP55Address
	if addr, err = ethkey.NewEIP55Address(args.ContractID); err != nil {
		return nil, err
	}
	var relayConfig types.RelayConfig
	if err = json.Unmarshal(args.RelayConfig, &relayConfig); err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if relayConfig.FeedID != nil {
		filterNames = []string{mercury.FilterName(addr.Address(), *relayConfig.FeedID)}
	} else {
		filterNames = []string{configPollerFilterName(addr.Address()), transmitterFilterName(addr.Address())}
	}
	return filterNames, err
}

type configWatcher struct {
	services.StateMachine
	lggr             logger.Logger
	contractAddress  common.Address
	contractABI      abi.ABI
	offchainDigester ocrtypes.OffchainConfigDigester
	configPoller     types.ConfigPoller
	chain            legacyevm.Chain
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
	configPoller types.ConfigPoller,
	chain legacyevm.Chain,
	fromBlock uint64,
	runReplay bool,
) *configWatcher {
	replayCtx, replayCancel := context.WithCancel(context.Background())
	return &configWatcher{
		lggr:             lggr.Named("ConfigWatcher").Named(contractAddress.String()),
		contractAddress:  contractAddress,
		contractABI:      contractABI,
		offchainDigester: offchainDigester,
		configPoller:     configPoller,
		chain:            chain,
		runReplay:        runReplay,
		fromBlock:        fromBlock,
		replayCtx:        replayCtx,
		replayCancel:     replayCancel,
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
		c.configPoller.Start()
		return nil
	})
}

func (c *configWatcher) Close() error {
	return c.StopOnce(fmt.Sprintf("configWatcher %x", c.contractAddress), func() error {
		c.replayCancel()
		c.wg.Wait()
		return c.configPoller.Close()
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

func newConfigProvider(lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts) (*configWatcher, error) {
	if !common.IsHexAddress(opts.ContractID) {
		return nil, pkgerrors.Errorf("invalid contractID, expected hex address")
	}

	aggregatorAddress := common.HexToAddress(opts.ContractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		return nil, pkgerrors.Wrap(err, "could not get contract ABI JSON")
	}
	var cp types.ConfigPoller

	relayConfig, err := opts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	if relayConfig.FeedID != nil {
		cp, err = mercury.NewConfigPoller(
			lggr.Named(relayConfig.FeedID.String()),
			chain.LogPoller(),
			aggregatorAddress,
			*relayConfig.FeedID,
			// TODO: Does mercury need to support config contract? DF-19182
		)
	} else {
		cp, err = NewConfigPoller(
			lggr,
			chain.Client(),
			chain.LogPoller(),
			aggregatorAddress,
			relayConfig.ConfigContractAddress,
		)
	}
	if err != nil {
		return nil, err
	}

	var offchainConfigDigester ocrtypes.OffchainConfigDigester
	if relayConfig.FeedID != nil {
		// Mercury
		offchainConfigDigester = mercury.NewOffchainConfigDigester(*relayConfig.FeedID, chain.Config().EVM().ChainID(), aggregatorAddress)
	} else {
		// Non-mercury
		offchainConfigDigester = evmutil.EVMOffchainConfigDigester{
			ChainID:         chain.Config().EVM().ChainID().Uint64(),
			ContractAddress: aggregatorAddress,
		}
	}
	return newConfigWatcher(lggr, aggregatorAddress, contractABI, offchainConfigDigester, cp, chain, relayConfig.FromBlock, opts.New), nil
}

type configTransmitterOpts struct {
	// override the gas limit default provided in the config watcher
	pluginGasLimit *uint32
}

func newContractTransmitter(lggr logger.Logger, rargs commontypes.RelayArgs, transmitterID string, ethKeystore keystore.Eth, configWatcher *configWatcher, opts configTransmitterOpts, reportToEthMeta ReportToEthMetadata) (*contractTransmitter, error) {
	var relayConfig types.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	var fromAddresses []common.Address
	sendingKeys := relayConfig.SendingKeys
	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, pkgerrors.New("EffectiveTransmitterID must be specified")
	}
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterID.String)

	sendingKeysLength := len(sendingKeys)
	if sendingKeysLength == 0 {
		return nil, pkgerrors.New("no sending keys provided")
	}

	// If we are using multiple sending keys, then a forwarder is needed to rotate transmissions.
	// Ensure that this forwarder is not set to a local sending key, and ensure our sending keys are enabled.
	for _, s := range sendingKeys {
		if sendingKeysLength > 1 && s == effectiveTransmitterAddress.String() {
			return nil, pkgerrors.New("the transmitter is a local sending key with transaction forwarding enabled")
		}
		if err := ethKeystore.CheckEnabled(common.HexToAddress(s), configWatcher.chain.Config().EVM().ChainID()); err != nil {
			return nil, pkgerrors.Wrap(err, "one of the sending keys given is not enabled")
		}
		fromAddresses = append(fromAddresses, common.HexToAddress(s))
	}

	scoped := configWatcher.chain.Config()
	strategy := txmgrcommon.NewQueueingTxStrategy(rargs.ExternalJobID, scoped.OCR2().DefaultTransactionQueueDepth(), scoped.Database().DefaultQueryTimeout())

	var checker txm.TransmitCheckerSpec
	if configWatcher.chain.Config().OCR2().SimulateTransactions() {
		checker.CheckerType = txm.TransmitCheckerTypeSimulate
	}

	gasLimit := configWatcher.chain.Config().EVM().GasEstimator().LimitDefault()
	ocr2Limit := configWatcher.chain.Config().EVM().GasEstimator().LimitJobType().OCR2()
	if ocr2Limit != nil {
		gasLimit = *ocr2Limit
	}
	if opts.pluginGasLimit != nil {
		gasLimit = *opts.pluginGasLimit
	}

	transmitter, err := ocrcommon.NewTransmitter(
		configWatcher.chain.TxManager(),
		fromAddresses,
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		checker,
		configWatcher.chain.ID(),
		ethKeystore,
	)

	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create transmitter")
	}

	return NewOCRContractTransmitter(
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		configWatcher.contractABI,
		transmitter,
		configWatcher.chain.LogPoller(),
		lggr,
		reportToEthMeta,
	)
}

func (r *Relayer) NewMedianProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.MedianProvider, error) {
	lggr := r.lggr.Named("MedianProvider").Named(rargs.ExternalJobID.String())
	relayOpts := types.NewRelayOpts(rargs)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	expectedChainID := relayConfig.ChainID.String()
	if expectedChainID != r.chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), r.chain.ID().String())
	}
	if !common.IsHexAddress(relayOpts.ContractID) {
		return nil, fmt.Errorf("invalid contractID %s, expected hex address", relayOpts.ContractID)
	}
	contractID := common.HexToAddress(relayOpts.ContractID)

	configWatcher, err := newConfigProvider(lggr, r.chain, relayOpts)
	if err != nil {
		return nil, err
	}

	reportCodec := evmreportcodec.ReportCodec{}
	contractTransmitter, err := newContractTransmitter(lggr, rargs, pargs.TransmitterID, r.ks.Eth(), configWatcher, configTransmitterOpts{}, nil)
	if err != nil {
		return nil, err
	}

	medianContract, err := newMedianContract(configWatcher.ContractConfigTracker(), configWatcher.contractAddress, configWatcher.chain, rargs.JobID, r.db, lggr)
	if err != nil {
		return nil, err
	}

	medianProvider := medianProvider{
		lggr:                lggr.Named("MedianProvider"),
		configWatcher:       configWatcher,
		reportCodec:         reportCodec,
		contractTransmitter: contractTransmitter,
		medianContract:      medianContract,
	}

	// allow fallback until chain reader is default and median contract is removed, but still log just in case
	var chainReaderService ChainReaderService
	if relayConfig.ChainReader != nil {
		if chainReaderService, err = NewChainReaderService(lggr, r.chain.LogPoller(), r.chain, *relayConfig.ChainReader); err != nil {
			return nil, err
		}

		boundContracts := []commontypes.BoundContract{{Name: "median", Pending: true, Address: contractID.String()}}
		if err = chainReaderService.Bind(context.Background(), boundContracts); err != nil {
			return nil, err
		}
	} else {
		lggr.Info("ChainReader missing from RelayConfig; falling back to internal MedianContract")
	}
	medianProvider.chainReader = chainReaderService

	if relayConfig.Codec != nil {
		medianProvider.codec, err = NewCodec(*relayConfig.Codec)
		if err != nil {
			return nil, err
		}
	} else {
		lggr.Info("Codec missing from RelayConfig; falling back to internal MedianContract")
	}

	return &medianProvider, nil
}

func (r *Relayer) NewAutomationProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.AutomationProvider, error) {
	lggr := r.lggr.Named("AutomationProvider").Named(rargs.ExternalJobID.String())
	ocr2keeperRelayer := NewOCR2KeeperRelayer(r.db, r.chain, lggr.Named("OCR2KeeperRelayer"), r.ks.Eth(), r.pgCfg)

	return ocr2keeperRelayer.NewOCR2KeeperProvider(rargs, pargs)
}

var _ commontypes.MedianProvider = (*medianProvider)(nil)

type medianProvider struct {
	lggr                logger.Logger
	configWatcher       *configWatcher
	contractTransmitter ContractTransmitter
	reportCodec         median.ReportCodec
	medianContract      *medianContract
	chainReader         ChainReaderService
	codec               commontypes.Codec
	ms                  services.MultiStart
}

func (p *medianProvider) Name() string { return p.lggr.Name() }

func (p *medianProvider) Start(ctx context.Context) error {
	srvcs := []services.StartClose{p.configWatcher, p.contractTransmitter, p.medianContract}
	if p.chainReader != nil {
		srvcs = append(srvcs, p.chainReader)
	}

	return p.ms.Start(ctx, srvcs...)
}

func (p *medianProvider) Close() error { return p.ms.Close() }

func (p *medianProvider) Ready() error { return nil }

func (p *medianProvider) HealthReport() map[string]error {
	hp := map[string]error{p.Name(): p.Ready()}
	services.CopyHealth(hp, p.configWatcher.HealthReport())
	services.CopyHealth(hp, p.contractTransmitter.HealthReport())
	services.CopyHealth(hp, p.medianContract.HealthReport())
	return hp
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

func (p *medianProvider) ChainReader() commontypes.ChainReader {
	return p.chainReader
}

func (p *medianProvider) Codec() commontypes.Codec {
	return p.codec
}
