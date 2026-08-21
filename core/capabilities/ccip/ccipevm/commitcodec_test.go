package ccipevm

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

var randomCommitReport = func() ccipocr3common.CommitPluginReport {
	return ccipocr3common.CommitPluginReport{
		BlessedMerkleRoots: []ccipocr3common.MerkleRootChain{
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      ccipocr3common.ChainSelector(rand.Uint64()),
				SeqNumsRange: ccipocr3common.NewSeqNumRange(
					ccipocr3common.SeqNum(rand.Uint64()),
					ccipocr3common.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      ccipocr3common.ChainSelector(rand.Uint64()),
				SeqNumsRange: ccipocr3common.NewSeqNumRange(
					ccipocr3common.SeqNum(rand.Uint64()),
					ccipocr3common.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
		},
		UnblessedMerkleRoots: []ccipocr3common.MerkleRootChain{
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      ccipocr3common.ChainSelector(rand.Uint64()),
				SeqNumsRange: ccipocr3common.NewSeqNumRange(
					ccipocr3common.SeqNum(rand.Uint64()),
					ccipocr3common.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      ccipocr3common.ChainSelector(rand.Uint64()),
				SeqNumsRange: ccipocr3common.NewSeqNumRange(
					ccipocr3common.SeqNum(rand.Uint64()),
					ccipocr3common.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
		},
		PriceUpdates: ccipocr3common.PriceUpdates{
			TokenPriceUpdates: []ccipocr3common.TokenPrice{
				{
					TokenID: ccipocr3common.UnknownEncodedAddress(utils.RandomAddress().String()),
					Price:   ccipocr3common.NewBigInt(utils.RandUint256()),
				},
			},
			GasPriceUpdates: []ccipocr3common.GasPriceChain{
				{GasPrice: ccipocr3common.NewBigInt(utils.RandUint256()), ChainSel: ccipocr3common.ChainSelector(rand.Uint64())},
				{GasPrice: ccipocr3common.NewBigInt(utils.RandUint256()), ChainSel: ccipocr3common.ChainSelector(rand.Uint64())},
				{GasPrice: ccipocr3common.NewBigInt(utils.RandUint256()), ChainSel: ccipocr3common.ChainSelector(rand.Uint64())},
			},
		},
		RMNSignatures: []ccipocr3common.RMNECDSASignature{
			{R: utils.RandomBytes32(), S: utils.RandomBytes32()},
			{R: utils.RandomBytes32(), S: utils.RandomBytes32()},
		},
	}
}

func TestCommitPluginCodecV1(t *testing.T) {
	testCases := []struct {
		name   string
		report func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport
		expErr bool
	}{
		{
			name: "base report",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				return report
			},
		},
		{
			name: "empty token address",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.PriceUpdates.TokenPriceUpdates[0].TokenID = ""
				return report
			},
			expErr: true,
		},
		{
			name: "empty merkle root",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.BlessedMerkleRoots[0].MerkleRoot = ccipocr3common.Bytes32{}
				report.UnblessedMerkleRoots[0].MerkleRoot = ccipocr3common.Bytes32{}
				return report
			},
		},
		{
			name: "zero token price",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.PriceUpdates.TokenPriceUpdates[0].Price = ccipocr3common.NewBigInt(big.NewInt(0))
				return report
			},
		},
		{
			name: "zero gas price",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.PriceUpdates.GasPriceUpdates[0].GasPrice = ccipocr3common.NewBigInt(big.NewInt(0))
				return report
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report := tc.report(randomCommitReport())
			commitCodec := NewCommitPluginCodecV1()
			ctx := t.Context()
			encodedReport, err := commitCodec.Encode(ctx, report)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			decodedReport, err := commitCodec.Decode(ctx, encodedReport)
			require.NoError(t, err)
			require.Equal(t, report, decodedReport)
		})
	}
}

func BenchmarkCommitPluginCodecV1_Encode(b *testing.B) {
	commitCodec := NewCommitPluginCodecV1()
	ctx := b.Context()

	rep := randomCommitReport()
	for b.Loop() {
		_, err := commitCodec.Encode(ctx, rep)
		require.NoError(b, err)
	}
}

func BenchmarkCommitPluginCodecV1_Decode(b *testing.B) {
	commitCodec := NewCommitPluginCodecV1()
	ctx := b.Context()
	encodedReport, err := commitCodec.Encode(ctx, randomCommitReport())
	require.NoError(b, err)

	for b.Loop() {
		_, err := commitCodec.Decode(ctx, encodedReport)
		require.NoError(b, err)
	}
}
