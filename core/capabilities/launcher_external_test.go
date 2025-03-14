package capabilities_test

import (
	"context"
	"crypto/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	corecapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func TestLauncher_UpdatesReceiverWithNewDON(t *testing.T) {
	m, err := values.NewMap(map[string]any{"response": "response1"})
	require.NoError(t, err)
	capabilityResponse := commoncap.CapabilityResponse{
		Value: m,
	}
	rawResponse, err := pb.MarshalCapabilityResponse(capabilityResponse)
	require.NoError(t, err)

	executeReceiveSafetly := func(ctx context.Context, receiver remotetypes.Receiver, msgBody *remotetypes.MessageBody) {
		// Create a wait group to ensure message processing completes
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			receiver.Receive(ctx, msgBody)
		}()

		// Wait for message processing with timeout
		doneCh := make(chan struct{})
		go func() {
			wg.Wait()
			close(doneCh)
		}()

		select {
		case <-doneCh:
			// Message processing completed
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for message processing")
		}
	}

	t.Run("receiver receives request from registered don", func(t *testing.T) {
		ctx := tests.Context(t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		var receiver remotetypes.Receiver

		// setup will create and start a launcher with a capabilities registry reflecting a current state of a WorkflowDon
		// and a CapabilitiesDon running a custom compute capability.
		// The receiver will be updated with the created CapabilitiesDon's ID
		th := setup(ctx, t, lggr, &receiver)
		msgBody := &remotetypes.MessageBody{
			Method:      remotetypes.MethodExecute,
			MessageId:   []byte("message_id"),
			Payload:     rawResponse,
			CallerDonId: th.workflowDonID,
			Sender:      th.workflowsNodes[0][:],
		}

		// we will now send a request to the receiver with a registered don's ID
		executeReceiveSafetly(ctx, receiver, msgBody)
		assert.Empty(t, observedLogs.FilterMessage("received request from unregistered don").All())
	})

	t.Run("receiver receives request from an unregistered don", func(t *testing.T) {
		ctx := tests.Context(t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		var receiver remotetypes.Receiver

		// setup will create and start a launcher with a capabilities registry reflecting a current state of a WorkflowDon
		// and a CapabilitiesDon running a custom compute capability.
		// The receiver will be updated with the created CapabilitiesDon's ID
		th := setup(ctx, t, lggr, &receiver)

		unregisteredDonID := uint32(3)
		msgBody := &remotetypes.MessageBody{
			Method:      remotetypes.MethodExecute,
			MessageId:   []byte("message_id"),
			Payload:     rawResponse,
			CallerDonId: unregisteredDonID,
			Sender:      th.workflowsNodes[0][:],
		}

		// we will now send a request to the receiver with an unregistered don's ID
		executeReceiveSafetly(ctx, receiver, msgBody)
		assert.Len(t, observedLogs.FilterMessage("received request from unregistered don").All(), 1)
	})

	t.Run("receivers gets updated when adding a new don", func(t *testing.T) {
		ctx := tests.Context(t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		var receiver remotetypes.Receiver

		// setup will create and start a launcher with a capabilities registry reflecting a current state of a WorkflowDon
		// and a CapabilitiesDon running a custom compute capability.
		// The receiver will be updated with the created CapabilitiesDon's ID
		th := setup(ctx, t, lggr, &receiver)

		// we will now add a new workflow don to the state which reflects a new workflow don being added to the registrySyncer.
		// this emulates what happens in the registrySyncer.Sync
		newWorkflowDonID := uint32(3)
		newWorkflowsNodes := []ragetypes.PeerID{
			randomWord(),
			randomWord(),
			randomWord(),
			randomWord(),
		}
		th.state.IDsToDONs[registrysyncer.DonID(newWorkflowDonID)] = registrysyncer.DON{
			DON: capabilities.DON{
				ID:               newWorkflowDonID,
				ConfigVersion:    uint32(0),
				F:                uint8(1),
				IsPublic:         true,
				AcceptsWorkflows: true,
				Members:          newWorkflowsNodes,
			},
		}

		th.state.IDsToNodes[newWorkflowsNodes[0]] = kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      1,
			Signer:              randomWord(),
			P2pId:               newWorkflowsNodes[0],
			EncryptionPublicKey: randomWord(),
			HashedCapabilityIds: [][32]byte{th.computeCapID},
		}
		th.state.IDsToNodes[newWorkflowsNodes[1]] = kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      1,
			Signer:              randomWord(),
			P2pId:               newWorkflowsNodes[1],
			EncryptionPublicKey: randomWord(),
			HashedCapabilityIds: [][32]byte{th.computeCapID},
		}
		th.state.IDsToNodes[newWorkflowsNodes[2]] = kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      1,
			Signer:              randomWord(),
			P2pId:               newWorkflowsNodes[2],
			EncryptionPublicKey: randomWord(),
			HashedCapabilityIds: [][32]byte{th.computeCapID},
		}
		th.state.IDsToNodes[newWorkflowsNodes[3]] = kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      1,
			Signer:              randomWord(),
			P2pId:               newWorkflowsNodes[3],
			EncryptionPublicKey: randomWord(),
			HashedCapabilityIds: [][32]byte{th.computeCapID},
		}

		// the dispatcher will return a remote.ErrReceiverExists error when trying to update the receiver with the new workflow don's ID
		// this is because currently the key is formed by the fullComputeCapID and the capabilitiesDonID and a receiver already exists for this key
		th.dispatcher.On("SetReceiver", th.fullComputeCapID, th.capabilitiesDonID, mock.AnythingOfType("*executable.server")).Return(remote.ErrReceiverExists).Once()
		err = th.launcher.Launch(ctx, &th.state)
		require.NoError(t, err)

		// we will now send a request to the receiver with the new workflow don's ID
		// the receiver should not log an error as the new workflow don is now registered
		// given that the receiver was updated with the new workflow don's ID when the launcher.Launch was called
		msgBody := &remotetypes.MessageBody{
			Method:      remotetypes.MethodExecute,
			MessageId:   []byte("message_id"),
			Payload:     rawResponse,
			CallerDonId: newWorkflowDonID,
			Sender:      newWorkflowsNodes[0][:],
		}

		executeReceiveSafetly(ctx, receiver, msgBody)
		assert.Empty(t, observedLogs.FilterMessage("received request from unregistered don").All())
	})
}

