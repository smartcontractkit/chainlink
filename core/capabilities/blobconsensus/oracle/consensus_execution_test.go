package oracle

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/cre-sdk-go/cre"

	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

func Test_CalculateOutcomeForObservations(t *testing.T) {
	type testCase struct {
		name            string
		observations    []*valuespb.Value
		descriptor      *sdk.ConsensusDescriptor
		defaultValue    *valuespb.Value
		f               int
		expectedOutcome *valuespb.Value
		expectedError   error
	}

	testCases := []testCase{
		{
			name: "insufficient observations (initial check)",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(10)),
				values.Proto(values.NewInt64(20)),
			},
			descriptor: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN,
				},
			},
			f:               2,
			expectedOutcome: nil,
			expectedError:   errors.New("insufficient observations"),
		},
		{
			name: "median aggregation: happy path (int64)",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(30)),
				values.Proto(values.NewInt64(40)),
				values.Proto(values.NewInt64(10)),
				values.Proto(values.NewInt64(20)),
				values.Proto(values.NewInt64(50)),
			},
			descriptor: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN,
				},
			},
			f:               4,
			expectedOutcome: values.Proto(values.NewInt64(30)),
			expectedError:   nil,
		},
		{
			name: "identical aggregation: happy path (int64)",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(42)),
				values.Proto(values.NewInt64(42)),
				values.Proto(values.NewInt64(42)),
				values.Proto(values.NewInt64(42)),
				values.Proto(values.NewInt64(50)), // spurious
				values.Proto(values.NewString("malicious")),
			},
			descriptor: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL,
				},
			},
			f:               3,
			expectedOutcome: values.Proto(values.NewInt64(42)),
			expectedError:   nil,
		},
		{
			name: "median: mixed types, two eligible types (int64, float64) - error returned",
			f:    1,
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(10)), values.Proto(values.NewFloat64(1.0)),
				values.Proto(values.NewInt64(20)), values.Proto(values.NewFloat64(2.0)),
				values.Proto(values.NewInt64(30)), values.Proto(values.NewInt64(40)),
			},
			descriptor: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN,
				},
			},
			expectedError: ErrMoreThanOneValidOutcomeForIdenticalConsensus,
		},
		{
			name: "median: mixed types, one eligible type (int64) - handled by filtering",
			f:    2,
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(10)), values.Proto(values.NewFloat64(1.0)),
				values.Proto(values.NewInt64(20)), values.Proto(values.NewFloat64(2.0)),
				values.Proto(values.NewInt64(30)), values.Proto(values.NewInt64(40)),
			},
			descriptor: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN,
				},
			},
			expectedOutcome: values.Proto(values.NewInt64(20)),
			expectedError:   nil,
		},
		{
			name: "common prefix",
			observations: []*valuespb.Value{
				mustNewList("1", "2", "3", "4"),
				mustNewList("1", "2", "3", "5"),
				mustNewList("1", "2", "3", "6"),
				mustNewList("1", "2", "3", "7"),
				mustNewList(),
				mustNewList(42, 43, 44, 45),
			},
			f: 3,
			descriptor: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_COMMON_PREFIX,
				},
			},
			expectedOutcome: mustNewList("1", "2", "3"),
		},
		{
			name: "fields map",
			observations: []*valuespb.Value{
				mustWrap(s{Val: 42}),
				mustWrap(s{Val: 43}),
				mustWrap(s{Val: 43}),
				mustWrap(s{Val: 44}),
				mustWrap(s{Val: 44}),
			},
			f:               2,
			descriptor:      cre.ConsensusAggregationFromTags[s]().Descriptor(),
			expectedOutcome: mustWrap(s{Val: 43}),
		},
		{
			name: "fields map one field succeeds one field returns default",
			observations: []*valuespb.Value{ // identical consensus fails for OtherField
				mustWrap(s{Val: 42, OtherField: "A"}),
				mustWrap(s{Val: 43, OtherField: "B"}),
				mustWrap(s{Val: 43, OtherField: "C"}),
				mustWrap(s{Val: 44, OtherField: "D"}),
				mustWrap(s{Val: 44, OtherField: "E"}),
			},
			f:               2,
			descriptor:      cre.ConsensusAggregationFromTags[s]().Descriptor(),
			defaultValue:    mustWrap(s{OtherField: "Z"}),
			expectedOutcome: mustWrap(s{Val: 43, OtherField: "Z"}),
		},
		{
			name: "fields map: all fields succeed with diverse types and aggregations",
			observations: []*valuespb.Value{
				mustWrap(s{Val: 10, OtherField: "common", PrefixSlice: []int64{1, 2, 3}, Nest: s1{Val: 100}, SuffixSlice: []int64{0, 2, 3}}),
				mustWrap(s{Val: 20, OtherField: "common", PrefixSlice: []int64{1, 2, 4}, Nest: s1{Val: 100}, SuffixSlice: []int64{10, 2, 3}}),
				mustWrap(s{Val: 30, OtherField: "common", PrefixSlice: []int64{1, 2, 5}, Nest: s1{Val: 100}, SuffixSlice: []int64{100, 2, 3}}),
				mustWrap(s{Val: 40, OtherField: "common", PrefixSlice: []int64{1, 9, 8}, Nest: s1{Val: 101}, SuffixSlice: []int64{99, 2, 3}}),
				mustWrap(s{Val: 50, OtherField: "common", PrefixSlice: []int64{1, 2, 6}, Nest: s1{Val: 102}, SuffixSlice: []int64{42, 2, 3}}),
			},
			descriptor: cre.ConsensusAggregationFromTags[s]().Descriptor(),
			f:          2,
			expectedOutcome: mustWrap(s{
				Val:         30,
				OtherField:  "common",
				PrefixSlice: []int64{1, 2},
				SuffixSlice: []int64{2, 3},
				Nest:        s1{Val: 100},
			}),
			expectedError: nil,
		},
		{
			name: "common suffix",
			observations: []*valuespb.Value{
				mustNewList("1", "2", "3", "4", "5", "6", "7", "8", "9"),
				mustNewList("1", "2", "3", "4", "11", "42", "42", "7", "8", "9"),
				mustNewList("1", "2", "3", "4", "8", "9", "42", "7", "8", "9"),
				mustNewList("1", "2", "3", "100", "99", "7", "8", "9"),
				mustNewList("1", "2", "3", "10", "99", "7", "8", "9"),
				mustNewList("1", "2", "3", "110", "99", "7", "8", "9"),
				mustNewList("1", "2", "3", "1000", "99", "7", "8", "9"),
				mustNewList("1", "2", "3", "x", "99"),
				mustNewList("1", "2", "3", "err", "99"),
				mustNewList("1", "2", "3", "4", "err", "7", "8", "9"),
				mustNewList(),
				mustNewList(42, 44, 45, 46),
				mustNewList("bad", "values", "bad", "values"),
			},
			f: 7,
			descriptor: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_COMMON_SUFFIX,
				},
			},
			expectedOutcome: mustNewList("7", "8", "9"),
		},
		{
			name: "unknown aggregation type (UNSPECIFIED)",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(10)),
			},
			descriptor: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_UNSPECIFIED,
				},
			},
			expectedOutcome: nil,
			expectedError:   errors.New("unknown aggregation type"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outcome, err := CalculateOutcomeForObservations(
				logger.Test(t),
				tc.observations,
				tc.descriptor,
				tc.defaultValue,
				tc.f,
			)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error(), "Error message mismatch")
			} else {
				require.NoError(t, err, "Unexpected error for test case %s", tc.name)
			}

			if tc.expectedError == nil {
				require.True(t, proto.Equal(outcome, tc.expectedOutcome),
					"Outcome mismatch for test: %s\nExpected: %+v\nActual:   %+v", tc.name, tc.expectedOutcome, outcome)
			}
		})
	}
}

