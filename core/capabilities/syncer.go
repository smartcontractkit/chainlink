package capabilities

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type reader interface {
	state(ctx context.Context) (state, error)
}

type registrySyncer struct {
	peerWrapper  p2ptypes.PeerWrapper
	registry     core.CapabilitiesRegistry
	dispatcher   remotetypes.Dispatcher
	stopCh       services.StopChan
	subServices  []services.Service
	networkSetup HardcodedDonNetworkSetup
	reader       reader

	wg   sync.WaitGroup
	lggr logger.Logger
}

var _ services.Service = &registrySyncer{}

var (
	defaultTickInterval = time.Duration(12 * time.Second)
)

var defaultStreamConfig = p2ptypes.StreamConfig{
	IncomingMessageBufferSize: 1000000,
	OutgoingMessageBufferSize: 1000000,
	MaxMessageLenBytes:        100000,
	MessageRateLimiter: ragep2p.TokenBucketParams{
		Rate:     100.0,
		Capacity: 1000,
	},
	BytesRateLimiter: ragep2p.TokenBucketParams{
		Rate:     100000.0,
		Capacity: 1000000,
	},
}

const maxRetryCount = 60

// RegistrySyncer updates local Registry to match its onchain counterpart
func NewRegistrySyncer(
	peerWrapper p2ptypes.PeerWrapper,
	registry core.CapabilitiesRegistry,
	dispatcher remotetypes.Dispatcher,
	lggr logger.Logger,
	networkSetup HardcodedDonNetworkSetup,
	relayer contractReaderFactory,
	registryAddress string,
) (*registrySyncer, error) {
	stopCh := make(services.StopChan)
	ctx, _ := stopCh.NewCtx()
	reader, err := newRemoteRegistryReader(ctx, relayer, registryAddress)
	if err != nil {
		return nil, err
	}

	return newRegistrySyncer(
		stopCh,
		peerWrapper,
		registry,
		dispatcher,
		lggr,
		networkSetup,
		reader,
	), nil
}

func newRegistrySyncer(
	stopCh services.StopChan,
	peerWrapper p2ptypes.PeerWrapper,
	registry core.CapabilitiesRegistry,
	dispatcher remotetypes.Dispatcher,
	lggr logger.Logger,
	networkSetup HardcodedDonNetworkSetup,
	reader reader,
) *registrySyncer {
	return &registrySyncer{
		stopCh:       stopCh,
		peerWrapper:  peerWrapper,
		registry:     registry,
		dispatcher:   dispatcher,
		networkSetup: networkSetup,
		lggr:         lggr,
		reader:       reader,
	}
}

func (s *registrySyncer) Start(ctx context.Context) error {
	s.wg.Add(2)
	go s.launch()
	go s.syncLoop()
	return nil
}

func (s *registrySyncer) syncLoop() {
	defer s.wg.Done()

	ctx, cancel := s.stopCh.NewCtx()
	defer cancel()

	ticker := time.NewTicker(defaultTickInterval)
	defer ticker.Stop()

	// Sync for a first time outside the loop; this means we'll start a remote
	// sync immediately once spinning up syncLoop, as by default a ticker will
	// fire for the first time at T+N, where N is the interval.
	s.lggr.Debug("starting initial sync with remote registry")
	err := s.sync(ctx)
	if err != nil {
		s.lggr.Errorw("failed to sync with remote registry", "error", err)
	}

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.lggr.Debug("starting regular sync with the remote registry")
			err := s.sync(ctx)
			if err != nil {
				s.lggr.Errorw("failed to sync with remote registry", "error", err)
			}
		}
	}
}