type testHarness struct {
	workflowDonID     uint32
	capabilitiesDonID uint32
	workflowsNodes    []ragetypes.PeerID
	launcher          registrysyncer.Launcher
	dispatcher        *remoteMocks.Dispatcher
	state             registrysyncer.LocalRegistry
	computeCapID      [32]byte
	fullComputeCapID  string
}

func setup(ctx context.Context, t *testing.T, lggr logger.Logger, receiver *remotetypes.Receiver) testHarness {
	registry := corecapabilities.NewRegistry(lggr)
	fullComputeCapID := "custom-compute@1.0.0"
	mt := newMockAction(capabilities.MustNewCapabilityInfo(
		fullComputeCapID,
		capabilities.CapabilityTypeAction,
		"custom compute",
	))
	computeCapID := randomWord()
	require.NoError(t, registry.Add(ctx, mt))

	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	dispatcher := remoteMocks.NewDispatcher(t)

	launcher := corecapabilities.NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	defer launcher.Close()

	workflowDonID := uint32(1)
	workflowsNodes := []ragetypes.PeerID{
		randomWord(),
		randomWord(),
		randomWord(),
		randomWord(),
	}

	capabilitiesDonID := uint32(2)
	capabilityNodes := []ragetypes.PeerID{
		wrapper.GetPeer().ID(),
		randomWord(),
		randomWord(),
		randomWord(),
	}

	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// And a Capabilities DON that exposes the Compute capability.
	// We expect the launcher to use a deep copy of the state,
	// making it able to modify the state without affecting the original state.
	state := &registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(workflowDonID): {
				DON: capabilities.DON{
					ID:               workflowDonID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowsNodes,
				},
			},
			registrysyncer.DonID(capabilitiesDonID): {
				DON: capabilities.DON{
					ID:               capabilitiesDonID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: false,
					Members:          capabilityNodes,
				},
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					fullComputeCapID: {},
				},
			},
		},
		IDsToCapabilities: map[string]registrysyncer.Capability{
			fullComputeCapID: {
				ID:             "custom-compute@1.0.0",
				CapabilityType: capabilities.CapabilityTypeAction,
			},
		},
		IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
			capabilityNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityNodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{computeCapID},
			},
			capabilityNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityNodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{computeCapID},
			},
			capabilityNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityNodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{computeCapID},
			},
			capabilityNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityNodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{computeCapID},
			},
			workflowsNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowsNodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{computeCapID},
			},
			workflowsNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowsNodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{computeCapID},
			},
			workflowsNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowsNodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{computeCapID},
			},
			workflowsNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowsNodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{computeCapID},
			},
		},
	}

	dispatcher.On("SetReceiver", fullComputeCapID, capabilitiesDonID, mock.AnythingOfType("*executable.server")).Run(func(args mock.Arguments) {
		*receiver = args.Get(2).(remotetypes.Receiver)
	}).Return(nil).Once()

	err = launcher.Launch(ctx, state)
	require.NoError(t, err)
	stateCopy := registrysyncer.DeepCopyLocalRegistry(state)

	return testHarness{
		workflowDonID:     workflowDonID,
		capabilitiesDonID: capabilitiesDonID,
		workflowsNodes:    workflowsNodes,
		launcher:          launcher,
		dispatcher:        dispatcher,
		state:             stateCopy,
		computeCapID:      computeCapID,
		fullComputeCapID:  fullComputeCapID,
	}
}

var _ capabilities.ActionCapability = (*mockAction)(nil)

type mockAction struct {
	capabilities.CapabilityInfo
}

func (m *mockAction) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}
func (m *mockAction) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
func (m *mockAction) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	return capabilities.CapabilityResponse{}, nil
}

func newMockAction(info capabilities.CapabilityInfo) *mockAction {
	return &mockAction{CapabilityInfo: info}
}

type mockDonNotifier struct {
}

func (m *mockDonNotifier) NotifyDonSet(don capabilities.DON) {
}

func randomWord() [32]byte {
	word := make([]byte, 32)
	_, err := rand.Read(word)
	if err != nil {
		panic(err)
	}
	return [32]byte(word)
}
