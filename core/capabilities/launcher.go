package capabilities

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

var defaultStreamConfig = p2ptypes.StreamConfig{
	IncomingMessageBufferSize: 500,
	OutgoingMessageBufferSize: 500,
	MaxMessageLenBytes:        500000, // 500 KB;  max capacity = 500 * 500000 = 250 MB
	MessageRateLimiter: ragep2p.TokenBucketParams{
		Rate:     100.0,
		Capacity: 500,
	},
	BytesRateLimiter: ragep2p.TokenBucketParams{
		Rate:     5000000.0, // 5 MB/s
		Capacity: 10000000,  // 10 MB
	},
}

type launcher struct {
	services.StateMachine
	lggr                logger.Logger
	myPeerID            p2ptypes.PeerID
	peerWrapper         p2ptypes.PeerWrapper
	dispatcher          remotetypes.Dispatcher
	cachedShims         cachedShims
	registry            *Registry
	subServices         []services.Service
	workflowDonNotifier donNotifier
	don2donSharedPeer   p2ptypes.SharedPeer
	p2pStreamConfig     p2ptypes.StreamConfig
}

// For V2 capabilities, shims are created once and their config is updated dynamically.
type cachedShims struct {
	combinedClients    map[string]remote.CombinedClient
	triggerSubscribers map[string]remote.TriggerSubscriber
	triggerPublishers  map[string]remote.TriggerPublisher
	executableClients  map[string]executable.Client
	executableServers  map[string]executable.Server
}

func shimKey(capID string, donID uint32, method string) string {
	return fmt.Sprintf("%s:%d:%s", capID, donID, method)
}

type donNotifier interface {
	NotifyDonSet(don capabilities.DON)
}

// TODO: add metric handler and instrument all the internal log.Error calls

// NewLauncher creates a new capabilities launcher.
// If peerWrapper is nil, no p2p connections will be managed by the launcher.
// If don2donSharedPeer is nil, no DON-to-DON connections will be managed by the launcher.
func NewLauncher(
	lggr logger.Logger,
	peerWrapper p2ptypes.PeerWrapper,
	don2donSharedPeer p2ptypes.SharedPeer,
	streamConfig config.StreamConfig,
	dispatcher remotetypes.Dispatcher,
	registry *Registry,
	workflowDonNotifier donNotifier,
) *launcher {
	p2pStreamConfig := defaultStreamConfig
	if streamConfig != nil {
		p2pStreamConfig.IncomingMessageBufferSize = streamConfig.IncomingMessageBufferSize()
		p2pStreamConfig.OutgoingMessageBufferSize = streamConfig.OutgoingMessageBufferSize()
		p2pStreamConfig.MaxMessageLenBytes = streamConfig.MaxMessageLenBytes()
		p2pStreamConfig.MessageRateLimiter = ragep2p.TokenBucketParams{
			Rate:     streamConfig.MessageRateLimiterRate(),
			Capacity: streamConfig.MessageRateLimiterCapacity(),
		}
		p2pStreamConfig.BytesRateLimiter = ragep2p.TokenBucketParams{
			Rate:     streamConfig.BytesRateLimiterRate(),
			Capacity: streamConfig.BytesRateLimiterCapacity(),
		}
	}
	return &launcher{
		lggr:        logger.Named(lggr, "CapabilitiesLauncher"),
		peerWrapper: peerWrapper,
		dispatcher:  dispatcher,
		cachedShims: cachedShims{
			combinedClients:    make(map[string]remote.CombinedClient),
			triggerSubscribers: make(map[string]remote.TriggerSubscriber),
			triggerPublishers:  make(map[string]remote.TriggerPublisher),
			executableClients:  make(map[string]executable.Client),
			executableServers:  make(map[string]executable.Server),
		},
		registry:            registry,
		subServices:         []services.Service{},
		workflowDonNotifier: workflowDonNotifier,
		don2donSharedPeer:   don2donSharedPeer,
		p2pStreamConfig:     p2pStreamConfig,
	}
}

// Maintain only necessary Don2Don connections:
//   - Workflow DONs connect only to other DONs that have at least one remote capability
//   - Capability DONs connect only to workflow DONs
//
// Returns boolean as:
//   - true: filter out
//   - false: keep
func filterDon2Don(
	lggr logger.Logger,
	belongsToACapabilityDON bool,
	belongsToAWorkflowDON bool,
	candidatePeerDON registrysyncer.DON,
) bool {
	// Below logic is based on identification who is who using a workflow acceptance flag
	// and does it support any capabilities
	candidatePeerBelongsToWorkflowDON := candidatePeerDON.AcceptsWorkflows
	candidatePeerBelongsToCapabilityDON := len(candidatePeerDON.CapabilityConfigurations) > 0

	// We identify few cases from the perspective of the node:
	if belongsToACapabilityDON && belongsToAWorkflowDON {
		// as both workflow & capability DON let's just connect to anything
		return false // keep
	}
	if !belongsToACapabilityDON && !belongsToAWorkflowDON {
		// as none of workflow & capability DON don't use bandwidth
		lggr.Warn("filterDon2Don: node does not belong to workflow or capability DON; misconfiguration")
		return true // filter out
	}
	if belongsToAWorkflowDON && !candidatePeerBelongsToCapabilityDON {
		lggr.Debugw(
			"filterDon2Don: as a workflow DON my peers should be only capability DONs - filtering out",
			"DON.ID",
			candidatePeerDON.ID,
		)
		return true // filter out
	}
	if belongsToACapabilityDON && !candidatePeerBelongsToWorkflowDON {
		lggr.Debugw(
			"filterDon2Don: as a capability DON my peers should only be workflow DONs - filtering out",
			"DON.ID",
			candidatePeerDON.ID,
		)
		return true // filter out
	}
	return false // keep
}

