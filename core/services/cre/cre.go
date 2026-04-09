package cre

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"google.golang.org/grpc/credentials"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
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
	linkingclient "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"
	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/localcapmgr"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr/capregconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2pmain "github.com/smartcontractkit/chainlink/v2/core/services/p2p"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	p2pwrapper "github.com/smartcontractkit/chainlink/v2/core/services/p2p/wrapper"
	registrysyncerV1 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	registrysyncerV2 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	artifactsV1 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"
	artifactsV2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	wfmonitoring "github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/shardownership"
	workflowstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	syncerV1 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer"
	syncerV2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	wftypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
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

	WorkflowKey  workflowkey.Key
	JWTGenerator nodeauthjwt.JWTGenerator

	ShardOrchestratorClient shardorchestrator.ClientInterface
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

	// callback to wire Delegates into CRE services (e.g. Launcher) when ready
	SetDelegatesDeps func(*standardcapabilities.Delegate) (commonsrv.Service, error)
}

func (s *Services) close() error {
	return s.WorkflowLimits.Close()
}

// newSubservices initializes and returns all CRE child services
func (s *Services) newSubservices(
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

	if capCfg.GatewayConnector().DonID() != "" {
		lggr.Debugw("Creating GatewayConnector wrapper", "donID", capCfg.GatewayConnector().DonID())
		gatewayConnectorWrapper, ierr := newGatewayConnectorWrapper(capCfg, keyStore, lggr)
		if ierr != nil {
			return nil, fmt.Errorf("could not create gateway connector wrapper: %w", ierr)
		}
		s.GatewayConnectorWrapper = gatewayConnectorWrapper
		srvs = append(srvs, gatewayConnectorWrapper)
	}

	if cfg.CRE().Linking().URL() != "" {
		lggr.Debugw("Creating OrgResolver")
		orgResolver, ierr := newOrgResolver(cfg, capCfg, opts, lggr)
		if ierr != nil {
			return nil, fmt.Errorf("could not create org resolver: %w", ierr)
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

	if dispatcherWrapper.dispatcher == nil {
		lggr.Warn("Skipping capabilities and workflow registry syncer, no dispatcher configured (peering disabled)")
		return srvs, nil
	}

	if capCfg.ExternalRegistry().Address() == "" {
		lggr.Warn("Skipping capabilities and workflow registry syncer, none configured")
		return srvs, nil
	}

	registrySyncerServices, donNotifier, err := s.newRegistrySyncer(
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
		cfg,
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
	ocrConfigService capregconfig.OCRConfigService,
	wfLauncher registrysyncerV1.Listener,
) ([]commonsrv.Service, error) {
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

	registrySyncer.AddListener(wfLauncher, ocrConfigService)
	return []commonsrv.Service{registrySyncer, ocrConfigService}, nil
}

func newRegistrySyncerV2(
	lggr logger.Logger,
	getPeerID func() (p2ptypes.PeerID, error),
	relayer loop.Relayer,
	registryAddress string,
	ds sqlutil.DataSource,
	ocrConfigService capregconfig.OCRConfigService,
	wfLauncher registrysyncerV1.Listener,
) ([]commonsrv.Service, error) {
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

	registrySyncer.AddListener(wfLauncher, ocrConfigService)
	return []commonsrv.Service{registrySyncer, ocrConfigService}, nil
}

// newRegistrySyncer creates a registry syncer based on the external registry version
func (s *Services) newRegistrySyncer(
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

	ocrConfigService, ocrErr := newOCRConfigService(lggr, rid, registryAddress, dispatcherWrapper)
	if ocrErr != nil {
		return nil, nil, ocrErr
	}
	s.OCRConfigService = ocrConfigService

	externalRegistryVersion, err := semver.NewVersion(capCfg.ExternalRegistry().ContractVersion())
	if err != nil {
		return nil, nil, err
	}

	wfLauncher, err := capabilities.NewLauncher(
		lggr,
		dispatcherWrapper.externalPeerWrapper,
		dispatcherWrapper.don2DonSharedPeer,
		streamConfig,
		dispatcherWrapper.dispatcher,
		opts.CapabilitiesRegistry,
		donNotifier,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create workflow launcher: %w", err)
	}
	srvcs = append(srvcs, wfLauncher)

	// callback to wire LocalCapabilityManager into the launcher if local capabilities are configured.
	localCfg := cfg.Capabilities().Local()
	if localCfg != nil && len(localCfg.RegistryBasedLaunchAllowlist()) > 0 {
		s.SetDelegatesDeps = func(stdcapDelegate *standardcapabilities.Delegate) (commonsrv.Service, error) {
			newServicesFn := func(ctx context.Context, capID string, command string, configJSON string) ([]job.ServiceCtx, error) {
				return stdcapDelegate.NewServices(ctx, command, configJSON, 0, capID, uuid.New(), job.OracleFactoryConfig{})
			}
			localCapMgr, lcmErr := localcapmgr.NewLocalCapabilityManager(lggr, localCfg, newServicesFn)
			if lcmErr != nil {
				return nil, fmt.Errorf("could not create local capability manager: %w", lcmErr)
			}
			wfLauncher.SetLocalCapabilityManager(localCapMgr)
			return localCapMgr, nil
		}
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
			ocrConfigService,
			wfLauncher,
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
			ocrConfigService,
			wfLauncher,
		)
		if err != nil {
			return nil, nil, err
		}
		srvcs = append(srvcs, srvs...)
		return srvcs, donNotifier, nil
	}

	return nil, nil, fmt.Errorf("could not configure capability registry syncer with version: %d", externalRegistryVersion.Major())
}

func newOCRConfigService(
	lggr logger.Logger,
	rid commontypes.RelayID,
	registryAddress string,
	dispatcherWrapper *dispatcherWrapper,
) (capregconfig.OCRConfigService, error) {
	registryChainID, err := strconv.ParseUint(rid.ChainID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse registry chain ID for OCRConfigService: %w", err)
	}
	return capregconfig.NewOCRConfigService(
		lggr,
		dispatcherWrapper.GetPeerID,
		registryChainID,
		registryAddress,
	), nil
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
		opts.CapabilitiesRegistry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
		return nil, nil
	}

	if opts.CapabilitiesDispatcher != nil {
		w.dispatcher = opts.CapabilitiesDispatcher
		w.externalPeerWrapper = opts.CapabilitiesPeerWrapper
		return []commonsrv.Service{w.externalPeerWrapper, w.dispatcher}, nil
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
	cfg Config,
	capCfg config.Capabilities,
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
		JWTGenerator:                  opts.JWTGenerator,
	}

	var (
		resolver orgresolver.OrgResolver
		err      error
	)
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

func newBillingClient(lggr logger.Logger, cfg Config, opts Opts) (metering.BillingClient, error) {
	if opts.BillingClient != nil {
		return opts.BillingClient, nil
	}

	if cfg.Billing().URL() == "" {
		return nil, nil
	}

	workflowOpts := []billing.WorkflowClientOpt{
		billing.WithJWTGenerator(opts.JWTGenerator),
	}
	if cfg.Billing().TLSEnabled() {
		workflowOpts = append(workflowOpts, billing.WithWorkflowTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}
	return billing.NewWorkflowClient(lggr, cfg.Billing().URL(), workflowOpts...)
}

func newShardOrchestratorClient(cfg Config, lggr logger.Logger) (*shardorchestrator.Client, error) {
	shardID := cfg.Sharding().ShardIndex()
	if shardID == 0 {
		return nil, nil
	}

	address := cfg.Sharding().ShardOrchestratorAddress()
	if address == nil {
		return nil, fmt.Errorf("shard %d requires ShardOrchestratorAddress configuration", shardID)
	}

	client, err := shardorchestrator.NewClient(address.String(), lggr)
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
) ([]commonsrv.Service, error) {
	var srvcs []commonsrv.Service

	fetcherFunc, srvs, err := newFetcherFuncV1(lggr, opts.FetcherFunc, gatewayConnectorWrapper)
	if err != nil {
		return nil, err
	}
	srvcs = append(srvcs, srvs...)

	key := opts.WorkflowKey

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

	featureFlags, err := v2.NewFeatureFlags(lf, nil)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate engine feature flags: %w", err)
	}

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
		featureFlags,
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
		return nil, fmt.Errorf("failed to instantiate contract reader factory: %w", err)
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
	opts Opts,
	capCfg config.Capabilities,
	lggr logger.Logger,
	gatewayConnectorWrapper *gatewayconnector.ServiceWrapper,
) (wftypes.FetcherFunc, wftypes.LocationRetrieverFunc, []commonsrv.Service, error) {
	if opts.FetcherFunc != nil {
		return opts.FetcherFunc, nil, []commonsrv.Service{}, nil
	}

	if gatewayConnectorWrapper == nil {
		return nil, nil, nil, errors.New("unable to create workflow registry syncer without gateway connector")
	}

	wfStorage := capCfg.WorkflowRegistry().WorkflowStorage()
	storageClient := opts.StorageClient
	if wfStorage.URL() != "" {
		workflowOpts := []storage.WorkflowClientOpt{
			storage.WithJWTGenerator(opts.JWTGenerator),
		}
		if wfStorage.TLSEnabled() {
			workflowOpts = append(workflowOpts, storage.WithWorkflowTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}

		sc, err := storage.NewWorkflowClient(lggr, wfStorage.URL(), workflowOpts...)
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
) (syncerV2.WorkflowRegistrySyncer, []commonsrv.Service, error) {
	capCfg := cfg.Capabilities()
	wfReg := capCfg.WorkflowRegistry()
	key := opts.WorkflowKey

	fetcherFunc, retrieverFunc, srvcs, err := newFetcherServiceV2(opts, capCfg, lggr, gatewayConnectorWrapper)
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
				MaxBinarySize:  uint64(wfReg.MaxBinarySize()),
				MaxSecretsSize: uint64(wfReg.MaxEncryptedSecretsSize()),
				MaxConfigSize:  uint64(wfReg.MaxConfigSize()),
			},
		),
		artifactsV2.WithConfig(artifactsV2.StoreConfig{
			ArtifactStorageHost: wfReg.WorkflowStorage().ArtifactStorageHost(),
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

	featureFlags, err := v2.NewFeatureFlags(lf, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not instantiate engine feature flags: %w", err)
	}

	selector, err := chainSelector(wfReg.ChainID(), wfReg.NetworkID())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get workflow registry chain details by chain ID and network ID: %w", err)
	}

	crFactory, err := newContractReaderFactory(capCfg, relayerChainInterops)
	if err != nil {
		return nil, nil, errors.New("failed to instantiate contract reader factory")
	}

	var shardOrchestratorClient shardorchestrator.ClientInterface
	if opts.ShardOrchestratorClient != nil {
		shardOrchestratorClient = opts.ShardOrchestratorClient
	} else {
		var c shardorchestrator.ClientInterface
		c, err = newShardOrchestratorClient(cfg, lggr)
		if err != nil {
			return nil, nil, err
		}
		shardOrchestratorClient = c
	}

	shardingEnabled := cfg.Sharding().ShardingEnabled()
	shardIndex := uint32(cfg.Sharding().ShardIndex())

	var shardRoutingSteady *shardownership.SteadySignal
	if shardingEnabled {
		steadyMetrics, errSteady := wfmonitoring.GlobalSteadySignalMetrics()
		if errSteady != nil {
			lggr.Warnw("Failed to register shard routing steady signal metrics; continuing without steady instrumentation", "err", errSteady)
		}
		shardRoutingSteady = shardownership.NewSteadySignal(shardownership.WithSteadySignalMetrics(steadyMetrics))
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
		featureFlags,
		workflowRateLimiter,
		workflowLimits,
		artifactsStore,
		key,
		workflowDonNotifier,
		syncerV2.WithBillingClient(billingClient),
		syncerV2.WithWorkflowRegistry(wfReg.Address(), selector),
		syncerV2.WithOrgResolver(orgResolver),
		syncerV2.WithDebugMode(cfg.CRE().DebugMode()),
		syncerV2.WithLocalSecrets(lggr, cfg.CRE().LocalSecrets()),
		syncerV2.WithShardExecutionGuard(shardOrchestratorClient, shardingEnabled, shardIndex),
		syncerV2.WithShardRoutingSteady(shardRoutingSteady),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create workflow registry event handler: %w", err)
	}

	addSources := wfReg.AdditionalSources()
	addSourceConfigs := make([]syncerV2.AdditionalSourceConfig, len(addSources))
	for i, src := range addSources {
		addSourceConfigs[i] = syncerV2.AdditionalSourceConfig{
			URL:          src.GetURL(),
			Name:         src.GetName(),
			TLSEnabled:   src.GetTLSEnabled(),
			JWTGenerator: opts.JWTGenerator,
		}
	}

	registryOpts := []syncerV2.Option{
		syncerV2.WithAdditionalSources(addSourceConfigs),
		syncerV2.WithShardOrchestratorClient(shardOrchestratorClient),
		syncerV2.WithMaxConcurrency(wfReg.MaxConcurrency()),
	}
	if cfg.Sharding().ShardingEnabled() {
		registryOpts = append(registryOpts,
			syncerV2.WithShardEnabled(true),
			syncerV2.WithShardID(uint32(cfg.Sharding().ShardIndex())),
		)
		if shardRoutingSteady != nil {
			registryOpts = append(registryOpts, syncerV2.WithRegistryShardRoutingObserver(shardRoutingSteady))
		}
	}

	workflowRegistrySyncerV2, err := syncerV2.NewWorkflowRegistry(
		lggr,
		crFactory,
		wfReg.Address(),
		selector,
		syncerV2.Config{
			QueryCount:   100,
			SyncStrategy: syncerV2.SyncStrategy(wfReg.SyncStrategy()),
		},
		eventHandler,
		workflowDonNotifier,
		engineRegistry,
		registryOpts...,
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
	cfg Config,
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

	billingClient, err := newBillingClient(lggr, cfg, opts)
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
		)
		return nil, billingClient, srvcs, err
	case 2:
		syncer, srvcs, err := newWorkflowRegistrySyncerV2(
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
		)
		return syncer, billingClient, srvcs, err
	default:
		return nil, nil, nil, fmt.Errorf("unsupported WorkflowRegistry contract version %s", wrVersion)
	}
}

// NewServices creates and initializes all CRE services
func NewServices(
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
		Close: s.close,
	}.NewServiceEngine(lggr)

	if subservicesErr != nil {
		return nil, subservicesErr
	}

	return s, nil
}
