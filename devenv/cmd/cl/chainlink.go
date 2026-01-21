package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// --- Observability Endpoints ---
// Predefined Grafana dashboards for monitoring performance and load tests
const (
	LocalWASPPerformanceDashboard = "http://localhost:3000/d/WASPLoadTests/wasp-load-test?orgId=1&refresh=5s"
	LocalCLDashboard              = "http://localhost:3000/d/cl-load-test?orgId=1&refresh=5s"
)

[Image of Chainlink node observability stack with Grafana and Prometheus]

/**
 * upCmd: Provisions the local development environment.
 * It configures environment variables and initializes testcontainers.
 */
var upCmd = &cobra.Command{
	Use:     "up",
	Short:   "Spin up the development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile := "env.toml,products/ocr2/basic.toml" // Default config
		if len(args) > 0 { configFile = args[0] }
		
		framework.L.Info().Str("Config", configFile).Msg("Initializing Environment")
		os.Setenv("CTF_CONFIGS", configFile)
		os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true") // Optimization: Disable Ryuk for faster local teardown
		
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return de.NewEnvironment(ctx)
	},
}

[Image of Docker Testcontainers lifecycle in a Go application]

/**
 * testCmd: Executes specific test suites (Soak, Gas, Chaos).
 * Uses 'go test' with long timeouts to simulate real-world conditions.
 */
var testCmd = &cobra.Command{
	Use:   "test [suite]",
	Short: "Run load or chaos tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Logic to map suite name to regex pattern (e.g., soak -> TestLoad/clean)
		// ...
		testCmd := exec.Command("go", "test", "-v", "-timeout", "4h", "-run", testPattern, "./...")
		testCmd.Stdout = os.Stdout
		return testCmd.Run()
	},
}

[Image of Software Testing Pyramid: Unit, Integration, and Load Testing]

/**
 * checkDockerIsRunning: Critical system check before execution.
 * Ensures the Docker daemon is responsive via Ping.
 */
func checkDockerIsRunning() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil || cli.Ping(context.Background()) != nil {
		fmt.Println("Error: Docker is not running!")
		os.Exit(1)
	}
}

func main() {
	checkDockerIsRunning()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