func (w *launcher) peers(
	belongsToACapabilityDON bool,
	belongsToAWorkflowDON bool,
	isBootstrap bool,
	localRegistry *registrysyncer.LocalRegistry,
) map[ragetypes.PeerID]p2ptypes.StreamConfig {
	allPeers := make(map[ragetypes.PeerID]p2ptypes.StreamConfig)
	for _, id := range w.allDONs(localRegistry) {
		candidatePeerDON := localRegistry.IDsToDONs[id]
		if !candidatePeerDON.IsPublic {
			w.lggr.Debugw("skipping non-public DON for peer connections", "DON.ID", candidatePeerDON.ID)
			continue
		}
		if !isBootstrap && filterDon2Don(w.lggr, belongsToACapabilityDON, belongsToAWorkflowDON, candidatePeerDON) {
			continue
		}
		for _, nid := range candidatePeerDON.Members {
			allPeers[nid] = defaultStreamConfig
		}
	}
	return allPeers
}

func (w *launcher) publicDONs(
	allDONIDs []registrysyncer.DonID,
	localRegistry *registrysyncer.LocalRegistry,
) []registrysyncer.DON {
	publicDONs := make([]registrysyncer.DON, 0)
	for _, id := range allDONIDs {
		candidatePeerDON := localRegistry.IDsToDONs[id]
		if !candidatePeerDON.IsPublic {
			continue
		}
		publicDONs = append(publicDONs, candidatePeerDON)
	}
	return publicDONs
}

func (w *launcher) allDONs(localRegistry *registrysyncer.LocalRegistry) []registrysyncer.DonID {
	allDONIDs := make([]registrysyncer.DonID, 0)
	for id, don := range localRegistry.IDsToDONs {
		if len(don.Members) > 0 {
			// only non-empty DONs
			allDONIDs = append(allDONIDs, id)
		}
	}
	slices.Sort(allDONIDs) // ensure deterministic order
	return allDONIDs
}

func (w *launcher) Start(ctx context.Context) error {
	if w.peerWrapper != nil && w.peerWrapper.GetPeer() != nil {
		w.myPeerID = w.peerWrapper.GetPeer().ID()
		return nil
	}
	if w.don2donSharedPeer != nil {
		w.myPeerID = w.don2donSharedPeer.ID()
		return nil
	}
	return errors.New("could not get peer ID from any source")
}

func (w *launcher) Close() error {
	for _, s := range w.subServices {
		if err := s.Close(); err != nil {
			w.lggr.Errorw("failed to close a sub-service", "name", s.Name(), "error", err)
		}
	}
	if w.peerWrapper != nil {
		return w.peerWrapper.GetPeer().UpdateConnections(map[ragetypes.PeerID]p2ptypes.StreamConfig{})
	}
	return nil
}

func (w *launcher) Ready() error {
	return nil
}

func (w *launcher) HealthReport() map[string]error {
	return nil
}

func (w *launcher) Name() string {
	return w.lggr.Name()
}

func (w *launcher) donPairsToUpdate(myID ragetypes.PeerID, localRegistry *registrysyncer.LocalRegistry) []p2ptypes.DonPair {
	allDONIds := w.allDONs(localRegistry)
	donPairs := []p2ptypes.DonPair{}
	isBootstrap := w.don2donSharedPeer.IsBootstrap()
	for i, idA := range allDONIds {
		donA := localRegistry.IDsToDONs[idA]
		nodeBelongsToA := slices.Contains(donA.Members, myID)
		for _, idB := range allDONIds[i+1:] {
			donB := localRegistry.IDsToDONs[idB]
			pairAB := p2ptypes.DonPair{donA.DON, donB.DON}
			nodeBelongsToB := slices.Contains(donB.Members, myID)
			if !nodeBelongsToA && !nodeBelongsToB && !isBootstrap { // bootstrap adds all allowed DON pairs
				continue // skip if node doesn't belong to either DON
			}
			if donA.AcceptsWorkflows && len(donB.CapabilityConfigurations) > 0 || // add DON pair if A is workflow and B is capability
				donB.AcceptsWorkflows && len(donA.CapabilityConfigurations) > 0 { // add DON pair if B is workflow and A is capability
				if !donFamiliesOverlap(donA.Families, donB.Families) {
					w.lggr.Debugw("donPairsToUpdate: filtering out DON pair due to family mismatch", "donA.ID", donA.ID, "donB.ID", donB.ID, "donA.Families", donA.Families, "donB.Families", donB.Families)
					continue
				}
				donPairs = append(donPairs, pairAB)
			}
		}
	}
	return donPairs
}

