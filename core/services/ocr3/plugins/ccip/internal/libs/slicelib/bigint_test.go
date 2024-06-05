package slicelib

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func TestBigIntSortedMiddle(t *testing.T) {
	tests := []struct {
		name string
		vals []cciptypes.BigInt
		want cciptypes.BigInt
	}{
		{
			name: "base case",
			vals: []cciptypes.BigInt{
				cciptypes.NewBigIntFromInt64(1),
				cciptypes.NewBigIntFromInt64(2),
				cciptypes.NewBigIntFromInt64(4),
				cciptypes.NewBigIntFromInt64(5),
			},
			want: cciptypes.NewBigIntFromInt64(4),
		},
		{
			name: "not sorted",
			vals: []cciptypes.BigInt{
				cciptypes.NewBigIntFromInt64(100),
				cciptypes.NewBigIntFromInt64(50),
				cciptypes.NewBigIntFromInt64(30),
				cciptypes.NewBigIntFromInt64(110),
			},
			want: cciptypes.NewBigIntFromInt64(100),
		},
		{
			name: "empty slice",
			vals: []cciptypes.BigInt{},
			want: cciptypes.BigInt{},
		},
		{
			name: "one item",
			vals: []cciptypes.BigInt{
				cciptypes.NewBigIntFromInt64(123),
			},
			want: cciptypes.NewBigIntFromInt64(123),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, BigIntSortedMiddle(tt.vals), "BigIntSortedMiddle(%v)", tt.vals)
		})
	}
}
