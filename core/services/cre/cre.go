package cre

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/jonboulle/clockwork"
	"google.golang.org/grpc/credentials"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/billing"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	nodeauthjwt "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/jwt"
	commonsrv "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/storage"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"
	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr/capregconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2pmain "github.com/smartcontractkit/chainlink/v2/core/services/p2p"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	p2pwrapper "github.com/smartcontractkit/chainlink/v2/core/services/p2p/wrapper"
	registrysyncerV1 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	registrysyncerV2 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	artifactsV1 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"
	artifactsV2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	workflowstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	syncerV1 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer"
	syncerV2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	wftypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"

	linkingclient "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"
)

// Keystore is the minimal interface needed from keystore for CRE
type Keystore interface {
	CSA() keystore.CSA
	Eth() keystore.Eth
	P2P() keystore.P2P
	Workflow() keystore.Workflow
}

// Opts are the options for the CRE services that are exposed by the application
type Opts struct {
	CapabilitiesRegistry    *capabilities.Registry
	CapabilitiesDispatcher  remotetypes.Dispatcher
	CapabilitiesPeerWrapper p2ptypes.PeerWrapper

	FetcherFunc      wftypes.FetcherFunc
	FetcherFactoryFn compute.FetcherFactory

	BillingClient metering.BillingClient
	LinkingClient linkingclient.LinkingServiceClient

	StorageClient storage.WorkflowClient

	DonTimeStore  *dontime.Store
	LimitsFactory limits.Factory

	UseLocalTimeProvider bool

	JWTGenerator nodeauthjwt.JWTGenerator // optional override, created automatically if nil
}

// Services contains all CRE-related services
type Services struct {
	commonsrv.Service
	eng *commonsrv.Engine

	BillingClient metering.BillingClient

	WorkflowRateLimiter *ratelimiter.RateLimiter

	WorkflowLimits limits.ResourceLimiter[int]

	GatewayConnectorWrapper *gatewayconnector.ServiceWrapper

	GetPeerID func() (p2ptypes.PeerID, error)

	WorkflowRegistrySyncer syncerV2.WorkflowRegistrySyncer

	OrgResolver orgresolver.OrgResolver

	OCRConfigService capregconfig.OCRConfigService
}

// newSubservices initializes and returns all CRE child services
func (s *Services) newSubservices(
	ctx context.Context,
	lggr logger.Logger,
	ds sqlutil.DataSource,
	keyStore Keystore,
	cfg Config,
	relayerChainInterops RelayerChainInterops,
	singletonPeerWrapper *ocrcommon.SingletonPeerWrapper,
	opts Opts,
) ([]commonsrv.Service, error) {
	var srvs []commonsrv.Service

	capCfg := cfg.Capabilities()

	workflowRateLimiter, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
		GlobalRPS:      capCfg.RateLimit().GlobalRPS(),
		GlobalBurst:    capCfg.RateLimit().GlobalBurst(),
		PerSenderRPS:   capCfg.RateLimit().PerSenderRPS(),
		PerSenderBurst: capCfg.RateLimit().PerSenderBurst(),
	})
	if err != nil {
		return nil, fmt.Errorf("could not instantiate workflow rate limiter: %w", err)
	}
	s.WorkflowRateLimiter = workflowRateLimiter
	wCfg := cfg.Workflows()
	if len(wCfg.Limits().PerOwnerOverrides()) > 0 {
		lggr.Debugw("loaded per owner overrides", "overrides", wCfg.Limits().PerOwnerOverrides())
	}

	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{
		Global:            wCfg.Limits().Global(),
		PerOwner:          wCfg.Limits().PerOwner(),
		PerOwnerOverrides: wCfg.Limits().PerOwnerOverrides(),
	}, opts.LimitsFactory)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate workflow syncer limiter: %w", err)
	}
	s.WorkflowLimits = workflowLimits
	srvs = append(srvs, closerService{name: "WorkflowExecutionLimiter", Closer: workflowLimits})

	if capCfg.GatewayConnector().DonID() != "" {
		lggr.Debugw("Creating GatewayConnector wrapper", "donID", capCfg.GatewayConnector().DonID())
		gatewayConnectorWrapper, err := newGatewayConnectorWrapper(capCfg, keyStore, lggr)
		if err != nil {
			return nil, fmt.Errorf("could not create gateway connector wrapper: %w", err)
		}
		s.GatewayConnectorWrapper = gatewayConnectorWrapper
		srvs = append(srvs, gatewayConnectorWrapper)
	}

	if cfg.CRE().Linking().URL() != "" {
		lggr.Debugw("Creating OrgResolver")
		orgResolver, err := newOrgResolver(ctx, cfg, capCfg, keyStore, opts, lggr)
		if err != nil {
			return nil, fmt.Errorf("could not create org resolver: %w", err)
		}
		s.OrgResolver = orgResolver
		srvs = append(srvs, orgResolver)
	} else {
		lggr.Warn("Skipping orgResolver, no linking service configured")
	}

	dispatcherWrapper, err := newDispatcherWrapper(cfg, opts, keyStore, ds, singletonPeerWrapper, lggr)
	if err != nil {
		return nil, fmt.Errorf("could not create dispatcher: %w", err)
	}
	s.GetPeerID = dispatcherWrapper.GetPeerID
	srvs = append(srvs, dispatcherWrapper)

	if capCfg.ExternalRegistry().Address() == "" {
		lggr.Warn("Skipping capabilities and workflow registry syncer, none configured")
		return srvs, nil
	}

	registrySyncerServices, donNotifier, err := newRegistrySyncer(
		lggr,
		cfg,
		relayerChainInterops,
		ds,
		opts,
		dispatcherWrapper,
	)
	if err != nil {
		return nil, err
	}
	srvs = append(srvs, registrySyncerServices...)

	if capCfg.WorkflowRegistry().Address() == "" {
		lggr.Warn("Skipping capabilities and workflow registry syncer, none configured")
		return srvs, nil
	}

	wfSyncer, billingClient, wfSyncerSrvcs, err := newWorkflowRegistrySyncer(
		ctx,
		cfg,
		keyStore,
		relayerChainInterops,
		opts,
		lggr,
		ds,
		opts.DonTimeStore,
		workflowRateLimiter,
		workflowLimits,
		donNotifier,
		opts.LimitsFactory,
		s.OrgResolver,
		s.GatewayConnectorWrapper,
	)
	if err != nil {
		return nil, err
	}
	s.BillingClient = billingClient
	s.WorkflowRegistrySyncer = wfSyncer
	srvs = append(srvs, wfSyncerSrvcs...)

	return srvs, nil
}

