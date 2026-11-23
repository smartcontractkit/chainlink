package evm

import (
	"testing"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	// Config is created during initialization, the following functions may panic:
	// MustGetABI
	// mustGetMethodName
	// mustGetEventName
}
