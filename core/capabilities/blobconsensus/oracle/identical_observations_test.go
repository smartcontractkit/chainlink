package oracle

import (
	"errors"
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandleIdenticalAggregation(t *testing.T) {
	nowTime := timestamppb.Now()

	tests := []struct {
		name        string
		inputValues []*valuespb.Value
		wantValue   *valuespb.Value
		wantErr     string
		f           int
	}{
		{
			name:        "NOK - Empty slice",
			inputValues: []*valuespb.Value{},
			wantErr:     "input slice cannot be empty",
		},
		{
			name: "NOK - Single string value",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_StringValue{StringValue: "hello"}},
			},
			f:       1,
			wantErr: "no values met f+1 threshold",
		},
		{
			name: "OK - Multiple identical string values",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_StringValue{StringValue: "apple"}},
				{Value: &valuespb.Value_StringValue{StringValue: "apple"}},
				{Value: &valuespb.Value_StringValue{StringValue: "apple"}},
			},
			f: 1,
			wantValue: &valuespb.Value{
				Value: &valuespb.Value_StringValue{StringValue: "apple"},
			},
		},
		{
			name: "OK - Multiple identical int64 values",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_Int64Value{Int64Value: 100}},
				{Value: &valuespb.Value_Int64Value{Int64Value: 100}},
			},
			f: 1,
			wantValue: &valuespb.Value{
				Value: &valuespb.Value_Int64Value{Int64Value: 100},
			},
		},
		{
			name: "OK - Multiple identical boolean values",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_BoolValue{BoolValue: true}},
				{Value: &valuespb.Value_BoolValue{BoolValue: true}},
			},
			f: 1,
			wantValue: &valuespb.Value{
				Value: &valuespb.Value_BoolValue{BoolValue: true},
			},
		},
		{
			name: "OK - Multiple identical float64 values",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_Float64Value{Float64Value: 3.14}},
				{Value: &valuespb.Value_Float64Value{Float64Value: 3.14}},
			},
			f: 1,
			wantValue: &valuespb.Value{
				Value: &valuespb.Value_Float64Value{Float64Value: 3.14},
			},
		},
		{
			name: "OK - Multiple identical time values",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_TimeValue{TimeValue: nowTime}},
				{Value: &valuespb.Value_TimeValue{TimeValue: nowTime}},
			},
			f: 1,
			wantValue: &valuespb.Value{
				Value: &valuespb.Value_TimeValue{TimeValue: nowTime},
			},
		},
		{
			name: "NOK - Multiple different string values",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_StringValue{StringValue: "alpha"}},
				{Value: &valuespb.Value_StringValue{StringValue: "beta"}},
			},
			f:       1,
			wantErr: "no values met f+1 threshold",
		},
		{
			name: "NOK - Mixed types",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_Int64Value{Int64Value: 1}},
				{Value: &valuespb.Value_StringValue{StringValue: "two"}},
				{Value: &valuespb.Value_BoolValue{BoolValue: false}},
			},
			f:       1,
			wantErr: "no values met f+1 threshold",
		},
		{
			name: "NOK - Mixed types and multiple options",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_Int64Value{Int64Value: 1}},
				{Value: &valuespb.Value_Int64Value{Int64Value: 1}},
				{Value: &valuespb.Value_StringValue{StringValue: "two"}},
				{Value: &valuespb.Value_StringValue{StringValue: "two"}},
				{Value: &valuespb.Value_BoolValue{BoolValue: false}},
			},
			f:       1,
			wantErr: "not identical",
		},
		{
			name: "OK - Slice with map values",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_MapValue{MapValue: &valuespb.Map{
					Fields: map[string]*valuespb.Value{
						"key":  {Value: &valuespb.Value_StringValue{StringValue: "value"}},
						"key2": {Value: &valuespb.Value_Float64Value{Float64Value: 3.14}},
						"key3": {Value: &valuespb.Value_Int64Value{Int64Value: 42}},
					},
				}}},
				{Value: &valuespb.Value_MapValue{MapValue: &valuespb.Map{
					Fields: map[string]*valuespb.Value{
						"key":  {Value: &valuespb.Value_StringValue{StringValue: "value"}},
						"key2": {Value: &valuespb.Value_Float64Value{Float64Value: 3.14}},
						"key3": {Value: &valuespb.Value_Int64Value{Int64Value: 42}},
					},
				}}},
			},
			f: 1,
			wantValue: &valuespb.Value{
				Value: &valuespb.Value_MapValue{MapValue: &valuespb.Map{
					Fields: map[string]*valuespb.Value{
						"key":  {Value: &valuespb.Value_StringValue{StringValue: "value"}},
						"key2": {Value: &valuespb.Value_Float64Value{Float64Value: 3.14}},
						"key3": {Value: &valuespb.Value_Int64Value{Int64Value: 42}},
					},
				}},
			},
		},
		{
			name: "OK - Slice with list values",
			inputValues: []*valuespb.Value{
				{Value: &valuespb.Value_ListValue{ListValue: &valuespb.List{
					Fields: []*valuespb.Value{{Value: &valuespb.Value_Int64Value{Int64Value: 1}}},
				}}},
				{Value: &valuespb.Value_ListValue{ListValue: &valuespb.List{
					Fields: []*valuespb.Value{{Value: &valuespb.Value_Int64Value{Int64Value: 1}}},
				}}},
			},
			f: 1,
			wantValue: &valuespb.Value{
				Value: &valuespb.Value_ListValue{
					ListValue: &valuespb.List{
						Fields: []*valuespb.Value{
							{
								Value: &valuespb.Value_Int64Value{Int64Value: 1},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)
			got, err := handleIdenticalAggregation(lggr, tc.inputValues, tc.f)

			if tc.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			require.True(t, proto.Equal(tc.wantValue, got))
		})
	}
}

func assertDeepEqualValuesSlice(t *testing.T, expected, actual []*valuespb.Value) {
	require.Len(t, actual, len(expected), "Slice length mismatch")
	for i := range expected {
		assert.True(t, proto.Equal(actual[i], expected[i]))
	}
}

func Test_filterObservations(t *testing.T) {
	type testCase struct {
		name                 string
		observationProtos    []*valuespb.Value
		minObservations      int
		expectedObservations []*valuespb.Value
		expectedType         reflect.Type
		expectedError        error
	}

	testCases := []testCase{
		{
			name: "sufficient observations of dominant type (int64)",
			observationProtos: []*valuespb.Value{
				values.Proto(values.NewInt64(10)),
				values.Proto(values.NewFloat64(1.1)),
				values.Proto(values.NewInt64(20)),
				values.Proto(values.NewFloat64(2.2)),
				values.Proto(values.NewInt64(30)),
			},
			minObservations: 3,
			expectedObservations: []*valuespb.Value{
				values.Proto(values.NewInt64(10)),
				values.Proto(values.NewInt64(20)),
				values.Proto(values.NewInt64(30)),
			},
			expectedType:  typeInt64,
			expectedError: nil,
		},
		{
			name: "insufficient total observations (initial check)",
			observationProtos: []*valuespb.Value{
				values.Proto(values.NewInt64(10)),
				values.Proto(values.NewInt64(20)),
			},
			minObservations:      3,
			expectedObservations: nil,
			expectedType:         nil,
			expectedError:        errors.New("insufficient observations (2) to meet minimum (3)"),
		},
		{
			name: "no dominant type meeting threshold",
			observationProtos: []*valuespb.Value{
				values.Proto(values.NewInt64(10)),
				values.Proto(values.NewInt64(20)),
				values.Proto(values.NewFloat64(1.1)),
				values.Proto(values.NewFloat64(2.2)),
			},
			minObservations:      3,
			expectedObservations: nil,
			expectedType:         nil,
			expectedError:        errors.New("no values met f+1 threshold"),
		},
		{
			name: "dominant type is TypeNil",
			observationProtos: []*valuespb.Value{
				values.Proto(nil),
				values.Proto(nil),
				values.Proto(values.NewInt64(10)),
			},
			minObservations:      2,
			expectedObservations: nil,
			expectedType:         nil,
			expectedError:        errors.New("no values met f+1 threshold"),
		},
		{
			name: "all observations are of dominant type",
			observationProtos: []*valuespb.Value{
				values.Proto(values.NewFloat64(1.1)),
				values.Proto(values.NewFloat64(2.2)),
				values.Proto(values.NewFloat64(3.3)),
			},
			minObservations: 3,
			expectedObservations: []*valuespb.Value{
				values.Proto(values.NewFloat64(1.1)),
				values.Proto(values.NewFloat64(2.2)),
				values.Proto(values.NewFloat64(3.3)),
			},
			expectedType:  typeFloat64,
			expectedError: nil,
		},
		{
			name: "mixed types but only dominant passes filter",
			observationProtos: []*valuespb.Value{
				values.Proto(values.NewInt64(100)),
				values.Proto(values.NewString("test")),
				values.Proto(values.NewInt64(200)),
				values.Proto(values.NewBool(true)),
			},
			minObservations: 2,
			expectedObservations: []*valuespb.Value{
				values.Proto(values.NewInt64(100)),
				values.Proto(values.NewInt64(200)),
			},
			expectedType:  typeInt64,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualObservations, actualType, err := filterObservations(tc.observationProtos, tc.minObservations)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				assert.Nil(t, actualObservations)
				assert.Nil(t, actualType)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedType.Name(), actualType.Name())
				assertDeepEqualValuesSlice(t, tc.expectedObservations, actualObservations)
			}
		})
	}
}
