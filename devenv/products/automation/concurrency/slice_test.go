package concurrency_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"
)

func TestDivideSlice(t *testing.T) {
	tests := []struct {
		name    string
		slice   []int
		parts   int
		want    [][]int
		wantErr bool
	}{
		{
			name:    "empty slice",
			slice:   []int{},
			parts:   3,
			want:    [][]int{{}, {}, {}},
			wantErr: false,
		},
		{
			name:    "single element",
			slice:   []int{1},
			parts:   1,
			want:    [][]int{{1}},
			wantErr: false,
		},
		{
			name:    "two elements three parts",
			slice:   []int{1, 2},
			parts:   3,
			want:    [][]int{{1}, {2}, {}},
			wantErr: false,
		},
		{
			name:    "equal division",
			slice:   []int{1, 2, 3, 4},
			parts:   2,
			want:    [][]int{{1, 2}, {3, 4}},
			wantErr: false,
		},
		{
			name:    "non-equal division",
			slice:   []int{1, 2, 3, 4, 5},
			parts:   3,
			want:    [][]int{{1, 2}, {3, 4}, {5}},
			wantErr: false,
		},
		{
			name:    "more parts than elements",
			slice:   []int{1, 2, 3},
			parts:   5,
			want:    [][]int{{1}, {2}, {3}, {}, {}},
			wantErr: false,
		},
		{
			name:    "zero parts",
			slice:   []int{1, 2, 3},
			parts:   0,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative parts",
			slice:   []int{1, 2, 3},
			parts:   -1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("DivideSlice() caused a panic for input %v divided into %d parts", tt.slice, tt.parts)
				}
			}()

			got := concurrency.DivideSlice(tt.slice, tt.parts)

			if !tt.wantErr {
				require.Equal(t, tt.want, got, "DivideSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