func (w *launcher) OnNewRegistry(ctx context.Context, localRegistry *registrysyncer.LocalRegistry) error {
	w.lggr.Debug("CapabilitiesLauncher triggered...")
	w.registry.SetLocalRegistry(localRegistry)

	allDONIDs := w.allDONs(localRegistry)
	w.lggr.Debugw("All DONs in the local registry", "allDONIDs", allDONIDs)

	// Let's start by identifying public DONs
	publicDONs := w.publicDONs(allDONIDs, localRegistry)

	// Next, we need to split the DONs into the following:
	// - workflow DONs the current node is a part of.
	// These will need remote shims to all remote capabilities on other DONs.
	//
	// We'll also construct a set to record what DONs the current node is a part of,
	// regardless of any modifiers (public/acceptsWorkflows etc).
	myWorkflowDONs := []registrysyncer.DON{}
	remoteWorkflowDONs := []registrysyncer.DON{}
	myDONs := map[uint32]bool{}
	myDONFamiliesSet := map[string]bool{}
	myDONFamilies := []string{}
	for _, id := range allDONIDs {
		d := localRegistry.IDsToDONs[id]
		for _, peerID := range d.Members {
			if peerID == w.myPeerID {
				myDONs[d.ID] = true
				for _, family := range d.Families {
					myDONFamiliesSet[family] = true
				}
			}
		}

		if d.AcceptsWorkflows {
			if myDONs[d.ID] {
				myWorkflowDONs = append(myWorkflowDONs, d)
			} else {
				remoteWorkflowDONs = append(remoteWorkflowDONs, d)
			}
		}
	}
	for family := range myDONFamiliesSet {
		myDONFamilies = append(myDONFamilies, family)
	}
	w.lggr.Debugw("Found my DON families", "count", len(myDONFamilies), "myDONFamilies", myDONFamilies)
	w.lggr.Debugw("Found my workflow DONs", "count", len(myWorkflowDONs), "myWorkflowDONs", myWorkflowDONs)
	w.lggr.Debugw("Found all remote workflow DONs", "count", len(remoteWorkflowDONs), "remoteWorkflowDONs", remoteWorkflowDONs)

	// Capability DONs (with IsPublic = true) the current node is a part of.
	// These need server-side shims to expose my own capabilities externally.
	myCapabilityDONs := []registrysyncer.DON{}
	remoteCapabilityDONs := []registrysyncer.DON{}
	for _, d := range publicDONs {
		if len(d.CapabilityConfigurations) > 0 {
			if myDONs[d.ID] {
				myCapabilityDONs = append(myCapabilityDONs, d)
			} else {
				remoteCapabilityDONs = append(remoteCapabilityDONs, d)
			}
		}
	}
	w.lggr.Debugw("Found my capability DONs", "count", len(myCapabilityDONs), "myCapabilityDONs", myCapabilityDONs)
	w.lggr.Debugw("Found all remote capability DONs", "count", len(remoteCapabilityDONs), "remoteCapabilityDONs", remoteCapabilityDONs)

	if len(myDONFamilies) > 0 {
		remoteWorkflowDONs = filterDONsByFamilies(remoteWorkflowDONs, myDONFamilies)
		remoteCapabilityDONs = filterDONsByFamilies(remoteCapabilityDONs, myDONFamilies)
		w.lggr.Debugw("Filtered remote workflow DONs to my families", "count", len(remoteWorkflowDONs), "remoteWorkflowDONs", remoteWorkflowDONs)
		w.lggr.Debugw("Filtered remote capability DONs to my families", "count", len(remoteCapabilityDONs), "remoteCapabilityDONs", remoteCapabilityDONs)
	} else {
		// legacy / Keystone setting
		w.lggr.Debug("My node doesn't belong to any DON families. No filtering will be applied.")
	}

	belongsToAWorkflowDON := len(myWorkflowDONs) > 0
	if belongsToAWorkflowDON {
		myDON := myWorkflowDONs[0]

		// NOTE: this is enforced on-chain and so should never happen.
		if len(myWorkflowDONs) > 1 {
			return errors.New("invariant violation: node is part of more than one workflowDON")
		}

		w.lggr.Debug("Notifying DON set...")
		w.workflowDonNotifier.NotifyDonSet(myDON.DON)

		for _, rcd := range remoteCapabilityDONs {
			w.addRemoteCapabilities(ctx, myDON, rcd, localRegistry)
		}
	}

	belongsToACapabilityDON := len(myCapabilityDONs) > 0
	if belongsToACapabilityDON {
		for _, myDON := range myCapabilityDONs {
			w.serveCapabilities(ctx, w.myPeerID, myDON, localRegistry, remoteWorkflowDONs)
		}
	}

	// Lastly, we identify peers to connect to, based on their DONs functions
	w.lggr.Debugw("Updating peer connections", "peerWrapperEnabled", w.peerWrapper != nil, "don2donSharedPeerEnabled", w.don2donSharedPeer != nil)
	if w.peerWrapper != nil { // legacy / Keystone setting
		peer := w.peerWrapper.GetPeer()
		myPeers := w.peers(belongsToACapabilityDON, belongsToAWorkflowDON, peer.IsBootstrap(), localRegistry)
		err := peer.UpdateConnections(myPeers)
		if err != nil {
			return fmt.Errorf("failed to update peer connections: %w", err)
		}
	}
	if w.don2donSharedPeer != nil {
		donPairs := w.donPairsToUpdate(w.myPeerID, localRegistry)
		err := w.don2donSharedPeer.UpdateConnectionsByDONs(ctx, donPairs, defaultStreamConfig)
		if err != nil {
			return fmt.Errorf("failed to update peer connections: %w", err)
		}
	}
	return nil
}

func filterDONsByFamilies(donList []registrysyncer.DON, myDONFamilies []string) []registrysyncer.DON {
	filteredDONs := []registrysyncer.DON{}
	for _, d := range donList {
		if donFamiliesOverlap(d.Families, myDONFamilies) {
			filteredDONs = append(filteredDONs, d)
		}
	}
	return filteredDONs
}

