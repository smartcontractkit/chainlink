package oracle

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

type s struct {
	Val         int64   `consensus_aggregation:"median"`
	OtherField  string  `consensus_aggregation:"identical"`
	PrefixSlice []int64 `consensus_aggregation:"common_prefix"`
	SuffixSlice []int64 `consensus_aggregation:"common_suffix"`
	Nest        s1      `consensus_aggregation:"nested"`
	InlineSlice []struct {
		Counts []int64
		Data   map[string]string
	} `consensus_aggregation:"identical"`
}

type s1 struct {
	Val  uint8 `consensus_aggregation:"median"`
	Nest s2    `consensus_aggregation:"nested"`
}

type s2 struct {
	Names []string `consensus_aggregation:"identical"`
}

func Test_handleFieldsMapAggregation(t *testing.T) {
	type testCase struct {
		name            string
		observations    []*valuespb.Value
		descriptor      map[string]*sdk.ConsensusDescriptor
		defaultValue    *valuespb.Value
		f               int
		expectedOutcome *valuespb.Value
		expectedError   error
	}

	lggr := logger.Test(t)
	desc := cre.ConsensusAggregationFromTags[s]()
	require.NoError(t, desc.Err(), "failed to get aggregation desc from tags")

	testCases := []testCase{
		{
			name: "single aggregation of single field",
			observations: []*valuespb.Value{
				mustWrap(s{Val: 10}),
				mustWrap(s{Val: 20}),
				mustWrap(s{Val: 30}),
				mustWrap(s{Val: 40}),
				mustWrap(s{Val: 50}),
			},
			descriptor:      desc.Descriptor().GetFieldsMap().GetFields(),
			f:               2,
			expectedOutcome: mustWrap(s{Val: 30}),
			expectedError:   nil,
		},
		{
			name: "mixed aggregation of multiple fields",
			observations: []*valuespb.Value{
				mustWrap(s{Val: 10, OtherField: "abc", PrefixSlice: []int64{1, 2, 3}, SuffixSlice: []int64{7, 8, 9}}),
				mustWrap(s{Val: 20, OtherField: "abc", PrefixSlice: []int64{1, 2, 4}, SuffixSlice: []int64{6, 8, 9}}),
				mustWrap(s{Val: 30, OtherField: "abc", PrefixSlice: []int64{1, 5, 6}, SuffixSlice: []int64{5, 8, 9}}),
				mustWrap(s{Val: 40, OtherField: "abc", PrefixSlice: []int64{1, 2, 7}, SuffixSlice: []int64{4, 8, 9}}),
				mustWrap(s{Val: 50, OtherField: "ghi", PrefixSlice: []int64{1, 2, 8}, SuffixSlice: []int64{3, 8, 9}}),
			},
			descriptor: desc.Descriptor().GetFieldsMap().GetFields(),
			f:          2,
			expectedOutcome: mustWrap(s{
				Val:         30,
				OtherField:  "abc",
				PrefixSlice: []int64{1, 2},
				SuffixSlice: []int64{8, 9},
			}),
			expectedError: nil,
		},
		{
			name: "all nested structs reach consensus",
			observations: []*valuespb.Value{
				mustWrap(s{Nest: s1{Val: 10, Nest: s2{Names: []string{"n1", "n2"}}}}),
				mustWrap(s{Nest: s1{Val: 20, Nest: s2{Names: []string{"n1", "n2"}}}}),
				mustWrap(s{Nest: s1{Val: 30, Nest: s2{Names: []string{"n1", "n2"}}}}),
				mustWrap(s{Nest: s1{Val: 40, Nest: s2{Names: []string{"n1", "n2"}}}}),
				mustWrap(s{Nest: s1{Val: 50, Nest: s2{Names: []string{"n1", "n2"}}}}),
			},
			descriptor: desc.Descriptor().GetFieldsMap().GetFields(),
			f:          2,
			expectedOutcome: mustWrap(s{
				Nest: s1{
					Val:  30,
					Nest: s2{Names: []string{"n1", "n2"}},
				},
			}),
			expectedError: nil,
		},
		{
			name: "error from child aggregation with no default value",
			observations: []*valuespb.Value{
				mustWrap(s{OtherField: "A"}),
				mustWrap(s{OtherField: "B"}),
				mustWrap(s{OtherField: "C"}),
				mustWrap(s{OtherField: "D"}),
				mustWrap(s{OtherField: "E"}),
			},
			descriptor:      desc.Descriptor().GetFieldsMap().GetFields(),
			f:               2,
			expectedOutcome: nil,
			expectedError:   errors.New("aggregation for field failed"),
		},
		{
			name: "child aggregation fails and returns default value",
			observations: []*valuespb.Value{
				mustWrap(s{OtherField: "A"}),
				mustWrap(s{OtherField: "B"}),
				mustWrap(s{OtherField: "C"}),
				mustWrap(s{OtherField: "D"}),
				mustWrap(s{OtherField: "E"}),
			},
			descriptor:      desc.Descriptor().GetFieldsMap().GetFields(),
			f:               2,
			defaultValue:    mustWrap(s{OtherField: "Z"}),
			expectedOutcome: mustWrap(s{OtherField: "Z"}),
		},
		{
			name: "insufficient observations fails",
			observations: []*valuespb.Value{
				mustWrap(s{}),
			},
			descriptor:    nil,
			f:             1,
			expectedError: ErrInsufficientObservations,
		},
		{
			name: "all observations are skipped with no defaults returns error",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(1)),
				values.Proto(values.NewBool(false)),
				values.Proto(values.NewString("test")),
				mustWrap(struct{}{}),
			},
			descriptor:      nil,
			f:               1,
			expectedOutcome: mustWrap(map[string]any{}),
		},
		{
			name: "all observations are skipped with description and no default errors",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(1)),
				values.Proto(values.NewBool(false)),
				values.Proto(values.NewString("test")),
				mustWrap(struct{}{}),
			},
			descriptor: desc.Descriptor().GetFieldsMap().GetFields(),
			f:          1,
			// non-deterministic which field will error first
			expectedError: errors.New("aggregation for field failed"),
		},
		{
			name: "all observations skipped with description with defaults returns default",
			observations: []*valuespb.Value{
				values.Proto(values.NewInt64(1)),
				values.Proto(values.NewBool(false)),
				values.Proto(values.NewString("test")),
				mustWrap(struct{}{}),
			},
			descriptor:      desc.Descriptor().GetFieldsMap().GetFields(),
			f:               3,
			defaultValue:    mustWrap(s{}),
			expectedOutcome: mustWrap(s{}),
		},
		{
			name: "nested s2 fields map with slices",
			observations: []*valuespb.Value{
				mustWrap(s{
					InlineSlice: inlineSlice(),
					Nest:        s1{Nest: s2{Names: []string{"n1", "n2"}}},
				}),
				mustWrap(s{
					InlineSlice: inlineSlice(),
					Nest:        s1{Nest: s2{Names: []string{"n1", "n2"}}},
				}),
				mustWrap(s{
					InlineSlice: inlineSlice(),
					Nest:        s1{Nest: s2{Names: []string{"n1", "n2"}}},
				}),
				mustWrap(s{
					InlineSlice: inlineSlice(),
					Nest:        s1{Nest: s2{Names: []string{"n1", "n2"}}},
				}),
				mustWrap(s{
					InlineSlice: inlineSlice(),
					Nest:        s1{Nest: s2{Names: []string{"n1", "n2"}}},
				}),
				mustWrap(s{
					InlineSlice: inlineSlice(),
					Nest:        s1{Nest: s2{Names: []string{"n1", "n2"}}},
				}),
			},
			descriptor: desc.Descriptor().GetFieldsMap().GetFields(),
			f:          2,
			expectedOutcome: mustWrap(s{
				InlineSlice: inlineSlice(),
				Nest: s1{Nest: s2{
					Names: []string{"n1", "n2"},
				}},
			},
			),
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := handleFieldsMapAggregation(lggr, tc.observations, tc.descriptor, tc.defaultValue, tc.f)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			require.True(t, proto.Equal(tc.expectedOutcome, got), "expected %v, got %v", tc.expectedOutcome, got)
		})
	}
}

func inlineSlice() []struct {
	Counts []int64
	Data   map[string]string
} {
	return []struct {
		Counts []int64
		Data   map[string]string
	}{
		{
			Counts: []int64{1, 2, 3},
			Data: map[string]string{
				"k1": "v1",
			},
		},

		{
			Counts: []int64{1, 2, 3},
			Data: map[string]string{
				"k1": "v1",
			},
		},

		{
			Counts: []int64{1, 2, 3},
			Data: map[string]string{
				"k1": "v1",
			},
		},
	}
}
