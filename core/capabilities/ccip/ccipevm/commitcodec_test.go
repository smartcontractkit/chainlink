package ccipevm

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

var randomCommitReport = func() ccipocr3.CommitPluginReport {
	return ccipocr3.CommitPluginReport{
		BlessedMerkleRoots: []ccipocr3.MerkleRootChain{
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      ccipocr3.ChainSelector(rand.Uint64()),
				SeqNumsRange: ccipocr3.NewSeqNumRange(
					ccipocr3.SeqNum(rand.Uint64()),
					ccipocr3.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      ccipocr3.ChainSelector(rand.Uint64()),
				SeqNumsRange: ccipocr3.NewSeqNumRange(
					ccipocr3.SeqNum(rand.Uint64()),
					ccipocr3.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
		},
		UnblessedMerkleRoots: []ccipocr3.MerkleRootChain{
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      ccipocr3.ChainSelector(rand.Uint64()),
				SeqNumsRange: ccipocr3.NewSeqNumRange(
					ccipocr3.SeqNum(rand.Uint64()),
					ccipocr3.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      ccipocr3.ChainSelector(rand.Uint64()),
				SeqNumsRange: ccipocr3.NewSeqNumRange(
					ccipocr3.SeqNum(rand.Uint64()),
					ccipocr3.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
		},
		PriceUpdates: ccipocr3.PriceUpdates{
			TokenPriceUpdates: []ccipocr3.TokenPrice{
				{
					TokenID: ccipocr3.UnknownEncodedAddress(utils.RandomAddress().String()),
					Price:   ccipocr3.NewBigInt(utils.RandUint256()),
				},
			},
			GasPriceUpdates: []ccipocr3.GasPriceChain{
				{GasPrice: ccipocr3.NewBigInt(utils.RandUint256()), ChainSel: ccipocr3.ChainSelector(rand.Uint64())},
				{GasPrice: ccipocr3.NewBigInt(utils.RandUint256()), ChainSel: ccipocr3.ChainSelector(rand.Uint64())},
				{GasPrice: ccipocr3.NewBigInt(utils.RandUint256()), ChainSel: ccipocr3.ChainSelector(rand.Uint64())},
			},
		},
		RMNSignatures: []ccipocr3.RMNECDSASignature{
			{R: utils.RandomBytes32(), S: utils.RandomBytes32()},
			{R: utils.RandomBytes32(), S: utils.RandomBytes32()},
		},
	}
}

func TestCommitPluginCodecV1(t *testing.T) {
	testCases := []struct {
		name   string
		report func(report ccipocr3.CommitPluginReport) ccipocr3.CommitPluginReport
		expErr bool
	}{
		{
			name: "base report",
			report: func(report ccipocr3.CommitPluginReport) ccipocr3.CommitPluginReport {
				return report
			},
		},
		{
			name: "empty token address",
			report: func(report ccipocr3.CommitPluginReport) ccipocr3.CommitPluginReport {
				report.PriceUpdates.TokenPriceUpdates[0].TokenID = ""
				return report
			},
			expErr: true,
		},
		{
			name: "empty merkle root",
			report: func(report ccipocr3.CommitPluginReport) ccipocr3.CommitPluginReport {
				report.BlessedMerkleRoots[0].MerkleRoot = ccipocr3.Bytes32{}
				report.UnblessedMerkleRoots[0].MerkleRoot = ccipocr3.Bytes32{}
				return report
			},
		},
		{
			name: "zero token price",
			report: func(report ccipocr3.CommitPluginReport) ccipocr3.CommitPluginReport {
				report.PriceUpdates.TokenPriceUpdates[0].Price = ccipocr3.NewBigInt(big.NewInt(0))
				return report
			},
		},
		{
			name: "zero gas price",
			report: func(report ccipocr3.CommitPluginReport) ccipocr3.CommitPluginReport {
				report.PriceUpdates.GasPriceUpdates[0].GasPrice = ccipocr3.NewBigInt(big.NewInt(0))
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
