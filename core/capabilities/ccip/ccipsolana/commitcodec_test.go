package ccipsolana

import (
	"bytes"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

var randomBlessedCommitReport = func() cciptypes.CommitPluginReport {
	pubkey, err := solanago.NewRandomPrivateKey()
	if err != nil {
		panic(err)
	}

	return cciptypes.CommitPluginReport{
		BlessedMerkleRoots: []cciptypes.MerkleRootChain{
			{
				OnRampAddress: cciptypes.UnknownAddress(pubkey.PublicKey().String()),
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
					TokenID: "C8WSPj3yyus1YN3yNB6YA5zStYtbjQWtpmKadmvyUXq8",
					Price:   cciptypes.NewBigInt(big.NewInt(rand.Int63())),
				},
			},
			GasPriceUpdates: []cciptypes.GasPriceChain{
				{GasPrice: cciptypes.NewBigInt(big.NewInt(rand.Int63())), ChainSel: cciptypes.ChainSelector(rand.Uint64())},
				{GasPrice: cciptypes.NewBigInt(big.NewInt(rand.Int63())), ChainSel: cciptypes.ChainSelector(rand.Uint64())},
				{GasPrice: cciptypes.NewBigInt(big.NewInt(rand.Int63())), ChainSel: cciptypes.ChainSelector(rand.Uint64())},
			},
		},
		RMNSignatures: []cciptypes.RMNECDSASignature{
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
			name: "base report blessed",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				return report
			},
		},
		{
			name: "base report unblessed",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.RMNSignatures = nil
				report.UnblessedMerkleRoots = report.BlessedMerkleRoots
				report.BlessedMerkleRoots = nil
				return report
			},
		},
		{
			name: "blessed report with no rmn signatures",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.RMNSignatures = nil
				return report
			},
			expErr: true,
		},
		{
			name: "rmn signature included without any blessed root",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.UnblessedMerkleRoots = report.BlessedMerkleRoots
				report.BlessedMerkleRoots = nil
				return report
			},
			expErr: true,
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
				return report
			},
		},
		{
			name: "both blessed and unblessed merkle roots",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.UnblessedMerkleRoots = []cciptypes.MerkleRootChain{
					report.BlessedMerkleRoots[0]}
				return report
			},
			expErr: true,
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
		{
			name: "empty gas price",
			report: func(report cciptypes.CommitPluginReport) cciptypes.CommitPluginReport {
				report.PriceUpdates.GasPriceUpdates[0].GasPrice = cciptypes.NewBigInt(nil)
				return report
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report := tc.report(randomBlessedCommitReport())
			commitCodec := NewCommitPluginCodecV1()
			ctx := testutils.Context(t)
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
	ctx := testutils.Context(b)

	rep := randomBlessedCommitReport()
	for i := 0; i < b.N; i++ {
		_, err := commitCodec.Encode(ctx, rep)
		require.NoError(b, err)
	}
}

func BenchmarkCommitPluginCodecV1_Decode(b *testing.B) {
	commitCodec := NewCommitPluginCodecV1()
	ctx := testutils.Context(b)
	encodedReport, err := commitCodec.Encode(ctx, randomBlessedCommitReport())
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_, err := commitCodec.Decode(ctx, encodedReport)
		require.NoError(b, err)
	}
}

func BenchmarkCommitPluginCodecV1_Encode_Decode(b *testing.B) {
	commitCodec := NewCommitPluginCodecV1()
	ctx := testutils.Context(b)

	rep := randomBlessedCommitReport()
	for i := 0; i < b.N; i++ {
		encodedReport, err := commitCodec.Encode(ctx, rep)
		require.NoError(b, err)
		decodedReport, err := commitCodec.Decode(ctx, encodedReport)
		require.NoError(b, err)
		require.Equal(b, rep, decodedReport)
	}
}

