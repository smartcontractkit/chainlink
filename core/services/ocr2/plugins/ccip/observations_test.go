package ccip

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

func TestObservationFilter(t *testing.T) {
	lggr := logger.TestLogger(t)
	obs1 := CommitObservation{Interval: cciptypes.CommitStoreInterval{Min: 1, Max: 10}}
	b1, err := obs1.Marshal()
	require.NoError(t, err)
	nonEmpty := GetParsableObservations[CommitObservation](lggr, []types.AttributedObservation{{Observation: b1}, {Observation: []byte{}}})
	require.Equal(t, 1, len(nonEmpty))
	assert.Equal(t, nonEmpty[0].Interval, obs1.Interval)
}

// After 1.2, the observation struct is version agnostic
// so only need to verify the 1.0->1.2 transition.
type CommitObservationV1_0_0 struct {
	Interval          commit_store_1_0_0.CommitStoreInterval `json:"interval"`
	TokenPricesUSD    map[common.Address]*big.Int            `json:"tokensPerFeeCoin"`
	SourceGasPriceUSD *big.Int                               `json:"sourceGasPrice"`
}

func TestObservationCompat100_120(t *testing.T) {
	v10 := CommitObservationV1_0_0{
		Interval: commit_store_1_0_0.CommitStoreInterval{
			Min: 1,
			Max: 12,
		},
		TokenPricesUSD:    map[common.Address]*big.Int{common.HexToAddress("0x1"): big.NewInt(1)},
		SourceGasPriceUSD: big.NewInt(3)}
	b10, err := json.Marshal(v10)
	require.NoError(t, err)
	v12 := CommitObservation{
		Interval: cciptypes.CommitStoreInterval{
			Min: 1,
			Max: 12,
		},
		TokenPricesUSD:    map[cciptypes.Address]*big.Int{ccipcalc.HexToAddress("0x1"): big.NewInt(1)},
		SourceGasPriceUSD: big.NewInt(3),
	}
	b12, err := json.Marshal(v12)
	require.NoError(t, err)
	// Assert identical json.
	assert.Equal(t, b10, b12)
}

func TestCommitObservationJsonDeserialization(t *testing.T) {
	expectedObservation := CommitObservation{
		Interval: cciptypes.CommitStoreInterval{
			Min: 1,
			Max: 12,
		},
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			ccipcalc.HexToAddress("0x1"): big.NewInt(1)},
		SourceGasPriceUSD: big.NewInt(3),
	}

	json := `{
		"interval": {
			"Min":1,
			"Max":12
		},
		"tokensPerFeeCoin": {
			"0x0000000000000000000000000000000000000001": 1
		},
		"sourceGasPrice": 3
	}`

	observations := GetParsableObservations[CommitObservation](logger.TestLogger(t), []types.AttributedObservation{{Observation: []byte(json)}})
	assert.Equal(t, 1, len(observations))
	assert.Equal(t, expectedObservation, observations[0])
}

func TestExecutionObservationJsonDeserialization(t *testing.T) {
	expectedObservation := ExecutionObservation{Messages: map[uint64]MsgData{
		2: {TokenData: tokenData("c")},
		1: {TokenData: tokenData("c")},
	}}

	// ["YQ=="] is "a"
	// ["Yw=="] is "c"
	json := `{
		"messages": {
			"2":{"tokenData":["YQ=="]},
			"1":{"tokenData":["Yw=="]},
			"2":{"tokenData":["Yw=="]}
		}
	}`

	observations := GetParsableObservations[ExecutionObservation](logger.TestLogger(t), []types.AttributedObservation{{Observation: []byte(json)}})
	assert.Equal(t, 1, len(observations))
	assert.Equal(t, 2, len(observations[0].Messages))
	assert.Equal(t, expectedObservation, observations[0])
}

func TestObservationSize(t *testing.T) {
	testParams := gopter.DefaultTestParameters()
	testParams.MinSuccessfulTests = 100
	p := gopter.NewProperties(testParams)
	p.Property("bounded observation size", prop.ForAll(func(min, max uint64) bool {
		o := NewExecutionObservation(
			[]ObservedMessage{
				{
					SeqNr:   min,
					MsgData: MsgData{},
				},
				{
					SeqNr:   max,
					MsgData: MsgData{},
				},
			},
		)
		b, err := o.Marshal()
		require.NoError(t, err)
		return len(b) <= MaxObservationLength
	}, gen.UInt64(), gen.UInt64()))
	p.TestingRun(t)
}

func TestNewExecutionObservation(t *testing.T) {
	tests := []struct {
		name         string
		observations []ObservedMessage
		want         ExecutionObservation
	}{
		{
			name:         "nil observations",
			observations: nil,
			want:         ExecutionObservation{Messages: map[uint64]MsgData{}},
		},
		{
			name:         "empty observations",
			observations: []ObservedMessage{},
			want:         ExecutionObservation{Messages: map[uint64]MsgData{}},
		},
		{
			name: "observations with different sequence numbers",
			observations: []ObservedMessage{
				NewObservedMessage(1, tokenData("a")),
				NewObservedMessage(2, tokenData("b")),
				NewObservedMessage(3, tokenData("c")),
			},
			want: ExecutionObservation{
				Messages: map[uint64]MsgData{
					1: {TokenData: tokenData("a")},
					2: {TokenData: tokenData("b")},
					3: {TokenData: tokenData("c")},
				},
			},
		},
		{
			name: "last one wins in case of duplicates",
			observations: []ObservedMessage{
				NewObservedMessage(1, tokenData("a")),
				NewObservedMessage(1, tokenData("b")),
				NewObservedMessage(1, tokenData("c")),
			},
			want: ExecutionObservation{
				Messages: map[uint64]MsgData{
					1: {TokenData: tokenData("c")},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewExecutionObservation(tt.observations), "NewExecutionObservation(%v)", tt.observations)
		})
	}
}

func tokenData(value string) [][]byte {
	return [][]byte{[]byte(value)}
}
