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

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_offramp"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

var randomBlessedCommitReport = func() ccipocr3common.CommitPluginReport {
	pubkey, err := solanago.NewRandomPrivateKey()
	if err != nil {
		panic(err)
	}

	return ccipocr3common.CommitPluginReport{
		BlessedMerkleRoots: []ccipocr3common.MerkleRootChain{
			{
				OnRampAddress: ccipocr3common.UnknownAddress(pubkey.PublicKey().String()),
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
					TokenID: "C8WSPj3yyus1YN3yNB6YA5zStYtbjQWtpmKadmvyUXq8",
					Price:   ccipocr3common.NewBigInt(big.NewInt(rand.Int63())),
				},
			},
			GasPriceUpdates: []ccipocr3common.GasPriceChain{
				{GasPrice: ccipocr3common.NewBigInt(big.NewInt(rand.Int63())), ChainSel: ccipocr3common.ChainSelector(rand.Uint64())},
				{GasPrice: ccipocr3common.NewBigInt(big.NewInt(rand.Int63())), ChainSel: ccipocr3common.ChainSelector(rand.Uint64())},
				{GasPrice: ccipocr3common.NewBigInt(big.NewInt(rand.Int63())), ChainSel: ccipocr3common.ChainSelector(rand.Uint64())},
			},
		},
		RMNSignatures: []ccipocr3common.RMNECDSASignature{
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
			name: "base report blessed",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				return report
			},
		},
		{
			name: "base report unblessed",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.RMNSignatures = nil
				report.UnblessedMerkleRoots = report.BlessedMerkleRoots
				report.BlessedMerkleRoots = nil
				return report
			},
		},
		{
			name: "blessed report with no rmn signatures",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.RMNSignatures = nil
				return report
			},
			expErr: true,
		},
		{
			name: "rmn signature included without any blessed root",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.UnblessedMerkleRoots = report.BlessedMerkleRoots
				report.BlessedMerkleRoots = nil
				return report
			},
			expErr: true,
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
				return report
			},
		},
		{
			name: "both blessed and unblessed merkle roots",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.UnblessedMerkleRoots = []ccipocr3common.MerkleRootChain{
					report.BlessedMerkleRoots[0],
				}
				return report
			},
			expErr: true,
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
		{
			name: "empty gas price",
			report: func(report ccipocr3common.CommitPluginReport) ccipocr3common.CommitPluginReport {
				report.PriceUpdates.GasPriceUpdates[0].GasPrice = ccipocr3common.NewBigInt(nil)
				return report
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report := tc.report(randomBlessedCommitReport())
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

	rep := randomBlessedCommitReport()
	for b.Loop() {
		_, err := commitCodec.Encode(ctx, rep)
		require.NoError(b, err)
	}
}

func BenchmarkCommitPluginCodecV1_Decode(b *testing.B) {
	commitCodec := NewCommitPluginCodecV1()
	ctx := b.Context()
	encodedReport, err := commitCodec.Encode(ctx, randomBlessedCommitReport())
	require.NoError(b, err)

	for b.Loop() {
		_, err := commitCodec.Decode(ctx, encodedReport)
		require.NoError(b, err)
	}
}

func BenchmarkCommitPluginCodecV1_Encode_Decode(b *testing.B) {
	commitCodec := NewCommitPluginCodecV1()
	ctx := b.Context()

	rep := randomBlessedCommitReport()
	for b.Loop() {
		encodedReport, err := commitCodec.Encode(ctx, rep)
		require.NoError(b, err)
		decodedReport, err := commitCodec.Decode(ctx, encodedReport)
		require.NoError(b, err)
		require.Equal(b, rep, decodedReport)
	}
}

func Test_DecodingCommitReport(t *testing.T) {
	t.Run("decode on-chain commit report", func(t *testing.T) {
		chainSel := ccipocr3common.ChainSelector(rand.Uint64())
		minSeqNr := rand.Uint64()
		maxSeqNr := minSeqNr + 10
		onRampAddr, err := solanago.NewRandomPrivateKey()
		require.NoError(t, err)

		tokenSource := solanago.MustPublicKeyFromBase58("C8WSPj3yyus1YN3yNB6YA5zStYtbjQWtpmKadmvyUXq8")
		tokenPrice := encodeBigIntToFixedLengthBE(big.NewInt(rand.Int63()), 28)
		gasPrice := encodeBigIntToFixedLengthBE(big.NewInt(rand.Int63()), 28)
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
			MerkleRoot: &ccip_offramp.MerkleRoot{
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
		decode, err := commitCodec.Decode(t.Context(), buf.Bytes())
		require.NoError(t, err)
		mr := decode.UnblessedMerkleRoots[0]

		// check decoded ocr report merkle root matches with on-chain report
		require.Equal(t, strconv.FormatUint(minSeqNr, 10), mr.SeqNumsRange.Start().String())
		require.Equal(t, strconv.FormatUint(maxSeqNr, 10), mr.SeqNumsRange.End().String())
		require.Equal(t, ccipocr3common.UnknownAddress(onRampAddr.PublicKey().Bytes()), mr.OnRampAddress)
		require.Equal(t, ccipocr3common.Bytes32(merkleRoot), mr.MerkleRoot)

		// check decoded ocr report token price update matches with on-chain report
		pu := decode.PriceUpdates.TokenPriceUpdates[0]
		require.Equal(t, decodeBEToBigInt(tokenPrice), pu.Price)
		require.Equal(t, ccipocr3common.UnknownEncodedAddress(tokenSource.String()), pu.TokenID)

		// check decoded ocr report gas price update matches with on-chain report
		gu := decode.PriceUpdates.GasPriceUpdates[0]
		require.Equal(t, decodeBEToBigInt(gasPrice), gu.GasPrice)
		require.Equal(t, chainSel, gu.ChainSel)
	})

	t.Run("decode on-chain commit report with no MerkleRoot", func(t *testing.T) {
		chainSel := ccipocr3common.ChainSelector(rand.Uint64())

		tokenSource := solanago.MustPublicKeyFromBase58("C8WSPj3yyus1YN3yNB6YA5zStYtbjQWtpmKadmvyUXq8")
		tokenPrice := encodeBigIntToFixedLengthBE(big.NewInt(rand.Int63()), 28)
		gasPrice := encodeBigIntToFixedLengthBE(big.NewInt(rand.Int63()), 28)

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
			MerkleRoot: nil,
			PriceUpdates: ccip_offramp.PriceUpdates{
				TokenPriceUpdates: tpu,
				GasPriceUpdates:   gpu,
			},
		}

		var buf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&buf)
		err := onChainReport.MarshalWithEncoder(encoder)
		require.NoError(t, err)

		commitCodec := NewCommitPluginCodecV1()
		decode, err := commitCodec.Decode(t.Context(), buf.Bytes())
		require.NoError(t, err)
		require.Nilf(t, decode.UnblessedMerkleRoots, "UnblessedMerkleRoots should be nil")
		require.Nilf(t, decode.BlessedMerkleRoots, "BlessedMerkleRoots should be nil")

		// check decoded ocr report token price update matches with on-chain report
		pu := decode.PriceUpdates.TokenPriceUpdates[0]
		require.Equal(t, decodeBEToBigInt(tokenPrice), pu.Price)
		require.Equal(t, ccipocr3common.UnknownEncodedAddress(tokenSource.String()), pu.TokenID)

		// check decoded ocr report gas price update matches with on-chain report
		gu := decode.PriceUpdates.GasPriceUpdates[0]
		require.Equal(t, decodeBEToBigInt(gasPrice), gu.GasPrice)
		require.Equal(t, chainSel, gu.ChainSel)
	})

	t.Run("decode Borsh encoded commit report", func(t *testing.T) {
		rep := randomBlessedCommitReport()
		commitCodec := NewCommitPluginCodecV1()
		decode, err := commitCodec.Encode(t.Context(), rep)
		require.NoError(t, err)

		decoder := agbinary.NewBorshDecoder(decode)
		decodedReport := ccip_offramp.CommitInput{}
		err = decodedReport.UnmarshalWithDecoder(decoder)
		require.NoError(t, err)

		reportMerkleRoot := rep.BlessedMerkleRoots[0]
		require.Equal(t, reportMerkleRoot.MerkleRoot, ccipocr3common.Bytes32(decodedReport.MerkleRoot.MerkleRoot))

		tu := rep.PriceUpdates.TokenPriceUpdates[0]
		require.Equal(t, tu.TokenID, ccipocr3common.UnknownEncodedAddress(decodedReport.PriceUpdates.TokenPriceUpdates[0].SourceToken.String()))
		require.Equal(t, tu.Price, decodeBEToBigInt(decodedReport.PriceUpdates.TokenPriceUpdates[0].UsdPerToken[:]))

		gu := rep.PriceUpdates.GasPriceUpdates[0]
		require.Equal(t, gu.ChainSel, ccipocr3common.ChainSelector(decodedReport.PriceUpdates.GasPriceUpdates[0].DestChainSelector))
		require.Equal(t, gu.GasPrice, decodeBEToBigInt(decodedReport.PriceUpdates.GasPriceUpdates[0].UsdPerUnitGas[:]))
	})
}
