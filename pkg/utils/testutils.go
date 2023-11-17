package utils

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

// Deprecated: use tests.Context
func Context(t *testing.T) context.Context {
	return tests.Context(t)
}
