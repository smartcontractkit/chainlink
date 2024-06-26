package remote_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	capId1   = "cap1"
	capId2   = "cap2"
	donId1   = "donA"
	payload1 = "hello world"
	payload2 = "goodbye world"
)

func TestValidateMessage(t *testing.T) {
	privKey1, peerId1 := newKeyPair(t)
	_, peerId2 := newKeyPair(t)

	// valid
	p2pMsg := encodeAndSign(t, privKey1, peerId1, peerId2, capId1, donId1, []byte(payload1))
	body, err := remote.ValidateMessage(p2pMsg, peerId2)
	require.NoError(t, err)
	require.Equal(t, peerId1[:], body.Sender)
	require.Equal(t, payload1, string(body.Payload))

	// invalid sender
	p2pMsg = encodeAndSign(t, privKey1, peerId1, peerId2, capId1, donId1, []byte(payload1))
	p2pMsg.Sender = peerId2
	_, err = remote.ValidateMessage(p2pMsg, peerId2)
	require.Error(t, err)

	// invalid receiver
	p2pMsg = encodeAndSign(t, privKey1, peerId1, peerId2, capId1, donId1, []byte(payload1))
	_, err = remote.ValidateMessage(p2pMsg, peerId1)
	require.Error(t, err)
}

func newKeyPair(t *testing.T) (ed25519.PrivateKey, ragetypes.PeerID) {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	peerID, err := ragetypes.PeerIDFromPrivateKey(privKey)
	require.NoError(t, err)
	return privKey, peerID
}

func encodeAndSign(t *testing.T, senderPrivKey ed25519.PrivateKey, senderId p2ptypes.PeerID, receiverId p2ptypes.PeerID, capabilityId string, donId string, payload []byte) p2ptypes.Message {
	body := remotetypes.MessageBody{
		Sender:          senderId[:],
		Receiver:        receiverId[:],
		CapabilityId:    capabilityId,
		CapabilityDonId: donId,
		Payload:         payload,
	}
	rawBody, err := proto.Marshal(&body)
	require.NoError(t, err)
	signature := ed25519.Sign(senderPrivKey, rawBody)

	msg := remotetypes.Message{
		Signature: signature,
		Body:      rawBody,
	}
	rawMsg, err := proto.Marshal(&msg)
	require.NoError(t, err)

	return p2ptypes.Message{
		Sender:  senderId,
		Payload: rawMsg,
	}
}

func TestToPeerID(t *testing.T) {
	id := remote.ToPeerID([]byte("12345678901234567890123456789012"))
	require.Equal(t, "12D3KooWD8QYTQVYjB6oog4Ej8PcPpqTrPRnxLQap8yY8KUQRVvq", id.String())
}

func TestDefaultModeAggregator_Aggregate(t *testing.T) {
	val, err := values.Wrap(triggerEvent1)
	require.NoError(t, err)
	capResponse1 := commoncap.CapabilityResponse{
		Value: val,
		Err:   nil,
	}
	marshaled1, err := pb.MarshalCapabilityResponse(capResponse1)
	require.NoError(t, err)

	val2, err := values.Wrap(triggerEvent2)
	require.NoError(t, err)
	capResponse2 := commoncap.CapabilityResponse{
		Value: val2,
		Err:   nil,
	}
	marshaled2, err := pb.MarshalCapabilityResponse(capResponse2)
	require.NoError(t, err)

	agg := remote.NewDefaultModeAggregator(2)
	_, err = agg.Aggregate("", [][]byte{marshaled1})
	require.Error(t, err)

	_, err = agg.Aggregate("", [][]byte{marshaled1, marshaled2})
	require.Error(t, err)

	res, err := agg.Aggregate("", [][]byte{marshaled1, marshaled2, marshaled1})
	require.NoError(t, err)
	require.Equal(t, res, capResponse1)
}
