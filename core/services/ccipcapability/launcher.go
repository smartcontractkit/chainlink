package ccipcapability

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

var (
	_ job.ServiceCtx = (*launcher)(nil)
)

const (
	tickInterval = 10 * time.Second
)

type ccipOracle struct {
	oracle     libocr2.Oracle
	pluginType PluginType
	configChan chan OCRConfig
	// TODO: full OCR config for this oracle?
	// some way to figure out what instance is blue and what instance is green?
	// can just use config counts?
}

func (o ccipOracle) Shutdown() error {
	return o.oracle.Close()
}

// launcher manages the lifecycles of the CCIP capability on all chains.
type launcher struct {
	capabilityVersion      string
	capabilityLabelledName string
	ocrKeyBundles          map[chaintype.ChainType]ocr2key.KeyBundle
	transmitters           map[types.RelayID]string
	p2pID                  p2pkey.KeyV2
	relayers               map[types.RelayID]loop.Relayer
	capRegistry            CapabilityRegistry
	lggr                   logger.Logger
	homeChainReader        HomeChainReader
	stopChan               chan struct{}
	regState               RegistryState
	peerWrapper            *ocrcommon.SingletonPeerWrapper
	externalJobID          uuid.UUID
	jobID                  int32
	isNewlyCreatedJob      bool
	relayConfigs           map[string]job.JSONConfig
	pluginConfig           job.JSONConfig
	db                     ocr3types.Database

	// dons is a map of CCIP DON IDs to the OCR instances that are running on them.
	// we can have up to two OCR instances per CCIP plugin, since we are running two plugins,
	// thats four OCR instances per CCIP DON maximum.
	dons map[uint32][]ccipOracle
}

// Close implements job.ServiceCtx.
func (l *launcher) Close() error {
	close(l.stopChan)
	return nil
}

// Start implements job.ServiceCtx.
func (l *launcher) Start(context.Context) error {
	l.stopChan = make(chan struct{})
	go l.monitor()
	return nil
}

func (l *launcher) monitor() {
	ticker := time.NewTicker(tickInterval)
	for {
		select {
		case <-l.stopChan:
			return
		case <-ticker.C:
			// Ensure that the home chain reader is healthy.
			// For new jobs it may be possible that the home chain reader is not yet ready
			// so we won't be able to fetch configs and start any OCR instances.
			if !l.homeChainReader.IsHealthy() {
				l.lggr.Error("Home chain reader is unhealthy, retrying in 10 seconds")
				continue
			}

			// Fetch the latest state from the capability registry and determine if we need to
			// launch or update any OCR instances.
			latestState, err := l.capRegistry.LatestState()
			if err != nil {
				l.lggr.Errorw("Failed to fetch latest state from capability registry", "err", err)
				continue
			}

			diffResult, err := diff(l.capabilityVersion, l.capabilityLabelledName, l.regState, latestState)
			if err != nil {
				l.lggr.Errorw("Failed to diff capability registry states", "err", err)
				continue
			}

			err = l.processDiff(diffResult)
			if err != nil {
				l.lggr.Errorw("Failed to process diff", "err", err)
				continue
			}

			// if diff is correctly processed, update latest state.
			l.regState = latestState
		}
	}
}

