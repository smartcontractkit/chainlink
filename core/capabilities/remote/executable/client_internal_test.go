package executable

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// Regression tests for the direct-read response race: a response arriving before
// the request is registered in requestIDToCallerRequest is dropped as an unknown
// message ID, starving the 2F+1 response quorum.

// Test_Client_RequestRegisteredBeforeSend asserts the request is present in
// requestIDToCallerRequest before the first dispatcher.Send call.
func Test_Client_RequestRegisteredBeforeSend(t *testing.T) {
	t.Parallel()

	var registeredAtSend bool
	cl, _, capabilityRequest := newInlineDispatcherFixture(t, func(c *client, msg *remotetypes.MessageBody) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		_, registeredAtSend = c.requestIDToCallerRequest[string(msg.MessageId)]
	})

	resp, err := cl.Execute(t.Context(), capabilityRequest)
	require.NoError(t, err)
	assert.True(t, registeredAtSend, "request must be registered in requestIDToCallerRequest before the first Send")

	unwrapped, err := resp.Value.Unwrap()
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"response": "response1"}, unwrapped)
}

// Test_Client_ResponseArrivingBeforeRegistrationIsCounted asserts that a response
// delivered synchronously from inside Send is still counted towards the quorum.
func Test_Client_ResponseArrivingBeforeRegistrationIsCounted(t *testing.T) {
	t.Parallel()

	cl, _, capabilityRequest := newInlineDispatcherFixture(t, nil)

	resp, err := cl.Execute(t.Context(), capabilityRequest)
	require.NoError(t, err)

	unwrapped, err := resp.Value.Unwrap()
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"response": "response1"}, unwrapped)
}

// newInlineDispatcherFixture returns a started client wired to an
// inlineResponseDispatcher, plus a capability request to execute against it.
func newInlineDispatcherFixture(t *testing.T, sendFn func(c *client, msg *remotetypes.MessageBody)) (*client, *inlineResponseDispatcher, commoncap.CapabilityRequest) {
	t.Helper()

	capabilityPeer := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeer.UnmarshalText([]byte(newInternalPeerID())))

	capDonInfo := commoncap.DON{ID: 2, Members: []p2ptypes.PeerID{capabilityPeer}, F: 0}
	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}
	workflowPeer := p2ptypes.PeerID{}
	require.NoError(t, workflowPeer.UnmarshalText([]byte(newInternalPeerID())))
	workflowDonInfo := commoncap.DON{ID: 1, Members: []p2ptypes.PeerID{workflowPeer}, F: 0}

	m, err := values.NewMap(map[string]any{"response": "response1"})
	require.NoError(t, err)
	responsePayload, err := pb.MarshalCapabilityResponse(commoncap.CapabilityResponse{Value: m})
	require.NoError(t, err)

	d := &inlineResponseDispatcher{sendFn: sendFn, responsePayload: responsePayload}
	cl := NewClient(capInfo.ID, "", d, logger.Test(t))
	d.cl = cl
	require.NoError(t, cl.SetConfig(capInfo, workflowDonInfo, 200*time.Millisecond, nil, 0))
	servicetest.Run(t, cl)

	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0",
			WorkflowExecutionID: "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed",
			ReferenceID:         "step1",
		},
		Inputs: m,
	}
	return cl, d, capabilityRequest
}

// inlineResponseDispatcher invokes sendFn (if set) and feeds a response back into
// the client from inside Send, before Send returns -- reproducing the fastest
// possible remote reply.
type inlineResponseDispatcher struct {
	cl              *client
	sendFn          func(c *client, msg *remotetypes.MessageBody)
	responsePayload []byte
}

func (d *inlineResponseDispatcher) Send(peerID p2ptypes.PeerID, msg *remotetypes.MessageBody) error {
	if d.sendFn != nil {
		d.sendFn(d.cl, msg)
	}
	resp := &remotetypes.MessageBody{
		CapabilityId:    msg.CapabilityId,
		CapabilityDonId: msg.CapabilityDonId,
		CallerDonId:     msg.CallerDonId,
		Method:          msg.Method,
		MessageId:       msg.MessageId,
		Sender:          peerID[:],
		Payload:         d.responsePayload,
	}
	d.cl.Receive(context.Background(), resp)
	return nil
}

func (d *inlineResponseDispatcher) SetReceiver(capabilityID string, donID uint32, receiver remotetypes.Receiver) error {
	return nil
}

func (d *inlineResponseDispatcher) RemoveReceiver(capabilityID string, donID uint32) {}

func (d *inlineResponseDispatcher) SetReceiverForMethod(capabilityID string, donID uint32, method string, receiver remotetypes.Receiver) error {
	return nil
}

func (d *inlineResponseDispatcher) RemoveReceiverForMethod(capabilityID string, donID uint32, method string) {
}

func (d *inlineResponseDispatcher) Start(ctx context.Context) error { return nil }

func (d *inlineResponseDispatcher) Close() error { return nil }

func (d *inlineResponseDispatcher) Ready() error { return nil }

func (d *inlineResponseDispatcher) HealthReport() map[string]error { return nil }

func (d *inlineResponseDispatcher) Name() string { return "inlineResponseDispatcher" }

func newInternalPeerID() string {
	var privKey [32]byte
	if _, err := rand.Read(privKey[:]); err != nil {
		panic(err)
	}
	peerID := append([]byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}, privKey[:]...)
	return base58.Encode(peerID)
}
