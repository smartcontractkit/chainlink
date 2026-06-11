package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForInstall(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want string
	}{
		{"2.53.0", "v2.53.0"},
		{"2.12.2", "v2.12.2"},
		{"0.1.7", "v0.1.7"},
		{"29.3", "v29.3"},
		{"1.26.4", "v1.26.4"},
		{"42dc7da8c2874db550e91c656f98d05fca3c2f98", "42dc7da8c2874db550e91c656f98d05fca3c2f98"},
		{"v1.2.3", "v1.2.3"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, ForInstall(tt.in), "ForInstall(%q)", tt.in)
	}
}
