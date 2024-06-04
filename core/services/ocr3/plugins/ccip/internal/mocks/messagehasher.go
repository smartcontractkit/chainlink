package mocks

import (
	"github.com/smartcontractkit/ccipocr3/internal/model"
)

type MessageHasher struct{}

func NewMessageHasher() *MessageHasher {
	return &MessageHasher{}
}

func (m *MessageHasher) Hash(msg model.CCIPMsg) (model.Bytes32, error) {
	return msg.ID, nil
}
