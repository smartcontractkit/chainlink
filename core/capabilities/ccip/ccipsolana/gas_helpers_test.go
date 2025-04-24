package ccipsolana

import (
	"encoding/binary"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

func Test_calculateMessageMaxGas(t *testing.T) {
	type args struct {
		dataLen          int
		numTokens        int
		extraArgs        []byte
		tokenGasOverhead uint32
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "allowOOO set to true makes no difference to final gas estimate",
			args: args{
				dataLen:          5,
				numTokens:        2,
				extraArgs:        makeExtraArgsV2(200_000, true),
				tokenGasOverhead: 100,
			},
			want: 550200,
		},
		{
			name: "allowOOO set to false makes no difference to final gas estimate",
			args: args{
				dataLen:          5,
				numTokens:        2,
				extraArgs:        makeExtraArgsV2(200_000, false),
				tokenGasOverhead: 100,
			},
			want: 550200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ccipocr3.Message{
				Data:         make([]byte, tt.args.dataLen),
				TokenAmounts: getTokenAmounts(t, tt.args.numTokens, tt.args.tokenGasOverhead),
				ExtraArgs:    tt.args.extraArgs,
			}
			// Set the source chain selector to be EVM for now
			msg.Header.SourceChainSelector = ccipocr3.ChainSelector(chainsel.SOLANA_TESTNET.Selector)
			ep := EstimateProvider{extraDataCodec: ExtraDataCodec}
			got := ep.CalculateMessageMaxGas(msg)
			t.Log(got)
			assert.Equalf(t, tt.want, got, "calculateMessageMaxGas(%v, %v)", tt.args.dataLen, tt.args.numTokens)
		})
	}
}

// TestCalculateMaxGas is taken from the ccip repo where the CalculateMerkleTreeGas and CalculateMessageMaxGas values
// are combined to one function.
func TestCalculateMaxGas(t *testing.T) {
	tests := []struct {
		name             string
		numRequests      int
		dataLength       int
		numberOfTokens   int
		extraArgs        []byte
		tokenGasOverhead uint32
		want             uint64
	}{
		{
			name:             "v2 extra args",
			numRequests:      3,
			dataLength:       len([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
			numberOfTokens:   1,
			extraArgs:        makeExtraArgsV2(200_000, true),
			tokenGasOverhead: 10,
			want:             550011,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ccipocr3.Message{
				Data:         make([]byte, tt.dataLength),
				TokenAmounts: getTokenAmounts(t, tt.numberOfTokens, tt.tokenGasOverhead),
				ExtraArgs:    tt.extraArgs,
			}

			msg.Header.SourceChainSelector = ccipocr3.ChainSelector(chainsel.SOLANA_TESTNET.Selector)
			ep := EstimateProvider{extraDataCodec: ExtraDataCodec}
			gotTree := ep.CalculateMerkleTreeGas(tt.numRequests)
			gotMsg := ep.CalculateMessageMaxGas(msg)
			t.Log("want", tt.want, "got", gotTree+gotMsg)
			assert.Equal(t, tt.want, gotTree+gotMsg)
		})
	}
}

func makeExtraArgsV2(computeUnits uint32, allowOOO bool) []byte {
	extraArgs, err := SerializeExtraArgs(svmExtraArgsV1Tag, fee_quoter.SVMExtraArgsV1{
		AllowOutOfOrderExecution: allowOOO,
		ComputeUnits:             computeUnits,
		AccountIsWritableBitmap:  0,
		TokenReceiver:            [32]uint8{},
		Accounts:                 [][32]uint8{},
	})
	if err != nil {
		panic(err)
	}
	return extraArgs
}

func getTokenAmounts(t *testing.T, numTokens int, tokenGasOverhead uint32) []ccipocr3.RampTokenAmount {
	tokenDestGasOverhead := binary.BigEndian.AppendUint32(nil, tokenGasOverhead)

	tokenAmounts := make([]ccipocr3.RampTokenAmount, numTokens)
	for i := range numTokens {
		tokenAmounts[i] = ccipocr3.RampTokenAmount{
			DestExecData: tokenDestGasOverhead,
		}
	}
	return tokenAmounts
}
