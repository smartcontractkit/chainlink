package cmd

import (
	"bytes"
	"testing"

	"github.com/pelletier/go-toml/v2"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateTestConfigCmd_LoggingLogTargets(t *testing.T) {
	// Setup
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(createTestConfigCmd)

	// Set the flag
	args := []string{"create", "--logging-log-targets=target1", "--logging-log-targets=target2"}
	rootCmd.SetArgs(args)

	// Capture the output
	var out bytes.Buffer
	rootCmd.SetOutput(&out)
	// createTestConfigCmd.SetOutput(&out)

	// Run the command
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Check the output
	var tc ctf_config.TestConfig
	err = toml.Unmarshal(out.Bytes(), &tc)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	// Assertions
	assert.NotNil(t, tc.Logging)
	assert.NotNil(t, tc.Logging.LogStream)
	assert.Equal(t, []string{"target1", "target2"}, tc.Logging.LogStream.LogTargets)
}
