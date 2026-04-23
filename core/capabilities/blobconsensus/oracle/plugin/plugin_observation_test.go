package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeduplicateRequestIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []string{}, want: []string{}},
		{name: "single", in: []string{"a"}, want: []string{"a"}},
		{name: "no duplicates", in: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		{name: "adjacent duplicates", in: []string{"a", "a", "b"}, want: []string{"a", "b"}},
		{name: "non adjacent duplicate drops later occurrence", in: []string{"a", "b", "a", "c"}, want: []string{"a", "b", "c"}},
		{name: "all same id", in: []string{"x", "x", "x"}, want: []string{"x"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := deduplicateRequestIDs(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
