package mocks

import (
	"github.com/smartcontractkit/ccipocr3/internal/model"
)

type NopMessageHasher struct{}

func NewNopMessageHasher() *NopMessageHasher {
	return &NopMessageHasher{}
}

func (m *NopMessageHasher) Hash(msg model.CCIPMsg) (model.Bytes32, error) {
	return msg.ID, nil
}
