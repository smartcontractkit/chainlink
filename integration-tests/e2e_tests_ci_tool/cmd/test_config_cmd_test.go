package cmd

import (
	"bytes"
	"testing"

	"github.com/pelletier/go-toml/v2"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateTestConfigCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    interface{}
		check   func(t *testing.T, tc *ctf_config.TestConfig)
		wantErr bool
	}{
		{
			name: "LoggingLogTargets",
			args: []string{"create", "--logging-log-targets=target1", "--logging-log-targets=target2"},
			check: func(t *testing.T, tc *ctf_config.TestConfig) {
				assert.NotNil(t, tc.Logging)
				assert.NotNil(t, tc.Logging.LogStream)
				assert.Equal(t, []string{"target1", "target2"}, tc.Logging.LogStream.LogTargets)
			},
		},
		{
			name: "PrivateEthereumNetworkExecutionLayerFlag",
			args: []string{"create", "--private-ethereum-network-execution-layer=geth", "--private-ethereum-network-ethereum-version=1.10.0"},
			check: func(t *testing.T, tc *ctf_config.TestConfig) {
				assert.NotNil(t, tc.PrivateEthereumNetwork)
				assert.NotNil(t, tc.PrivateEthereumNetwork.ExecutionLayer)
				assert.Equal(t, ctf_config.ExecutionLayer("geth"), *tc.PrivateEthereumNetwork.ExecutionLayer)
				assert.Equal(t, ctf_config.EthereumVersion("1.10.0"), *tc.PrivateEthereumNetwork.EthereumVersion)
			},
		},
		{
			name: "PrivateEthereumNetworkCustomDockerImageFlag",
			args: []string{"create", "--private-ethereum-network-execution-layer=geth", "--private-ethereum-network-ethereum-version=1.10.0", "--private-ethereum-network-custom-docker-image=custom-image:v1.0"},
			check: func(t *testing.T, tc *ctf_config.TestConfig) {
				assert.NotNil(t, tc.PrivateEthereumNetwork)
				assert.NotNil(t, tc.PrivateEthereumNetwork.ExecutionLayer)
				assert.Equal(t, map[ctf_config.ContainerType]string{"execution_layer": "custom-image:v1.0"}, tc.PrivateEthereumNetwork.CustomDockerImages)
			},
		},
	}

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(createTestConfigCmd)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)
			var out bytes.Buffer
			rootCmd.SetOutput(&out)
			err := rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			var tc ctf_config.TestConfig
			err = toml.Unmarshal(out.Bytes(), &tc)
			if err != nil {
				t.Fatalf("Failed to unmarshal output: %v", err)
			}
			if tt.check != nil {
				tt.check(t, &tc)
			}
		})
	}
}
