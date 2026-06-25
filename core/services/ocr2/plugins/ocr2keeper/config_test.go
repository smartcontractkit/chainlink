package ocr2keeper

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalDuration(t *testing.T) {
	t.Parallel()
	raw := `"2s"`

	var value Duration
	err := json.Unmarshal([]byte(raw), &value)

	require.NoError(t, err)
	assert.Equal(t, 2*time.Second, value.Value())
}

func TestUnmarshalConfig(t *testing.T) {
	t.Parallel()
	raw := `{"cacheExpiration":"2s","maxServiceWorkers":42}`

	var config PluginConfig
	err := json.Unmarshal([]byte(raw), &config)

	require.NoError(t, err)
	assert.Equal(t, 2*time.Second, config.CacheExpiration.Value())
	assert.Equal(t, 42, config.MaxServiceWorkers)
}