func donFamiliesOverlap(donA []string, donB []string) bool {
	if len(donA) == 0 && len(donB) == 0 {
		return true // legacy setting with empty families - ignore filtering
	}
	for _, family := range donA {
		if slices.Contains(donB, family) {
			return true
		}
	}
	return false
}

// addRemoteCapabilities adds remote capabilities from a remote DON to the local node,
// allowing the local node to use these capabilities in its workflows.
// it is best effort to ensure that valid capabilities are added even if some fail
func (w *launcher) addRemoteCapabilities(ctx context.Context, myDON registrysyncer.DON, remoteDON registrysyncer.DON, localRegistry *registrysyncer.LocalRegistry) {
	for cid, c := range remoteDON.CapabilityConfigurations {
		err := w.addRemoteCapability(ctx, cid, c, myDON, remoteDON, localRegistry)
		if err != nil {
			// TODO CRE-1021 metrics for failures
			w.lggr.Errorw("failed to add remote capability ", "myDON", myDON, "remoteDON", remoteDON, "capabilityID", cid, "err", err)
		}
	}
}

func (w *launcher) addRemoteCapability(ctx context.Context, cid string, c registrysyncer.CapabilityConfiguration, myDON registrysyncer.DON, remoteDON registrysyncer.DON, localRegistry *registrysyncer.LocalRegistry) error {
	capability, ok := localRegistry.IDsToCapabilities[cid]
	if !ok {
		return fmt.Errorf("could not find capability matching id %s", cid)
	}

	capabilityConfig, err := c.Unmarshal()
	if err != nil {
		return fmt.Errorf("could not unmarshal capability config for id %s: %w", cid, err)
	}

	methodConfig := capabilityConfig.CapabilityMethodConfig
	if methodConfig != nil { // v2 capability - handle via CombinedClient
		errAdd := w.addRemoteCapabilityV2(ctx, capability.ID, methodConfig, myDON, remoteDON)
		if errAdd != nil {
			return fmt.Errorf("failed to add remote v2 capability %s: %w", capability.ID, errAdd)
		}
	}

	switch capability.CapabilityType {
	case capabilities.CapabilityTypeTrigger:
		newTriggerFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
			var aggregator remotetypes.Aggregator
			switch {
			case strings.HasPrefix(info.ID, "streams-trigger"):
				v := info.ID[strings.LastIndexAny(info.ID, "@")+1:] // +1 to skip the @; also gracefully handle the case where there is no @ (which should not happen)
				version, err := semver.NewVersion(v)
				if err != nil {
					return nil, fmt.Errorf("could not extract version from %s (%s): %w", info.ID, v, err)
				}
				switch version.Major() {
				case 1: // legacy streams trigger
					codec := streams.NewCodec(w.lggr)

					signers, err := signersFor(remoteDON, localRegistry)
					if err != nil {
						return nil, fmt.Errorf("failed to get signers for streams-trigger: %w", err)
					}

					// deprecated pre-LLO Mercury aggregator
					aggregator = triggers.NewMercuryRemoteAggregator(
						codec,
						signers,
						int(remoteDON.F+1),
						info.ID,
						w.lggr,
					)
				case 2: // LLO
					// TODO: add a flag in capability onchain config to indicate whether it's OCR based
					// the "SignedReport" aggregator is generic
					signers, err := signersFor(remoteDON, localRegistry)
					if err != nil {
						return nil, fmt.Errorf("failed to get signers for llo-trigger: %w", err)
					}

					const maxAgeSec = 120 // TODO move to capability onchain config
					aggregator = aggregation.NewSignedReportRemoteAggregator(
						signers,
						int(remoteDON.F+1),
						info.ID,
						maxAgeSec,
						w.lggr,
					)
				default:
					return nil, fmt.Errorf("unsupported stream trigger %s", info.ID)
				}
			default:
				aggregator = aggregation.NewDefaultModeAggregator(uint32(remoteDON.F) + 1)
			}

			shimKey := shimKey(capability.ID, remoteDON.ID, "") // empty method name for V1
			triggerCap, alreadyExists := w.cachedShims.triggerSubscribers[shimKey]
			if !alreadyExists {
				triggerCap = remote.NewTriggerSubscriber(
					capability.ID,
					"", // empty method name for v1
					w.dispatcher,
					w.lggr,
				)
				w.cachedShims.triggerSubscribers[shimKey] = triggerCap
			}
			if errCfg := triggerCap.SetConfig(capabilityConfig.RemoteTriggerConfig, info, myDON.ID, remoteDON.DON, aggregator); errCfg != nil {
				return nil, fmt.Errorf("failed to set trigger config: %w", errCfg)
			}
			return triggerCap.(capabilityService), nil
		}
		err := w.addToRegistryAndSetDispatcher(ctx, capability, remoteDON, newTriggerFn)
		if err != nil {
			return fmt.Errorf("failed to add trigger shim: %w", err)
		}
	case capabilities.CapabilityTypeAction:
		newActionFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
			shimKey := shimKey(capability.ID, remoteDON.ID, "") // empty method name for V1
			execCap, alreadyExists := w.cachedShims.executableClients[shimKey]
			if !alreadyExists {
				execCap = executable.NewClient(
					info.ID,
					"", // empty method name for v1
					w.dispatcher,
					w.lggr,
				)
				w.cachedShims.executableClients[shimKey] = execCap
			}
			// V1 capabilities read transmission schedule from every request
			if errCfg := execCap.SetConfig(info, myDON.DON, defaultTargetRequestTimeout, nil); errCfg != nil {
				return nil, fmt.Errorf("failed to set trigger config: %w", errCfg)
			}
			return execCap.(capabilityService), nil
		}

		err := w.addToRegistryAndSetDispatcher(ctx, capability, remoteDON, newActionFn)
		if err != nil {
			return fmt.Errorf("failed to add action shim: %w", err)
		}
	case capabilities.CapabilityTypeConsensus:
		// nothing to do; we don't support remote consensus capabilities for now
	case capabilities.CapabilityTypeTarget:
		newTargetFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
			shimKey := shimKey(capability.ID, remoteDON.ID, "") // empty method name for V1
			execCap, alreadyExists := w.cachedShims.executableClients[shimKey]
			if !alreadyExists {
				execCap = executable.NewClient(
					info.ID,
					"", // empty method name for v1
					w.dispatcher,
					w.lggr,
				)
				w.cachedShims.executableClients[shimKey] = execCap
			}
			// V1 capabilities read transmission schedule from every request
			if errCfg := execCap.SetConfig(info, myDON.DON, defaultTargetRequestTimeout, nil); errCfg != nil {
				return nil, fmt.Errorf("failed to set trigger config: %w", errCfg)
			}
			return execCap.(capabilityService), nil
		}

		err := w.addToRegistryAndSetDispatcher(ctx, capability, remoteDON, newTargetFn)
		if err != nil {
			return fmt.Errorf("failed to add target shim: %w", err)
		}
	default:
		w.lggr.Warnf("unknown capability type, skipping configuration: %+v", capability)
	}

	return nil
}

