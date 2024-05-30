package ccip

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
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

func TestCommitObservationJsonSerializationDeserialization(t *testing.T) {
	jsonEncoded := `{
		"interval": {
			"Min":1,
			"Max":12
		},
		"tokensPerFeeCoin": {
			"0x0000000000000000000000000000000000000001": 1,
			"0x507877C2E26f1387432D067D2DaAfa7d0420d90a": 2,
			"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa": 3
		},
		"sourceGasPrice": 3
	}`

	expectedObservation := CommitObservation{
		Interval: cciptypes.CommitStoreInterval{
			Min: 1,
			Max: 12,
		},
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			cciptypes.Address("0x0000000000000000000000000000000000000001"): big.NewInt(1),
			cciptypes.Address("0x507877C2E26f1387432D067D2DaAfa7d0420d90a"): big.NewInt(2), // json eip55->eip55 parsed
			cciptypes.Address("0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa"): big.NewInt(3), // json lower->eip55 parsed
		},
		SourceGasPriceUSD: big.NewInt(3),
	}

	observations := GetParsableObservations[CommitObservation](logger.TestLogger(t), []types.AttributedObservation{
		{Observation: []byte(jsonEncoded)},
	})
	assert.Equal(t, 1, len(observations))
	assert.Equal(t, expectedObservation, observations[0])

	backToJson, err := expectedObservation.Marshal()
	// we expect the json encoded addresses to be lower-case
	exp := strings.ReplaceAll(
		jsonEncoded, "0x507877C2E26f1387432D067D2DaAfa7d0420d90a", strings.ToLower("0x507877C2E26f1387432D067D2DaAfa7d0420d90a"))
	assert.NoError(t, err)
	assert.JSONEq(t, exp, string(backToJson))

	// and we expect to get the same results after we parse the lower-case addresses
	observations = GetParsableObservations[CommitObservation](logger.TestLogger(t), []types.AttributedObservation{
		{Observation: []byte(jsonEncoded)},
	})
	assert.Equal(t, 1, len(observations))
	assert.Equal(t, expectedObservation, observations[0])
}

func TestAddressEncodingBackwardsCompatibility(t *testing.T) {
	// The intention of this test is to remind including proper formatting of addresses after config is updated.
	//
	// The following tests will fail when a new cciptypes.Address field is added or removed.
	// If you notice that the test is failing, make sure to apply proper address formatting
	// after the struct is marshalled/unmarshalled and then include your new field in the expected fields slice to
	// make this test pass or if you removed a field, remove it from the expected fields slice.

	t.Run("job spec config", func(t *testing.T) {
		exp := []string{"ccip.Address OffRamp"}

		fields := testhelpers.FindStructFieldsOfCertainType(
			"ccip.Address",
			config.CommitPluginJobSpecConfig{PriceGetterConfig: &config.DynamicPriceGetterConfig{}},
		)
		assert.Equal(t, exp, fields)
	})

	t.Run("commit observation", func(t *testing.T) {
		exp := []string{"map[ccip.Address]*big.Int TokenPricesUSD"}

		fields := testhelpers.FindStructFieldsOfCertainType(
			"ccip.Address",
			CommitObservation{SourceGasPriceUSD: big.NewInt(0)},
		)
		assert.Equal(t, exp, fields)
	})

}
