package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestTransformHostDockerInternalReferences(t *testing.T) {
	t.Parallel()

	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")
	cfg := &Config{
		NodeSets: []*cre.NodeSet{
			{
				NodeSpecs: []*cre.NodeSpecWithRole{
					{
						Input: &clnode.Input{
							Node: &clnode.NodeInput{},
						},
					},
				},
				CapabilityConfigs: map[cre.CapabilityFlag]cre.CapabilityConfig{
					cre.VaultCapability: {
						Values: map[string]any{
							"auth0": map[string]any{
								"issuerURL": "http://host.docker.internal:18123/",
								"urls":      []any{"host.docker.internal:18124"},
							},
						},
					},
				},
			},
		},
		CapabilityConfigs: map[string]cre.CapabilityConfig{
			cre.VaultCapability: {
				Values: map[string]any{
					"endpoint": "host.docker.internal:9999",
				},
			},
		},
	}
	cfg.NodeSets[0].NodeSpecs[0].Node.UserConfigOverrides = "[CRE.Linking]\nURL = \"host.docker.internal:18124\"\n"

	transformHostDockerInternalReferences(cfg)

	require.Contains(t, cfg.NodeSets[0].NodeSpecs[0].Node.UserConfigOverrides, dockerHost+":18124")

	auth0 := cfg.NodeSets[0].CapabilityConfigs[cre.VaultCapability].Values["auth0"].(map[string]any)
	require.Equal(t, framework.HostDockerInternal()+":18123/", auth0["issuerURL"])
	require.Equal(t, []any{dockerHost + ":18124"}, auth0["urls"])

	require.Equal(t, dockerHost+":9999", cfg.CapabilityConfigs[cre.VaultCapability].Values["endpoint"])
}
