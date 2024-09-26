package cache

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache/mocks"
)

var address = cciptypes.Address(common.HexToAddress("0x1234567890123456789012345678901234567890").String())

func Test_ObservedChainStateSkipErrors(t *testing.T) {
	mockedHealthcheck := mocks.NewChainHealthcheck(t)
	mockedHealthcheck.On("IsHealthy", mock.Anything).Return(false, fmt.Errorf("error"))

	observedChainState := NewObservedChainHealthCheck(
		mockedHealthcheck,
		"plugin",
		10,
		20,
		address,
	)

	_, err := observedChainState.IsHealthy(tests.Context(t))
	assert.Error(t, err)
	assert.Equal(t, float64(0), testutil.ToFloat64(laneHealthStatus.WithLabelValues("plugin", "10", "20", "0x1234567890123456789012345678901234567890")))
}

func Test_ObservedChainStateReportsStatus(t *testing.T) {
	mockedHealthcheck := mocks.NewChainHealthcheck(t)
	mockedHealthcheck.On("IsHealthy", mock.Anything).Return(true, nil).Once()

	observedChainState := NewObservedChainHealthCheck(
		mockedHealthcheck,
		"plugin",
		10,
		20,
		address,
	)

	health, err := observedChainState.IsHealthy(tests.Context(t))
	require.NoError(t, err)
	assert.True(t, health)
	assert.Equal(t, float64(1), testutil.ToFloat64(laneHealthStatus.WithLabelValues("plugin", "10", "20", "0x1234567890123456789012345678901234567890")))

	// Mark as unhealthy
	mockedHealthcheck.On("IsHealthy", mock.Anything).Return(false, nil).Once()

	health, err = observedChainState.IsHealthy(tests.Context(t))
	require.NoError(t, err)
	assert.False(t, health)
	assert.Equal(t, float64(0), testutil.ToFloat64(laneHealthStatus.WithLabelValues("plugin", "10", "20", "0x1234567890123456789012345678901234567890")))
}
