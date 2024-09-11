package remote

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"unicode"

	"google.golang.org/protobuf/proto"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	maxLoggedStringLen = 256
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

func ToPeerID(peerID []byte) (p2ptypes.PeerID, error) {
	if len(peerID) != p2ptypes.PeerIDLength {
		return p2ptypes.PeerID{}, fmt.Errorf("invalid peer ID length: %d", len(peerID))
	}

	var id p2ptypes.PeerID
	copy(id[:], peerID)
	return id, nil
}

// Default MODE Aggregator needs a configurable number of identical responses for aggregation to succeed
type defaultModeAggregator struct {
	minIdenticalResponses uint32
}

var _ remotetypes.Aggregator = &defaultModeAggregator{}

func NewDefaultModeAggregator(minIdenticalResponses uint32) *defaultModeAggregator {
	return &defaultModeAggregator{
		minIdenticalResponses: minIdenticalResponses,
	}
}

func (a *defaultModeAggregator) Aggregate(_ string, responses [][]byte) (commoncap.TriggerResponse, error) {
	found, err := AggregateModeRaw(responses, a.minIdenticalResponses)
	if err != nil {
		return commoncap.TriggerResponse{}, fmt.Errorf("failed to aggregate responses, err: %w", err)
	}

	unmarshaled, err := pb.UnmarshalTriggerResponse(found)
	if err != nil {
		return commoncap.TriggerResponse{}, fmt.Errorf("failed to unmarshal aggregated responses, err: %w", err)
	}
	return unmarshaled, nil
}

func AggregateModeRaw(elemList [][]byte, minIdenticalResponses uint32) ([]byte, error) {
	hashToCount := make(map[string]uint32)
	var found []byte
	for _, elem := range elemList {
		hasher := sha256.New()
		hasher.Write(elem)
		sha := hex.EncodeToString(hasher.Sum(nil))
		hashToCount[sha]++
		if hashToCount[sha] >= minIdenticalResponses {
			found = elem
			// update in case we find another elem with an even higher count
			minIdenticalResponses = hashToCount[sha]
		}
	}
	if found == nil {
		return nil, errors.New("not enough identical responses found")
	}
	return found, nil
}

func SanitizeLogString(s string) string {
	tooLongSuffix := ""
	if len(s) > maxLoggedStringLen {
		s = s[:maxLoggedStringLen]
		tooLongSuffix = " [TRUNCATED]"
	}
	for i := 0; i < len(s); i++ {
		if !unicode.IsPrint(rune(s[i])) {
			return "[UNPRINTABLE] " + hex.EncodeToString([]byte(s)) + tooLongSuffix
		}
	}
	return s + tooLongSuffix
}