func Test_DecodingCommitReport(t *testing.T) {
	t.Run("decode on-chain commit report", func(t *testing.T) {
		chainSel := cciptypes.ChainSelector(rand.Uint64())
		minSeqNr := rand.Uint64()
		maxSeqNr := minSeqNr + 10
		onRampAddr, err := solanago.NewRandomPrivateKey()
		require.NoError(t, err)

		tokenSource := solanago.MustPublicKeyFromBase58("C8WSPj3yyus1YN3yNB6YA5zStYtbjQWtpmKadmvyUXq8")
		tokenPrice := encodeBigIntToFixedLengthLE(big.NewInt(rand.Int63()), 28)
		gasPrice := encodeBigIntToFixedLengthLE(big.NewInt(rand.Int63()), 28)
		merkleRoot := utils.RandomBytes32()

		tpu := []ccip_offramp.TokenPriceUpdate{
			{
				SourceToken: tokenSource,
				UsdPerToken: [28]uint8(tokenPrice),
			},
		}

		gpu := []ccip_offramp.GasPriceUpdate{
			{UsdPerUnitGas: [28]uint8(gasPrice), DestChainSelector: uint64(chainSel)},
			{UsdPerUnitGas: [28]uint8(gasPrice), DestChainSelector: uint64(chainSel)},
			{UsdPerUnitGas: [28]uint8(gasPrice), DestChainSelector: uint64(chainSel)},
		}

		onChainReport := ccip_offramp.CommitInput{
			MerkleRoot: ccip_offramp.MerkleRoot{
				SourceChainSelector: uint64(chainSel),
				OnRampAddress:       onRampAddr.PublicKey().Bytes(),
				MinSeqNr:            minSeqNr,
				MaxSeqNr:            maxSeqNr,
				MerkleRoot:          merkleRoot,
			},
			PriceUpdates: ccip_offramp.PriceUpdates{
				TokenPriceUpdates: tpu,
				GasPriceUpdates:   gpu,
			},
		}

		var buf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&buf)
		err = onChainReport.MarshalWithEncoder(encoder)
		require.NoError(t, err)

		commitCodec := NewCommitPluginCodecV1()
		decode, err := commitCodec.Decode(testutils.Context(t), buf.Bytes())
		require.NoError(t, err)
		mr := decode.UnblessedMerkleRoots[0]

		// check decoded ocr report merkle root matches with on-chain report
		require.Equal(t, strconv.FormatUint(minSeqNr, 10), mr.SeqNumsRange.Start().String())
		require.Equal(t, strconv.FormatUint(maxSeqNr, 10), mr.SeqNumsRange.End().String())
		require.Equal(t, cciptypes.UnknownAddress(onRampAddr.PublicKey().Bytes()), mr.OnRampAddress)
		require.Equal(t, cciptypes.Bytes32(merkleRoot), mr.MerkleRoot)

		// check decoded ocr report token price update matches with on-chain report
		pu := decode.PriceUpdates.TokenPriceUpdates[0]
		require.Equal(t, decodeLEToBigInt(tokenPrice), pu.Price)
		require.Equal(t, cciptypes.UnknownEncodedAddress(tokenSource.String()), pu.TokenID)

		// check decoded ocr report gas price update matches with on-chain report
		gu := decode.PriceUpdates.GasPriceUpdates[0]
		require.Equal(t, decodeLEToBigInt(gasPrice), gu.GasPrice)
		require.Equal(t, chainSel, gu.ChainSel)
	})

	t.Run("decode Borsh encoded commit report", func(t *testing.T) {
		rep := randomBlessedCommitReport()
		commitCodec := NewCommitPluginCodecV1()
		decode, err := commitCodec.Encode(testutils.Context(t), rep)
		require.NoError(t, err)

		decoder := agbinary.NewBorshDecoder(decode)
		decodedReport := ccip_offramp.CommitInput{}
		err = decodedReport.UnmarshalWithDecoder(decoder)
		require.NoError(t, err)

		reportMerkleRoot := rep.BlessedMerkleRoots[0]
		require.Equal(t, reportMerkleRoot.MerkleRoot, cciptypes.Bytes32(decodedReport.MerkleRoot.MerkleRoot))

		tu := rep.PriceUpdates.TokenPriceUpdates[0]
		require.Equal(t, tu.TokenID, cciptypes.UnknownEncodedAddress(decodedReport.PriceUpdates.TokenPriceUpdates[0].SourceToken.String()))
		require.Equal(t, tu.Price, decodeLEToBigInt(decodedReport.PriceUpdates.TokenPriceUpdates[0].UsdPerToken[:]))

		gu := rep.PriceUpdates.GasPriceUpdates[0]
		require.Equal(t, gu.ChainSel, cciptypes.ChainSelector(decodedReport.PriceUpdates.GasPriceUpdates[0].DestChainSelector))
		require.Equal(t, gu.GasPrice, decodeLEToBigInt(decodedReport.PriceUpdates.GasPriceUpdates[0].UsdPerUnitGas[:]))
	})
}
