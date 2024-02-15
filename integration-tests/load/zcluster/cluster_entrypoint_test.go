package zcluster

import (
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestClusterEntrypoint(t *testing.T) {
	config, err := tc.GetConfig("Load", tc.OCR)
	require.NoError(t, err)
	cfgBase64, err := config.AsBase64()
	require.NoError(t, err)

	p, err := wasp.NewClusterProfile(&wasp.ClusterConfig{
		Namespace:    "wasp",
		UpdateImage:  true,
		BuildCtxPath: "..",
		HelmValues: map[string]string{
			"env.loki.url":        os.Getenv("LOKI_URL"),
			"env.loki.token":      os.Getenv("LOKI_TOKEN"),
			"env.loki.basic_auth": os.Getenv("LOKI_BASIC_AUTH"),
			"env.loki.tenant_id":  os.Getenv("LOKI_TENANT_ID"),
			"image":               os.Getenv("WASP_TEST_IMAGE"),
			"test.binaryName":     os.Getenv("WASP_TEST_BIN"),
			"test.name":           os.Getenv("WASP_TEST_NAME"),
			"env.wasp.log_level":  "debug",
			"jobs":                "3",
			// other test vars pass through
			"test.MY_CUSTOM_VAR":          "abc",
			"test.BASE64_CONFIG_OVERRIDE": cfgBase64,
		},
	})
	require.NoError(t, err)
	err = p.Run()
	require.NoError(t, err)
}