type capabilityService interface {
	capabilities.BaseCapability
	remotetypes.Receiver
	services.Service
}

func (w *launcher) addToRegistryAndSetDispatcher(ctx context.Context, capability registrysyncer.Capability, don registrysyncer.DON, newCapFn func(info capabilities.CapabilityInfo) (capabilityService, error)) error {
	capabilityID := capability.ID
	info, err := capabilities.NewRemoteCapabilityInfo(
		capabilityID,
		capability.CapabilityType,
		"Remote Capability for "+capabilityID,
		&don.DON,
	)
	if err != nil {
		return fmt.Errorf("failed to create remote capability info: %w", err)
	}
	w.lggr.Debugw("Adding remote capability to registry", "id", info.ID, "don", info.DON)
	cp, err := newCapFn(info)
	if err != nil {
		return fmt.Errorf("failed to instantiate capability %q: %w", capabilityID, err)
	}

	err = w.registry.Add(ctx, cp)
	if err != nil {
		// If the capability already exists, then it's either local
		// or we've handled this in a previous syncer iteration,
		// let's skip and move on to other capabilities.
		if errors.Is(err, registry.ErrCapabilityAlreadyExists) {
			return nil
		}

		return fmt.Errorf("failed to add capability to registry: %w", err)
	}

	err = w.dispatcher.SetReceiver(
		capabilityID,
		don.ID,
		cp,
	)
	if err != nil {
		return err
	}
	w.lggr.Debugw("Setting receiver for capability", "id", capabilityID, "donID", don.ID)
	err = cp.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start capability: %w", err)
	}
	w.subServices = append(w.subServices, cp)
	return nil
}

var (
	// TODO: make this configurable
	defaultTargetRequestTimeout                 = 8 * time.Minute
	defaultMaxParallelCapabilityExecuteRequests = uint32(1000)
)

// serveCapabilities exposes capabilities that are available on this node, as part of the given DON.
// It is best effort, ensuring that valid capabilities are exposed even if some fail
func (w *launcher) serveCapabilities(ctx context.Context, myPeerID p2ptypes.PeerID, don registrysyncer.DON, localRegistry *registrysyncer.LocalRegistry, remoteWorkflowDONs []registrysyncer.DON) {
	idsToDONs := map[uint32]capabilities.DON{}
	for _, d := range remoteWorkflowDONs {
		idsToDONs[d.ID] = d.DON
	}

	for cid, c := range don.CapabilityConfigurations {
		err := w.serveCapability(ctx, cid, c, myPeerID, don, localRegistry, idsToDONs)
		if err != nil {
			// TODO CRE-1021 metrics for failures
			w.lggr.Errorw("failed to serve capability", "myPeerID", myPeerID, "don", don, "capabilityID", cid, "err", err)
		}
	}
}

