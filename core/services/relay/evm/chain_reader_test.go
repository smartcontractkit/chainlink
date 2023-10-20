package evm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocklogpoller "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestChainReaderStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	lp := mocklogpoller.NewLogPoller(t)
	chainReader, err := evm.NewChainReaderService(lggr, lp)
	require.NoError(t, err)
	require.NotNil(t, chainReader)
	err = chainReader.Start(testutils.Context(t))
	assert.NoError(t, err)
	err = chainReader.Close()
	assert.NoError(t, err)
}
