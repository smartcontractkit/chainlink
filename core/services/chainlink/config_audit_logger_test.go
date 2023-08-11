package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuditLoggerConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	auditConfig := cfg.AuditLogger()

	require.Equal(t, true, auditConfig.Enabled())
	require.Equal(t, "event", auditConfig.JsonWrapperKey())

	fUrl, err := auditConfig.ForwardToUrl()
	require.NoError(t, err)
	require.Equal(t, "http", fUrl.Scheme)
	require.Equal(t, "localhost:9898", fUrl.Host)

	headers, err := auditConfig.Headers()
	require.NoError(t, err)
	require.Len(t, headers, 2)
	require.Equal(t, "Authorization", headers[0].Header)
	require.Equal(t, "token", headers[0].Value)
	require.Equal(t, "X-SomeOther-Header", headers[1].Header)
	require.Equal(t, "value with spaces | and a bar+*", headers[1].Value)
}
