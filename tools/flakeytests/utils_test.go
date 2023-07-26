package flakeytests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigString(t *testing.T) {
	in := map[string]interface{}{
		"pull_request": map[string]interface{}{
			"url": "some-url",
		},
	}
	out, err := DigString(in, []string{"pull_request", "url"})
	require.NoError(t, err)
	assert.Equal(t, "some-url", out)
}
