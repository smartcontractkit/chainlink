package capabilities

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
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

type workflowHandler struct {
	services.StateMachine
	lggr        logger.Logger
	peerWrapper p2ptypes.PeerWrapper
	dispatcher  remotetypes.Dispatcher
	registry    core.CapabilitiesRegistry
	localNode   capabilities.Node
	subServices []services.Service
}

func NewWorkflowSyncerHandler(
	lggr logger.Logger,
	peerWrapper p2ptypes.PeerWrapper,
	dispatcher remotetypes.Dispatcher,
	registry core.CapabilitiesRegistry,
) *workflowHandler {
	return &workflowHandler{
		lggr:        lggr,
		peerWrapper: peerWrapper,
		dispatcher:  dispatcher,
		registry:    registry,
		subServices: []services.Service{},
	}
}

func (w *workflowHandler) Start(ctx context.Context) error {
	return nil
}

func (w *workflowHandler) Close() error {
	for _, s := range w.subServices {
		if err := s.Close(); err != nil {
			w.lggr.Errorw("failed to close a sub-service", "name", s.Name(), "error", err)
		}
	}

	return w.peerWrapper.GetPeer().UpdateConnections(map[ragetypes.PeerID]p2ptypes.StreamConfig{})
}

func (w *workflowHandler) Ready() error {
	return nil
}

func (w *workflowHandler) HealthReport() map[string]error {
	return nil
}

func (w *workflowHandler) Name() string {
	return "WorkflowSyncerHandler"
}

func (w *workflowHandler) LocalNode(ctx context.Context) (capabilities.Node, error) {
	if w.peerWrapper.GetPeer() == nil {
		return w.localNode, errors.New("unable to get local node: peerWrapper hasn't started yet")
	}

	if w.localNode.WorkflowDON.ID == "" {
		return w.localNode, errors.New("unable to get local node: waiting for initial call from syncer")
	}

	return w.localNode, nil
}

func (w *workflowHandler) updateLocalNode(state registrysyncer.State) {
	pid := w.peerWrapper.GetPeer().ID()

	var workflowDON capabilities.DON
	capabilityDONs := []capabilities.DON{}
	for _, d := range state.IDsToDONs {
		for _, p := range d.NodeP2PIds {
			if p == pid {
				if d.AcceptsWorkflows {
					if workflowDON.ID == "" {
						workflowDON = *toDONInfo(d)
					} else {
						w.lggr.Errorf("Configuration error: node %s belongs to more than one workflowDON", pid)
					}
				}

				capabilityDONs = append(capabilityDONs, *toDONInfo(d))
			}
		}
	}

	w.localNode = capabilities.Node{
		PeerID:         &pid,
		WorkflowDON:    workflowDON,
		CapabilityDONs: capabilityDONs,
	}

}

func (w *workflowHandler) Handle(ctx context.Context, state registrysyncer.State) error {
	w.updateLocalNode(state)

	// Let's start by updating the list of Peers
	// We do this by creating a new entry for each node belonging
	// to a public DON.
	// We also add the hardcoded peers determined by the NetworkSetup.
	allPeers := make(map[ragetypes.PeerID]p2ptypes.StreamConfig)

	publicDONs := []kcr.CapabilitiesRegistryDONInfo{}
	for _, d := range state.IDsToDONs {
		if !d.IsPublic {
			continue
		}

		publicDONs = append(publicDONs, d)

		for _, nid := range d.NodeP2PIds {
			allPeers[nid] = defaultStreamConfig
		}
	}

	// TODO: be a bit smarter about who we connect to; we should ideally only
	// be connecting to peers when we need to.
	// https://smartcontract-it.atlassian.net/browse/KS-330
	err := w.peerWrapper.GetPeer().UpdateConnections(allPeers)
	if err != nil {
		return fmt.Errorf("failed to update peer connections: %w", err)
	}

	// Next, we need to split the DONs into the following:
	// - workflow DONs the current node is a part of.
	// These will need remote shims to all remote capabilities on other DONs.
	//
	// We'll also construct a set to record what DONs the current node is a part of,
	// regardless of any modifiers (public/acceptsWorkflows etc).
	myID := w.peerWrapper.GetPeer().ID()
	myWorkflowDONs := []kcr.CapabilitiesRegistryDONInfo{}
	remoteWorkflowDONs := []kcr.CapabilitiesRegistryDONInfo{}
	myDONs := map[uint32]bool{}
	for _, d := range state.IDsToDONs {
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
	myCapabilityDONs := []kcr.CapabilitiesRegistryDONInfo{}
	remoteCapabilityDONs := []kcr.CapabilitiesRegistryDONInfo{}
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
			w.lggr.Warn("node is part of more than one workflow DON; assigning first DON as caller")
		}

		for _, rcd := range remoteCapabilityDONs {
			err := w.addRemoteCapabilities(ctx, myDON, rcd, state)
			if err != nil {
				return err
			}
		}
	}

	// Finally, if I'm a capability DON, let's enable external access
	// to the capability.
	if len(myCapabilityDONs) > 0 {
		for _, mcd := range myCapabilityDONs {
			err := w.enableExternalAccess(ctx, myID, mcd, state, remoteWorkflowDONs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *workflowHandler) addRemoteCapabilities(ctx context.Context, myDON kcr.CapabilitiesRegistryDONInfo, remoteDON kcr.CapabilitiesRegistryDONInfo, state registrysyncer.State) error {
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

				codec := streams.NewCodec(w.lggr)

				signers, err := signersFor(remoteDON, state)
				if err != nil {
					return nil, err
				}

				aggregator := triggers.NewMercuryRemoteAggregator(
					codec,
					signers,
					int(remoteDON.F+1),
					w.lggr,
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
					w.dispatcher,
					aggregator,
					w.lggr,
				)
				return triggerCap, nil
			}
			err := w.addToRegistryAndSetDispatcher(ctx, capability, remoteDON, newTriggerFn)
			if err != nil {
				return fmt.Errorf("failed to add trigger shim: %w", err)
			}
		case capabilities.CapabilityTypeAction:
			w.lggr.Warn("no remote client configured for capability type action, skipping configuration")
		case capabilities.CapabilityTypeConsensus:
			w.lggr.Warn("no remote client configured for capability type consensus, skipping configuration")
		case capabilities.CapabilityTypeTarget:
			newTargetFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
				client := target.NewClient(
					info,
					*toDONInfo(myDON),
					w.dispatcher,
					defaultTargetRequestTimeout,
					w.lggr,
				)
				return client, nil
			}

			err := w.addToRegistryAndSetDispatcher(ctx, capability, remoteDON, newTargetFn)
			if err != nil {
				return fmt.Errorf("failed to add target shim: %w", err)
			}
		default:
			w.lggr.Warnf("unknown capability type, skipping configuration: %+v", capability)
		}
	}
	return nil
}