// serveCapability exposes a single capability.
// trigger capabilities are exposed via a TriggerPublisher local execution
// all other capabilities are exposed via an Executable for remote execution
func (w *launcher) serveCapability(ctx context.Context, cid string, c registrysyncer.CapabilityConfiguration, myPeerID p2ptypes.PeerID, don registrysyncer.DON, localRegistry *registrysyncer.LocalRegistry, idsToDONs map[uint32]capabilities.DON) error {
	capability, ok := localRegistry.IDsToCapabilities[cid]
	if !ok {
		return fmt.Errorf("could not find capability matching id %s", cid)
	}

	capabilityConfig, err := c.Unmarshal()
	if err != nil {
		return fmt.Errorf("could not unmarshal capability config for id %s: %w", cid, err)
	}

	methodConfig := capabilityConfig.CapabilityMethodConfig
	if methodConfig != nil { // v2 capability
		errExpose := w.exposeCapabilityV2(ctx, cid, methodConfig, myPeerID, don, idsToDONs)
		if errExpose != nil {
			return fmt.Errorf("failed to expose v2 capability remotely %s: %w", cid, errExpose)
		}
		return nil
	}

	switch capability.CapabilityType {
	case capabilities.CapabilityTypeTrigger:
		newTriggerPublisher := func(bc capabilities.BaseCapability, info capabilities.CapabilityInfo) (remotetypes.ReceiverService, error) {
			triggerCapability, ok := (bc).(capabilities.TriggerCapability)
			if !ok {
				return nil, errors.New("capability does not implement TriggerCapability")
			}
			shimKey := shimKey(capability.ID, don.ID, "") // empty method name for V1
			publisher, alreadyExists := w.cachedShims.triggerPublishers[shimKey]
			if !alreadyExists {
				publisher = remote.NewTriggerPublisher(
					capability.ID,
					"", // empty method name for v1
					w.dispatcher,
					w.lggr,
				)
				w.cachedShims.triggerPublishers[shimKey] = publisher
			}
			if errCfg := publisher.SetConfig(capabilityConfig.RemoteTriggerConfig, triggerCapability, don.DON, idsToDONs); errCfg != nil {
				return nil, fmt.Errorf("failed to set config for trigger publisher: %w", errCfg)
			}
			return publisher, nil
		}

		if err = w.addReceiver(ctx, capability, don, newTriggerPublisher); err != nil {
			return fmt.Errorf("failed to add server-side receiver for a trigger capability '%s' - it won't be exposed remotely: %w", cid, err)
		}
	case capabilities.CapabilityTypeAction:
		newActionServer := func(bc capabilities.BaseCapability, info capabilities.CapabilityInfo) (remotetypes.ReceiverService, error) {
			actionCapability, ok := (bc).(capabilities.ActionCapability) //nolint:staticcheck //SA1019
			if !ok {
				return nil, errors.New("capability does not implement ActionCapability")
			}
			shimKey := shimKey(capability.ID, don.ID, "") // empty method name for V1
			server, alreadyExists := w.cachedShims.executableServers[shimKey]
			if !alreadyExists {
				server = executable.NewServer(
					info.ID,
					"", // empty method name for v1
					myPeerID,
					w.dispatcher,
					w.lggr,
				)
				w.cachedShims.executableServers[shimKey] = server
			}

			remoteConfig := &capabilities.RemoteExecutableConfig{
				// deprecated defaults - v2 reads these from onchain config
				RequestTimeout:            defaultTargetRequestTimeout,
				ServerMaxParallelRequests: defaultMaxParallelCapabilityExecuteRequests,
			}
			if capabilityConfig.RemoteTargetConfig != nil {
				remoteConfig.RequestHashExcludedAttributes = capabilityConfig.RemoteTargetConfig.RequestHashExcludedAttributes
			}
			errCfg := server.SetConfig(
				remoteConfig,
				actionCapability,
				info,
				don.DON,
				idsToDONs,
				nil,
			)
			if errCfg != nil {
				return nil, fmt.Errorf("failed to set server config: %w", errCfg)
			}

			return server, nil
		}

		if err = w.addReceiver(ctx, capability, don, newActionServer); err != nil {
			return fmt.Errorf("failed to add action server-side receiver '%s' - it won't be exposed remotely: %w", cid, err)
		}
	case capabilities.CapabilityTypeConsensus:
		w.lggr.Debug("no remote client configured for capability type consensus, skipping configuration")
	case capabilities.CapabilityTypeTarget: // TODO: unify Target and Action into Executable
		newTargetServer := func(bc capabilities.BaseCapability, info capabilities.CapabilityInfo) (remotetypes.ReceiverService, error) {
			targetCapability, ok := (bc).(capabilities.TargetCapability) //nolint:staticcheck //SA1019
			if !ok {
				return nil, errors.New("capability does not implement TargetCapability")
			}

			shimKey := shimKey(capability.ID, don.ID, "") // empty method name for V1
			server, alreadyExists := w.cachedShims.executableServers[shimKey]
			if !alreadyExists {
				server = executable.NewServer(
					info.ID,
					"", // empty method name for v1
					myPeerID,
					w.dispatcher,
					w.lggr,
				)
				w.cachedShims.executableServers[shimKey] = server
			}

			remoteConfig := &capabilities.RemoteExecutableConfig{
				// deprecated defaults - v2 reads these from onchain config
				RequestTimeout:            defaultTargetRequestTimeout,
				ServerMaxParallelRequests: defaultMaxParallelCapabilityExecuteRequests,
			}
			if capabilityConfig.RemoteTargetConfig != nil {
				remoteConfig.RequestHashExcludedAttributes = capabilityConfig.RemoteTargetConfig.RequestHashExcludedAttributes
			}
			errCfg := server.SetConfig(
				remoteConfig,
				targetCapability,
				info,
				don.DON,
				idsToDONs,
				nil,
			)
			if errCfg != nil {
				return nil, fmt.Errorf("failed to set server config: %w", errCfg)
			}

			return server, nil
		}

		if err = w.addReceiver(ctx, capability, don, newTargetServer); err != nil {
			return fmt.Errorf("failed to add server-side receiver for a target capability '%s' - it won't be exposed remotely: %w", cid, err)
		}
	default:
		w.lggr.Warnf("unknown capability type, skipping configuration: %+v", capability)
	}
	return nil
}

