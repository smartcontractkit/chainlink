package ccipton

import (
	"context"
	"encoding/base64"
	"math/big"
	"math/rand"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
)

func randomTONExecuteReport(t *testing.T, sourceChainSelector uint64) cciptypes.ExecutePluginReport {
	const numChainReports = 2
	const msgsPerReport = 2
	const numTokensPerMsg = 2

	chainReports := make([]cciptypes.ExecutePluginReportSingleChain, numChainReports)
	for i := 0; i < numChainReports; i++ {
		reportMessages := make([]cciptypes.Message, msgsPerReport)
		for j := 0; j < msgsPerReport; j++ {
			addr, err := address.ParseAddr("EQDtFpEwcFAEcRe5mLVh2N6C0x-_hJEM7W61_JLnSF74p4q2")
			require.NoError(t, err)
			require.NoError(t, err)
			extraData := []byte{0x12, 0x34}

			addrBytes, err := base64.RawURLEncoding.DecodeString(addr.String())
			require.NoError(t, err)
			tokenAmounts := make([]cciptypes.RampTokenAmount, numTokensPerMsg)
			for z := 0; z < numTokensPerMsg; z++ {
				tokenAmounts[z] = cciptypes.RampTokenAmount{
					SourcePoolAddress: cciptypes.UnknownAddress(addr.String()),
					DestTokenAddress:  addrBytes, // pad to 36 bytes
					ExtraData:         extraData,
					Amount:            cciptypes.NewBigInt(big.NewInt(rand.Int63())),
					DestExecData:      []byte{0, 0, 0, 0},
				}
			}
			reportMessages[j] = cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           [32]byte{},
					SourceChainSelector: cciptypes.ChainSelector(sourceChainSelector),
					DestChainSelector:   cciptypes.ChainSelector(rand.Uint64()),
					SequenceNumber:      cciptypes.SeqNum(rand.Uint64()),
					Nonce:               rand.Uint64(),
				},
				Sender:       cciptypes.UnknownAddress(addr.String()),
				Data:         extraData,
				Receiver:     addrBytes,
				ExtraArgs:    []byte{0, 0, 0, 0},
				TokenAmounts: tokenAmounts,
			}
		}
		chainReports[i] = cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(sourceChainSelector),
			Messages:            reportMessages,
			OffchainTokenData:   [][][]byte{{{0x1}, {0x2, 0x3}}},
			Proofs:              []cciptypes.Bytes32{},
			ProofFlagBits:       cciptypes.BigInt{Int: big.NewInt(1)},
		}
	}
	return cciptypes.ExecutePluginReport{ChainReports: chainReports}
}

func TestExecutePluginCodecV1_TON(t *testing.T) {
	ctx := context.Background()
	mockExtraDataCodec := new(mocks.SourceChainExtraDataCodec)
	edc := common.ExtraDataCodec(map[string]common.SourceChainExtraDataCodec{
		chainsel.FamilyEVM:    mockExtraDataCodec,
		chainsel.FamilySolana: mockExtraDataCodec,
		chainsel.FamilyTon:    mockExtraDataCodec,
	})

	mockExtraDataCodec.On("DecodeDestExecDataToMap", mock.Anything).Return(map[string]any{
		"destgasamount": uint32(1000),
	}, nil)
	mockExtraDataCodec.On("DecodeExtraArgsToMap", mock.Anything).Return(map[string]any{
		"gasLimit": big.NewInt(1000),
	}, nil)
	codec := NewExecutePluginCodecV1(edc)

	t.Run("encode/decode roundtrip", func(t *testing.T) {
		report := randomTONExecuteReport(t, 5009297550715157269) // evm selector for TON
		encoded, err := codec.Encode(ctx, report)
		require.NoError(t, err)
		decoded, err := codec.Decode(ctx, encoded)
		require.NoError(t, err)
		assert.Equal(t, report.ChainReports[0].SourceChainSelector, decoded.ChainReports[0].SourceChainSelector)
		assert.Equal(t, report.ChainReports[0].Messages[0].TokenAmounts[0].Amount, decoded.ChainReports[0].Messages[0].TokenAmounts[0].Amount)
	})

	t.Run("empty report", func(t *testing.T) {
		encoded, err := codec.Encode(ctx, cciptypes.ExecutePluginReport{})
		require.NoError(t, err)
		assert.Nil(t, encoded)
	})
}
