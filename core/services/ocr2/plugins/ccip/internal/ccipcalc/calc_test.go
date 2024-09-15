package ccipcalc

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestMergeEpochAndRound(t *testing.T) {
	type args struct {
		epoch uint32
		round uint8
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "zero round and epoch",
			args: args{epoch: 0, round: 0},
			want: 0,
		},
		{
			name: "avg case",
			args: args{
				epoch: 243,
				round: 15,
			},
			want: 62223,
		},
		{
			name: "largest epoch and round",
			args: args{
				epoch: math.MaxUint32,
				round: math.MaxUint8,
			},
			want: 1099511627775,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want,
				MergeEpochAndRound(tt.args.epoch, tt.args.round),
				"mergeEpochAndRound(%v, %v)", tt.args.epoch, tt.args.round)
		})
	}
}

func TestContiguousReqs(t *testing.T) {
	testCases := []struct {
		min    uint64
		max    uint64
		seqNrs []uint64
		exp    bool
	}{
		{min: 5, max: 10, seqNrs: []uint64{5, 6, 7, 8, 9, 10}, exp: true},
		{min: 5, max: 10, seqNrs: []uint64{5, 7, 8, 9, 10}, exp: false},
		{min: 5, max: 10, seqNrs: []uint64{5, 6, 7, 8, 9, 10, 11}, exp: false},
		{min: 5, max: 10, seqNrs: []uint64{}, exp: false},
		{min: 1, max: 1, seqNrs: []uint64{1}, exp: true},
		{min: 6, max: 10, seqNrs: []uint64{5, 7, 8, 9, 10}, exp: false},
	}

	for _, tc := range testCases {
		res := ContiguousReqs(logger.NullLogger, tc.min, tc.max, tc.seqNrs)
		assert.Equal(t, tc.exp, res)
	}
}

func TestCalculateUsdPerUnitGas(t *testing.T) {
	testCases := []struct {
		name           string
		sourceGasPrice *big.Int
		usdPerFeeCoin  *big.Int
		exp            *big.Int
	}{
		{
			name:           "base case",
			sourceGasPrice: big.NewInt(2e18),
			usdPerFeeCoin:  big.NewInt(3e18),
			exp:            big.NewInt(6e18),
		},
		{
			name:           "small numbers",
			sourceGasPrice: big.NewInt(1000),
			usdPerFeeCoin:  big.NewInt(2000),
			exp:            big.NewInt(0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := CalculateUsdPerUnitGas(tc.sourceGasPrice, tc.usdPerFeeCoin)
			assert.Zero(t, tc.exp.Cmp(res))
		})
	}
}

func TestBigIntSortedMiddle(t *testing.T) {
	tests := []struct {
		name string
		vals []*big.Int
		want *big.Int
	}{
		{
			name: "base case",
			vals: []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(4), big.NewInt(5)},
			want: big.NewInt(4),
		},
		{
			name: "not sorted",
			vals: []*big.Int{big.NewInt(100), big.NewInt(50), big.NewInt(30), big.NewInt(110)},
			want: big.NewInt(100),
		},
		{
			name: "empty slice",
			vals: []*big.Int{},
			want: nil,
		},
		{
			name: "one item",
			vals: []*big.Int{big.NewInt(123)},
			want: big.NewInt(123),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, BigIntSortedMiddle(tt.vals), "BigIntSortedMiddle(%v)", tt.vals)
		})
	}
}

func TestDeviates(t *testing.T) {
	type args struct {
		x1  *big.Int
		x2  *big.Int
		ppb int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "base case",
			args: args{x1: big.NewInt(1e9), x2: big.NewInt(2e9), ppb: 1},
			want: true,
		},
		{
			name: "x1 is zero and x1 neq x2",
			args: args{x1: big.NewInt(0), x2: big.NewInt(1), ppb: 999},
			want: true,
		},
		{
			name: "x2 is zero and x1 neq x2",
			args: args{x1: big.NewInt(1), x2: big.NewInt(0), ppb: 999},
			want: true,
		},
		{
			name: "x1 and x2 are both zero",
			args: args{x1: big.NewInt(0), x2: big.NewInt(0), ppb: 999},
			want: false,
		},
		{
			name: "deviates when ppb is 0",
			args: args{x1: big.NewInt(0), x2: big.NewInt(1), ppb: 0},
			want: true,
		},
		{
			name: "does not deviate when x1 eq x2",
			args: args{x1: big.NewInt(5), x2: big.NewInt(5), ppb: 1},
			want: false,
		},
		{
			name: "does not deviate with high ppb when x2 is greater",
			args: args{x1: big.NewInt(5), x2: big.NewInt(10), ppb: 2e9},
			want: false,
		},
		{
			name: "does not deviate with high ppb when x1 is greater",
			args: args{x1: big.NewInt(10), x2: big.NewInt(5), ppb: 2e9},
			want: false,
		},
		{
			name: "deviates with low ppb when x2 is greater",
			args: args{x1: big.NewInt(5), x2: big.NewInt(10), ppb: 9e8},
			want: true,
		},
		{
			name: "deviates with low ppb when x1 is greater",
			args: args{x1: big.NewInt(10), x2: big.NewInt(5), ppb: 9e8},
			want: true,
		},
		{
			name: "near deviation limit but deviates",
			args: args{x1: big.NewInt(10), x2: big.NewInt(5), ppb: 1e9 - 1},
			want: true,
		},
		{
			name: "at deviation limit but does not deviate",
			args: args{x1: big.NewInt(10), x2: big.NewInt(5), ppb: 1e9},
			want: false,
		},
		{
			name: "near deviation limit but does not deviate",
			args: args{x1: big.NewInt(10), x2: big.NewInt(5), ppb: 1e9 + 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Deviates(tt.args.x1, tt.args.x2, tt.args.ppb), "Deviates(%v, %v, %v)", tt.args.x1, tt.args.x2, tt.args.ppb)
		})
	}
}