// Config is the minimal interface needed from GeneralConfig for CRE
type Config interface {
	Billing() config.Billing
	Capabilities() config.Capabilities
	Workflows() config.Workflows
	CRE() config.CRE
	P2P() config.P2P
	Sharding() config.Sharding
}

// RelayerChainInterops is the minimal interface needed for relayer chain interops
type RelayerChainInterops interface {
	Get(commontypes.RelayID) (loop.Relayer, error)
}

// newGatewayConnectorWrapper creates a new GatewayConnector service wrapper if configured
func newGatewayConnectorWrapper(
	capCfg config.Capabilities,
	keyStore Keystore,
	lggr logger.Logger,
) (*gatewayconnector.ServiceWrapper, error) {
	chainID, ok := new(big.Int).SetString(capCfg.GatewayConnector().ChainIDForNodeKey(), 0)
	if !ok {
		return nil, fmt.Errorf("failed to parse gateway connector chain ID as integer: %s", capCfg.GatewayConnector().ChainIDForNodeKey())
	}

	wrapper := gatewayconnector.NewGatewayConnectorServiceWrapper(
		capCfg.GatewayConnector(),
		keys.NewStore(keystore.NewEthSigner(keyStore.Eth(), chainID)),
		keyStore.Eth(),
		chainID,
		clockwork.NewRealClock(),
		lggr)

	return wrapper, nil
}

// dispatcherWrapper is a service that encapsulates the dispatcher and its peer dependencies.
// It manages the lifecycle of the external peer wrapper, shared peer, and dispatcher as subservices.
type dispatcherWrapper struct {
	commonsrv.Service
	eng *commonsrv.Engine

	dispatcher          remotetypes.Dispatcher
	externalPeerWrapper p2ptypes.PeerWrapper
	don2DonSharedPeer   p2ptypes.SharedPeer
}

// GetPeerID returns the peer ID from either the shared peer or external peer wrapper
func (w *dispatcherWrapper) GetPeerID() (p2ptypes.PeerID, error) {
	if w.don2DonSharedPeer != nil {
		return w.don2DonSharedPeer.ID(), nil
	}
	if w.externalPeerWrapper != nil {
		p := w.externalPeerWrapper.GetPeer()
		if p == nil {
			return p2ptypes.PeerID{}, errors.New("could not get peer from externalPeerWrapper")
		}
		return p.ID(), nil
	}
	return p2ptypes.PeerID{}, errors.New("could not get peer from any source")
}