type capabilityService interface {
	capabilities.BaseCapability
	remotetypes.Receiver
	services.Service
}

func (w *workflowHandler) addToRegistryAndSetDispatcher(ctx context.Context, capabilityInfo kcr.CapabilitiesRegistryCapabilityInfo, don kcr.CapabilitiesRegistryDONInfo, newCapFn func(info capabilities.CapabilityInfo) (capabilityService, error)) error {
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
	w.lggr.Debugw("Adding remote capability to registry", "id", info.ID, "don", info.DON)
	capability, err := newCapFn(info)
	if err != nil {
		return err
	}

	err = w.registry.Add(ctx, capability)
	if err != nil {
		// If the capability already exists, then it's either local
		// or we've handled this in a previous syncer iteration,
		// let's skip and move on to other capabilities.
		if errors.Is(err, ErrCapabilityAlreadyExists) {
			return nil
		}

		return err
	}

	err = w.dispatcher.SetReceiver(
		fullCapID,
		fmt.Sprint(don.Id),
		capability,
	)
	if err != nil {
		return err
	}
	w.lggr.Debugw("Setting receiver for capability", "id", fullCapID, "donID", don.Id)
	err = capability.Start(ctx)
	if err != nil {
		return err
	}
	w.subServices = append(w.subServices, capability)
	return nil
}

var (
	defaultTargetRequestTimeout = time.Minute
)

func (w *workflowHandler) enableExternalAccess(ctx context.Context, myPeerID p2ptypes.PeerID, don kcr.CapabilitiesRegistryDONInfo, state registrysyncer.State, remoteWorkflowDONs []kcr.CapabilitiesRegistryDONInfo) error {
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
					w.dispatcher,
					w.lggr,
				)
				return publisher, nil
			}

			err := w.addReceiver(ctx, capability, don, newTriggerPublisher)
			if err != nil {
				return fmt.Errorf("failed to add server-side receiver: %w", err)
			}
		case capabilities.CapabilityTypeAction:
			w.lggr.Warn("no remote client configured for capability type action, skipping configuration")
		case capabilities.CapabilityTypeConsensus:
			w.lggr.Warn("no remote client configured for capability type consensus, skipping configuration")
		case capabilities.CapabilityTypeTarget:
			newTargetServer := func(capability capabilities.BaseCapability, info capabilities.CapabilityInfo) (receiverService, error) {
				return target.NewServer(
					myPeerID,
					capability.(capabilities.TargetCapability),
					info,
					*toDONInfo(don),
					idsToDONs,
					w.dispatcher,
					defaultTargetRequestTimeout,
					w.lggr,
				), nil
			}

			err := w.addReceiver(ctx, capability, don, newTargetServer)
			if err != nil {
				return fmt.Errorf("failed to add server-side receiver: %w", err)
			}
		default:
			w.lggr.Warnf("unknown capability type, skipping configuration: %+v", capability)
		}
	}
	return nil
}

type receiverService interface {
	services.Service
	remotetypes.Receiver
}

func (w *workflowHandler) addReceiver(ctx context.Context, capability kcr.CapabilitiesRegistryCapabilityInfo, don kcr.CapabilitiesRegistryDONInfo, newReceiverFn func(capability capabilities.BaseCapability, info capabilities.CapabilityInfo) (receiverService, error)) error {
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
	underlying, err := w.registry.Get(ctx, fullCapID)
	if err != nil {
		return err
	}

	receiver, err := newReceiverFn(underlying, info)
	if err != nil {
		return err
	}

	w.lggr.Debugw("Enabling external access for capability", "id", fullCapID, "donID", don.Id)
	err = w.dispatcher.SetReceiver(fullCapID, fmt.Sprint(don.Id), receiver)
	if err != nil {
		return err
	}

	err = receiver.Start(ctx)
	if err != nil {
		return err
	}

	w.subServices = append(w.subServices, receiver)
	return nil
}

func signersFor(don kcr.CapabilitiesRegistryDONInfo, state registrysyncer.State) ([][]byte, error) {
	s := [][]byte{}
	for _, nodeID := range don.NodeP2PIds {
		node, ok := state.IDsToNodes[nodeID]
		if !ok {
			return nil, fmt.Errorf("could not find node for id %s", nodeID)
		}

		// NOTE: the capability registry stores signers as [32]byte,
		// but we only need the first [20], as the rest is padded.
		s = append(s, node.Signer[0:20])
	}

	return s, nil
}

func toDONInfo(don kcr.CapabilitiesRegistryDONInfo) *capabilities.DON {
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
