package slicelib

import (
	"math/big"
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
				{Int: big.NewInt(1)},
				{Int: big.NewInt(2)},
				{Int: big.NewInt(4)},
				{Int: big.NewInt(5)},
			},
			want: model.BigInt{Int: big.NewInt(4)},
		},
		{
			name: "not sorted",
			vals: []model.BigInt{
				{Int: big.NewInt(100)},
				{Int: big.NewInt(50)},
				{Int: big.NewInt(30)},
				{Int: big.NewInt(110)},
			},
			want: model.BigInt{Int: big.NewInt(100)},
		},
		{
			name: "empty slice",
			vals: []model.BigInt{},
			want: model.BigInt{},
		},
		{
			name: "one item",
			vals: []model.BigInt{
				{Int: big.NewInt(123)},
			},
			want: model.BigInt{Int: big.NewInt(123)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, BigIntSortedMiddle(tt.vals), "BigIntSortedMiddle(%v)", tt.vals)
		})
	}
}