func (s *registrySyncer) sync(ctx context.Context) error {
	state, err := s.reader.state(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync with remote registry: %w", err)
	}

	// Let's start by updating the list of Peers
	// We do this by creating a new entry for each node belonging
	// to a public DON.
	// We also add the hardcoded peers determined by the NetworkSetup.
	allPeers := make(map[ragetypes.PeerID]p2ptypes.StreamConfig)
	// TODO: Remove this when we're no longer hard-coding
	// a `networkSetup`.
	for p, cfg := range s.networkSetup.allPeers {
		allPeers[p] = cfg
	}

	publicDONs := []kcr.CapabilityRegistryDONInfo{}
	for _, d := range state.DONs {
		if !d.IsPublic {
			continue
		}

		publicDONs = append(publicDONs, d)

		for _, nid := range d.NodeP2PIds {
			allPeers[nid] = defaultStreamConfig
		}
	}
	s.peerWrapper.GetPeer().UpdateConnections(allPeers)

	// Next, we need to split the DONs into the following:
	// - workflow DONs the current node is a part of.
	// These will need remote shims to all remote capabilities on other DONs.
	//
	// We'll also construct a set to record what DONs the current node is a part of,
	// regardless of any modifiers (public/accetptsWorkflows etc).
	myID := s.peerWrapper.GetPeer().ID()
	myWorkflowDONs := []kcr.CapabilityRegistryDONInfo{}
	remoteWorkflowDONs := []kcr.CapabilityRegistryDONInfo{}
	myDONs := map[uint32]bool{}
	for _, d := range state.DONs {
		for _, peerID := range d.NodeP2PIds {
			if peerID == myID {
				myDONs[d.Id] = true
			}
		}

		if d.AcceptsWorkflows {
			if myDONs[d.Id] {
				myWorkflowDONs = append(myWorkflowDONs, d)
			} else {
				remoteWorkflowDONs = append(remoteWorkflowDONs, d)
			}
		}

	}

	// - remote capability DONs (with IsPublic = true) the current node is a part of.
	// These need server-side shims.
	myCapabilityDONs := []kcr.CapabilityRegistryDONInfo{}
	remoteCapabilityDONs := []kcr.CapabilityRegistryDONInfo{}
	for _, d := range publicDONs {
		if len(d.CapabilityConfigurations) > 0 {
			if myDONs[d.Id] {
				myCapabilityDONs = append(myCapabilityDONs, d)
			} else {
				remoteCapabilityDONs = append(remoteCapabilityDONs, d)
			}
		}
	}

	// Now, if my node is a workflow DON, let's setup any shims
	// to external capabilities.
	if len(myWorkflowDONs) > 0 {
		myDON := myWorkflowDONs[0]

		// TODO: this is a bit nasty; figure out how best to handle this.
		if len(myWorkflowDONs) > 1 {
			s.lggr.Warn("node is part of more than one workflow DON; assigning first DON as caller")
		}

		for _, rcd := range remoteCapabilityDONs {
			err := s.addRemoteCapabilities(ctx, myDON, rcd, state)
			if err != nil {
				return err
			}
		}
	}

	// Finally, if I'm a capability DON, let's enable external access
	// to the capability.
	if len(myCapabilityDONs) > 0 {
		for _, mcd := range myCapabilityDONs {
			err := s.enableExternalAccess(ctx, myID, mcd, state, remoteWorkflowDONs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func signersFor(don kcr.CapabilityRegistryDONInfo, state state) ([][]byte, error) {
	s := [][]byte{}
	for _, nodeID := range don.NodeP2PIds {
		node, ok := state.IDsToNodes[nodeID]
		if !ok {
			return nil, fmt.Errorf("could not find node for id %s", nodeID)
		}

		s = append(s, node.Signer[:])
	}

	return s, nil
}

func toDONInfo(don kcr.CapabilityRegistryDONInfo) *capabilities.DON {
	peerIDs := []p2ptypes.PeerID{}
	for _, p := range don.NodeP2PIds {
		peerIDs = append(peerIDs, p)
	}

	return &capabilities.DON{
		ID:      fmt.Sprint(don.Id),
		Members: peerIDs,
		F:       don.F,
	}
}

func toCapabilityType(capabilityType uint8) capabilities.CapabilityType {
	switch capabilityType {
	case 0:
		return capabilities.CapabilityTypeTrigger
	case 1:
		return capabilities.CapabilityTypeAction
	case 2:
		return capabilities.CapabilityTypeConsensus
	case 3:
		return capabilities.CapabilityTypeTarget
	default:
		// Not found
		return capabilities.CapabilityType(-1)
	}

}

func (s *registrySyncer) addRemoteCapabilities(ctx context.Context, myDON kcr.CapabilityRegistryDONInfo, remoteDON kcr.CapabilityRegistryDONInfo, state state) error {
	for _, c := range remoteDON.CapabilityConfigurations {
		capability, ok := state.IDsToCapabilities[c.CapabilityId]
		if !ok {
			return fmt.Errorf("could not find capability matching id %s", c.CapabilityId)
		}

		switch toCapabilityType(capability.CapabilityType) {
		case capabilities.CapabilityTypeTrigger:
			newTriggerFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
				if !strings.HasPrefix(info.ID, "streams-trigger") {
					return nil, errors.New("not supported: trigger capability does not have id = streams-trigger")
				}

				codec := streams.NewCodec(s.lggr)

				signers, err := signersFor(remoteDON, state)
				if err != nil {
					return nil, err
				}

				aggregator := triggers.NewMercuryRemoteAggregator(
					codec,
					signers,
					int(remoteDON.F+1),
					s.lggr,
				)
				cfg := &remotetypes.RemoteTriggerConfig{}
				cfg.ApplyDefaults()
				err = proto.Unmarshal(c.Config, cfg)
				if err != nil {
					return nil, err
				}
				// TODO: We need to implement a custom, Mercury-specific
				// aggregator here, because there is no guarantee that
				// all trigger events in the workflow will have the same
				// payloads. As a workaround, we validate the signatures.
				// When this is solved, we can move to a generic aggregator
				// and remove this.
				triggerCap := remote.NewTriggerSubscriber(
					cfg,
					info,
					*toDONInfo(remoteDON),
					*toDONInfo(myDON),
					s.dispatcher,
					aggregator,
					s.lggr,
				)
				return triggerCap, nil
			}
			err := s.addToRegistryAndSetDispatcher(ctx, capability, remoteDON, newTriggerFn)
			if err != nil {
				return fmt.Errorf("failed to add trigger shim: %w", err)
			}
		case capabilities.CapabilityTypeAction:
			s.lggr.Warn("no remote client configured for capability type action, skipping configuration")
		case capabilities.CapabilityTypeConsensus:
			s.lggr.Warn("no remote client configured for capability type consensus, skipping configuration")
		case capabilities.CapabilityTypeTarget:
			newTargetFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
				client := target.NewClient(
					info,
					*toDONInfo(remoteDON),
					s.dispatcher,
					defaultTargetRequestTimeout,
					s.lggr,
				)
				return client, nil
			}

			err := s.addToRegistryAndSetDispatcher(ctx, capability, remoteDON, newTargetFn)
			if err != nil {
				return fmt.Errorf("failed to add target shim: %w", err)
			}
		default:
			s.lggr.Warnf("unknown capability type, skipping configuration: %+v", capability)
		}
	}
	return nil
}

type capabilityService interface {
	capabilities.BaseCapability
	remotetypes.Receiver
	services.Service
}

func (s *registrySyncer) addToRegistryAndSetDispatcher(ctx context.Context, capabilityInfo kcr.CapabilityRegistryCapability, don kcr.CapabilityRegistryDONInfo, newCapFn func(info capabilities.CapabilityInfo) (capabilityService, error)) error {
	fullCapID := fmt.Sprintf("%s@%s", capabilityInfo.LabelledName, capabilityInfo.Version)
	info, err := capabilities.NewRemoteCapabilityInfo(
		fullCapID,
		toCapabilityType(capabilityInfo.CapabilityType),
		fmt.Sprintf("Remote Capability for %s", fullCapID),
		toDONInfo(don),
	)
	if err != nil {
		return err
	}
	capability, err := newCapFn(info)
	if err != nil {
		return err
	}

	err = s.registry.Add(ctx, capability)
	if err != nil {
		// If the capability already exists, we've handled this in
		// a previous syncer iteration, let's skip and move on
		// to other capabilities.
		if errors.Is(err, ErrCapabilityAlreadyExists) {
			return nil
		}

		return err
	}

	err = s.dispatcher.SetReceiver(
		fullCapID,
		fmt.Sprint(don.Id),
		capability,
	)
	if err != nil {
		return err
	}
	err = capability.Start(ctx)
	if err != nil {
		return err
	}
	s.subServices = append(s.subServices, capability)
	return nil
}

var (
	defaultTargetRequestTimeout = time.Minute
)

func (s *registrySyncer) enableExternalAccess(ctx context.Context, myPeerID p2ptypes.PeerID, don kcr.CapabilityRegistryDONInfo, state state, remoteWorkflowDONs []kcr.CapabilityRegistryDONInfo) error {
	idsToDONs := map[string]capabilities.DON{}
	for _, d := range remoteWorkflowDONs {
		idsToDONs[fmt.Sprint(d.Id)] = *toDONInfo(d)
	}

	for _, c := range don.CapabilityConfigurations {
		capability, ok := state.IDsToCapabilities[c.CapabilityId]
		if !ok {
			return fmt.Errorf("could not find capability matching id %s", c.CapabilityId)
		}

		switch toCapabilityType(capability.CapabilityType) {
		case capabilities.CapabilityTypeTrigger:
			newTriggerPublisher := func(capability capabilities.BaseCapability, info capabilities.CapabilityInfo) (receiverService, error) {
				cfg := &remotetypes.RemoteTriggerConfig{}
				cfg.ApplyDefaults()
				err := proto.Unmarshal(c.Config, cfg)
				if err != nil {
					return nil, err
				}
				publisher := remote.NewTriggerPublisher(
					cfg,
					capability.(capabilities.TriggerCapability),
					info,
					*toDONInfo(don),
					idsToDONs,
					s.dispatcher,
					s.lggr,
				)
				return publisher, nil
			}

			err := s.addReceiver(ctx, capability, don, newTriggerPublisher)
			if err != nil {
				return fmt.Errorf("failed to add server-side receiver: %w", err)
			}
		case capabilities.CapabilityTypeAction:
			s.lggr.Warn("no remote client configured for capability type action, skipping configuration")
		case capabilities.CapabilityTypeConsensus:
			s.lggr.Warn("no remote client configured for capability type consensus, skipping configuration")
		case capabilities.CapabilityTypeTarget:
			newTargetServer := func(capability capabilities.BaseCapability, info capabilities.CapabilityInfo) (receiverService, error) {
				return target.NewServer(
					myPeerID,
					capability.(capabilities.TargetCapability),
					info,
					*toDONInfo(don),
					idsToDONs,
					s.dispatcher,
					defaultTargetRequestTimeout,
					s.lggr,
				), nil
			}

			err := s.addReceiver(ctx, capability, don, newTargetServer)
			if err != nil {
				return fmt.Errorf("failed to add server-side receiver: %w", err)
			}
		default:
			s.lggr.Warnf("unknown capability type, skipping configuration: %+v", capability)
		}
	}
	return nil
}

type receiverService interface {
	services.Service
	remotetypes.Receiver
}

func (s *registrySyncer) addReceiver(ctx context.Context, capability kcr.CapabilityRegistryCapability, don kcr.CapabilityRegistryDONInfo, newReceiverFn func(capability capabilities.BaseCapability, info capabilities.CapabilityInfo) (receiverService, error)) error {
	fullCapID := fmt.Sprintf("%s@%s", capability.LabelledName, capability.Version)
	info, err := capabilities.NewRemoteCapabilityInfo(
		fullCapID,
		toCapabilityType(capability.CapabilityType),
		fmt.Sprintf("Remote Capability for %s", fullCapID),
		toDONInfo(don),
	)
	if err != nil {
		return err
	}
	underlying, err := s.registry.Get(ctx, fullCapID)
	if err != nil {
		return err
	}

	receiver, err := newReceiverFn(underlying, info)
	if err != nil {
		return err
	}

	err = s.dispatcher.SetReceiver(fullCapID, fmt.Sprint(don.Id), receiver)
	if err != nil {
		return err
	}

	err = receiver.Start(ctx)
	if err != nil {
		return err
	}

	s.subServices = append(s.subServices, receiver)
	return nil
}

// NOTE: this implementation of the Syncer is temporary and will be replaced by one
// that reads the configuration from chain (KS-117).
func (s *registrySyncer) launch() {
	ctx, _ := s.stopCh.NewCtx()
	defer s.wg.Done()
	capId := "streams-trigger@1.0.0"
	triggerInfo, err := capabilities.NewRemoteCapabilityInfo(
		capId,
		capabilities.CapabilityTypeTrigger,
		"Remote Trigger",
		&s.networkSetup.TriggerCapabilityDonInfo,
	)
	if err != nil {
		s.lggr.Errorw("failed to create capability info for streams-trigger", "error", err)
		return
	}

	targetCapId := "write_ethereum-testnet-sepolia@1.0.0"
	targetInfo, err := capabilities.NewRemoteCapabilityInfo(
		targetCapId,
		capabilities.CapabilityTypeTarget,
		"Remote Target",
		&s.networkSetup.TargetCapabilityDonInfo,
	)
	if err != nil {
		s.lggr.Errorw("failed to create capability info for write_ethereum-testnet-sepolia", "error", err)
		return
	}

	myId := s.peerWrapper.GetPeer().ID()
	config := &remotetypes.RemoteTriggerConfig{
		RegistrationRefreshMs:   20000,
		RegistrationExpiryMs:    60000,
		MinResponsesToAggregate: uint32(s.networkSetup.TriggerCapabilityDonInfo.F) + 1,
	}
	err = s.peerWrapper.GetPeer().UpdateConnections(s.networkSetup.allPeers)
	if err != nil {
		s.lggr.Errorw("failed to update connections", "error", err)
		return
	}
	if s.networkSetup.IsWorkflowDon(myId) {
		s.lggr.Info("member of a workflow DON - starting remote subscribers")
		codec := streams.NewCodec(s.lggr)
		aggregator := triggers.NewMercuryRemoteAggregator(codec, hexStringsToBytes(s.networkSetup.triggerDonSigners), int(s.networkSetup.TriggerCapabilityDonInfo.F+1), s.lggr)
		triggerCap := remote.NewTriggerSubscriber(config, triggerInfo, s.networkSetup.TriggerCapabilityDonInfo, s.networkSetup.WorkflowsDonInfo, s.dispatcher, aggregator, s.lggr)
		err = s.registry.Add(ctx, triggerCap)
		if err != nil {
			s.lggr.Errorw("failed to add remote trigger capability to registry", "error", err)
			return
		}
		err = s.dispatcher.SetReceiver(capId, s.networkSetup.TriggerCapabilityDonInfo.ID, triggerCap)
		if err != nil {
			s.lggr.Errorw("workflow DON failed to set receiver for trigger", "capabilityId", capId, "donId", s.networkSetup.TriggerCapabilityDonInfo.ID, "error", err)
			return
		}
		s.subServices = append(s.subServices, triggerCap)

		s.lggr.Info("member of a workflow DON - starting remote targets")
		targetCap := target.NewClient(targetInfo, s.networkSetup.WorkflowsDonInfo, s.dispatcher, 60*time.Second, s.lggr)
		err = s.registry.Add(ctx, targetCap)
		if err != nil {
			s.lggr.Errorw("failed to add remote target capability to registry", "error", err)
			return
		}
		err = s.dispatcher.SetReceiver(targetCapId, s.networkSetup.TargetCapabilityDonInfo.ID, targetCap)
		if err != nil {
			s.lggr.Errorw("workflow DON failed to set receiver for target", "capabilityId", capId, "donId", s.networkSetup.TargetCapabilityDonInfo.ID, "error", err)
			return
		}
		s.subServices = append(s.subServices, targetCap)
	}
	if s.networkSetup.IsTriggerDon(myId) {
		s.lggr.Info("member of a capability DON - starting remote publishers")

		/*{
			// ---- This is for local tests only, until a full-blown Syncer is implemented
			// ---- Normally this is set up asynchronously (by the Relayer + job specs in Mercury's case)
			localTrigger := triggers.NewMercuryTriggerService(1000, s.lggr)
			mockMercuryDataProducer := NewMockMercuryDataProducer(localTrigger, s.lggr)
			err = s.registry.Add(ctx, localTrigger)
			if err != nil {
				s.lggr.Errorw("failed to add local trigger capability to registry", "error", err)
				return err
			}
			s.subServices = append(s.subServices, localTrigger)
			s.subServices = append(s.subServices, mockMercuryDataProducer)
			// ----
		}*/

		count := 0
		for {
			count++
			if count > maxRetryCount {
				s.lggr.Error("failed to get Streams Trigger from the Registry")
				return
			}
			underlying, err2 := s.registry.GetTrigger(ctx, capId)
			if err2 != nil {
				// NOTE: it's possible that the jobs are not launched yet at this moment.
				// If not found yet, Syncer won't add to Registry but retry on the next tick.
				s.lggr.Infow("trigger not found yet ...", "capabilityId", capId, "error", err2)
				time.Sleep(1 * time.Second)
				continue
			}
			workflowDONs := map[string]capabilities.DON{
				s.networkSetup.WorkflowsDonInfo.ID: s.networkSetup.WorkflowsDonInfo,
			}
			triggerCap := remote.NewTriggerPublisher(config, underlying, triggerInfo, s.networkSetup.TriggerCapabilityDonInfo, workflowDONs, s.dispatcher, s.lggr)
			err = s.dispatcher.SetReceiver(capId, s.networkSetup.TriggerCapabilityDonInfo.ID, triggerCap)
			if err != nil {
				s.lggr.Errorw("capability DON failed to set receiver", "capabilityId", capId, "donId", s.networkSetup.TriggerCapabilityDonInfo.ID, "error", err)
				return
			}
			s.subServices = append(s.subServices, triggerCap)
			break
		}
	}
	if s.networkSetup.IsTargetDon(myId) {
		s.lggr.Info("member of a target DON - starting remote shims")
		underlying, err2 := s.registry.GetTarget(ctx, targetCapId)
		if err2 != nil {
			s.lggr.Errorw("target not found yet", "capabilityId", targetCapId, "error", err2)
			return
		}
		workflowDONs := map[string]capabilities.DON{
			s.networkSetup.WorkflowsDonInfo.ID: s.networkSetup.WorkflowsDonInfo,
		}
		targetCap := target.NewServer(myId, underlying, targetInfo, *targetInfo.DON, workflowDONs, s.dispatcher, 60*time.Second, s.lggr)
		err = s.dispatcher.SetReceiver(targetCapId, s.networkSetup.TargetCapabilityDonInfo.ID, targetCap)
		if err != nil {
			s.lggr.Errorw("capability DON failed to set receiver", "capabilityId", capId, "donId", s.networkSetup.TargetCapabilityDonInfo.ID, "error", err)
			return
		}
		s.subServices = append(s.subServices, targetCap)
	}
	// NOTE: temporary service start - should be managed by capability creation
	for _, srv := range s.subServices {
		err = srv.Start(ctx)
		if err != nil {
			s.lggr.Errorw("failed to start remote trigger caller", "error", err)
			return
		}
	}
	s.lggr.Info("registry syncer started")
}

func (s *registrySyncer) Close() error {
	close(s.stopCh)
	s.wg.Wait()
	for _, subService := range s.subServices {
		err := subService.Close()
		if err != nil {
			s.lggr.Errorw("failed to close a sub-service", "name", subService.Name(), "error", err)
		}
	}
	return s.peerWrapper.GetPeer().UpdateConnections(map[ragetypes.PeerID]p2ptypes.StreamConfig{})
}

func (s *registrySyncer) Ready() error {
	return nil
}

func (s *registrySyncer) HealthReport() map[string]error {
	return nil
}

func (s *registrySyncer) Name() string {
	return "RegistrySyncer"
}

// HardcodedDonNetworkSetup is a temporary setup for testing purposes
type HardcodedDonNetworkSetup struct {
	workflowDonPeers  []string
	triggerDonPeers   []string
	targetDonPeers    []string
	triggerDonSigners []string
	allPeers          map[ragetypes.PeerID]p2ptypes.StreamConfig

	WorkflowsDonInfo         capabilities.DON
	TriggerCapabilityDonInfo capabilities.DON
	TargetCapabilityDonInfo  capabilities.DON
}

func NewHardcodedDonNetworkSetup() (HardcodedDonNetworkSetup, error) {
	result := HardcodedDonNetworkSetup{}

	result.workflowDonPeers = []string{
		"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N",
		"12D3KooWG1AyvwmCpZ93J8pBQUE1SuzrjDXnT4BeouncHR3jWLCG",
		"12D3KooWGeUKZBRMbx27FUTgBwZa9Ap9Ym92mywwpuqkEtz8XWyv",
		"12D3KooW9zYWQv3STmDeNDidyzxsJSTxoCTLicafgfeEz9nhwhC4",
		"12D3KooWG1AeBnSJH2mdcDusXQVye2jqodZ6pftTH98HH6xvrE97",
		"12D3KooWBf3PrkhNoPEmp7iV291YnPuuTsgEDHTscLajxoDvwHGA",
		"12D3KooWP3FrMTFXXRU2tBC8aYvEBgUX6qhcH9q2JZCUi9Wvc2GX",
	}
	result.triggerDonPeers = []string{
		"12D3KooWBaiTbbRwwt2fbNifiL7Ew9tn3vds9AJE3Nf3eaVBX36m",
		"12D3KooWS7JSY9fzSfWgbCE1S3W2LNY6ZVpRuun74moVBkKj6utE",
		"12D3KooWMMTDXcWhpVnwrdAer1jnVARTmnr3RyT3v7Djg8ZuoBh9",
		"12D3KooWGzVXsKxXsF4zLgxSDM8Gzx1ywq2pZef4PrHMKuVg4K3P",
		"12D3KooWSyjmmzjVtCzwN7bXzZQFmWiJRuVcKBerNjVgL7HdLJBW",
		"12D3KooWLGz9gzhrNsvyM6XnXS3JRkZoQdEzuAvysovnSChNK5ZK",
		"12D3KooWAvZnvknFAfSiUYjATyhzEJLTeKvAzpcLELHi4ogM3GET",
	}
	result.triggerDonSigners = []string{
		"0x9CcE7293a4Cc2621b61193135A95928735e4795F",
		"0x3c775F20bCB2108C1A818741Ce332Bb5fe0dB925",
		"0x50314239e2CF05555ceeD53E7F47eB2A8Eab0dbB",
		"0xd76A4f98898c3b9A72b244476d7337b50D54BCd8",
		"0x656A873f6895b8a03Fb112dE927d43FA54B2c92A",
		"0x5d1e87d87bF2e0cD4Ea64F381a2dbF45e5f0a553",
		"0x91d9b0062265514f012Eb8fABA59372fD9520f56",
	}
	result.targetDonPeers = []string{ // "cap-one"
		"12D3KooWJrthXtnPHw7xyHFAxo6NxifYTvc8igKYaA6wRRRqtsMb",
		"12D3KooWFQekP9sGex4XhqEJav5EScjTpDVtDqJFg1JvrePBCEGJ",
		"12D3KooWFLEq4hYtdyKWwe47dXGEbSiHMZhmr5xLSJNhpfiEz8NF",
		"12D3KooWN2hztiXNNS1jMQTTvvPRYcarK1C7T3Mdqk4x4gwyo5WS",
	}

	result.allPeers = make(map[ragetypes.PeerID]p2ptypes.StreamConfig)
	addPeersToDONInfo := func(peers []string, donInfo *capabilities.DON) error {
		for _, peerID := range peers {
			var p ragetypes.PeerID
			err := p.UnmarshalText([]byte(peerID))
			if err != nil {
				return err
			}
			result.allPeers[p] = defaultStreamConfig
			donInfo.Members = append(donInfo.Members, p)
		}
		return nil
	}
	result.WorkflowsDonInfo = capabilities.DON{ID: "workflowDon1", F: 2}
	if err := addPeersToDONInfo(result.workflowDonPeers, &result.WorkflowsDonInfo); err != nil {
		return HardcodedDonNetworkSetup{}, fmt.Errorf("failed to add peers to workflow DON info: %w", err)
	}
	result.TriggerCapabilityDonInfo = capabilities.DON{ID: "capabilityDon1", F: 1} // NOTE: misconfiguration - should be 2
	if err := addPeersToDONInfo(result.triggerDonPeers, &result.TriggerCapabilityDonInfo); err != nil {
		return HardcodedDonNetworkSetup{}, fmt.Errorf("failed to add peers to trigger DON info: %w", err)
	}

	result.TargetCapabilityDonInfo = capabilities.DON{ID: "targetDon1", F: 1}
	if err := addPeersToDONInfo(result.targetDonPeers, &result.TargetCapabilityDonInfo); err != nil {
		return HardcodedDonNetworkSetup{}, fmt.Errorf("failed to add peers to target DON info: %w", err)
	}

	return result, nil
}

func (h HardcodedDonNetworkSetup) IsWorkflowDon(id p2ptypes.PeerID) bool {
	return slices.Contains(h.workflowDonPeers, id.String())
}

func (h HardcodedDonNetworkSetup) IsTriggerDon(id p2ptypes.PeerID) bool {
	return slices.Contains(h.triggerDonPeers, id.String())
}

func (h HardcodedDonNetworkSetup) IsTargetDon(id p2ptypes.PeerID) bool {
	return slices.Contains(h.targetDonPeers, id.String())
}

type mockMercuryDataProducer struct {
	trigger *triggers.MercuryTriggerService
	wg      sync.WaitGroup
	closeCh chan struct{}
	lggr    logger.Logger
}

var _ services.Service = &mockMercuryDataProducer{}

func NewMockMercuryDataProducer(trigger *triggers.MercuryTriggerService, lggr logger.Logger) *mockMercuryDataProducer {
	return &mockMercuryDataProducer{
		trigger: trigger,
		closeCh: make(chan struct{}),
		lggr:    lggr,
	}
}

func (m *mockMercuryDataProducer) Start(ctx context.Context) error {
	m.wg.Add(1)
	go m.loop()
	return nil
}

func (m *mockMercuryDataProducer) loop() {
	defer m.wg.Done()

	sleepSec := 60
	ticker := time.NewTicker(time.Duration(sleepSec) * time.Second)
	defer ticker.Stop()

	prices := []*big.Int{big.NewInt(300000), big.NewInt(40000), big.NewInt(5000000)}

	for range ticker.C {
		for i := range prices {
			prices[i].Add(prices[i], big.NewInt(1))
		}

		reports := []datastreams.FeedReport{
			{
				FeedID:               "0x0003fbba4fce42f65d6032b18aee53efdf526cc734ad296cb57565979d883bdd",
				FullReport:           []byte{0x11, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[0].Bytes(),
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               "0x0003c317fec7fad514c67aacc6366bf2f007ce37100e3cddcacd0ccaa1f3746d",
				FullReport:           []byte{0x22, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[1].Bytes(),
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               "0x0003da6ab44ea9296674d80fe2b041738189103d6b4ea9a4d34e2f891fa93d12",
				FullReport:           []byte{0x33, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[2].Bytes(),
				ObservationTimestamp: time.Now().Unix(),
			},
		}

		m.lggr.Infow("New set of Mercury reports", "timestamp", time.Now().Unix(), "payload", reports)
		err := m.trigger.ProcessReport(reports)
		if err != nil {
			m.lggr.Errorw("failed to process Mercury reports", "err", err, "timestamp", time.Now().Unix(), "payload", reports)
		}
	}
}

func (m *mockMercuryDataProducer) Close() error {
	close(m.closeCh)
	m.wg.Wait()
	return nil
}

func (m *mockMercuryDataProducer) HealthReport() map[string]error {
	return nil
}

func (m *mockMercuryDataProducer) Ready() error {
	return nil
}

func (m *mockMercuryDataProducer) Name() string {
	return "mockMercuryDataProducer"
}

func hexStringsToBytes(strs []string) (res [][]byte) {
	for _, s := range strs {
		b, _ := hex.DecodeString(s[2:])
		res = append(res, b)
	}
	return res
}
