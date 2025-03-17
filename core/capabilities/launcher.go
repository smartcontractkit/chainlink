package capabilities

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	peerWrapper         p2ptypes.PeerWrapper
	dispatcher          remotetypes.Dispatcher
	registry            *Registry
	subServices         []services.Service
	workflowDonNotifier donNotifier
}

func unmarshalCapabilityConfig(data []byte) (capabilities.CapabilityConfiguration, error) {
	cconf := &capabilitiespb.CapabilityConfig{}
	err := proto.Unmarshal(data, cconf)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, err
	}

	var remoteTriggerConfig *capabilities.RemoteTriggerConfig
	var remoteTargetConfig *capabilities.RemoteTargetConfig

	switch cconf.GetRemoteConfig().(type) {
	case *capabilitiespb.CapabilityConfig_RemoteTriggerConfig:
		prtc := cconf.GetRemoteTriggerConfig()
		remoteTriggerConfig = &capabilities.RemoteTriggerConfig{}
		remoteTriggerConfig.RegistrationRefresh = prtc.RegistrationRefresh.AsDuration()
		remoteTriggerConfig.RegistrationExpiry = prtc.RegistrationExpiry.AsDuration()
		remoteTriggerConfig.MinResponsesToAggregate = prtc.MinResponsesToAggregate
		remoteTriggerConfig.MessageExpiry = prtc.MessageExpiry.AsDuration()
	case *capabilitiespb.CapabilityConfig_RemoteTargetConfig:
		prtc := cconf.GetRemoteTargetConfig()
		remoteTargetConfig = &capabilities.RemoteTargetConfig{}
		remoteTargetConfig.RequestHashExcludedAttributes = prtc.RequestHashExcludedAttributes
	}

	dc, err := values.FromMapValueProto(cconf.DefaultConfig)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, err
	}

	rc, err := values.FromMapValueProto(cconf.RestrictedConfig)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, err
	}

	return capabilities.CapabilityConfiguration{
		DefaultConfig:       dc,
		RestrictedKeys:      cconf.RestrictedKeys,
		RestrictedConfig:    rc,
		RemoteTriggerConfig: remoteTriggerConfig,
		RemoteTargetConfig:  remoteTargetConfig,
	}, nil
}

type donNotifier interface {
	NotifyDonSet(don capabilities.DON)
}

func NewLauncher(
	lggr logger.Logger,
	peerWrapper p2ptypes.PeerWrapper,
	dispatcher remotetypes.Dispatcher,
	registry *Registry,
	workflowDonNotifier donNotifier,
) *launcher {
	return &launcher{
		lggr:                lggr.Named("CapabilitiesLauncher"),
		peerWrapper:         peerWrapper,
		dispatcher:          dispatcher,
		registry:            registry,
		subServices:         []services.Service{},
		workflowDonNotifier: workflowDonNotifier,
	}
}

func (w *launcher) Start(ctx context.Context) error {
	return nil
}

