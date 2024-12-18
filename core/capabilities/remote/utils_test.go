package remote_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	capID1   = "cap1"
	capID2   = "cap2"
	donID1   = uint32(1)
	payload1 = "hello world"
	payload2 = "goodbye world"
)

func TestValidateMessage(t *testing.T) {
	privKey1, peerID1 := newKeyPair(t)
	_, peerID2 := newKeyPair(t)

	// valid
	p2pMsg := encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	body, err := remote.ValidateMessage(p2pMsg, peerID2)
	require.NoError(t, err)
	require.Equal(t, peerID1[:], body.Sender)
	require.Equal(t, payload1, string(body.Payload))

	// invalid sender
	p2pMsg = encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	p2pMsg.Sender = peerID2
	_, err = remote.ValidateMessage(p2pMsg, peerID2)
	require.Error(t, err)

	// invalid receiver
	p2pMsg = encodeAndSign(t, privKey1, peerID1, peerID2, capID1, donID1, []byte(payload1))
	_, err = remote.ValidateMessage(p2pMsg, peerID1)
	require.Error(t, err)
}

func newKeyPair(t *testing.T) (ed25519.PrivateKey, ragetypes.PeerID) {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	peerID, err := ragetypes.PeerIDFromPrivateKey(privKey)
	require.NoError(t, err)
	return privKey, peerID
}

func encodeAndSign(t *testing.T, senderPrivKey ed25519.PrivateKey, senderID p2ptypes.PeerID, receiverID p2ptypes.PeerID, capabilityID string, donID uint32, payload []byte) p2ptypes.Message {
	body := remotetypes.MessageBody{
		Sender:          senderID[:],
		Receiver:        receiverID[:],
		CapabilityId:    capabilityID,
		CapabilityDonId: donID,
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
		Sender:  senderID,
		Payload: rawMsg,
	}
}

func TestToPeerID(t *testing.T) {
	id, err := remote.ToPeerID([]byte("12345678901234567890123456789012"))
	require.NoError(t, err)
	require.Equal(t, "12D3KooWD8QYTQVYjB6oog4Ej8PcPpqTrPRnxLQap8yY8KUQRVvq", id.String())
}

func TestSanitizeLogString(t *testing.T) {
	require.Equal(t, "hello", remote.SanitizeLogString("hello"))
	require.Equal(t, "[UNPRINTABLE] 0a", remote.SanitizeLogString("\n"))

	longString := ""
	for i := 0; i < 100; i++ {
		longString += "aa-aa-aa-"
	}
	require.Equal(t, longString[:256]+" [TRUNCATED]", remote.SanitizeLogString(longString))
}