func newRegistrySyncerV1(
	lggr logger.Logger,
	getPeerID func() (p2ptypes.PeerID, error),
	relayer loop.Relayer,
	registryAddress string,
	ds sqlutil.DataSource,
	externalPeerWrapper p2ptypes.PeerWrapper,
	don2donSharedPeer p2ptypes.SharedPeer,
	streamConfig config.StreamConfig,
	dispatcher remotetypes.Dispatcher,
	capabilitiesRegistry *capabilities.Registry,
	donNotifier capabilities.DonNotifier,
) ([]commonsrv.Service, error) {
	wfLauncher, err := capabilities.NewLauncher(
		lggr,
		externalPeerWrapper,
		don2donSharedPeer,
		streamConfig,
		dispatcher,
		capabilitiesRegistry,
		donNotifier,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create workflow launcher: %w", err)
	}

	registrySyncer, err := registrysyncerV1.New(
		lggr,
		getPeerID,
		relayer,
		registryAddress,
		registrysyncerV1.NewORM(ds, lggr),
	)
	if err != nil {
		return nil, fmt.Errorf("could not configure syncer: %w", err)
	}

	registrySyncer.AddListener(wfLauncher)
	return []commonsrv.Service{wfLauncher, registrySyncer}, nil
}

func newRegistrySyncerV2(
	lggr logger.Logger,
	getPeerID func() (p2ptypes.PeerID, error),
	relayer loop.Relayer,
	registryAddress string,
	ds sqlutil.DataSource,
	rid commontypes.RelayID,
	don2donSharedPeer p2ptypes.SharedPeer,
	externalPeerWrapper p2ptypes.PeerWrapper,
	streamConfig config.StreamConfig,
	dispatcher remotetypes.Dispatcher,
	capabilitiesRegistry *capabilities.Registry,
	donNotifier capabilities.DonNotifier,
) ([]commonsrv.Service, error) {
	wfLauncher, err := capabilities.NewLauncher(
		lggr,
		externalPeerWrapper,
		don2donSharedPeer,
		streamConfig,
		dispatcher,
		capabilitiesRegistry,
		donNotifier,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create workflow launcher: %w", err)
	}

	registrySyncer, err := registrysyncerV2.New(
		lggr,
		getPeerID,
		relayer,
		registryAddress,
		registrysyncerV1.NewORM(ds, lggr),
	)
	if err != nil {
		return nil, fmt.Errorf("could not configure syncer: %w", err)
	}

	registrySyncer.AddListener(wfLauncher)

	registryChainID, parseErr := strconv.ParseUint(rid.ChainID, 10, 64)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse registry chain ID for OCRConfigService: %w", parseErr)
	}
	ocrConfigService := capregconfig.NewOCRConfigService(
		lggr,
		func() ragetypes.PeerID { return don2donSharedPeer.ID() },
		registryChainID,
		registryAddress,
	)

	registrySyncer.AddListener(ocrConfigService)
	return []commonsrv.Service{wfLauncher, registrySyncer, ocrConfigService}, nil
}

// newRegistrySyncer creates a registry syncer based on the external registry version
func newRegistrySyncer(
	lggr logger.Logger,
	cfg Config,
	relayerChainInterops RelayerChainInterops,
	ds sqlutil.DataSource,
	opts Opts,
	dispatcherWrapper *dispatcherWrapper,
) ([]commonsrv.Service, capabilities.DonNotifyWaitSubscriber, error) {
	var srvcs []commonsrv.Service

	capCfg := cfg.Capabilities()

	rid := capCfg.ExternalRegistry().RelayID()
	registryAddress := capCfg.ExternalRegistry().Address()
	relayer, err := relayerChainInterops.Get(rid)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch relayer %s configured for capabilities registry: %w", rid, err)
	}

	var streamConfig config.StreamConfig
	if capCfg.SharedPeering().Enabled() {
		streamConfig = capCfg.SharedPeering().StreamConfig()
	}

	donNotifier := capabilities.NewDonNotifier()

	externalRegistryVersion, err := semver.NewVersion(capCfg.ExternalRegistry().ContractVersion())
	if err != nil {
		return nil, nil, err
	}

	switch externalRegistryVersion.Major() {
	case 1:
		srvs, err := newRegistrySyncerV1(
			lggr,
			dispatcherWrapper.GetPeerID,
			relayer,
			registryAddress,
			ds,
			dispatcherWrapper.externalPeerWrapper,
			dispatcherWrapper.don2DonSharedPeer,
			streamConfig,
			dispatcherWrapper.dispatcher,
			opts.CapabilitiesRegistry,
			donNotifier,
		)
		if err != nil {
			return nil, nil, err
		}
		srvcs = append(srvcs, srvs...)
		return srvcs, donNotifier, nil
	case 2:
		srvs, err := newRegistrySyncerV2(
			lggr,
			dispatcherWrapper.GetPeerID,
			relayer,
			registryAddress,
			ds,
			rid,
			dispatcherWrapper.don2DonSharedPeer,
			dispatcherWrapper.externalPeerWrapper,
			streamConfig,
			dispatcherWrapper.dispatcher,
			opts.CapabilitiesRegistry,
			donNotifier,
		)
		if err != nil {
			return nil, nil, err
		}
		srvcs = append(srvcs, srvs...)
		return srvcs, donNotifier, nil
	}

	return nil, nil, fmt.Errorf("could not configure capability registry syncer with version: %d", externalRegistryVersion.Major())
}

