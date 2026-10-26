package static

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_short(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		ver, sha       string
		expVer, expSha string
	}{
		{"1.0", "1234567890", "1.0", "1234567"},
		{"1", "a", "1", "a"},
		{"", "", "unset", "unset"},
		{"1.0", "", "1.0", "unset"},
		{"", "1234567890", "unset", "1234567"},
	} {
		t.Run(tt.ver+":"+tt.sha, func(t *testing.T) {
			t.Parallel()
			sha, ver := short(tt.sha, tt.ver)
			assert.Equal(t, tt.expSha, sha)
			assert.Equal(t, tt.expVer, ver)
		})
	}
}
