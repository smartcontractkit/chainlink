package slicelib

import (
	"testing"

	"github.com/smartcontractkit/ccipocr3/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestBigIntSortedMiddle(t *testing.T) {
	tests := []struct {
		name string
		vals []model.BigInt
		want model.BigInt
	}{
		{
			name: "base case",
			vals: []model.BigInt{
				model.NewBigIntFromInt64(1),
				model.NewBigIntFromInt64(2),
				model.NewBigIntFromInt64(4),
				model.NewBigIntFromInt64(5),
			},
			want: model.NewBigIntFromInt64(4),
		},
		{
			name: "not sorted",
			vals: []model.BigInt{
				model.NewBigIntFromInt64(100),
				model.NewBigIntFromInt64(50),
				model.NewBigIntFromInt64(30),
				model.NewBigIntFromInt64(110),
			},
			want: model.NewBigIntFromInt64(100),
		},
		{
			name: "empty slice",
			vals: []model.BigInt{},
			want: model.BigInt{},
		},
		{
			name: "one item",
			vals: []model.BigInt{
				model.NewBigIntFromInt64(123),
			},
			want: model.NewBigIntFromInt64(123),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, BigIntSortedMiddle(tt.vals), "BigIntSortedMiddle(%v)", tt.vals)
		})
	}
}
