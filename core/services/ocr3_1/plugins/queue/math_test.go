package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ts(sec int64) *timestamppb.Timestamp { return &timestamppb.Timestamp{Seconds: sec} }

func TestMedianTimestamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*timestamppb.Timestamp
		want time.Time
	}{
		{"single", []*timestamppb.Timestamp{ts(7)}, time.Unix(7, 0)},
		{"odd sorted", []*timestamppb.Timestamp{ts(1), ts(2), ts(3)}, time.Unix(2, 0)},
		{"odd unsorted", []*timestamppb.Timestamp{ts(9), ts(1), ts(5)}, time.Unix(5, 0)},
		{"even averages middle two", []*timestamppb.Timestamp{ts(10), ts(0), ts(4), ts(6)}, time.Unix(5, 0)},
		{"even with nanos", []*timestamppb.Timestamp{{Seconds: 1, Nanos: 0}, {Seconds: 1, Nanos: 10}}, time.Unix(1, 5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := medianTimestamp(tt.in).AsTime()
			assert.True(t, tt.want.Equal(got), "want %v got %v", tt.want, got)
		})
	}

	t.Run("does not reorder input", func(t *testing.T) {
		t.Parallel()
		in := []*timestamppb.Timestamp{ts(3), ts(1), ts(2)}
		medianTimestamp(in)
		require.Equal(t, []int64{3, 1, 2}, []int64{in[0].Seconds, in[1].Seconds, in[2].Seconds})
	})

	t.Run("empty panics", func(t *testing.T) {
		t.Parallel()
		require.Panics(t, func() { medianTimestamp(nil) })
	})
}