// processDiff processes the diff between the current and latest capability registry states.
// for any added OCR instances, it will launch them.
// for any removed OCR instances, it will shut them down.
// for any updated OCR instances, it will restart them with the new configuration.
func (l *launcher) processDiff(diff diffResult) error {
	for id := range diff.removed {
		oracles, ok := l.dons[id]
		if !ok {
			// not running this particular DON.
			continue
		}
		for _, oracle := range oracles {
			if err := oracle.Shutdown(); err != nil {
				l.lggr.Errorw("Failed to shutdown OCR instance", "err", err, "donId", id)
				continue
			}
		}
	}

	for id, don := range diff.added {
		// check if we are a member of this DON.
		// if we are, spin up the OCR instances specified by the OCR configs.
		var found bool
		for _, node := range don.NodeP2PIds {
			if node == l.p2pID.PeerID() {
				found = true
				break
			}
		}
		if !found {
			l.lggr.Infow("Not a member of this DON, skipping", "donId", id, "p2pId", l.p2pID.ID())
			continue
		}

		// this should be a retryable error.
		commitOCRConfigs, err := l.homeChainReader.GetOCRConfigs(context.Background(), don.Id, PluginTypeCCIPCommit)
		if err != nil {
			return fmt.Errorf("failed to fetch OCR configs for CCIP commit plugin (don id: %d) from home chain config contract: %w",
				don.Id, err)
		}

		execOCRConfigs, err := l.homeChainReader.GetOCRConfigs(context.Background(), don.Id, PluginTypeCCIPExec)
		if err != nil {
			return fmt.Errorf("failed to fetch OCR configs for CCIP exec plugin (don id: %d) from home chain config contract: %w",
				don.Id, err)
		}

		// upon creation we should only have one OCR config per plugin type.
		if len(commitOCRConfigs) != 1 {
			return fmt.Errorf("expected exactly one OCR config for CCIP commit plugin (don id: %d), got %d", don.Id, len(commitOCRConfigs))
		}

		if len(execOCRConfigs) != 1 {
			return fmt.Errorf("expected exactly one OCR config for CCIP exec plugin (don id: %d), got %d", don.Id, len(execOCRConfigs))
		}

		commitOracle, err := l.createCommitOracle(commitOCRConfigs[0])
		if err != nil {
			return fmt.Errorf("failed to create CCIP commit oracle: %w", err)
		}

		execOracle, err := l.createExecOracle(execOCRConfigs[0])
		if err != nil {
			return fmt.Errorf("failed to create CCIP exec oracle: %w", err)
		}
	}
}

func (l *launcher) createExecOracle(config OCRConfig) (ccipOracle, error) {
	// Create the exec provider.
	// TODO: get the right relay ID from chain selector.
	relayID := types.NewRelayID("evm", "1")
	relayer, ok := l.relayers[relayID]
	if !ok {
		return ccipOracle{}, fmt.Errorf("unable to find relayer for chain selector %d", config.ChainSelector())
	}

	execProvider, err := relayer.NewPluginProvider(
		context.Background(),
		types.RelayArgs{
			ExternalJobID: l.externalJobID,
			JobID:         l.jobID,
			ContractID:    config.OfframpAddress(),
			New:           l.isNewlyCreatedJob,
			RelayConfig:   l.relayConfigs[relayID.Network].Bytes(),
			ProviderType:  "CCIPOCR3ExecProvider",
		},
		types.PluginArgs{
			TransmitterID: l.transmitters[relayID],
			PluginConfig:  l.pluginConfig.Bytes(),
		},
	)
	if err != nil {
		return ccipOracle{}, fmt.Errorf("failed to create commit provider", err)
	}

	ocrLogger := ocrcommon.NewOCRWrapper(
		l.lggr.Named(fmt.Sprintf("CCIPOCR3Commit-%d", config.ChainSelector())),
		false, /* l.cfg.OCR2().TraceLogging() */
		func(ctx context.Context, msg string) {
			// l.lggr.ErrorIf(d.jobORM.RecordError(ctx, jb.ID, msg), "unable to record error") // TODO
		},
	)

	// TODO: ReportInfo generic type.
	configChan := make(chan OCRConfig)
	oracleArgs := libocr2.OCR3OracleArgs[[]byte]{
		BinaryNetworkEndpointFactory: l.peerWrapper.Peer2, // TODO
		V2Bootstrappers:              nil,                 // TODO how to get?
		ContractConfigTracker:        newConfigTracker(configChan),
		// ContractTransmitter:          commitProvider.ContractTransmitter(), // TODO need OCR3 transmitter
		Database: l.db,
		// how much of this local config is required/applicable for the role don?
		LocalConfig: ocrtypes.LocalConfig{
			BlockchainTimeout:                  0,    // TODO
			ContractConfigConfirmations:        0,    // TODO
			SkipContractConfigConfirmations:    true, // TODO
			ContractConfigTrackerPollInterval:  0,    // TODO
			ContractTransmitterTransmitTimeout: 0,    // TODO
			DatabaseTimeout:                    0,    // TODO
			MinOCR2MaxDurationQuery:            0,    // TODO
			DevelopmentMode:                    "",
		}, // TODO
		Logger:                 ocrLogger, // TODO
		MetricsRegisterer:      nil,       // TODO
		MonitoringEndpoint:     nil,       // TODO
		OffchainConfigDigester: &offchainConfigDigester{},
		// TODO: need to know the chain type from home chain config, perhaps chain config?
		OffchainKeyring:        l.ocrKeyBundles[chaintype.EVM],
		OnchainKeyring:         ocrcommon.NewOCR3OnchainKeyringAdapter(l.ocrKeyBundles[chaintype.EVM]),
		ReportingPluginFactory: nil, // TODO
	}
	oracle, err := libocr2.NewOracle(oracleArgs)
	if err != nil {
		return ccipOracle{}, fmt.Errorf("unable to create commit oracle, err: %w", err)
	}

	return ccipOracle{oracle: oracle, pluginType: PluginTypeCCIPExec, configChan: configChan}, nil
}

