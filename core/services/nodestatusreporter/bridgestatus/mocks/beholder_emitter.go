package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// BeholderEmitter is a mock implementation of beholder.Emitter for testing
type BeholderEmitter struct {
	mock.Mock
}

func (m *BeholderEmitter) Emit(ctx context.Context, body []byte, attrKVs ...any) error {
	args := m.Called(ctx, body, attrKVs)
	return args.Error(0)
}

// NewBeholderEmitter creates a new beholder emitter mock for testing
func NewBeholderEmitter() *BeholderEmitter {
	return &BeholderEmitter{}
}
