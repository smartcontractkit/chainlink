package evm

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/bm"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/channeldefinitions"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/grpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/mercurytransmitter"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	lloconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/llo/config"
	evmllo "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var _ commontypes.LLOProvider = (*lloProvider)(nil)

type LLOTransmitter interface {
	services.Service
	llotypes.Transmitter
}

type FilterRegisterer interface {
	Replay(ctx context.Context, fromBlock int64) error
	RegisterFilter(ctx context.Context, filter logpoller.Filter) error
}

type lloProvider struct {
	services.Service
	eng *services.Engine

	cps []evmllo.ConfigPollerService

	transmitter            LLOTransmitter
	logger                 logger.Logger
	channelDefinitionCache llotypes.ChannelDefinitionCache
	digester               ocrtypes.OffchainConfigDigester
	shouldRetireCache      evmllo.ShouldRetireCacheService
	csaSigner              *coretypes.Ed25519Signer

	lp              FilterRegisterer
	runReplay       bool
	replayFromBlock uint64

	ms services.MultiStart
}

func lloProviderConfiguratorFilterName(addr common.Address, donID uint32) string {
	return logpoller.FilterName("LLOProvider Configurator", addr.String(), fmt.Sprintf("%d", donID))
}

func NewLLOProvider(
	ctx context.Context,
	lggr logger.Logger,
	pargs commontypes.PluginArgs,
	cc evmllo.ConfigCache,
	chain legacyevm.Chain,
	configuratorAddress common.Address,
	cdcFactory channeldefinitions.ChannelDefinitionCacheFactory,
	relayConfig types.RelayConfig,
	relayOpts *types.RelayOpts,
	csaKeystore coretypes.Keystore,
	mercuryCfg MercuryConfig,
	retirementReportCache retirement.RetirementReportCache,
	ds sqlutil.DataSource,
	mercuryPool wsrpc.Pool,
	capabilitiesRegistry coretypes.CapabilitiesRegistry,
) (relaytypes.LLOProvider, error) {
	donID := relayConfig.LLODONID
	lp := chain.LogPoller()
	lggr = logger.Sugared(lggr).With("configMode", relayConfig.LLOConfigMode, "configuratorAddress", configuratorAddress, "donID", donID)
	lggr = logger.Named(lggr, fmt.Sprintf("LLO-%d", donID))

	cps, digester, err := newLLOConfigPollers(ctx, lggr, cc, lp, chain.Config().EVM().ChainID(), configuratorAddress, relayConfig)
	if err != nil {
		return nil, err
	}
	src := evmllo.NewShouldRetireCache(lggr, lp, configuratorAddress, donID)

	csaPub := relayConfig.EffectiveTransmitterID.String
	csaSigner, err := coretypes.NewEd25519Signer(csaPub, csaKeystore.Sign)
	if err != nil {
		return nil, fmt.Errorf("failed to create ed25519 signer: %w", err)
	}

	var lloCfg lloconfig.PluginConfig
	if err = json.Unmarshal(pargs.PluginConfig, &lloCfg); err != nil {
		return nil, pkgerrors.WithStack(err)
	}
	if err = lloCfg.Validate(); err != nil {
		return nil, err
	}

	var transmitter LLOTransmitter
	if lloCfg.BenchmarkMode {
		lggr.Info("Benchmark mode enabled, using dummy transmitter. NOTE: THIS WILL NOT TRANSMIT ANYTHING")
		transmitter = bm.NewTransmitter(lggr, csaPub)
	} else {
		clients := make(map[string]grpc.Client)

		mercuryServers := lloCfg.GetServers()

		for _, server := range mercuryServers {
			var client grpc.Client
			switch mercuryCfg.Transmitter().Protocol() {
			case config.MercuryTransmitterProtocolGRPC:
				client = grpc.NewClient(grpc.ClientOpts{
					Logger: logger.Sugared(lggr).
						Named(fmt.Sprintf("%q", server.URL)).
						With("serverURL", server.URL),
					ClientSigner: csaSigner,
					ServerPubKey: ed25519.PublicKey(server.PubKey),
					ServerURL:    server.URL,
				})
			case config.MercuryTransmitterProtocolWSRPC:
				wsrpcClient, checkoutErr := mercuryPool.Checkout(ctx, csaPub, csaSigner, server.PubKey, server.URL)
				if checkoutErr != nil {
					return nil, checkoutErr
				}
				client = wsrpc.GRPCCompatibilityWrapper{Client: wsrpcClient}
			default:
				return nil, fmt.Errorf("unsupported protocol %q", mercuryCfg.Transmitter().Protocol())
			}
			clients[server.URL] = client
		}

		var mercuryTransmitterOpts *mercurytransmitter.Opts = nil

		if len(mercuryServers) > 0 {
			mercuryTransmitterOpts = &mercurytransmitter.Opts{
				Lggr:                 lggr,
				VerboseLogging:       mercuryCfg.VerboseLogging(),
				Cfg:                  mercuryCfg.Transmitter(),
				Clients:              clients,
				FromAccount:          csaPub,
				DonID:                lloCfg.DonID,
				ORM:                  mercurytransmitter.NewORM(ds, relayConfig.LLODONID),
				CapabilitiesRegistry: capabilitiesRegistry,
			}
		}

		// FIXME: The transmitter instantiation really ought to be moved out of
		// the evm relay into llo package
		// https://smartcontract-it.atlassian.net/browse/MERC-6847
		transmitter, err = llo.NewTransmitter(llo.TransmitterOpts{
			Lggr:                   lggr,
			DonID:                  lloCfg.DonID,
			FromAccount:            csaPub, // NOTE: This may need to change if we support e.g. multiple tranmsmitters, to be a composite of all keys
			VerboseLogging:         mercuryCfg.VerboseLogging(),
			MercuryTransmitterOpts: mercuryTransmitterOpts,
			Subtransmitters:        lloCfg.Transmitters,
			RetirementReportCache:  retirementReportCache,
			CapabilitiesRegistry:   capabilitiesRegistry,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create LLO transmitter: %w", err)
		}
	}

	cdc, err := cdcFactory.NewCache(lloCfg)
	if err != nil {
		return nil, err
	}

	p := &lloProvider{
		nil,
		nil,
		cps,
		transmitter,
		logger.Named(lggr, "LLOProvider"),
		cdc,
		digester,
		src,
		csaSigner,
		lp,
		relayOpts.New,
		relayConfig.FromBlock,
		services.MultiStart{},
	}

	p.Service, p.eng = services.Config{
		Name:  "LLOProvider",
		Start: p.start,
		Close: p.close,
	}.NewServiceEngine(lggr)

	return p, nil
}

func (p *lloProvider) start(ctx context.Context) error {
	// NOTE: Remember that all filters must be registered first for this replay
	// to be effective
	// 1. Replay
	// 2. Start all services
	if p.runReplay && p.replayFromBlock != 0 {
		// Only replay if it's a brand new job.
		p.eng.Go(func(ctx context.Context) {
			p.eng.Infow("starting replay for config", "fromBlock", p.replayFromBlock)
			// #nosec G115
			if err := p.lp.Replay(ctx, int64(p.replayFromBlock)); err != nil {
				p.eng.Errorw("error replaying for config", "err", err)
			} else {
				p.eng.Infow("completed replaying for config", "replayFromBlock", p.replayFromBlock)
			}
		})
	}
	srvs := []services.StartClose{p.transmitter, p.channelDefinitionCache, p.shouldRetireCache, p.csaSigner}
	for _, cp := range p.cps {
		srvs = append(srvs, cp)
	}
	err := p.ms.Start(ctx, srvs...)
	return err
}

func (p *lloProvider) close() error {
	return p.ms.Close()
}

func (p *lloProvider) Ready() error {
	errs := make([]error, len(p.cps))
	for i, cp := range p.cps {
		errs[i] = cp.Ready()
	}
	errs = append(errs, p.transmitter.Ready(), p.channelDefinitionCache.Ready(), p.shouldRetireCache.Ready())
	return errors.Join(errs...)
}

func (p *lloProvider) Name() string {
	return p.logger.Name()
}

func (p *lloProvider) HealthReport() map[string]error {
	report := map[string]error{}
	for _, cp := range p.cps {
		services.CopyHealth(report, cp.HealthReport())
	}
	services.CopyHealth(report, p.transmitter.HealthReport())
	services.CopyHealth(report, p.channelDefinitionCache.HealthReport())
	services.CopyHealth(report, p.shouldRetireCache.HealthReport())
	return report
}

func (p *lloProvider) ContractConfigTrackers() (cps []ocrtypes.ContractConfigTracker) {
	cps = make([]ocrtypes.ContractConfigTracker, len(p.cps))
	for i, cp := range p.cps {
		cps[i] = cp
	}
	return
}

func (p *lloProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.digester
}

func (p *lloProvider) ContractTransmitter() llotypes.Transmitter {
	return p.transmitter
}

func (p *lloProvider) ChannelDefinitionCache() llotypes.ChannelDefinitionCache {
	return p.channelDefinitionCache
}

func (p *lloProvider) ShouldRetireCache() llotypes.ShouldRetireCache {
	return p.shouldRetireCache
}

// wrapper is needed to turn mercury config poller into a service
type mercuryConfigPollerWrapper struct {
	*mercury.ConfigPoller
	services.Service
	eng *services.Engine

	runReplay bool
	fromBlock uint64
}

func newMercuryConfigPollerWrapper(lggr logger.Logger, cp *mercury.ConfigPoller, fromBlock uint64, runReplay bool) *mercuryConfigPollerWrapper {
	w := &mercuryConfigPollerWrapper{cp, nil, nil, runReplay, fromBlock}
	w.Service, w.eng = services.Config{
		Name:  "LLOMercuryConfigWrapper",
		Start: w.start,
		Close: w.close,
	}.NewServiceEngine(lggr)
	return w
}

func (w *mercuryConfigPollerWrapper) Start(ctx context.Context) error {
	return w.Service.Start(ctx)
}

func (w *mercuryConfigPollerWrapper) start(ctx context.Context) error {
	w.ConfigPoller.Start()
	return nil
}

func (w *mercuryConfigPollerWrapper) Close() error {
	return w.Service.Close()
}

func (w *mercuryConfigPollerWrapper) close() error {
	return w.ConfigPoller.Close()
}

func newLLOConfigPollers(ctx context.Context, lggr logger.Logger, cc evmllo.ConfigCache, lp logpoller.LogPoller, chainID *big.Int, configuratorAddress common.Address, relayConfig types.RelayConfig) (cps []evmllo.ConfigPollerService, configDigester ocrtypes.OffchainConfigDigester, err error) {
	donID := relayConfig.LLODONID
	donIDHash := evmllo.DonIDToBytes32(donID)
	switch relayConfig.LLOConfigMode {
	case types.LLOConfigModeMercury:
		// NOTE: This uses the old config digest prefix for compatibility with legacy contracts
		configDigester = mercury.NewOffchainConfigDigester(donIDHash, chainID, configuratorAddress, ocrtypes.ConfigDigestPrefixMercuryV02)
		// Mercury config poller will register its own filter
		mcp, err := mercury.NewConfigPoller(
			ctx,
			lggr,
			lp,
			configuratorAddress,
			evmllo.DonIDToBytes32(donID),
		)
		if err != nil {
			return nil, nil, err
		}
		// don't need to replay in the wrapper since the provider will handle it
		w := newMercuryConfigPollerWrapper(lggr, mcp, relayConfig.FromBlock, false)
		cps = []evmllo.ConfigPollerService{w}
	case types.LLOConfigModeBlueGreen:
		// NOTE: Register filter here because the config poller doesn't do it on its own
		err := lp.RegisterFilter(ctx, logpoller.Filter{Name: lloProviderConfiguratorFilterName(configuratorAddress, donID), EventSigs: []common.Hash{evmllo.ProductionConfigSet, evmllo.StagingConfigSet, evmllo.PromoteStagingConfig}, Topic2: []common.Hash{donIDHash}, Addresses: []common.Address{configuratorAddress}})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to register filter: %w", err)
		}

		configDigester = evmllo.NewOffchainConfigDigester(donIDHash, chainID, configuratorAddress, ocrtypes.ConfigDigestPrefixLLO)
		blueCP := evmllo.NewConfigPoller(
			lggr,
			lp,
			cc,
			configuratorAddress,
			donID,
			evmllo.InstanceTypeBlue,
			relayConfig.FromBlock,
		)
		greenCP := evmllo.NewConfigPoller(
			lggr,
			lp,
			cc,
			configuratorAddress,
			donID,
			evmllo.InstanceTypeGreen,
			relayConfig.FromBlock,
		)
		cps = []evmllo.ConfigPollerService{blueCP, greenCP}
	}
	return cps, configDigester, nil
}