// Test_handleMedianAggregation tests the handleMedianAggregation function directly.
func Test_handleMedianAggregation(t *testing.T) {
	type testCase struct {
		name            string
		observations    []*valuespb.Value
		expectedOutcome *valuespb.Value
		expectedError   error
		f               int
	}

	testCases := []testCase{
		{
			name: "int64 median: basic five values",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(30)), values.Proto(values.NewInt64(40)), values.Proto(values.NewInt64(10)), values.Proto(values.NewInt64(20)), values.Proto(values.NewInt64(50)),
			},
			expectedOutcome: values.Proto(values.NewInt64(30)),
			expectedError:   nil,
			f:               2,
		},
		{
			name: "int64 median: even number of values returns left value",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(10)), values.Proto(values.NewInt64(20)), values.Proto(values.NewInt64(30)), values.Proto(values.NewInt64(40)),
			},
			expectedOutcome: values.Proto(values.NewInt64(20)),
			expectedError:   nil,
			f:               1,
		},
		{
			name: "uint64 median: basic five values",
			observations: []*valuespb.Value{
				values.Proto(values.NewUint64(30)), values.Proto(values.NewUint64(40)), values.Proto(values.NewUint64(10)), values.Proto(values.NewUint64(20)), values.Proto(values.NewUint64(50)),
			},
			expectedOutcome: values.Proto(values.NewUint64(30)),
			expectedError:   nil,
			f:               2,
		},
		{
			name: "uint64 median: even number of values returns left value",
			observations: []*valuespb.Value{
				values.Proto(values.NewUint64(10)), values.Proto(values.NewUint64(20)), values.Proto(values.NewUint64(30)), values.Proto(values.NewUint64(40)),
			},
			expectedOutcome: values.Proto(values.NewUint64(20)),
			expectedError:   nil,
			f:               1,
		},
		{
			name: "float64 median: basic five values",
			observations: []*valuespb.Value{
				values.Proto(values.NewFloat64(30.5)), values.Proto(values.NewFloat64(40.5)), values.Proto(values.NewFloat64(10.5)), values.Proto(values.NewFloat64(20.5)), values.Proto(values.NewFloat64(50.5)),
			},
			expectedOutcome: values.Proto(values.NewFloat64(30.5)),
			expectedError:   nil,
			f:               2,
		},
		{
			name: "float64 median: even number of values returns left value",
			observations: []*valuespb.Value{
				values.Proto(values.NewFloat64(10.5)), values.Proto(values.NewFloat64(20.5)), values.Proto(values.NewFloat64(30.5)), values.Proto(values.NewFloat64(40.5)),
			},
			expectedOutcome: values.Proto(values.NewFloat64(20.5)),
			expectedError:   nil,
			f:               1,
		},
		{
			name: "decimal median: basic five values",
			observations: []*valuespb.Value{
				values.Proto(values.NewDecimal(decimal.NewFromFloat(30.3))), values.Proto(values.NewDecimal(decimal.NewFromFloat(40.4))),
				values.Proto(values.NewDecimal(decimal.NewFromFloat(10.1))), values.Proto(values.NewDecimal(decimal.NewFromFloat(20.2))),
				values.Proto(values.NewDecimal(decimal.NewFromFloat(50.5))),
			},
			expectedOutcome: values.Proto(values.NewDecimal(decimal.NewFromFloat(30.3))),
			expectedError:   nil,
			f:               2,
		},
		{
			name: "decimal median: even number of values returns left value",
			observations: []*valuespb.Value{
				values.Proto(values.NewDecimal(decimal.NewFromFloat(10.1))), values.Proto(values.NewDecimal(decimal.NewFromFloat(20.2))),
				values.Proto(values.NewDecimal(decimal.NewFromFloat(30.3))), values.Proto(values.NewDecimal(decimal.NewFromFloat(40.4))),
			},

			expectedOutcome: values.Proto(values.NewDecimal(decimal.NewFromFloat(20.2))),
			expectedError:   nil,
			f:               1,
		},
		{
			name: "bigint median: basic five values",
			observations: []*valuespb.Value{
				values.Proto(values.NewBigInt(big.NewInt(300))), values.Proto(values.NewBigInt(big.NewInt(400))),
				values.Proto(values.NewBigInt(big.NewInt(100))), values.Proto(values.NewBigInt(big.NewInt(200))),
				values.Proto(values.NewBigInt(big.NewInt(500))),
			},
			expectedOutcome: values.Proto(values.NewBigInt(big.NewInt(300))),
			expectedError:   nil,
			f:               2,
		},
		{
			name: "bigint median: even number of values returns left value",
			observations: []*valuespb.Value{
				values.Proto(values.NewBigInt(big.NewInt(100))), values.Proto(values.NewBigInt(big.NewInt(200))),
				values.Proto(values.NewBigInt(big.NewInt(300))), values.Proto(values.NewBigInt(big.NewInt(400))),
			},
			expectedOutcome: values.Proto(values.NewBigInt(big.NewInt(200))),
			expectedError:   nil,
			f:               1,
		},
		{
			name: "time median: basic five values",
			observations: []*valuespb.Value{
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:30Z"))),
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:40Z"))),
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:10Z"))),
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:20Z"))),
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:50Z"))),
			},
			expectedOutcome: values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:30Z"))),
			expectedError:   nil,
			f:               2,
		},
		{
			name: "time median: even number of values returns left value",
			observations: []*valuespb.Value{
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:10Z"))),
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:20Z"))),
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:30Z"))),
				values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:40Z"))),
			},
			expectedOutcome: values.Proto(values.NewTime(parseTime(t, "2023-01-01T00:00:20Z"))),
			expectedError:   nil,
			f:               1,
		},
		{
			name: "median: unsupported type for median aggregation (string)",
			observations: []*valuespb.Value{
				values.Proto(values.NewString("foo")),
				values.Proto(values.NewString("bar")),
				values.Proto(values.NewString("baz")),
				values.Proto(values.NewString("bah")),
				values.Proto(values.NewString("cad")),
			},
			expectedOutcome: nil,
			expectedError:   errors.New("unsupported type for median aggregation: " + typeString.Name()),
			f:               2,
		},
		{
			name:            "empty filtered observations for median",
			observations:    []*valuespb.Value{},
			expectedOutcome: nil,
			expectedError:   errors.New("insufficient observations"),
			f:               2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outcome, err := handleMedianAggregation(
				logger.Test(t),
				tc.observations,
				tc.f,
			)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error(), "Error message mismatch")
			} else {
				require.NoError(t, err, "Unexpected error for test case %s", tc.name)
			}

			if tc.expectedError == nil {
				require.True(t, proto.Equal(outcome, tc.expectedOutcome),
					"Outcome mismatch for %s\nExpected: %+v\nActual:   %+v", tc.name, tc.expectedOutcome, outcome)
			}
		})
	}
}

