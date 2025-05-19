package chainlink

import (
	"testing"
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, 1*time.Minute, ws.HTTPWriteTimeout())
	assert.Equal(t, uint16(56), ws.HTTPPort())
	assert.True(t, ws.SecureCookies())
	assert.Equal(t, *commonconfig.MustNewDuration(1 * time.Hour), ws.SessionTimeout())
	assert.Equal(t, *commonconfig.MustNewDuration(168 * time.Hour), ws.SessionReaperExpiration())
	assert.Equal(t, int64(32770), ws.HTTPMaxSize())
	assert.Equal(t, 15*time.Second, ws.StartTimeout())
	tls := ws.TLS()
	assert.Equal(t, "test/root/dir/tls", tls.Dir())
	assert.Equal(t, "tls/cert/path", tls.(*tlsConfig).certPath())
	assert.True(t, tls.ForceRedirect())
	assert.Equal(t, "tls-host", tls.Host())
	assert.Equal(t, uint16(6789), tls.HTTPSPort())
	assert.Equal(t, "tls/key/path", tls.(*tlsConfig).keyPath())

	rl := ws.RateLimit()
	assert.Equal(t, int64(42), rl.Authenticated())
	assert.Equal(t, 1*time.Second, rl.AuthenticatedPeriod())
	assert.Equal(t, int64(7), rl.Unauthenticated())
	assert.Equal(t, 1*time.Minute, rl.UnauthenticatedPeriod())

	mf := ws.MFA()
	assert.Equal(t, "test-rpid", mf.RPID())
	assert.Equal(t, "test-rp-origin", mf.RPOrigin())
}
