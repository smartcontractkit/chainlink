package mocks

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type MessageHasher struct{}

func NewMessageHasher() *MessageHasher {
	return &MessageHasher{}
}

func (m *MessageHasher) Hash(ctx context.Context, msg cciptypes.CCIPMsg) (cciptypes.Bytes32, error) {
	// simply return the msg id as bytes32
	var b32 [32]byte
	copy(b32[:], msg.ID)
	return b32, nil
}