func (w *launcher) addReceiver(ctx context.Context, capability registrysyncer.Capability, don registrysyncer.DON, newReceiverFn func(capability capabilities.BaseCapability, info capabilities.CapabilityInfo) (remotetypes.ReceiverService, error)) error {
	capID := capability.ID
	info, err := capabilities.NewRemoteCapabilityInfo(
		capID,
		capability.CapabilityType,
		"Remote Capability for "+capability.ID,
		&don.DON,
	)
	if err != nil {
		return fmt.Errorf("failed to instantiate remote capability for receiver: %w", err)
	}
	underlying, err := w.registry.Get(ctx, capability.ID)
	if err != nil {
		return fmt.Errorf("failed to get capability from registry: %w", err)
	}

	receiver, err := newReceiverFn(underlying, info)
	if err != nil {
		return fmt.Errorf("failed to instantiate receiver: %w", err)
	}

	w.lggr.Debugw("Enabling external access for capability", "id", capID, "donID", don.ID)
	err = w.dispatcher.SetReceiver(capID, don.ID, receiver)
	if errors.Is(err, remote.ErrReceiverExists) {
		// If a receiver already exists, let's log the error for debug purposes, but
		// otherwise short-circuit here. We've handled this capability in a previous iteration.
		w.lggr.Debugf("receiver already exists for cap ID %s and don ID %d: %s", capID, don.ID, err)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to set receiver: %w", err)
	}

	err = receiver.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start receiver: %w", err)
	}

	w.subServices = append(w.subServices, receiver)
	return nil
}

func signersFor(don registrysyncer.DON, localRegistry *registrysyncer.LocalRegistry) ([][]byte, error) {
	s := [][]byte{}
	for _, nodeID := range don.Members {
		node, ok := localRegistry.IDsToNodes[nodeID]
		if !ok {
			return nil, fmt.Errorf("could not find node for id %s", nodeID)
		}

		// NOTE: the capability registry stores signers as [32]byte,
		// but we only need the first [20], as the rest is padded.
		s = append(s, node.Signer[0:20])
	}

	return s, nil
}

// Add a V2 capability with multiple methods, using CombinedClient.
func (w *launcher) addRemoteCapabilityV2(ctx context.Context, capID string, methodConfig map[string]capabilities.CapabilityMethodConfig, myDON registrysyncer.DON, remoteDON registrysyncer.DON) error {
	info, err := capabilities.NewRemoteCapabilityInfo(
		capID,
		capabilities.CapabilityTypeCombined,
		"Remote Capability for "+capID,
		&remoteDON.DON,
	)
	if err != nil {
		return fmt.Errorf("failed to create remote capability info: %w", err)
	}

	cc, isNewCC := w.getCombinedClient(info)
	for method, config := range methodConfig {
		w.lggr.Infow("addRemoteCapabilityV2", "capID", capID, "method", method)
		if config.RemoteTriggerConfig == nil && config.RemoteExecutableConfig == nil {
			// TODO CRE-1021 metrics
			w.lggr.Errorw("no remote config found", "method", method, "capID", capID)
			continue
		}

		shimKey := shimKey(capID, remoteDON.ID, method)
		if config.RemoteTriggerConfig != nil { // trigger
			sub, alreadyExists := w.cachedShims.triggerSubscribers[shimKey]
			if !alreadyExists {
				sub = remote.NewTriggerSubscriber(capID, method, w.dispatcher, w.lggr)
				cc.SetTriggerSubscriber(method, sub)
				// add to cachedShims later, only after startNewShim succeeds
			}
			// TODO(CRE-590): add support for SignedReportAggregator (needed by LLO Streams Trigger V2)
			agg := aggregation.NewDefaultModeAggregator(config.RemoteTriggerConfig.MinResponsesToAggregate)
			if errCfg := sub.SetConfig(config.RemoteTriggerConfig, info, myDON.ID, remoteDON.DON, agg); errCfg != nil {
				return fmt.Errorf("failed to set trigger config: %w", errCfg)
			}

			if !alreadyExists {
				if err2 := w.startNewShim(ctx, sub.(remotetypes.ReceiverService), capID, remoteDON.ID, method); err2 != nil {
					// TODO CRE-1021 metrics
					w.lggr.Errorw("failed to start receiver", "capID", capID, "method", method, "error", err2)
					continue
				}
				w.cachedShims.triggerSubscribers[shimKey] = sub
				w.lggr.Infow("added new remote trigger subscriber", "capID", capID, "method", method)
			}
		} else { // executable
			client, alreadyExists := w.cachedShims.executableClients[shimKey]
			if !alreadyExists {
				client = executable.NewClient(info.ID, method, w.dispatcher, w.lggr)
				cc.SetExecutableClient(method, client)
				// add to cachedShims later, only after startNewShim succeeds
			}
			// Update existing client config
			transmissionConfig := &transmission.TransmissionConfig{
				Schedule:   transmission.EnumToString(config.RemoteExecutableConfig.TransmissionSchedule),
				DeltaStage: config.RemoteExecutableConfig.DeltaStage,
			}
			err := client.SetConfig(info, myDON.DON, config.RemoteExecutableConfig.RequestTimeout, transmissionConfig)
			if err != nil {
				w.lggr.Errorw("failed to update client config", "capID", capID, "method", method, "error", err)
				continue
			}

			if !alreadyExists {
				if err2 := w.startNewShim(ctx, client.(remotetypes.ReceiverService), capID, remoteDON.ID, method); err2 != nil {
					// TODO CRE-1021 metrics
					w.lggr.Errorw("failed to start receiver", "capID", capID, "method", method, "error", err2)
					continue
				}
				w.cachedShims.executableClients[shimKey] = client
				w.lggr.Infow("added new remote executable client", "capID", capID, "method", method)
			}
		}
	}

	if isNewCC { // add new CombinedClient to registry, only after all methods are configured
		if err2 := w.registry.Add(ctx, cc); err2 != nil {
			return fmt.Errorf("failed to add CombinedClient for capability %s to registry: %w", capID, err2)
		}
	}
	return nil
}

