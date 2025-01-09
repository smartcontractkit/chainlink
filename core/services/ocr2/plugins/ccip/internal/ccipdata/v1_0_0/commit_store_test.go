package v1_0_0

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCommitReportEncoding(t *testing.T) {
	ctx := testutils.Context(t)
	report := cciptypes.CommitStoreReport{
		TokenPrices: []cciptypes.TokenPrice{
			{
				Token: cciptypes.Address(utils.RandomAddress().String()),
				Value: big.NewInt(9e18),
			},
		},
		GasPrices: []cciptypes.GasPrice{
			{
				DestChainSelector: rand.Uint64(),
				Value:             big.NewInt(2000e9),
			},
		},
		MerkleRoot: [32]byte{123},
		Interval:   cciptypes.CommitStoreInterval{Min: 1, Max: 10},
	}

	c, err := NewCommitStore(logger.TestLogger(t), utils.RandomAddress(), nil, mocks.NewLogPoller(t))
	assert.NoError(t, err)

	encodedReport, err := c.EncodeCommitReport(ctx, report)
	require.NoError(t, err)
	assert.Greater(t, len(encodedReport), 0)

	decodedReport, err := c.DecodeCommitReport(ctx, encodedReport)
	require.NoError(t, err)
	require.Equal(t, report.TokenPrices, decodedReport.TokenPrices)
	require.Equal(t, report, decodedReport)
}
