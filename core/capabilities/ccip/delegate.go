package ccip

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	configsevm "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/oraclecreator"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type RelayGetter interface {
	Get(types.RelayID) (loop.Relayer, error)
	GetIDToRelayerMap() (map[types.RelayID]loop.Relayer, error)
}

type Delegate struct {
	lggr                  logger.Logger
	registrarConfig       plugins.RegistrarConfig
	pipelineRunner        pipeline.Runner
	relayers              RelayGetter
	keystore              keystore.Master
	ds                    sqlutil.DataSource
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	capabilityConfig      config.Capabilities
	evmConfigs            toml.EVMConfigs

	isNewlyCreatedJob bool
}

func NewDelegate(
	lggr logger.Logger,
	registrarConfig plugins.RegistrarConfig,
	pipelineRunner pipeline.Runner,
	relayers RelayGetter,
	keystore keystore.Master,
	ds sqlutil.DataSource,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	capabilityConfig config.Capabilities,
	evmConfigs toml.EVMConfigs,
) *Delegate {
	return &Delegate{
		lggr:                  lggr,
		registrarConfig:       registrarConfig,
		pipelineRunner:        pipelineRunner,
		relayers:              relayers,
		ds:                    ds,
		keystore:              keystore,
		peerWrapper:           peerWrapper,
		monitoringEndpointGen: monitoringEndpointGen,
		capabilityConfig:      capabilityConfig,
		evmConfigs:            evmConfigs,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.CCIP
}

func (d *Delegate) BeforeJobCreated(job.Job) {
	// This is only called first time the job is created
	d.isNewlyCreatedJob = true
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	// In general there should only be one P2P key but the node may have multiple.
	// The job spec should specify the correct P2P key to use.
	peerID, err := p2pkey.MakePeerID(spec.CCIPSpec.P2PKeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to make peer ID from provided spec p2p id (%s): %w", spec.CCIPSpec.P2PKeyID, err)
	}

	p2pID, err := d.keystore.P2P().Get(peerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all p2p keys: %w", err)
	}

	cfg := d.capabilityConfig
	rid := cfg.ExternalRegistry().RelayID()
	homeChainRelayer, err := d.relayers.Get(rid)
	if err != nil {
		return nil, fmt.Errorf("could not fetch relayer %s configured for capabilities registry: %w", rid, err)
	}
	registrySyncer, err := registrysyncer.New(
		d.lggr,
		func() (p2ptypes.PeerID, error) {
			return p2ptypes.PeerID(p2pID.PeerID()), nil
		},
		homeChainRelayer,
		cfg.ExternalRegistry().Address(),
		registrysyncer.NewORM(d.ds, d.lggr),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create registry syncer: %w", err)
	}

	ocrKeys, err := d.getOCRKeys(spec.CCIPSpec.OCRKeyBundleIDs)
	if err != nil {
		return nil, err
	}

	allRelayers, err := d.relayers.GetIDToRelayerMap()
	if err != nil {
		return nil, fmt.Errorf("could not fetch all relayers: %w", err)
	}
	transmitterKeys, err := d.getTransmitterKeys(ctx, maps.Keys(allRelayers))
	if err != nil {
		return nil, err
	}

	bootstrapperLocators, err := ocrcommon.ParseBootstrapPeers(spec.CCIPSpec.P2PV2Bootstrappers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bootstrapper locators: %w", err)
	}

	// NOTE: we can use the same DB for all plugin instances,
	// since all queries are scoped by config digest.
	ocrDB := ocr2.NewDB(d.ds, spec.ID, 0, d.lggr)

	homeChainContractReader, ccipConfigBinding, err := d.getHomeChainContractReader(
		ctx,
		homeChainRelayer,
		spec.CCIPSpec.CapabilityLabelledName,
		spec.CCIPSpec.CapabilityVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get home chain contract reader: %w", err)
	}

	hcr := ccipreaderpkg.NewHomeChainReader(
		homeChainContractReader,
		d.lggr.Named("HomeChainReader"),
		100*time.Millisecond,
		ccipConfigBinding,
	)

	// get the chain selector for the home chain
	homeChainChainID, err := strconv.ParseUint(d.capabilityConfig.ExternalRegistry().RelayID().ChainID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chain ID %s: %w", d.capabilityConfig.ExternalRegistry().RelayID().ChainID, err)
	}
	homeChainChainSelector, err := chainsel.SelectorFromChainId(homeChainChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain selector from chain ID %d", homeChainChainID)
	}

	// if bootstrappers are provided we assume that the node is a plugin oracle.
	// the reason for this is that bootstrap oracles do not need to be aware
	// of other bootstrap oracles. however, plugin oracles, at least initially,
	// must be aware of available bootstrappers.
	var oracleCreator cctypes.OracleCreator
	if len(spec.CCIPSpec.P2PV2Bootstrappers) > 0 {
		oracleCreator = oraclecreator.NewPluginOracleCreator(
			ocrKeys,
			transmitterKeys,
			allRelayers,
			d.peerWrapper,
			spec.ExternalJobID,
			spec.ID,
			d.isNewlyCreatedJob,
			spec.CCIPSpec.PluginConfig,
			ocrDB,
			d.lggr,
			d.monitoringEndpointGen,
			bootstrapperLocators,
			hcr,
			cciptypes.ChainSelector(homeChainChainSelector),
			d.evmConfigs,
		)
	} else {
		oracleCreator = oraclecreator.NewBootstrapOracleCreator(
			d.peerWrapper,
			bootstrapperLocators,
			ocrDB,
			d.monitoringEndpointGen,
			d.lggr,
		)
	}

	capabilityID := fmt.Sprintf("%s@%s", spec.CCIPSpec.CapabilityLabelledName, spec.CCIPSpec.CapabilityVersion)
	capLauncher := launcher.New(
		capabilityID,
		ragep2ptypes.PeerID(p2pID.PeerID()),
		d.lggr,
		hcr,
		12*time.Second,
		oracleCreator,
	)

	// register the capability launcher with the registry syncer
	registrySyncer.AddLauncher(capLauncher)

	return []job.ServiceCtx{
		homeChainContractReader,
		registrySyncer,
		hcr,
		capLauncher,
	}, nil
}

func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, spec job.Job) error {
	// TODO: shut down needed services?
	return nil
}

