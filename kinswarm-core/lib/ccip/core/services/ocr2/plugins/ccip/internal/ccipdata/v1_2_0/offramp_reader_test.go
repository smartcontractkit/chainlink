package v1_2_0_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

func TestExecutionReportEncodingV120(t *testing.T) {
	// Note could consider some fancier testing here (fuzz/property)
	// but I think that would essentially be testing geth's abi library
	// as our encode/decode is a thin wrapper around that.
	report := cciptypes.ExecReport{
		Messages:          []cciptypes.EVM2EVMMessage{},
		OffchainTokenData: [][][]byte{{}},
		Proofs:            [][32]byte{testutils.Random32Byte()},
		ProofFlagBits:     big.NewInt(133),
	}

	feeEstimatorConfig := mocks.NewFeeEstimatorConfigReader(t)

	offRamp, err := v1_2_0.NewOffRamp(logger.TestLogger(t), utils.RandomAddress(), nil, lpmocks.NewLogPoller(t), nil, nil, feeEstimatorConfig)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	encodeExecutionReport, err := offRamp.EncodeExecutionReport(ctx, report)
	require.NoError(t, err)
	decodeCommitReport, err := offRamp.DecodeExecutionReport(ctx, encodeExecutionReport)
	require.NoError(t, err)
	require.Equal(t, report.Proofs, decodeCommitReport.Proofs)
	require.Equal(t, report, decodeCommitReport)
}