func (l *launcher) createCommitOracle(config OCRConfig) (ccipOracle, error) {
	// Create the commit provider.
	// TODO: get the right relay ID from chain selector.
	relayID := types.NewRelayID("evm", "1")
	relayer, ok := l.relayers[relayID]
	if !ok {
		return ccipOracle{}, fmt.Errorf("unable to find relayer for chain selector %d", config.ChainSelector())
	}

	// TODO: create chain readers for all enabled chains on this node.
	// TODO: create chain writers for all enabled chains on this node.
	// TODO: pass in the readers + writers through the factory somehow.
	commitProvider, err := relayer.NewPluginProvider(
		context.Background(),
		types.RelayArgs{
			ExternalJobID: l.externalJobID,
			JobID:         l.jobID,
			ContractID:    config.OfframpAddress(),
			New:           l.isNewlyCreatedJob,
			RelayConfig:   l.relayConfigs[relayID.Network].Bytes(),
			ProviderType:  "CCIPOCR3CommitProvider",
		},
		types.PluginArgs{
			TransmitterID: l.transmitters[relayID],
			PluginConfig:  l.pluginConfig.Bytes(),
		},
	)
	if err != nil {
		return ccipOracle{}, fmt.Errorf("failed to create commit provider", err)
	}

	commitProvider.ContractTransmitter()

	ocrLogger := ocrcommon.NewOCRWrapper(
		l.lggr.Named(fmt.Sprintf("CCIPOCR3Commit-%d", config.ChainSelector())),
		false, /* l.cfg.OCR2().TraceLogging() */
		func(ctx context.Context, msg string) {
			// l.lggr.ErrorIf(d.jobORM.RecordError(ctx, jb.ID, msg), "unable to record error") // TODO
		},
	)

	// TODO: ReportInfo generic type.
	configChan := make(chan OCRConfig)
	oracleArgs := libocr2.OCR3OracleArgs[[]byte]{
		BinaryNetworkEndpointFactory: l.peerWrapper.Peer2, // TODO
		V2Bootstrappers:              nil,                 // TODO how to get?
		ContractConfigTracker:        newConfigTracker(configChan),
		// ContractTransmitter:          commitProvider.ContractTransmitter(), // TODO need OCR3 transmitter
		Database: l.db,
		// how much of this local config is required/applicable for the role don?
		LocalConfig: ocrtypes.LocalConfig{
			BlockchainTimeout:                  0,    // TODO
			ContractConfigConfirmations:        0,    // TODO
			SkipContractConfigConfirmations:    true, // TODO
			ContractConfigTrackerPollInterval:  0,    // TODO
			ContractTransmitterTransmitTimeout: 0,    // TODO
			DatabaseTimeout:                    0,    // TODO
			MinOCR2MaxDurationQuery:            0,    // TODO
			DevelopmentMode:                    "",
		}, // TODO
		Logger:                 ocrLogger, // TODO
		MetricsRegisterer:      nil,       // TODO
		MonitoringEndpoint:     nil,       // TODO
		OffchainConfigDigester: &offchainConfigDigester{},
		// TODO: need to know the chain type from home chain config, perhaps chain config?
		OffchainKeyring:        l.ocrKeyBundles[chaintype.EVM],
		OnchainKeyring:         ocrcommon.NewOCR3OnchainKeyringAdapter(l.ocrKeyBundles[chaintype.EVM]),
		ReportingPluginFactory: nil, // TODO
	}
	oracle, err := libocr2.NewOracle(oracleArgs)
	if err != nil {
		return ccipOracle{}, fmt.Errorf("unable to create commit oracle, err: %w", err)
	}

	return ccipOracle{oracle: oracle, pluginType: PluginTypeCCIPCommit, configChan: configChan}, nil
}