func (w *launcher) Close() error {
	for _, s := range w.subServices {
		if err := s.Close(); err != nil {
			w.lggr.Errorw("failed to close a sub-service", "name", s.Name(), "error", err)
		}
	}

	return w.peerWrapper.GetPeer().UpdateConnections(map[ragetypes.PeerID]p2ptypes.StreamConfig{})
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

func (w *launcher) Launch(ctx context.Context, state *registrysyncer.LocalRegistry) error {
	w.lggr.Debug("CapabilitiesLauncher triggered...")
	w.registry.SetLocalRegistry(state)

	allDONIDs := []registrysyncer.DonID{}
	for id := range state.IDsToDONs {
		allDONIDs = append(allDONIDs, id)
	}
	slices.Sort(allDONIDs) // ensure deterministic order

	// Let's start by updating the list of Peers
	// We do this by creating a new entry for each node belonging
	// to a public DON.
	// We also add the hardcoded peers determined by the NetworkSetup.
	allPeers := make(map[ragetypes.PeerID]p2ptypes.StreamConfig)

	publicDONs := []registrysyncer.DON{}
	for _, id := range allDONIDs {
		d := state.IDsToDONs[id]
		if !d.DON.IsPublic {
			continue
		}

		publicDONs = append(publicDONs, d)

		for _, nid := range d.DON.Members {
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
	myWorkflowDONs := []registrysyncer.DON{}
	remoteWorkflowDONs := []registrysyncer.DON{}
	myDONs := map[uint32]bool{}
	for _, id := range allDONIDs {
		d := state.IDsToDONs[id]
		for _, peerID := range d.Members {
			if peerID == myID {
				myDONs[d.ID] = true
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

	// Now, if my node is a workflow DON, let's setup any shims
	// to external capabilities.
	if len(myWorkflowDONs) > 0 {
		myDON := myWorkflowDONs[0]

		// NOTE: this is enforced on-chain and so should never happen.
		if len(myWorkflowDONs) > 1 {
			return errors.New("invariant violation: node is part of more than one workflowDON")
		}

		w.lggr.Debug("Notifying DON set...")
		w.workflowDonNotifier.NotifyDonSet(myDON.DON)

		for _, rcd := range remoteCapabilityDONs {
			err := w.addRemoteCapabilities(ctx, myDON, rcd, state)
			if err != nil {
				return err
			}
		}
	}

	// Finally, if I'm in a capability DON, let's enable external access
	// to the capability.
	if len(myCapabilityDONs) > 0 {
		for _, myDON := range myCapabilityDONs {
			err := w.exposeCapabilities(ctx, myID, myDON, state, remoteWorkflowDONs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *launcher) addRemoteCapabilities(ctx context.Context, myDON registrysyncer.DON, remoteDON registrysyncer.DON, state *registrysyncer.LocalRegistry) error {
	for cid, c := range remoteDON.CapabilityConfigurations {
		capability, ok := state.IDsToCapabilities[cid]
		if !ok {
			return fmt.Errorf("could not find capability matching id %s", cid)
		}

		capabilityConfig, err := unmarshalCapabilityConfig(c.Config)
		if err != nil {
			return fmt.Errorf("could not unmarshal capability config for id %s", cid)
		}

		switch capability.CapabilityType {
		case capabilities.CapabilityTypeTrigger:
			newTriggerFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
				var aggregator remotetypes.Aggregator
				if strings.HasPrefix(info.ID, "streams-trigger") {
					codec := streams.NewCodec(w.lggr)

					signers, err := signersFor(remoteDON, state)
					if err != nil {
						return nil, err
					}

					aggregator = triggers.NewMercuryRemoteAggregator(
						codec,
						signers,
						int(remoteDON.F+1),
						info.ID,
						w.lggr,
					)
				} else {
					aggregator = aggregation.NewDefaultModeAggregator(uint32(remoteDON.F) + 1)
				}

				// TODO: We need to implement a custom, Mercury-specific
				// aggregator here, because there is no guarantee that
				// all trigger events in the workflow will have the same
				// payloads. As a workaround, we validate the signatures.
				// When this is solved, we can move to a generic aggregator
				// and remove this.
				triggerCap := remote.NewTriggerSubscriber(
					capabilityConfig.RemoteTriggerConfig,
					info,
					remoteDON.DON,
					myDON.DON,
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
			newActionFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
				client := executable.NewClient(
					info,
					myDON.DON,
					w.dispatcher,
					defaultTargetRequestTimeout,
					w.lggr,
				)
				return client, nil
			}

			err := w.addToRegistryAndSetDispatcher(ctx, capability, remoteDON, newActionFn)
			if err != nil {
				return fmt.Errorf("failed to add action shim: %w", err)
			}
		case capabilities.CapabilityTypeConsensus:
			// nothing to do; we don't support remote consensus capabilities for now
		case capabilities.CapabilityTypeTarget:
			newTargetFn := func(info capabilities.CapabilityInfo) (capabilityService, error) {
				client := executable.NewClient(
					info,
					myDON.DON,
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

func (w *launcher) addToRegistryAndSetDispatcher(ctx context.Context, capability registrysyncer.Capability, don registrysyncer.DON, newCapFn func(info capabilities.CapabilityInfo) (capabilityService, error)) error {
	capabilityID := capability.ID
	info, err := capabilities.NewRemoteCapabilityInfo(
		capabilityID,
		capability.CapabilityType,
		fmt.Sprintf("Remote Capability for %s", capabilityID),
		&don.DON,
	)
	if err != nil {
		return fmt.Errorf("failed to create remote capability info: %w", err)
	}
	w.lggr.Debugw("Adding remote capability to registry", "id", info.ID, "don", info.DON)
	cp, err := newCapFn(info)
	if err != nil {
		return fmt.Errorf("failed to instantiate capability: %w", err)
	}

	err = w.registry.Add(ctx, cp)
	if err != nil {
		// If the capability already exists, then it's either local
		// or we've handled this in a previous syncer iteration,
		// let's skip and move on to other capabilities.
		if errors.Is(err, ErrCapabilityAlreadyExists) {
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
	defaultTargetRequestTimeout = 8 * time.Minute
)

func (w *launcher) exposeCapabilities(ctx context.Context, myPeerID p2ptypes.PeerID, don registrysyncer.DON, state *registrysyncer.LocalRegistry, remoteWorkflowDONs []registrysyncer.DON) error {
	idsToDONs := map[uint32]capabilities.DON{}
	for _, d := range remoteWorkflowDONs {
		idsToDONs[d.ID] = d.DON
	}

	for cid, c := range don.CapabilityConfigurations {
		capability, ok := state.IDsToCapabilities[cid]
		if !ok {
			return fmt.Errorf("could not find capability matching id %s", cid)
		}

		capabilityConfig, err := unmarshalCapabilityConfig(c.Config)
		if err != nil {
			return fmt.Errorf("could not unmarshal capability config for id %s", cid)
		}

		switch capability.CapabilityType {
		case capabilities.CapabilityTypeTrigger:
			newTriggerPublisher := func(cap capabilities.BaseCapability, info capabilities.CapabilityInfo) (remotetypes.ReceiverService, error) {
				triggerCapability, ok := (cap).(capabilities.TriggerCapability)
				if !ok {
					return nil, errors.New("capability does not implement TriggerCapability")
				}

				publisher := remote.NewTriggerPublisher(
					capabilityConfig.RemoteTriggerConfig,
					triggerCapability,
					info,
					don.DON,
					idsToDONs,
					w.dispatcher,
					w.lggr,
				)
				return publisher, nil
			}

			err := w.addReceiver(ctx, capability, don, newTriggerPublisher)
			if err != nil {
				w.lggr.Errorw("failed to add server-side receiver for a trigger capability - it won't be exposed remotely", "id", cid, "error", err)
				// continue attempting other capabilities
			}
		case capabilities.CapabilityTypeAction:
			newActionServer := func(cap capabilities.BaseCapability, info capabilities.CapabilityInfo) (remotetypes.ReceiverService, error) {
				actionCapability, ok := (cap).(capabilities.ActionCapability)
				if !ok {
					return nil, errors.New("capability does not implement ActionCapability")
				}

				remoteConfig := &capabilities.RemoteExecutableConfig{}
				if capabilityConfig.RemoteTargetConfig != nil {
					remoteConfig.RequestHashExcludedAttributes = capabilityConfig.RemoteTargetConfig.RequestHashExcludedAttributes
				}

				return executable.NewServer(
					capabilityConfig.RemoteExecutableConfig,
					myPeerID,
					actionCapability,
					info,
					don.DON,
					idsToDONs,
					w.dispatcher,
					defaultTargetRequestTimeout,
					w.lggr,
				), nil
			}

			err = w.addReceiver(ctx, capability, don, newActionServer)
			if err != nil {
				w.lggr.Errorw("failed to add action server-side receiver - it won't be exposed remotely", "id", cid, "error", err)
				// continue attempting other capabilities
			}
		case capabilities.CapabilityTypeConsensus:
			w.lggr.Warn("no remote client configured for capability type consensus, skipping configuration")
		case capabilities.CapabilityTypeTarget:
			newTargetServer := func(cap capabilities.BaseCapability, info capabilities.CapabilityInfo) (remotetypes.ReceiverService, error) {
				targetCapability, ok := (cap).(capabilities.TargetCapability)
				if !ok {
					return nil, errors.New("capability does not implement TargetCapability")
				}

				remoteConfig := &capabilities.RemoteExecutableConfig{}
				if capabilityConfig.RemoteTargetConfig != nil {
					remoteConfig.RequestHashExcludedAttributes = capabilityConfig.RemoteTargetConfig.RequestHashExcludedAttributes
				}

				return executable.NewServer(
					remoteConfig,
					myPeerID,
					targetCapability,
					info,
					don.DON,
					idsToDONs,
					w.dispatcher,
					defaultTargetRequestTimeout,
					w.lggr,
				), nil
			}

			err := w.addReceiver(ctx, capability, don, newTargetServer)
			if err != nil {
				w.lggr.Errorw("failed to add server-side receiver for a target capability - it won't be exposed remotely", "id", cid, "error", err)
				// continue attempting other capabilities
			}
		default:
			w.lggr.Warnf("unknown capability type, skipping configuration: %+v", capability)
		}
	}
	return nil
}

func (w *launcher) addReceiver(ctx context.Context, capability registrysyncer.Capability, don registrysyncer.DON, newReceiverFn func(capability capabilities.BaseCapability, info capabilities.CapabilityInfo) (remotetypes.ReceiverService, error)) error {
	capID := capability.ID
	info, err := capabilities.NewRemoteCapabilityInfo(
		capID,
		capability.CapabilityType,
		fmt.Sprintf("Remote Capability for %s", capability.ID),
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

func signersFor(don registrysyncer.DON, state *registrysyncer.LocalRegistry) ([][]byte, error) {
	s := [][]byte{}
	for _, nodeID := range don.Members {
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
