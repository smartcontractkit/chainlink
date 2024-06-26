package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var maskTestConfigCmd = &cobra.Command{
	Use:   "mask-secrets",
	Short: "Run 'echo ::add-mask::${secret}' for all secrets in the test config",
	Run: func(cmd *cobra.Command, args []string) {
		mustMaskSecret("chainlink image", oc.ChainlinkImage)
		mustMaskSecret("chainlink version", oc.ChainlinkVersion)
		mustMaskSecret("pyroscope server url", oc.PyroscopeServerURL)
		mustMaskSecret("pyroscope key", oc.PyroscopeKey)
		mustMaskSecret("loki tenant id", oc.LoggingLokiTenantID)
		mustMaskSecret("loki endpoint", oc.LoggingLokiEndpoint)
		mustMaskSecret("loki basic auth", oc.LoggingLokiBasicAuth)
		mustMaskSecret("loki grafana token", oc.LoggingGrafanaToken)
		mustMaskSecret("pyroscope environment", oc.PyroscopeEnvironment)
		mustMaskSecret("pyroscope server url", oc.PyroscopeServerURL)
		mustMaskSecret("pyroscope key", oc.PyroscopeKey)
	},
}

func mustMaskSecret(description string, secret string) {
	if secret != "" {
		fmt.Printf("Mask '%s'\n", description)
		fmt.Printf("::add-mask::%s\n", secret)
		cmd := exec.Command("bash", "-c", fmt.Sprintf("echo ::add-mask::%s", secret))
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to mask secret '%s'", description)
			os.Exit(1)
		}
	}
}
