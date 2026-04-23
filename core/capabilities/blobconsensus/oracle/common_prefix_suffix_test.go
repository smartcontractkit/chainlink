package oracle

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

func Test_handleCommonPrefixAggregation(t *testing.T) {
	type testCase struct {
		name       string
		giveValues []*valuespb.Value
		f          int
		wantValue  *valuespb.Value
		wantErr    string
	}

	testCases := []testCase{
		{
			name: "OK - common prefix of f+1 lists",
			giveValues: []*valuespb.Value{
				mustNewList("1", "2", "3", "4"),
				mustNewList("1", "2", "3", "5"),
				mustNewList("1", "2", "3", "6"),
				mustNewList("1", "2", "3", "7"),
			},
			f:         3,
			wantValue: mustNewList("1", "2", "3"),
		},
		{
			name: "OK - common prefix with one list disregarded after first index",
			giveValues: []*valuespb.Value{
				mustNewList("1", "2", "3"),
				mustNewList("1", "3", "2"),
				mustNewList("3", "2", "1"),
			},
			f:         1,
			wantValue: mustNewList("1"),
		},
		{
			name: "OK - common prefix with one list disregarded after second index",
			giveValues: []*valuespb.Value{
				mustNewList("1", "2", "3"),
				mustNewList("1", "3", "2"),
				mustNewList("3", "2", "1"),
				mustNewList("1", "2", "err"),
			},
			f:         1,
			wantValue: mustNewList("1", "2"),
		},
		{
			name: "OK - common prefix of f+1 lists mixed",
			giveValues: []*valuespb.Value{
				mustNewList("1", "2", "3", "4", "5", "6", "7", "8", "9"),
				mustNewList("1", "2", "3", "4", "11", "42", "42"),
				mustNewList("1", "2", "3", "4", "8", "9", "42"),
				mustNewList("1", "2", "3", "100", "99"),
				mustNewList("1", "2", "3", "10", "99"),
				mustNewList("1", "2", "3", "110", "99"),
				mustNewList("1", "2", "3", "1000", "99"),
				mustNewList("1", "2", "3", "x", "99"),
				mustNewList("1", "2", "3", "err", "99"),
				mustNewList("1", "2", "3", "4", "err"),
				mustNewList(),
				mustNewList(42, 44, 45, 46),
				mustNewList("bad", "values", "bad", "values"),
			},
			f:         7,
			wantValue: mustNewList("1", "2", "3"),
		},
		{
			name: "OK - fails identical consensus on first check returns empty slice",
			giveValues: []*valuespb.Value{
				mustNewList("1", "2", "3", "4"),
				mustNewList("1", "2", "3", "5"),
				mustNewList("0", "1", "2", "5"),
				mustNewList("1", "2", "3", "7"),
			},
			f:         3,
			wantValue: mustNewList(),
		},
		{
			name:       "NOK - no lists provided returns error",
			giveValues: []*valuespb.Value{},
			f:          3,
			wantErr:    ErrInsufficientObservations.Error(),
		},
		{
			name: "NOK - less than f+1 lists provided errors",
			giveValues: []*valuespb.Value{
				mustNewList(),
				mustNewList(),
				mustNewList(),
			},
			f:       3,
			wantErr: ErrInsufficientObservations.Error(),
		},
		{
			name: "OK - f+1 empty lists returns empty list",
			giveValues: []*valuespb.Value{
				mustNewList(),
				mustNewList(),
				mustNewList(),
				mustNewList(),
				values.Proto(values.NewString("bad entry")), // dropped
			},
			f:         3,
			wantValue: mustNewList(),
		},
		{
			name: "NOK - all non-lists provided",
			giveValues: []*valuespb.Value{
				values.Proto(values.NewInt64(42)),
				values.Proto(values.NewInt64(99)),
				values.Proto(values.NewString("bad entry")),
				values.Proto(values.NewString("bad entry")),
			},
			f:       3,
			wantErr: ErrInsufficientObservations.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)
			got, err := handleCommonPrefixAggregation(lggr, tc.giveValues, tc.f)

			if tc.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			require.True(t, proto.Equal(tc.wantValue, got), "expected %v, got %v", tc.wantValue, got)
		})
	}
}

func Test_handleCommonSuffixAggregation(t *testing.T) {
	type testCase struct {
		name       string
		giveValues []*valuespb.Value
		f          int
		wantValue  *valuespb.Value
		wantErr    string
	}

	testCases := []testCase{
		{
			name: "OK - common suffix of f+1 lists",
			giveValues: []*valuespb.Value{
				mustNewList("4", "2", "3", "1"),
				mustNewList("5", "2", "3", "1"),
				mustNewList("6", "2", "3", "1"),
				mustNewList("7", "2", "3", "1"),
			},
			f:         3,
			wantValue: mustNewList("2", "3", "1"),
		},
		{
			name: "OK - common suffix of mixed lists",
			giveValues: []*valuespb.Value{
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
			f:         7,
			wantValue: mustNewList("7", "8", "9"),
		},
		{
			name: "OK - fails identical consensus on first check returns empty slice",
			giveValues: []*valuespb.Value{
				mustNewList("4", "2", "3", "1"),
				mustNewList("5", "1", "2", "0"),
				mustNewList("6", "2", "3", "1"),
				mustNewList("7", "2", "3", "1"),
			},
			f:         3,
			wantValue: mustNewList(),
		},
		{
			name:       "NOK - no lists provided returns error",
			giveValues: []*valuespb.Value{},
			f:          3,
			wantErr:    ErrInsufficientObservations.Error(),
		},
		{
			name: "OK - f+1 empty lists provided returns empty list",
			giveValues: []*valuespb.Value{
				mustNewList(),
				mustNewList(),
				mustNewList(),
				mustNewList(),
			},
			f:         3,
			wantValue: values.Proto(&values.List{}),
		},
		{
			name: "NOK - less than f+1 lists to select from",
			giveValues: []*valuespb.Value{
				mustNewList(),
				mustNewList(),
				values.Proto(values.NewString("bad entry")), // dropped
			},
			f:       3,
			wantErr: ErrInsufficientObservations.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)
			got, err := handleCommonSuffixAggregation(lggr, tc.giveValues, tc.f)

			if tc.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			require.True(t, proto.Equal(tc.wantValue, got), "expected %v, got %v", tc.wantValue, got)
		})
	}
}

func mustNewList(elems ...any) *valuespb.Value {
	l, err := values.NewList(elems)
	if err != nil {
		panic(err)
	}
	return values.Proto(l)
}
