package evm

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestConfig(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	// Config is created during initialization, the following functions may panic:
	// MustGetABI
	// mustGetMethodName
	// mustGetEventName
}
