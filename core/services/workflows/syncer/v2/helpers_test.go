package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandWorkflowFamilies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		families []string
		want     []string
	}{
		{
			name:     "no workflows suffix",
			families: []string{"zone-a"},
			want:     []string{"zone-a"},
		},
		{
			name:     "workflows suffix expands to base family",
			families: []string{"zone-a_workflows"},
			want:     []string{"zone-a_workflows", "zone-a"},
		},
		{
			name:     "base family already present is not duplicated",
			families: []string{"zone-a_workflows", "zone-a"},
			want:     []string{"zone-a_workflows", "zone-a"},
		},
		{
			name:     "mixed families",
			families: []string{"zone-a", "zone-b_workflows"},
			want:     []string{"zone-a", "zone-b_workflows", "zone-b"},
		},
		{
			name:     "empty input",
			families: []string{},
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, ExpandWorkflowFamilies(tt.families))
		})
	}
}
