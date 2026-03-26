package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestExecutionHelper_ConfidentialHTTPPerWorkflowLimit(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	lf := limits.Factory{Logger: lggr}

	// Configure per-workflow confidential-http call limit to 1
	cfgFn := func(w *cresettings.Workflows) {
		w.ConfidentialHTTP.CallLimit = settings.Int(1)
	}

	limiters, err := NewLimiters(lf, cfgFn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = limiters.Close() })

	// Build ExecutionHelper and initialize its call limiters from EngineLimiters
	exec := &ExecutionHelper{}
	exec.initLimiters(limiters)

	// Grab the configured limiter instance for confidential-http SendRequest
	capCallValue := capCall{name: "confidential-http", method: "SendRequest"}
	limiter, ok := exec.callLimiters[capCallValue]
	require.True(t, ok, "expected confidential-http limiter to be configured")

	// Prime the internal callCounts to simulate one prior call so the next call will exceed the configured limit (1)
	exec.callCounts = make(map[limits.Limiter[int]]int)
	exec.callCounts[limiter] = 1

	// Prepare a request that will parse to capName == "confidential-http" and method == "SendRequest"
	req := &sdk.CapabilityRequest{
		Id:         "confidential-http",
		Method:     "SendRequest",
		CallbackId: 1,
	}

	// Call and expect an error from the bound limiter (limit exceeded)
	_, err = exec.CallCapability(t.Context(), req)
	require.Error(t, err, "expected CallCapability to fail when per-workflow confidential-http call limit is exceeded")
}

func TestExecutionHelper_FailedCapabilityIDs(t *testing.T) {
	t.Parallel()

	t.Run("empty when no failures", func(t *testing.T) {
		exec := &ExecutionHelper{}
		assert.Empty(t, exec.FailedCapabilityIDs())
	})

	t.Run("records single failed capability", func(t *testing.T) {
		exec := &ExecutionHelper{}
		exec.failedCapabilityIDs.Store("evm@1.0.0", struct{}{})
		assert.Equal(t, []string{"evm@1.0.0"}, exec.FailedCapabilityIDs())
	})

	t.Run("deduplicates same capability", func(t *testing.T) {
		exec := &ExecutionHelper{}
		exec.failedCapabilityIDs.Store("evm@1.0.0", struct{}{})
		exec.failedCapabilityIDs.Store("evm@1.0.0", struct{}{})
		assert.Equal(t, []string{"evm@1.0.0"}, exec.FailedCapabilityIDs())
	})

	t.Run("records multiple failed capabilities sorted", func(t *testing.T) {
		exec := &ExecutionHelper{}
		exec.failedCapabilityIDs.Store("evm@1.0.0", struct{}{})
		exec.failedCapabilityIDs.Store("confidential-http@1.0.0", struct{}{})
		assert.Equal(t, []string{"confidential-http@1.0.0", "evm@1.0.0"}, exec.FailedCapabilityIDs())
	})
}