func (w *dispatcherWrapper) newSubservices(
	lggr logger.Logger,
	cfg Config,
	opts Opts,
	keyStore Keystore,
	ds sqlutil.DataSource,
	singletonPeerWrapper *ocrcommon.SingletonPeerWrapper,
) ([]commonsrv.Service, error) {
	capCfg := cfg.Capabilities()

	if !capCfg.Peering().Enabled() && !capCfg.SharedPeering().Enabled() {
		return nil, nil
	}

	if opts.CapabilitiesDispatcher != nil {
		w.dispatcher = opts.CapabilitiesDispatcher
		w.externalPeerWrapper = opts.CapabilitiesPeerWrapper
		return nil, nil
	}

	var subs []commonsrv.Service
	var signer p2ptypes.Signer
	if capCfg.Peering().Enabled() {
		w.externalPeerWrapper = p2pwrapper.NewExternalPeerWrapper(keyStore.P2P(), capCfg.Peering(), ds, lggr)
		subs = append(subs, w.externalPeerWrapper)

		signer = p2pmain.NewSigner(keyStore.P2P(), capCfg.Peering().PeerID())
	}

	if capCfg.SharedPeering().Enabled() {
		if !cfg.P2P().Enabled() {
			return nil, errors.New("top-level P2P must be enabled in order to use SharedPeering")
		}
		if singletonPeerWrapper == nil {
			return nil, errors.New("singleton peer wrapper is required for shared peering (are OCR and P2P enabled?)")
		}
		bootstrappers := capCfg.SharedPeering().Bootstrappers()
		if len(bootstrappers) == 0 {
			bootstrappers = cfg.P2P().V2().DefaultBootstrappers()
		}
		w.don2DonSharedPeer = p2pmain.NewDon2DonSharedPeer(singletonPeerWrapper, bootstrappers, lggr)
		subs = append(subs, w.don2DonSharedPeer)

		signer = p2pmain.NewSigner(keyStore.P2P(), cfg.P2P().PeerID())
	}

	remoteDispatcher, err := remote.NewDispatcher(capCfg.Dispatcher(), w.externalPeerWrapper, w.don2DonSharedPeer, signer, opts.CapabilitiesRegistry, lggr)
	if err != nil {
		return nil, fmt.Errorf("could not create dispatcher: %w", err)
	}
	w.dispatcher = remoteDispatcher
	subs = append(subs, remoteDispatcher)
	return subs, nil
}

// newDispatcherWrapper creates a new dispatcherWrapper service with peer wrappers if peering is enabled
func newDispatcherWrapper(
	cfg Config,
	opts Opts,
	keyStore Keystore,
	ds sqlutil.DataSource,
	singletonPeerWrapper *ocrcommon.SingletonPeerWrapper,
	lggr logger.Logger,
) (*dispatcherWrapper, error) {
	w := &dispatcherWrapper{}

	var initErr error
	w.Service, w.eng = commonsrv.Config{
		Name: "DispatcherWrapper",
		NewSubServices: func(lggr logger.Logger) []commonsrv.Service {
			subs, err := w.newSubservices(lggr, cfg, opts, keyStore, ds, singletonPeerWrapper)
			if err != nil {
				initErr = err
				return nil
			}
			return subs
		},
	}.NewServiceEngine(lggr)

	if initErr != nil {
		return nil, initErr
	}

	return w, nil
}

