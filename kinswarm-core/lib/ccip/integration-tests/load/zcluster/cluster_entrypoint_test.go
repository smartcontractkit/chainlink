package zcluster

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestClusterEntrypoint(t *testing.T) {
	config, err := tc.GetConfig([]string{"Load"}, tc.OCR)
	require.NoError(t, err)
	cfgBase64, err := config.AsBase64()
	require.NoError(t, err)
	p, err := wasp.NewClusterProfile(&wasp.ClusterConfig{
		// you set up these only once, no need to configure through TOML
		DockerCmdExecPath: "../../..",
		BuildCtxPath:      "integration-tests/load",
		Namespace:         *config.WaspConfig.Namespace,
		KeepJobs:          config.WaspConfig.KeepJobs,
		UpdateImage:       config.WaspConfig.UpdateImage,
		HelmValues: map[string]string{
			"env.loki.url":       *config.Logging.Loki.Endpoint,
			"env.loki.tenant_id": *config.Logging.Loki.TenantId,
			"image":              *config.WaspConfig.RepoImageVersionURI,
			"test.binaryName":    *config.WaspConfig.TestBinaryName,
			"test.name":          *config.WaspConfig.TestName,
			"test.timeout":       *config.WaspConfig.TestTimeout,
			"env.wasp.log_level": *config.WaspConfig.WaspLogLevel,
			"jobs":               *config.WaspConfig.WaspJobs,
			// other test vars pass through
			"test.BASE64_CONFIG_OVERRIDE": cfgBase64,
		},
	})
	require.NoError(t, err)
	err = p.Run()
	require.NoError(t, err)
}
