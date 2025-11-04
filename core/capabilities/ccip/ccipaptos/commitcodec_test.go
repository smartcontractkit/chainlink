package ccipaptos

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

var randomCommitReport = func() cciptypes.CommitPluginReport {
	return cciptypes.CommitPluginReport{
		BlessedMerkleRoots: []cciptypes.MerkleRootChain{
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      cciptypes.ChainSelector(rand.Uint64()),
				SeqNumsRange: cciptypes.NewSeqNumRange(
					cciptypes.SeqNum(rand.Uint64()),
					cciptypes.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      cciptypes.ChainSelector(rand.Uint64()),
				SeqNumsRange: cciptypes.NewSeqNumRange(
					cciptypes.SeqNum(rand.Uint64()),
					cciptypes.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
		},
		UnblessedMerkleRoots: []cciptypes.MerkleRootChain{
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      cciptypes.ChainSelector(rand.Uint64()),
				SeqNumsRange: cciptypes.NewSeqNumRange(
					cciptypes.SeqNum(rand.Uint64()),
					cciptypes.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
			{
				OnRampAddress: common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				ChainSel:      cciptypes.ChainSelector(rand.Uint64()),
				SeqNumsRange: cciptypes.NewSeqNumRange(
					cciptypes.SeqNum(rand.Uint64()),
					cciptypes.SeqNum(rand.Uint64()),
				),
				MerkleRoot: utils.RandomBytes32(),
			},
		},
		PriceUpdates: cciptypes.PriceUpdates{
			TokenPriceUpdates: []cciptypes.TokenPrice{
				{
					TokenID: cciptypes.UnknownEncodedAddress(generateAddressString()),
					Price:   cciptypes.NewBigInt(utils.RandUint256()),
				},
			},
			GasPriceUpdates: []cciptypes.GasPriceChain{
				{GasPrice: cciptypes.NewBigInt(utils.RandUint256()), ChainSel: cciptypes.ChainSelector(rand.Uint64())},
				{GasPrice: cciptypes.NewBigInt(utils.RandUint256()), ChainSel: cciptypes.ChainSelector(rand.Uint64())},
				{GasPrice: cciptypes.NewBigInt(utils.RandUint256()), ChainSel: cciptypes.ChainSelector(rand.Uint64())},
			},
		},
		RMNSignatures: []cciptypes.RMNECDSASignature{
			{R: utils.RandomBytes32(), S: utils.RandomBytes32()},
			{R: utils.RandomBytes32(), S: utils.RandomBytes32()},
		},
	}
}

func TestCommitPluginCodecV1(t *testing.T) {
	testCases := []struct {
		name   string
		report func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport
		expErr bool
	}{
		{
			name: "base report",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				return report
			},
		},
		{
			name: "empty token address",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.PriceUpdates.TokenPriceUpdates[0].TokenID = ""
				return report
			},
			expErr: true,
		},
		{
			name: "empty merkle root",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.BlessedMerkleRoots[0].MerkleRoot = cciptypes.Bytes32{}
				report.UnblessedMerkleRoots[0].MerkleRoot = cciptypes.Bytes32{}
				return report
			},
		},
		{
			name: "zero token price",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.PriceUpdates.TokenPriceUpdates[0].Price = cciptypes.NewBigInt(big.NewInt(0))
				return report
			},
		},
		{
			name: "zero gas price",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.PriceUpdates.GasPriceUpdates[0].GasPrice = cciptypes.NewBigInt(big.NewInt(0))
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

// Go equivalent of test_deserialize_commit_report
// https://github.com/smartcontractkit/chainlink-aptos/blob/cb70d13f90d16ea7fea7f0f52f02fbebc38d16a9/contracts/ccip/ccip_offramp/tests/offramp_test.move#L525
func TestCommitPluginCodecV1_Decode(t *testing.T) {
	expectedSourceToken := "0x000000000000000000000000000000000000000000000000000000000000000a"
	expectedUsdPerToken, ok := new(big.Int).SetString("500000000000000000000", 10)
	require.True(t, ok)
	expectedSourceChainSelector := cciptypes.ChainSelector(909606746561742123)
	expectedOnRampAddress, err := hexutil.Decode("0x47a1f0a819457f01153f35c6b6b0d42e2e16e91e")
	require.NoError(t, err)
	expectedMinSeqNr := cciptypes.SeqNum(1)
	expectedMaxSeqNr := cciptypes.SeqNum(1)
	expectedMerkleRootBytes, err := hexutil.Decode("0x258dc7f9ec033388ee50bf3e0debfc841a278054f5b2ce41728f7459267c719e")
	require.NoError(t, err)
	var expectedMerkleRoot cciptypes.Bytes32
	copy(expectedMerkleRoot[:], expectedMerkleRootBytes)

	commitReportBytes, err := hexutil.Decode("0x01000000000000000000000000000000000000000000000000000000000000000a000050efe2d6e41a1b00000000000000000000000000000000000000000000000000012b851c4684929f0c1447a1f0a819457f01153f35c6b6b0d42e2e16e91e01000000000000000100000000000000258dc7f9ec033388ee50bf3e0debfc841a278054f5b2ce41728f7459267c719e00")
	require.NoError(t, err)

	codec := NewCommitPluginCodecV1()
	ctx := t.Context()

	commitReport, err := codec.Decode(ctx, commitReportBytes)
	require.NoError(t, err)

	priceUpdates := commitReport.PriceUpdates
	require.Len(t, priceUpdates.TokenPriceUpdates, 1, "Expected one token price update")
	tokenPriceUpdate := priceUpdates.TokenPriceUpdates[0]
	assert.Equal(t, expectedSourceToken, string(tokenPriceUpdate.TokenID), "Source token mismatch")
	assert.Equal(t, expectedUsdPerToken, tokenPriceUpdate.Price.Int, "USD per token mismatch")
	assert.Empty(t, priceUpdates.GasPriceUpdates, "Expected no gas price updates")

	assert.Empty(t, commitReport.BlessedMerkleRoots, "Expected no blessed merkle roots")

	require.Len(t, commitReport.UnblessedMerkleRoots, 1, "Expected one unblessed merkle root")
	merkleRootStruct := commitReport.UnblessedMerkleRoots[0]
	assert.Equal(t, expectedSourceChainSelector, merkleRootStruct.ChainSel, "Source chain selector mismatch")
	assert.Equal(t, expectedOnRampAddress, []byte(merkleRootStruct.OnRampAddress), "On ramp address mismatch")
	assert.Equal(t, expectedMinSeqNr, merkleRootStruct.SeqNumsRange.Start(), "Min sequence number mismatch")
	assert.Equal(t, expectedMaxSeqNr, merkleRootStruct.SeqNumsRange.End(), "Max sequence number mismatch")
	assert.Equal(t, expectedMerkleRoot, merkleRootStruct.MerkleRoot, "Merkle root mismatch")

	assert.Empty(t, commitReport.RMNSignatures, "Expected no RMN signatures")
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