// newOrgResolver creates a new OrgResolver if configured
func newOrgResolver(
	ctx context.Context,
	cfg Config,
	capCfg config.Capabilities,
	keyStore Keystore,
	opts Opts,
	lggr logger.Logger,
) (orgresolver.OrgResolver, error) {
	var wrChainDetails chainselectors.ChainDetails
	if capCfg.WorkflowRegistry().Address() != "" {
		var err error
		wrChainDetails, err = chainselectors.GetChainDetailsByChainIDAndFamily(
			capCfg.WorkflowRegistry().ChainID(),
			capCfg.WorkflowRegistry().NetworkID(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get workflow registry chain details by chain ID and network ID: %w", err)
		}
	}

	orgResolverConfig := orgresolver.Config{
		URL:                           cfg.CRE().Linking().URL(),
		TLSEnabled:                    cfg.CRE().Linking().TLSEnabled(),
		WorkflowRegistryAddress:       capCfg.WorkflowRegistry().Address(),
		WorkflowRegistryChainSelector: wrChainDetails.ChainSelector,
	}

	jwtGenerator, err := newJWTGenerator(ctx, keyStore, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT generator for org resolver: %w", err)
	}
	orgResolverConfig.JWTGenerator = jwtGenerator

	var resolver orgresolver.OrgResolver
	if opts.LinkingClient != nil {
		resolver, err = orgresolver.NewOrgResolverWithClient(orgResolverConfig, opts.LinkingClient, lggr)
	} else {
		resolver, err = orgresolver.NewOrgResolver(orgResolverConfig, lggr)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create org resolver: %w", err)
	}

	return resolver, nil
}

func newJWTGenerator(ctx context.Context, keyStore Keystore, opts Opts) (nodeauthjwt.JWTGenerator, error) {
	if opts.JWTGenerator != nil {
		return opts.JWTGenerator, nil
	}
	csaKeystore := &keystore.CSASigner{CSA: keyStore.CSA()}
	signer, csaPubKey, err := keystore.BuildNodeAuth(ctx, csaKeystore)
	if err != nil {
		return nil, fmt.Errorf("failed to build node auth: %w", err)
	}
	return nodeauthjwt.NewNodeJWTGenerator(signer, csaPubKey), nil
}

func newBillingClient(ctx context.Context, lggr logger.Logger, cfg Config, keyStore Keystore, opts Opts) (metering.BillingClient, error) {
	if opts.BillingClient != nil {
		return opts.BillingClient, nil
	}

	if cfg.Billing().URL() == "" {
		return nil, nil
	}

	jwtGenerator, err := newJWTGenerator(ctx, keyStore, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT generator for billing client: %w", err)
	}

	workflowOpts := []billing.WorkflowClientOpt{
		billing.WithJWTGenerator(jwtGenerator),
	}
	if cfg.Billing().TLSEnabled() {
		workflowOpts = append(workflowOpts, billing.WithWorkflowTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}
	return billing.NewWorkflowClient(lggr, cfg.Billing().URL(), workflowOpts...)
}

func newShardOrchestratorClient(ctx context.Context, cfg Config, lggr logger.Logger) (*shardorchestrator.Client, error) {
	shardID := cfg.Sharding().ShardIndex()
	if shardID == 0 {
		return nil, nil
	}

	address := cfg.Sharding().ShardOrchestratorAddress()
	if address == nil {
		return nil, fmt.Errorf("shard %d requires ShardOrchestratorAddress configuration", shardID)
	}

	client, err := shardorchestrator.NewClient(ctx, address.String(), lggr)
	if err != nil {
		return nil, fmt.Errorf("failed to create ShardOrchestrator gRPC client: %w", err)
	}

	lggr.Infow("ShardOrchestrator gRPC client created", "shardID", shardID, "serverAddress", address)
	return client, nil
}

func newContractReaderFactory(capCfg config.Capabilities, relayerChainInterops RelayerChainInterops) (func(ctx context.Context, bytes []byte) (commontypes.ContractReader, error), error) {
	wfRegRid := capCfg.WorkflowRegistry().RelayID()
	wfRegRelayer, err := relayerChainInterops.Get(wfRegRid)
	if err != nil {
		return nil, fmt.Errorf("could not fetch relayer %s configured for workflow registry: %w", wfRegRid, err)
	}

	return func(ctx context.Context, bytes []byte) (commontypes.ContractReader, error) {
		return wfRegRelayer.NewContractReader(ctx, bytes)
	}, nil
}

func chainSelector(chainID, networkID string) (string, error) {
	wrChainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(
		chainID,
		networkID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get workflow registry chain details by chain ID and network ID: %w", err)
	}

	return strconv.FormatUint(wrChainDetails.ChainSelector, 10), err
}

func newFetcherFuncV1(lggr logger.Logger, optsFetcherFunc wftypes.FetcherFunc, gatewayConnectorWrapper *gatewayconnector.ServiceWrapper) (wftypes.FetcherFunc, []commonsrv.Service, error) {
	if optsFetcherFunc != nil {
		return optsFetcherFunc, nil, nil
	}

	if gatewayConnectorWrapper == nil {
		return nil, nil, errors.New("unable to create workflow registry syncer without gateway connector")
	}
	f := syncerV1.NewFetcherService(lggr, gatewayConnectorWrapper)
	return f.Fetch, []commonsrv.Service{f}, nil
}

func newWorkflowRegistrySyncerV1(
	ctx context.Context,
	capCfg config.Capabilities,
	relayerChainInterops RelayerChainInterops,
	billingClient metering.BillingClient,
	opts Opts,
	lggr logger.Logger,
	ds sqlutil.DataSource,
	dontimeStore *dontime.Store,
	workflowRateLimiter *ratelimiter.RateLimiter,
	workflowLimits limits.ResourceLimiter[int],
	workflowDonNotifier capabilities.DonNotifyWaitSubscriber,
	lf limits.Factory,
	gatewayConnectorWrapper *gatewayconnector.ServiceWrapper,
	keyStore Keystore,
) ([]commonsrv.Service, error) {
	var srvcs []commonsrv.Service

	fetcherFunc, srvs, err := newFetcherFuncV1(lggr, opts.FetcherFunc, gatewayConnectorWrapper)
	if err != nil {
		return nil, err
	}
	srvcs = append(srvcs, srvs...)

	key, err := keystore.GetDefault(ctx, keyStore.Workflow())
	if err != nil {
		return nil, fmt.Errorf("failed to get all workflow keys: %w", err)
	}

	artifactsStore := artifactsV1.NewStore(
		lggr,
		artifactsV1.NewWorkflowRegistryDS(ds, lggr),
		fetcherFunc,
		clockwork.NewRealClock(),
		key,
		custmsg.NewLabeler(),
		artifactsV1.WithMaxArtifactSize(
			artifactsV1.ArtifactConfig{
				MaxBinarySize:  uint64(capCfg.WorkflowRegistry().MaxBinarySize()),
				MaxSecretsSize: uint64(capCfg.WorkflowRegistry().MaxEncryptedSecretsSize()),
				MaxConfigSize:  uint64(capCfg.WorkflowRegistry().MaxConfigSize()),
			},
		),
	)

	engineRegistry := syncerV1.NewEngineRegistry()

	engineLimiters, err := v2.NewLimiters(lf, nil)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate engine limiters: %w", err)
	}
	srvcs = append(srvcs, closerService{name: "WorkflowEngineLimiters", Closer: engineLimiters})

	selector, err := chainSelector(capCfg.WorkflowRegistry().ChainID(), capCfg.WorkflowRegistry().NetworkID())
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow registry chain details by chain ID and network ID: %w", err)
	}

	eventHandler, err := syncerV1.NewEventHandler(
		lggr,
		workflowstore.NewInMemoryStore(lggr, clockwork.NewRealClock()),
		opts.CapabilitiesRegistry,
		dontimeStore,
		opts.UseLocalTimeProvider,
		engineRegistry,
		custmsg.NewLabeler(),
		engineLimiters,
		workflowRateLimiter,
		workflowLimits,
		artifactsStore,
		key,
		workflowDonNotifier,
		syncerV1.WithBillingClient(billingClient),
		syncerV1.WithWorkflowRegistry(capCfg.WorkflowRegistry().Address(), selector),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create workflow registry event handler: %w", err)
	}

	crFactory, err := newContractReaderFactory(capCfg, relayerChainInterops)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract reader factory")
	}

	wfSyncer, err := syncerV1.NewWorkflowRegistry(
		lggr,
		crFactory,
		capCfg.WorkflowRegistry().Address(),
		syncerV1.Config{
			QueryCount:   100,
			SyncStrategy: syncerV1.SyncStrategy(capCfg.WorkflowRegistry().SyncStrategy()),
		},
		eventHandler,
		workflowDonNotifier,
		engineRegistry,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create workflow registry syncer: %w", err)
	}

	srvcs = append(srvcs, wfSyncer)
	lggr.Debugw("Created WorkflowRegistrySyncer V1")
	return srvcs, nil
}