func (d *Delegate) getOCRKeys(ocrKeyBundleIDs job.JSONConfig) (map[string]ocr2key.KeyBundle, error) {
	ocrKeys := make(map[string]ocr2key.KeyBundle)
	for networkType, bundleIDRaw := range ocrKeyBundleIDs {
		if networkType != relay.NetworkEVM {
			return nil, fmt.Errorf("unsupported chain type: %s", networkType)
		}

		bundleID, ok := bundleIDRaw.(string)
		if !ok {
			return nil, fmt.Errorf("OCRKeyBundleIDs must be a map of chain types to OCR key bundle IDs, got: %T", bundleIDRaw)
		}

		bundle, err2 := d.keystore.OCR2().Get(bundleID)
		if err2 != nil {
			return nil, fmt.Errorf("OCR key bundle with ID %s not found: %w", bundleID, err2)
		}

		ocrKeys[networkType] = bundle
	}
	return ocrKeys, nil
}

func (d *Delegate) getTransmitterKeys(ctx context.Context, relayIDs []types.RelayID) (map[types.RelayID][]string, error) {
	transmitterKeys := make(map[types.RelayID][]string)
	for _, relayID := range relayIDs {
		chainID, ok := new(big.Int).SetString(relayID.ChainID, 10)
		if !ok {
			return nil, fmt.Errorf("error parsing chain ID, expected big int: %s", relayID.ChainID)
		}

		ethKeys, err := d.keystore.Eth().EnabledAddressesForChain(ctx, chainID)
		if err != nil {
			return nil, fmt.Errorf("error getting enabled addresses for chain: %s %w", chainID.String(), err)
		}

		transmitterKeys[relayID] = func() (r []string) {
			for _, key := range ethKeys {
				r = append(r, key.Hex())
			}
			return
		}()
	}
	return transmitterKeys, nil
}

func (d *Delegate) getHomeChainContractReader(
	ctx context.Context,
	homeChainRelayer loop.Relayer,
	capabilityLabelledName,
	capabilityVersion string,
) (types.ContractReader, types.BoundContract, error) {
	reader, err := homeChainRelayer.NewContractReader(ctx, configsevm.HomeChainReaderConfig)
	if err != nil {
		return nil, types.BoundContract{}, fmt.Errorf("failed to create home chain contract reader: %w", err)
	}

	reader, ccipConfigBinding, err := bindReader(ctx, reader, d.capabilityConfig.ExternalRegistry().Address(), capabilityLabelledName, capabilityVersion)
	if err != nil {
		return nil, types.BoundContract{}, fmt.Errorf("failed to bind home chain contract reader: %w", err)
	}

	return reader, ccipConfigBinding, nil
}

func bindReader(ctx context.Context,
	reader types.ContractReader,
	capRegAddress,
	capabilityLabelledName,
	capabilityVersion string,
) (boundReader types.ContractReader, ccipConfigBinding types.BoundContract, err error) {
	boundContract := types.BoundContract{
		Address: capRegAddress,
		Name:    consts.ContractNameCapabilitiesRegistry,
	}

	err = reader.Bind(ctx, []types.BoundContract{boundContract})
	if err != nil {
		return nil, types.BoundContract{}, fmt.Errorf("failed to bind home chain contract reader: %w", err)
	}

	hid, err := common.HashedCapabilityID(capabilityLabelledName, capabilityVersion)
	if err != nil {
		return nil, types.BoundContract{}, fmt.Errorf("failed to hash capability id: %w", err)
	}

	var ccipCapabilityInfo kcr.CapabilitiesRegistryCapabilityInfo
	err = reader.GetLatestValue(
		ctx,
		boundContract.ReadIdentifier(consts.MethodNameGetCapability),
		primitives.Unconfirmed,
		map[string]any{
			"hashedId": hid,
		},
		&ccipCapabilityInfo,
	)
	if err != nil {
		return nil, types.BoundContract{}, fmt.Errorf("failed to get CCIP capability info from chain reader: %w", err)
	}

	// bind the ccip capability configuration contract
	ccipConfigBinding = types.BoundContract{
		Address: ccipCapabilityInfo.ConfigurationContract.String(),
		Name:    consts.ContractNameCCIPConfig,
	}
	err = reader.Bind(ctx, []types.BoundContract{ccipConfigBinding})
	if err != nil {
		return nil, types.BoundContract{}, fmt.Errorf("failed to bind CCIP capability configuration contract: %w", err)
	}

	return reader, ccipConfigBinding, nil
}
