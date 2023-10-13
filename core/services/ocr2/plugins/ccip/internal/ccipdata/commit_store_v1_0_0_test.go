package ccipdata

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestCommitReportEncodingV1_0_0(t *testing.T) {
	report := CommitStoreReport{
		TokenPrices: []TokenPrice{
			{
				Token: utils.RandomAddress(),
				Value: big.NewInt(9e18),
			},
		},
		GasPrices: []GasPrice{
			{
				DestChainSelector: rand.Uint64(),
				Value:             big.NewInt(2000e9),
			},
		},
		MerkleRoot: [32]byte{123},
		Interval:   CommitStoreInterval{Min: 1, Max: 10},
	}

	lp := mocks.NewLogPoller(t)
	lp.On("RegisterFilter", mock.Anything).Return(nil)

	c, err := NewCommitStoreV1_0_0(logger.TestLogger(t), utils.RandomAddress(), nil, lp, nil)
	assert.NoError(t, err)

	encodedReport, err := c.EncodeCommitReport(report)
	require.NoError(t, err)
	assert.Greater(t, len(encodedReport), 0)

	decodedReport, err := c.DecodeCommitReport(encodedReport)
	require.NoError(t, err)
	require.Equal(t, report.TokenPrices, decodedReport.TokenPrices)
	require.Equal(t, report, decodedReport)
}