func newFetcherServiceV2(
	ctx context.Context,
	opts Opts,
	capCfg config.Capabilities,
	keyStore Keystore,
	lggr logger.Logger,
	gatewayConnectorWrapper *gatewayconnector.ServiceWrapper,
) (wftypes.FetcherFunc, wftypes.LocationRetrieverFunc, []commonsrv.Service, error) {
	if opts.FetcherFunc != nil {
		return opts.FetcherFunc, nil, []commonsrv.Service{}, nil
	}

	if gatewayConnectorWrapper == nil {
		return nil, nil, nil, errors.New("unable to create workflow registry syncer without gateway connector")
	}

	storageClient := opts.StorageClient
	if capCfg.WorkflowRegistry().WorkflowStorage().URL() != "" {
		jwtGenerator, err := newJWTGenerator(ctx, keyStore, opts)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create JWT generator for storage client: %w", err)
		}
		workflowOpts := []storage.WorkflowClientOpt{
			storage.WithJWTGenerator(jwtGenerator),
		}
		if capCfg.WorkflowRegistry().WorkflowStorage().TLSEnabled() {
			workflowOpts = append(workflowOpts, storage.WithWorkflowTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}

		sc, err := storage.NewWorkflowClient(lggr, capCfg.WorkflowRegistry().WorkflowStorage().URL(), workflowOpts...)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create storage client: %w", err)
		}

		storageClient = sc
	}

	if storageClient == nil {
		return nil, nil, nil, errors.New("must have a storage client")
	}

	fetcher := syncerV2.NewFetcherService(lggr, gatewayConnectorWrapper, storageClient)
	return fetcher.Fetch, fetcher.RetrieveURL, []commonsrv.Service{fetcher}, nil
}

