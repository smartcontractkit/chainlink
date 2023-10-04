package utils

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils/tests"
)

// Deprecated: use tests.Context
func Context(t *testing.T) context.Context {
	return tests.Context(t)
}
