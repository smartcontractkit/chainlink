package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestWebServerConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ws := cfg.WebServer()
	assert.Equal(t, "*", ws.AllowOrigins())
	assert.Equal(t, "https://bridge.response", ws.BridgeResponseURL().String())
	assert.Equal(t, 10*time.Second, ws.BridgeCacheTTL())
	assert.Equal(t, 1*time.Minute, ws.HTTPServerWriteTimeout())
	assert.Equal(t, uint16(56), ws.Port())
	assert.True(t, ws.SecureCookies())
	assert.Equal(t, *models.MustNewDuration(1 * time.Hour), ws.SessionTimeout())
	assert.Equal(t, *models.MustNewDuration(168 * time.Hour), ws.ReaperExpiration())
	assert.Equal(t, int64(32770), ws.WebServerHTTPMaxSize())
	assert.Equal(t, 15*time.Second, ws.WebServerStartTimeout())
	assert.Equal(t, "test-rpid", ws.RPID())
	assert.Equal(t, "test-rp-origin", ws.RPOrigin())
	assert.Equal(t, int64(42), ws.AuthenticatedRateLimit())
	assert.Equal(t, *models.MustNewDuration(1 * time.Second), ws.AuthenticatedRateLimitPeriod())
	assert.Equal(t, int64(7), ws.UnAuthenticatedRateLimit())
	assert.Equal(t, *models.MustNewDuration(1 * time.Minute), ws.UnAuthenticatedRateLimitPeriod())

	tls := ws.TLS()
	assert.Equal(t, "test/root/dir/tls", tls.Dir())
	assert.Equal(t, "tls/cert/path", tls.CertPath())
	assert.True(t, tls.ForceRedirect())
	assert.Equal(t, "tls-host", tls.Host())
	assert.Equal(t, uint16(6789), tls.TLSPort())
	assert.Equal(t, "tls/key/path", tls.TLSKeyPath())
}