func newWorkflowRegistrySyncerV2(
	ctx context.Context,
	cfg Config,
	relayerChainInterops RelayerChainInterops,
	billingClient metering.BillingClient,
	opts Opts,
	lggr logger.Logger,
	ds sqlutil.DataSource,
	dontimeStore *dontime.Store,
	workflowRateLimiter *ratelimiter.RateLimiter,
	workflowLimits limits.ResourceLimiter[int],
	workflowDonNotifier capabilities.DonNotifyWaitSubscriber,
	lf limits.Factory,
	orgResolver orgresolver.OrgResolver,
	gatewayConnectorWrapper *gatewayconnector.ServiceWrapper,
	keyStore Keystore,
) (syncerV2.WorkflowRegistrySyncer, []commonsrv.Service, error) {
	capCfg := cfg.Capabilities()
	key, err := keystore.GetDefault(ctx, keyStore.Workflow())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get all workflow keys: %w", err)
	}

	fetcherFunc, retrieverFunc, srvcs, err := newFetcherServiceV2(ctx, opts, capCfg, keyStore, lggr, gatewayConnectorWrapper)
	if err != nil {
		return nil, nil, err
	}

	artifactsStore, err := artifactsV2.NewStore(
		lggr,
		artifactsV2.NewWorkflowRegistryDS(ds, lggr),
		fetcherFunc,
		retrieverFunc,
		clockwork.NewRealClock(),
		key,
		custmsg.NewLabeler(),
		lf,
		artifactsV2.WithMaxArtifactSize(
			artifactsV2.ArtifactConfig{
				MaxBinarySize:  uint64(capCfg.WorkflowRegistry().MaxBinarySize()),
				MaxSecretsSize: uint64(capCfg.WorkflowRegistry().MaxEncryptedSecretsSize()),
				MaxConfigSize:  uint64(capCfg.WorkflowRegistry().MaxConfigSize()),
			},
		),
		artifactsV2.WithConfig(artifactsV2.StoreConfig{
			ArtifactStorageHost: capCfg.WorkflowRegistry().WorkflowStorage().ArtifactStorageHost(),
		}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create artifact store: %w", err)
	}

	engineRegistry := syncerV2.NewEngineRegistry()

	engineLimiters, err := v2.NewLimiters(lf, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not instantiate engine limiters: %w", err)
	}
	srvcs = append(srvcs, closerService{name: "WorkflowEngineLimiters", Closer: engineLimiters})

	selector, err := chainSelector(capCfg.WorkflowRegistry().ChainID(), capCfg.WorkflowRegistry().NetworkID())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get workflow registry chain details by chain ID and network ID: %w", err)
	}

	eventHandler, err := syncerV2.NewEventHandler(
		lggr,
		workflowstore.NewInMemoryStore(lggr, clockwork.NewRealClock()),
		dontimeStore,
		opts.UseLocalTimeProvider,
		opts.CapabilitiesRegistry,
		engineRegistry,
		custmsg.NewLabeler(),
		engineLimiters,
		workflowRateLimiter,
		workflowLimits,
		artifactsStore,
		key,
		workflowDonNotifier,
		syncerV2.WithBillingClient(billingClient),
		syncerV2.WithWorkflowRegistry(capCfg.WorkflowRegistry().Address(), selector),
		syncerV2.WithOrgResolver(orgResolver),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create workflow registry event handler: %w", err)
	}

	crFactory, err := newContractReaderFactory(capCfg, relayerChainInterops)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate contract reader factory")
	}

	shardOrchestratorClient, err := newShardOrchestratorClient(ctx, cfg, lggr)
	if err != nil {
		return nil, nil, err
	}
	if shardOrchestratorClient != nil {
		srvcs = append(srvcs, closerService{name: "ShardOrchestratorClient", Closer: shardOrchestratorClient})
	}

	addSources := capCfg.WorkflowRegistry().AdditionalSources()
	addSourceConfigs := make([]syncerV2.AdditionalSourceConfig, 0, len(addSources))
	if len(addSources) > 0 {
		jwtGenerator, jwtErr := newJWTGenerator(ctx, keyStore, opts)
		if jwtErr != nil {
			return nil, nil, fmt.Errorf("failed to create JWT generator for additional sources: %w", jwtErr)
		}
		for _, src := range addSources {
			addSourceConfigs = append(addSourceConfigs, syncerV2.AdditionalSourceConfig{
				URL:          src.GetURL(),
				Name:         src.GetName(),
				TLSEnabled:   src.GetTLSEnabled(),
				JWTGenerator: jwtGenerator,
			})
		}
	}

	workflowRegistrySyncerV2, err := syncerV2.NewWorkflowRegistry(
		lggr,
		crFactory,
		capCfg.WorkflowRegistry().Address(),
		selector,
		syncerV2.Config{
			QueryCount:   100,
			SyncStrategy: syncerV2.SyncStrategy(capCfg.WorkflowRegistry().SyncStrategy()),
		},
		eventHandler,
		workflowDonNotifier,
		engineRegistry,
		syncerV2.WithAdditionalSources(addSourceConfigs),
		syncerV2.WithShardOrchestratorClient(shardOrchestratorClient),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create workflow registry syncer: %w", err)
	}

	srvcs = append(srvcs, workflowRegistrySyncerV2)
	lggr.Debugw("Created WorkflowRegistrySyncer V2")
	return workflowRegistrySyncerV2, srvcs, nil
}

