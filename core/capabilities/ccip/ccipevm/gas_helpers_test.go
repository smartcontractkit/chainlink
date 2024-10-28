package ccipevm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

func Test_calculateMessageMaxGas(t *testing.T) {
	type args struct {
		dataLen   int
		numTokens int
		extraArgs []byte
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "base",
			args: args{dataLen: 5, numTokens: 2, extraArgs: makeExtraArgsV1(200_000)},
			want: 1_022_264,
		},
		{
			name: "large",
			args: args{dataLen: 1000, numTokens: 1000, extraArgs: makeExtraArgsV1(200_000)},
			want: 346_677_520,
		},
		{
			name: "overheadGas test 1",
			args: args{dataLen: 0, numTokens: 0, extraArgs: makeExtraArgsV1(200_000)},
			want: 319_920,
		},
		{
			name: "overheadGas test 2",
			args: args{dataLen: len([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}), numTokens: 1, extraArgs: makeExtraArgsV1(200_000)},
			want: 675_948,
		},
		{
			name: "allowOOO set to true makes no difference to final gas estimate",
			args: args{dataLen: 5, numTokens: 2, extraArgs: makeExtraArgsV2(200_000, true)},
			want: 1_022_264,
		},
		{
			name: "allowOOO set to false makes no difference to final gas estimate",
			args: args{dataLen: 5, numTokens: 2, extraArgs: makeExtraArgsV2(200_000, false)},
			want: 1_022_264,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ccipocr3.Message{
				Data:         make([]byte, tt.args.dataLen),
				TokenAmounts: make([]ccipocr3.RampTokenAmount, tt.args.numTokens),
				ExtraArgs:    tt.args.extraArgs,
			}
			ep := EstimateProvider{}
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
		name           string
		numRequests    int
		dataLength     int
		numberOfTokens int
		extraArgs      []byte
		want           uint64
	}{
		{
			name:           "maxGasOverheadGas 1",
			numRequests:    6,
			dataLength:     0,
			numberOfTokens: 0,
			extraArgs:      makeExtraArgsV1(200_000),
			want:           322_992,
		},
		{
			name:           "maxGasOverheadGas 2",
			numRequests:    3,
			dataLength:     len([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
			numberOfTokens: 1,
			extraArgs:      makeExtraArgsV1(200_000),
			want:           678_508,
		},
		{
			name:           "v2 extra args",
			numRequests:    3,
			dataLength:     len([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
			numberOfTokens: 1,
			extraArgs:      makeExtraArgsV2(200_000, true),
			want:           678_508,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ccipocr3.Message{
				Data:         make([]byte, tt.dataLength),
				TokenAmounts: make([]ccipocr3.RampTokenAmount, tt.numberOfTokens),
				ExtraArgs:    tt.extraArgs,
			}
			ep := EstimateProvider{}

			gotTree := ep.CalculateMerkleTreeGas(tt.numRequests)
			gotMsg := ep.CalculateMessageMaxGas(msg)
			t.Log("want", tt.want, "got", gotTree+gotMsg)
			assert.Equal(t, tt.want, gotTree+gotMsg)
		})
	}
}

func makeExtraArgsV1(gasLimit uint64) []byte {
	// extra args is the tag followed by the gas limit abi-encoded.
	var extraArgs []byte
	extraArgs = append(extraArgs, evmExtraArgsV1Tag...)
	gasLimitBytes := new(big.Int).SetUint64(gasLimit).Bytes()
	// pad from the left to 32 bytes
	gasLimitBytes = common.LeftPadBytes(gasLimitBytes, 32)
	extraArgs = append(extraArgs, gasLimitBytes...)
	return extraArgs
}

func makeExtraArgsV2(gasLimit uint64, allowOOO bool) []byte {
	// extra args is the tag followed by the gas limit and allowOOO abi-encoded.
	var extraArgs []byte
	extraArgs = append(extraArgs, evmExtraArgsV2Tag...)
	gasLimitBytes := new(big.Int).SetUint64(gasLimit).Bytes()
	// pad from the left to 32 bytes
	gasLimitBytes = common.LeftPadBytes(gasLimitBytes, 32)

	// abi-encode allowOOO
	var allowOOOBytes []byte
	if allowOOO {
		allowOOOBytes = append(allowOOOBytes, 1)
	} else {
		allowOOOBytes = append(allowOOOBytes, 0)
	}
	// pad from the left to 32 bytes
	allowOOOBytes = common.LeftPadBytes(allowOOOBytes, 32)

	extraArgs = append(extraArgs, gasLimitBytes...)
	extraArgs = append(extraArgs, allowOOOBytes...)
	return extraArgs
}
