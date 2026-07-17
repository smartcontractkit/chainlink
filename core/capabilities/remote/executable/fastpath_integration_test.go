package executable_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/vaultshare"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// TestVaultFastPathCapability returns a GetSecretsResponse with one distinct binary share per call.
// It is shared across all capability peers, so an atomic counter assigns each peer a unique share.
type TestVaultFastPathCapability struct {
	abstractTestCapability
	encKey string
	owner  string
	calls  atomic.Int64
}

func (t *TestVaultFastPathCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	n := t.calls.Add(1)
	vaultResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: &vaultcommon.SecretIdentifier{Owner: t.owner, Namespace: "main", Key: "secret"},
			Result: &vaultcommon.SecretResponse_Data{
				Data: &vaultcommon.SecretData{
					EncryptedValue: "abc123",
					EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
						EncryptionKey: t.encKey,
						BinaryShares:  [][]byte{[]byte{byte(n)}},
					}},
				},
			},
		}},
	}
	anyPayload, err := anypb.New(vaultResp)
	if err != nil {
		return commoncap.CapabilityResponse{}, err
	}
	return commoncap.CapabilityResponse{Payload: anyPayload}, nil
}

func Test_RemoteExecutable_FastPathAggregator_MergesFPlusOneShares(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4
	capabilityDonF := uint8(1)
	numWorkflowPeers := 1
	workflowDonF := uint8(0)
	capability := &TestVaultFastPathCapability{
		encKey: "enc-key",
		owner:  "owner",
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	method := func(ctx context.Context, caller commoncap.ExecutableCapability) {
		gsr := &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{{
				Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"},
			}},
		}
		anyPayload, err := anypb.New(gsr)
		require.NoError(t, err)

		response, err := caller.Execute(ctx, commoncap.CapabilityRequest{
			Method:   vaulttypes.MethodSecretsGet,
			Payload:  anyPayload,
			Config:   transmissionSchedule,
			Metadata: commoncap.RequestMetadata{WorkflowID: workflowID1, WorkflowExecutionID: workflowExecutionID1},
		})
		require.NoError(t, err)

		merged := &vaultcommon.GetSecretsResponse{}
		require.NoError(t, response.Payload.UnmarshalTo(merged))
		require.Len(t, merged.Responses, 1)
		require.Len(t, merged.Responses[0].GetData().EncryptedDecryptionKeyShares, 1)
		require.Len(t, merged.Responses[0].GetData().EncryptedDecryptionKeyShares[0].BinaryShares, 2)
	}

	// Use the shared harness from endtoend_test.go with an ungated vault aggregator factory.
	testFastPathRemoteExecutableCapability(ctx, t, capability, numWorkflowPeers, workflowDonF, 10*time.Second,
		numCapabilityPeers, capabilityDonF, 10*time.Second, method, true,
		withUngatedVaultAggregatorFactory(int(capabilityDonF)))
}

// aggregatorFactoryOpt allows the test harness to inject a custom aggregator factory into the
// workflow client SetConfig call. The default is nil (legacy hash-quorum behavior).
type aggregatorFactoryOpt func(*aggregatorFactoryCfg)

type aggregatorFactoryCfg struct {
	factory request.AggregatorFactory
}

func withAggregatorFactory(factory request.AggregatorFactory) aggregatorFactoryOpt {
	return func(cfg *aggregatorFactoryCfg) {
		cfg.factory = factory
	}
}

// withUngatedVaultAggregatorFactory returns an aggregator factory that always creates a vaultshare
// Aggregator without gating on VaultFastPathGetSecretsEnabled. The real vaultshare factory is gated,
// so tests that exercise the aggregation path need this ungated variant.
func withUngatedVaultAggregatorFactory(f int) aggregatorFactoryOpt {
	return func(cfg *aggregatorFactoryCfg) {
		cfg.factory = func(_ context.Context, _ commoncap.CapabilityRequest) request.ResponseAggregator {
			return vaultshare.NewAggregator(f+1, 2*f+1)
		}
	}
}

// testFastPathRemoteExecutableCapability is a variant of testRemoteExecutableCapability that
// accepts an aggregator factory override. It keeps the original harness unmodified.
func testFastPathRemoteExecutableCapability(ctx context.Context, t *testing.T, underlying commoncap.ExecutableCapability, numWorkflowPeers int, workflowDonF uint8, workflowNodeTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration,
	method func(ctx context.Context, caller commoncap.ExecutableCapability), waitForExecuteCalls bool, opts ...aggregatorFactoryOpt) {
	lggr := logger.Test(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(NewPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

	capDonInfo := commoncap.DON{
		ID:      2,
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "vault@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Vault Target",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeerID := p2ptypes.PeerID{}
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(NewPeerID())))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
		F:       workflowDonF,
	}

	broker := newTestAsyncMessageBroker(t, 1000)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	capabilityNodes := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := executable.NewServer(capInfo.ID, "", capabilityPeer, capabilityDispatcher, lggr)
		cfg := &commoncap.RemoteExecutableConfig{
			RequestHashExcludedAttributes: []string{},
			RequestTimeout:                capabilityNodeResponseTimeout,
			ServerMaxParallelRequests:     10,
		}
		require.NoError(t, capabilityNode.SetConfig(cfg, underlying, capInfo, capDonInfo, workflowDONs, nil))
		servicetest.Run(t, capabilityNode)
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
		capabilityNodes[i] = capabilityNode
	}

	cfg := &aggregatorFactoryCfg{}
	for _, opt := range opts {
		opt(cfg)
	}

	workflowNodes := make([]commoncap.ExecutableCapability, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := executable.NewClient(capInfo.ID, "", workflowPeerDispatcher, lggr)
		err := workflowNode.SetConfig(capInfo, workflowDonInfo, workflowNodeTimeout, nil, nil, 0, cfg.factory)
		require.NoError(t, err)
		servicetest.Run(t, workflowNode)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
	}

	servicetest.Run(t, broker)

	wg := &sync.WaitGroup{}
	wg.Add(len(workflowNodes))

	for _, caller := range workflowNodes {
		go func(caller commoncap.ExecutableCapability) {
			defer wg.Done()
			method(ctx, caller)
		}(caller)
	}
	if waitForExecuteCalls {
		wg.Wait()
	}
}
