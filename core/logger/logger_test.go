package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// no sampling
	assert.Nil(t, newZapConfigBase().Sampling)
	assert.Nil(t, newZapConfigTest().Sampling)
	assert.Nil(t, newZapConfigProd(false, false).Sampling)

	// not development, which would trigger panics for Critical level
	assert.False(t, newZapConfigBase().Development)
	assert.False(t, newZapConfigTest().Development)
	assert.False(t, newZapConfigProd(false, false).Development)
}

func Test_verShaName(t *testing.T) {
	for _, tt := range []struct {
		ver, sha string
		exp      string
	}{
		{"1.0", "1234567890", "1.0@1234567"},
		{"1", "a", "1@a"},
		{"", "", "unset@unset"},
		{"1.0", "", "1.0@unset"},
		{"", "1234567890", "unset@1234567"},
	} {
		t.Run(tt.ver+":"+tt.sha, func(t *testing.T) {
			if got := verShaName(tt.ver, tt.sha); got != tt.exp {
				t.Errorf("expected %q but got %q", tt.exp, got)
			}
		})
	}
}