// newWorkflowRegistrySyncer creates a workflow registry syncer based on the contract version
func newWorkflowRegistrySyncer(
	ctx context.Context,
	cfg Config,
	keyStore Keystore,
	relayerChainInterops RelayerChainInterops,
	opts Opts,
	lggr logger.Logger,
	ds sqlutil.DataSource,
	dontimeStore *dontime.Store,
	workflowRateLimiter *ratelimiter.RateLimiter,
	workflowLimits limits.ResourceLimiter[int],
	workflowDonNotifier capabilities.DonNotifyWaitSubscriber,
	lf limits.Factory,
	orgResolver orgresolver.OrgResolver,
	gatewayConnectorWrapper *gatewayconnector.ServiceWrapper,
) (syncerV2.WorkflowRegistrySyncer, metering.BillingClient, []commonsrv.Service, error) {
	capCfg := cfg.Capabilities()

	lggr.Debugw("Creating WorkflowRegistrySyncer")
	lggr = logger.Named(lggr, "WorkflowRegistrySyncer")

	billingClient, err := newBillingClient(ctx, lggr, cfg, keyStore, opts)
	if err != nil {
		lggr.Infof("failed to create billing client: %s", err)
	}

	wrVersion, vErr := semver.NewVersion(capCfg.WorkflowRegistry().ContractVersion())
	if vErr != nil {
		return nil, nil, nil, vErr
	}

	switch wrVersion.Major() {
	case 1:
		srvcs, err := newWorkflowRegistrySyncerV1(
			ctx,
			capCfg,
			relayerChainInterops,
			billingClient,
			opts,
			lggr,
			ds,
			dontimeStore,
			workflowRateLimiter,
			workflowLimits,
			workflowDonNotifier,
			lf,
			gatewayConnectorWrapper,
			keyStore,
		)
		return nil, billingClient, srvcs, err
	case 2:
		syncer, srvcs, err := newWorkflowRegistrySyncerV2(
			ctx,
			cfg,
			relayerChainInterops,
			billingClient,
			opts,
			lggr,
			ds,
			dontimeStore,
			workflowRateLimiter,
			workflowLimits,
			workflowDonNotifier,
			lf,
			orgResolver,
			gatewayConnectorWrapper,
			keyStore,
		)
		return syncer, billingClient, srvcs, err
	default:
		return nil, nil, nil, fmt.Errorf("unsupported WorkflowRegistry contract version %s", wrVersion)
	}
}

// NewServices creates and initializes all CRE services
func NewServices(
	ctx context.Context,
	lggr logger.Logger,
	ds sqlutil.DataSource,
	keyStore Keystore,
	cfg Config,
	relayerChainInterops RelayerChainInterops,
	singletonPeerWrapper *ocrcommon.SingletonPeerWrapper,
	opts Opts,
) (*Services, error) {
	s := &Services{}

	var subservicesErr error
	s.Service, s.eng = commonsrv.Config{
		Name: "CRE",
		NewSubServices: func(subLggr logger.Logger) []commonsrv.Service {
			srvs, err := s.newSubservices(
				ctx,
				subLggr,
				ds,
				keyStore,
				cfg,
				relayerChainInterops,
				singletonPeerWrapper,
				opts,
			)
			if err != nil {
				subservicesErr = err
				return nil
			}
			return srvs
		},
	}.NewServiceEngine(lggr)

	if subservicesErr != nil {
		return nil, subservicesErr
	}

	return s, nil
}

var _ commonsrv.Service = closerService{}

// closerService extends an io.Closer to implement [services.ServiceCtx]
type closerService struct {
	name string
	io.Closer
}

func (c closerService) Start(ctx context.Context) error { return nil }

func (c closerService) Ready() error { return nil }

func (c closerService) HealthReport() map[string]error { return map[string]error{c.Name(): nil} }

func (c closerService) Name() string { return c.name }
