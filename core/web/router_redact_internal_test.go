package web

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLogRedactsOIDCSecrets(t *testing.T) {
	t.Parallel()

	t.Run("signin query", func(t *testing.T) {
		t.Parallel()
		q, err := url.ParseQuery("code=secret-code&state=secret-state&other=kept")
		require.NoError(t, err)
		out := redact(q)
		assert.NotContains(t, out, "secret-code")
		assert.NotContains(t, out, "secret-state")
		assert.Contains(t, out, "other=kept")
	})

	t.Run("exchange body", func(t *testing.T) {
		t.Parallel()
		body := bytes.NewBufferString(`{"code":"secret-code","state":"secret-state","device_handle":"secret-handle","other":"kept"}`)
		out, err := readSanitizedJSON(body)
		require.NoError(t, err)
		assert.NotContains(t, out, "secret-code")
		assert.NotContains(t, out, "secret-state")
		assert.NotContains(t, out, "secret-handle")
		assert.Contains(t, out, `"other":"kept"`)
	})
}
