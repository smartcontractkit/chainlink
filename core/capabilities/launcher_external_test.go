package capabilities_test

import (
	"context"
	"crypto/rand"
	"testing"

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
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func TestLauncher_UpdatesReceiverWithNewDON(t *testing.T) {
	ctx := tests.Context(t)
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)

	m, err := values.NewMap(map[string]any{"response": "response1"})
	require.NoError(t, err)
	capabilityResponse := commoncap.CapabilityResponse{
		Value: m,
	}
	rawResponse, err := pb.MarshalCapabilityResponse(capabilityResponse)
	require.NoError(t, err)

	t.Run("receiver receives request from registered don", func(t *testing.T) {
		var receiver remotetypes.Receiver
		workflowDonID, workflowsNodes := testSetup(t, ctx, lggr, &receiver)
		msgBody := &remotetypes.MessageBody{
			Method:      remotetypes.MethodExecute,
			MessageId:   []byte("message_id"),
			Payload:     rawResponse,
			CallerDonId: workflowDonID,
			Sender:      workflowsNodes[0][:],
		}

		receiver.Receive(ctx, msgBody)
		assert.Equal(t, 0, len(observedLogs.FilterMessage("received request from unregistered don").All()))
	})

	t.Run("receiver receives request from an unregistered don", func(t *testing.T) {
		var receiver remotetypes.Receiver
		_, workflowsNodes := testSetup(t, ctx, lggr, &receiver)

		unregisteredDonID := uint32(3)
		msgBody := &remotetypes.MessageBody{
			Method:      remotetypes.MethodExecute,
			MessageId:   []byte("message_id"),
			Payload:     rawResponse,
			CallerDonId: unregisteredDonID,
			Sender:      workflowsNodes[0][:],
		}

		receiver.Receive(ctx, msgBody)
		assert.Equal(t, 1, len(observedLogs.FilterMessage("received request from unregistered don").All()))
	})
}

func testSetup(t *testing.T, ctx context.Context, lggr logger.Logger, receiver *remotetypes.Receiver) (uint32, []ragetypes.PeerID) {
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
	}).Return(nil)

	err = launcher.Launch(ctx, state)
	require.NoError(t, err)

	return workflowDonID, workflowsNodes
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