func (w *launcher) startNewShim(ctx context.Context, receiver remotetypes.ReceiverService, capID string, donID uint32, method string) error {
	w.lggr.Debugw("Starting new remote shim for capability method", "id", capID, "method", method, "donID", donID)
	if err := receiver.Start(ctx); err != nil {
		return fmt.Errorf("failed to start receiver for capability %s, method %s: %w", capID, method, err)
	}
	if err := w.dispatcher.SetReceiverForMethod(capID, donID, method, receiver); err != nil {
		_ = receiver.Close()
		return fmt.Errorf("failed to register receiver for capability %s, method %s: %w", capID, method, err)
	}
	w.subServices = append(w.subServices, receiver)
	w.lggr.Debugw("New remote shim started successfully for capability method", "id", capID, "method", method, "donID", donID)
	return nil
}

func (w *launcher) exposeCapabilityV2(ctx context.Context, capID string, methodConfig map[string]capabilities.CapabilityMethodConfig, myPeerID p2ptypes.PeerID, myDON registrysyncer.DON, idsToDONs map[uint32]capabilities.DON) error {
	info, err := capabilities.NewRemoteCapabilityInfo(
		capID,
		capabilities.CapabilityTypeCombined,
		"Remote Capability for "+capID,
		&myDON.DON,
	)
	if err != nil {
		return fmt.Errorf("failed to create remote capability info: %w", err)
	}
	underlying, err := w.registry.Get(ctx, capID)
	if err != nil {
		return fmt.Errorf("failed to get capability %s from registry: %w", capID, err)
	}
	for method, config := range methodConfig {
		if config.RemoteTriggerConfig != nil { // trigger
			underlyingTriggerCapability, ok := (underlying).(capabilities.TriggerCapability)
			if !ok {
				return fmt.Errorf("capability %s does not implement TriggerCapability", capID)
			}
			shimKey := shimKey(capID, myDON.ID, method)
			publisher, alreadyExists := w.cachedShims.triggerPublishers[shimKey]
			if !alreadyExists {
				publisher = remote.NewTriggerPublisher(
					capID,
					method,
					w.dispatcher,
					w.lggr,
				)
				// add to cachedShims later, only after startNewShim succeeds
			}
			if errCfg := publisher.SetConfig(config.RemoteTriggerConfig, underlyingTriggerCapability, myDON.DON, idsToDONs); errCfg != nil {
				return fmt.Errorf("failed to set config for trigger publisher: %w", errCfg)
			}

			if !alreadyExists {
				if err2 := w.startNewShim(ctx, publisher.(remotetypes.ReceiverService), capID, myDON.ID, method); err2 != nil {
					// TODO CRE-1021 metrics
					w.lggr.Errorw("failed to start receiver", "capID", capID, "method", method, "error", err2)
					continue
				}
				w.cachedShims.triggerPublishers[shimKey] = publisher
				w.lggr.Infow("added new remote trigger publisher", "capID", capID, "method", method)
			}
		} else { // executable
			underlyingExecutableCapability, ok := (underlying).(capabilities.ExecutableCapability)
			if !ok {
				return fmt.Errorf("capability %s does not implement ExecutableCapability", capID)
			}

			shimKey := shimKey(capID, myDON.ID, method)
			server, alreadyExists := w.cachedShims.executableServers[shimKey]
			if !alreadyExists {
				server = executable.NewServer(
					info.ID,
					method,
					myPeerID,
					w.dispatcher,
					w.lggr,
				)
				// add to cachedShims later, only after startNewShim succeeds
			}

			var requestHasher remotetypes.MessageHasher
			switch config.RemoteExecutableConfig.RequestHasherType {
			case capabilities.RequestHasherType_Simple:
				requestHasher = executable.NewSimpleHasher()
			case capabilities.RequestHasherType_WriteReportExcludeSignatures:
				requestHasher = executable.NewWriteReportExcludeSignaturesHasher()
			default:
				requestHasher = executable.NewSimpleHasher()
			}

			err := server.SetConfig(
				config.RemoteExecutableConfig,
				underlyingExecutableCapability,
				info,
				myDON.DON,
				idsToDONs,
				requestHasher,
			)
			if err != nil {
				return fmt.Errorf("failed to set server config: %w", err)
			}

			if !alreadyExists {
				if err2 := w.startNewShim(ctx, server.(remotetypes.ReceiverService), capID, myDON.ID, method); err2 != nil {
					// TODO CRE-1021 metrics
					w.lggr.Errorw("failed to start receiver", "capID", capID, "method", method, "error", err2)
					continue
				}
				w.cachedShims.executableServers[shimKey] = server
				w.lggr.Infow("added new remote execcutable server", "capID", capID, "method", method)
			}
		}
	}
	return nil
}

// retrieve or create a CombinedClient for the given capability
func (w *launcher) getCombinedClient(info capabilities.CapabilityInfo) (remote.CombinedClient, bool) {
	key := shimKey(info.ID, info.DON.ID, "") // empty method name - CombinedClient covers all methods
	cc, exists := w.cachedShims.combinedClients[key]
	if !exists { // create a new combined client and cache it
		cc = remote.NewCombinedClient(info)
		w.cachedShims.combinedClients[key] = cc
		return cc, true
	}
	return cc, false
}