func parseTime(t *testing.T, s string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, s)
	require.NoError(t, err)
	return parsedTime
}

func mustWrap(v any) *valuespb.Value {
	wrapped, err := values.Wrap(v)
	if err != nil {
		panic(err)
	}
	return values.Proto(wrapped)
}

// Test_FieldsMapAggregation_ErrorDeterminism verifies that handleFieldsMapAggregation
// always returns the same error when multiple fields fail aggregation with no defaults.
// The desc map keys are iterated in sorted order so the first-failing key is stable.
func Test_FieldsMapAggregation_ErrorDeterminism(t *testing.T) {
	lggr := logger.Test(t)
	f := 2

	observations := make([]*valuespb.Value, 2*f+1)
	for i := range observations {
		observations[i] = valuespb.NewMapValue(map[string]*valuespb.Value{
			"Alpha":   values.Proto(values.NewString("unique-" + string(rune('A'+i)))),
			"Beta":    values.Proto(values.NewString("unique-" + string(rune('Z'-i)))),
			"Gamma":   values.Proto(values.NewString("unique-" + string(rune('a'+i)))),
			"Delta":   values.Proto(values.NewString("unique-" + string(rune('z'-i)))),
			"Epsilon": values.Proto(values.NewString("unique-" + string(rune('0'+i)))),
		})
	}

	desc := map[string]*sdk.ConsensusDescriptor{
		"Alpha":   {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
		"Beta":    {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
		"Gamma":   {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
		"Delta":   {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
		"Epsilon": {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
	}

	seenErrors := map[string]bool{}
	for i := 0; i < 200; i++ {
		_, err := handleFieldsMapAggregation(lggr, observations, desc, nil, f)
		require.Error(t, err)
		seenErrors[err.Error()] = true
	}

	require.Equal(t, 1, len(seenErrors),
		"handleFieldsMapAggregation must return a deterministic error regardless of map iteration order")
	require.Contains(t, firstKey(seenErrors), "aggregation for field failed 'Alpha'",
		"sorted iteration should always fail on the alphabetically first key")
}

func firstKey(m map[string]bool) string {
	for k := range m {
		return k
	}
	return ""
}
