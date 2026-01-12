package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"
)

const (
	LocalWASPLoadDashboard = "http://localhost:3000/d/WASPLoadTests/wasp-load-test?orgId=1&from=now-5m&to=now&refresh=5s"
	LocalCLDashboard       = "http://localhost:3000/d/f8a04cef-653f-46d3-86df-87c532300672/cl-load-test?orgId=1&refresh=5s"
)

var rootCmd = &cobra.Command{
	Use:   "cl",
	Short: "A Chainlink local environment tool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		if debug {
			framework.L.Info().Msg("Debug mode enabled, setting CTF_CLNODE_DLV=true")
			os.Setenv("CTF_CLNODE_DLV", "true")
		}
		return nil
	},
}

var restartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"r"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Restart development environment, remove apps and apply default configuration again",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile string
		if len(args) > 0 {
			configFile = args[0]
		} else {
			configFile = "env.toml"
		}
		framework.L.Info().Str("Config", configFile).Msg("Reconfiguring development environment")
		_ = os.Setenv("CTF_CONFIGS", configFile)
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		framework.L.Info().Msg("Tearing down the development environment")
		err := framework.RemoveTestContainers()
		if err != nil {
			return fmt.Errorf("failed to clean Docker resources: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return de.NewEnvironment(ctx)
	},
}

var upCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"u"},
	Short:   "Spin up the development environment",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile string
		if len(args) > 0 {
			configFile = args[0]
		} else {
			configFile = "env.toml"
		}
		framework.L.Info().Str("Config", configFile).Msg("Creating development environment")
		_ = os.Setenv("CTF_CONFIGS", configFile)
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return de.NewEnvironment(ctx)
	},
}

var downCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"d"},
	Short:   "Tear down the development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		framework.L.Info().Msg("Tearing down the development environment")
		err := framework.RemoveTestContainers()
		if err != nil {
			return fmt.Errorf("failed to clean Docker resources: %w", err)
		}
		return nil
	},
}

var bsCmd = &cobra.Command{
	Use:   "bs",
	Short: "Manage the Blockscout EVM block explorer",
	Long:  "Spin up or down the Blockscout EVM block explorer",
}

var bsUpCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"u"},
	Short:   "Spin up Blockscout EVM block explorer",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := bsCmd.Flags().GetString("url")
		chainID, _ := bsCmd.Flags().GetString("chain-id")
		return framework.BlockScoutUp(url, chainID)
	},
}

var bsDownCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"d"},
	Short:   "Spin down Blockscout EVM block explorer",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := bsCmd.Flags().GetString("url")
		return framework.BlockScoutDown(url)
	},
}

var bsRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"r"},
	Short:   "Restart the Blockscout EVM block explorer",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := bsCmd.Flags().GetString("url")
		chainID, _ := bsCmd.Flags().GetString("chain-id")
		if err := framework.BlockScoutDown(url); err != nil {
			return err
		}
		return framework.BlockScoutUp(url, chainID)
	},
}

var obsCmd = &cobra.Command{
	Use:   "obs",
	Short: "Manage the observability stack",
	Long:  "Spin up or down the observability stack with subcommands 'up' and 'down'",
}

var obsUpCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"u"},
	Short:   "Spin up the observability stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		full, _ := cmd.Flags().GetBool("full")
		var err error
		if full {
			err = framework.ObservabilityUpFull()
		} else {
			err = framework.ObservabilityUp()
		}
		if err != nil {
			return fmt.Errorf("observability up failed: %w", err)
		}
		ocr2.L.Info().Msgf("OCR2 Dashboard: %s", LocalCLDashboard)
		ocr2.L.Info().Msgf("OCR2 Load Test Dashboard: %s", LocalWASPLoadDashboard)
		return nil
	},
}

var obsDownCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"d"},
	Short:   "Spin down the observability stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		return framework.ObservabilityDown()
	},
}

var obsRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"r"},
	Short:   "Restart the observability stack (data wipe)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := framework.ObservabilityDown(); err != nil {
			return fmt.Errorf("observability down failed: %w", err)
		}
		full, _ := cmd.Flags().GetBool("full")
		var err error
		if full {
			err = framework.ObservabilityUpFull()
		} else {
			err = framework.ObservabilityUp()
		}
		if err != nil {
			return fmt.Errorf("observability up failed: %w", err)
		}
		ocr2.L.Info().Msgf("OCR2 Dashboard: %s", LocalCLDashboard)
		ocr2.L.Info().Msgf("OCR2 Load Test Dashboard: %s", LocalWASPLoadDashboard)
		return nil
	},
}

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Run the tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("specify the test suite: smoke or load")
		}
		var testPattern string
		switch args[0] {
		case "load":
			testPattern = "TestLoad/clean"
		case "gas":
			testPattern = "TestLoad/gas_spikes"
		case "chaos":
			testPattern = "TestLoad/chaos"
		default:
			return fmt.Errorf("test suite %s is unknown, choose between smoke or load", args[0])
		}

		testCmd := exec.Command("go", "test", "-v", "-run", testPattern, "./...")
		testCmd.Dir = "./tests"
		testCmd.Stdout = os.Stdout
		testCmd.Stderr = os.Stderr
		testCmd.Stdin = os.Stdin

		if err := testCmd.Run(); err != nil {
			exitError := &exec.ExitError{}
			if errors.As(err, &exitError) {
				if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
				os.Exit(1)
			}
			return fmt.Errorf("failed to run test command: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable running services with dlv to allow remote debugging.")

	rootCmd.AddCommand(testCmd)

	// Blockscout, on-chain debug
	bsCmd.PersistentFlags().StringP("url", "u", "http://host.docker.internal:8555", "EVM RPC node URL (default to dst chain on 8555")
	bsCmd.PersistentFlags().StringP("chain-id", "c", "2337", "RPC's Chain ID")
	bsCmd.AddCommand(bsUpCmd)
	bsCmd.AddCommand(bsDownCmd)
	bsCmd.AddCommand(bsRestartCmd)
	rootCmd.AddCommand(bsCmd)

	// observability
	obsCmd.PersistentFlags().BoolP("full", "f", false, "Enable full observability stack with additional components")
	obsCmd.AddCommand(obsRestartCmd)
	obsCmd.AddCommand(obsUpCmd)
	obsCmd.AddCommand(obsDownCmd)
	rootCmd.AddCommand(obsCmd)

	// main env commands
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(downCmd)
}

func checkDockerIsRunning() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Can't create Docker client, please check if Docker daemon is running!")
		os.Exit(1)
	}
	_, err = cli.Ping(context.Background())
	if err != nil {
		fmt.Println("Docker is not running, please start Docker daemon first!")
		os.Exit(1)
	}
}

func main() {
	checkDockerIsRunning()
	if len(os.Args) == 2 && (os.Args[1] == "shell" || os.Args[1] == "sh") {
		_ = os.Setenv("CTF_CONFIGS", "env.toml") // Set default config for shell

		StartShell()
		return
	}
	if err := rootCmd.Execute(); err != nil {
		ocr2.L.Err(err).Send()
		os.Exit(1)
	}
}