type diffResult struct {
	added   map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo
	removed map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo
	updated map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo
}

func diff(
	capabilityVersion,
	capabilityLabelledName string,
	oldState,
	newState RegistryState,
) (diffResult, error) {
	ccipCapability, err := checkCapabilityPresence(capabilityVersion, capabilityLabelledName, newState)
	if err != nil {
		return diffResult{}, err
	}

	newCCIPDONs, err := filterCCIPDONs(ccipCapability, newState)
	if err != nil {
		return diffResult{}, err
	}

	currCCIPDONs, err := filterCCIPDONs(ccipCapability, oldState)
	if err != nil {
		return diffResult{}, err
	}

	// compare curr with new and launch or update OCR instances as needed
	added, removed, updated, err := compareDONs(currCCIPDONs, newCCIPDONs)
	if err != nil {
		return diffResult{}, err
	}

	return diffResult{
		added:   added,
		removed: removed,
		updated: updated,
	}, nil
}

func compareDONs(
	currCCIPDONs,
	newCCIPDONs map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo,
) (
	added, removed, updated map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo,
	err error,
) {
	added = make(map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo)
	removed = make(map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo)
	updated = make(map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo)

	for id, don := range newCCIPDONs {
		if _, ok := currCCIPDONs[id]; !ok {
			added[id] = don
		} else {
			updated[id] = don
		}
	}

	for id, don := range currCCIPDONs {
		if _, ok := newCCIPDONs[id]; !ok {
			removed[id] = don
		}
	}

	return added, removed, updated, nil
}

func filterCCIPDONs(
	ccipCapability keystone_capability_registry.CapabilityRegistryCapability,
	state RegistryState,
) (map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo, error) {
	ccipDONs := make(map[uint32]keystone_capability_registry.CapabilityRegistryDONInfo)
	for _, don := range state.DONs {
		// CCIP DONs should only have one capability, CCIP.
		var found bool
		for _, donCapabilities := range don.CapabilityConfigurations {
			if donCapabilities.CapabilityId == hashedCapabilityId(ccipCapability.Version, ccipCapability.LabelledName) {
				ccipDONs[don.Id] = don
				found = true
			}
		}
		if found && len(don.CapabilityConfigurations) > 1 {
			return nil, fmt.Errorf("found more than one capability (actual: %d) in the CCIP DON %d",
				len(don.CapabilityConfigurations), don.Id)
		}
	}

	return ccipDONs, nil
}

func checkCapabilityPresence(
	capabilityVersion,
	capabilityLabelledName string,
	state RegistryState,
) (keystone_capability_registry.CapabilityRegistryCapability, error) {
	// Sanity check to make sure the capability registry has the capability we are looking for.
	var ccipCapability keystone_capability_registry.CapabilityRegistryCapability
	for _, capability := range state.Capabilities {
		if string(capability.LabelledName[:]) == capabilityLabelledName &&
			string(capability.Version[:]) == capabilityVersion {
			ccipCapability = capability
			break
		}
	}

	if ccipCapability.LabelledName == [32]byte{} {
		return keystone_capability_registry.CapabilityRegistryCapability{},
			fmt.Errorf("unable to find capability with name %s and version %s in capability registry state",
				capabilityLabelledName, capabilityVersion)
	}

	return ccipCapability, nil
}

func hashedCapabilityId(capabilityVersion, capabilityLabelledName [32]byte) (r [32]byte) {
	h := crypto.Keccak256(capabilityVersion[:], capabilityLabelledName[:])
	copy(r[:], h)
	return r
}
