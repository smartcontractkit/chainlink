package remote

import (
	"bytes"
	"crypto/ed25519"
	"fmt"

	"google.golang.org/protobuf/proto"

	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func ValidateMessage(msg p2ptypes.Message, expectedReceiver p2ptypes.PeerID) (*remotetypes.MessageBody, error) {
	var topLevelMessage remotetypes.Message
	err := proto.Unmarshal(msg.Payload, &topLevelMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message, err: %v", err)
	}
	var body remotetypes.MessageBody
	err = proto.Unmarshal(topLevelMessage.Body, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message body, err: %v", err)
	}
	if len(body.Sender) != p2ptypes.PeerIDLength || len(body.Receiver) != p2ptypes.PeerIDLength {
		return &body, fmt.Errorf("invalid sender length (%d) or receiver length (%d)", len(body.Sender), len(body.Receiver))
	}
	if !ed25519.Verify(body.Sender, topLevelMessage.Body, topLevelMessage.Signature) {
		return &body, fmt.Errorf("failed to verify message signature")
	}
	// NOTE we currently don't support relaying messages so the p2p message sender needs to be the message author
	if !bytes.Equal(body.Sender, msg.Sender[:]) {
		return &body, fmt.Errorf("sender in message body does not match sender of p2p message")
	}
	if !bytes.Equal(body.Receiver, expectedReceiver[:]) {
		return &body, fmt.Errorf("receiver in message body does not match expected receiver")
	}
	return &body, nil
}
