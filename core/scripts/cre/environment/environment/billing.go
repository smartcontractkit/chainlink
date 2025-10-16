package environment

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	billingplatformservice "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/billing_platform_service"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

const DefaultBillingConfigFile = "configs/billing-platform-service.toml"

func billingCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billing",
		Short: "Billing management commands",
		Long:  `Commands to manage the billing platform service`,
	}

	cmd.AddCommand(startBillingCmds())
	cmd.AddCommand(stopBillingCmd)

	return cmd
}

func startBillingCmds() *cobra.Command {
	timeout := 15 * time.Second

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the billing platform service",
		Long:  `Starts the billing platform service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			initDxTracker()
			var startBillingErr error

			defer func() {
				metaData := map[string]any{}
				if startBillingErr != nil {
					metaData["result"] = "failure"
					metaData["error"] = oneLineErrorMessage(startBillingErr)
				} else {
					metaData["result"] = "success"
				}

				trackingErr := dxTracker.Track(MetricBillingStart, metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track billing start: %s\n", trackingErr)
				}
			}()

			// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
			setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
			if setErr != nil {
				return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
			}

			startBillingErr = startBilling(cmd.Context(), timeout, nil)
			if startBillingErr != nil {
				waitToCleanUp(timeout)
				billingRemoveErr := framework.RemoveTestStack(billingplatformservice.DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)
				if billingRemoveErr != nil {
					fmt.Fprint(os.Stderr, errors.Wrap(billingRemoveErr, manualBillingCleanupMsg).Error())
				}

				return errors.Wrap(startBillingErr, "failed to start Billing Platform Service")
			}

			return nil
		},
	}

	return cmd
}

var stopBillingCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop the billing platform service",
	Long:  `stop the billing platform service and clean up resources`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stopBilling()
	},
}

func stopBilling() error {
	setErr := os.Setenv("CTF_CONFIGS", DefaultBillingConfigFile)
	if setErr != nil {
		return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
	}

	removeCacheErr := removeBillingStateFiles(relativePathToRepoRoot)
	if removeCacheErr != nil {
		framework.L.Warn().Msgf("failed to remove cache files: %s\n", removeCacheErr)
	}

	return framework.RemoveTestStack(billingplatformservice.DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)
}

func removeBillingStateFiles(relativePathToRepoRoot string) error {
	path := filepath.Join(relativePathToRepoRoot, envconfig.StateDirname, envconfig.BillingStateFilename)
	absPath, absErr := filepath.Abs(path)
	if absErr != nil {
		return errors.Wrap(absErr, "error getting absolute path for billing platform service state file")
	}

	return os.Remove(absPath)
}

func startBilling(_ context.Context, cleanupWait time.Duration, setupOutput *environment.SetupOutput) (startupErr error) {
	// just in case, remove the stack if it exists
	_ = framework.RemoveTestStack(billingplatformservice.DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)

	defer func() {
		p := recover()

		if p != nil {
			fmt.Println("Panicked when starting Billing Platform Service")

			if err, ok := p.(error); ok {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

				startupErr = err
			} else {
				fmt.Fprintf(os.Stderr, "panic: %v\n", p)
				fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

				startupErr = fmt.Errorf("panic: %v", p)
			}

			time.Sleep(cleanupWait)

			billingRemoveErr := framework.RemoveTestStack(billingplatformservice.DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)
			if billingRemoveErr != nil {
				fmt.Fprint(os.Stderr, errors.Wrap(billingRemoveErr, manualBillingCleanupMsg).Error())
			}

			os.Exit(1)
		}
	}()

	stageGen := stagegen.NewStageGen(1, "STAGE")
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting Billing Platform Service")))

	previousCTFConfig := os.Getenv("CTF_CONFIGS")
	defer func() {
		setErr := os.Setenv("CTF_CONFIGS", previousCTFConfig)
		if setErr != nil {
			framework.L.Warn().Msgf("failed to restore previous CTF_CONFIGS environment variable: %s", setErr)
		}
	}()

	setErr := os.Setenv("CTF_CONFIGS", DefaultBillingConfigFile)
	if setErr != nil {
		return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
	}

	// Load and validate test configuration
	in, err := framework.Load[envconfig.BillingConfig](nil)
	if err != nil {
		return errors.Wrap(err, "failed to load test configuration")
	}

	if setupOutput != nil {
		in.BillingService.ChainSelector = setupOutput.WorkflowRegistryConfigurationOutput.ChainSelector
		addressRefs, err := setupOutput.CldEnvironment.DataStore.Addresses().Fetch()
		if err != nil {
			return errors.Wrap(err, "failed to fetch address references")
		}

		for _, ref := range addressRefs {
			// TODO CRE-878 fail test if test relies on v1 registries
			switch ref.Type {
			case "WorkflowRegistry":
				if in.BillingService.ChainSelector == ref.ChainSelector {
					in.BillingService.WorkflowRegistryAddress = ref.Address
				}
			case "CapabilitiesRegistry":
				if in.BillingService.ChainSelector == ref.ChainSelector {
					in.BillingService.CapabilitiesRegistryAddress = ref.Address
				}
			default:
				continue
			}
		}

		// Select the appropriate chain for billing service from available chains in the environment.
		// otherwise, if RPCURL is defined, billing service can be used standalone
		if len(setupOutput.BlockchainOutput) != 0 {
			var selectedChain *blockchain.Output

			for _, chain := range setupOutput.BlockchainOutput {
				if chain.ChainSelector == in.BillingService.ChainSelector {
					selectedChain = chain.BlockchainOutput
				}
			}

			if selectedChain == nil || len(selectedChain.Nodes) == 0 {
				return errors.Wrap(err, fmt.Sprintf("configured chain selector does not exist in the current topology: %d", in.BillingService.ChainSelector))
			}

			in.BillingService.RPCURL = strings.Replace(selectedChain.Nodes[0].ExternalHTTPUrl, "127.0.0.1", "host.docker.internal", 1)
		}

		in.BillingService.WorkflowOwners = make([]string, len(setupOutput.WorkflowRegistryConfigurationOutput.WorkflowOwners))
		for idx, owner := range setupOutput.WorkflowRegistryConfigurationOutput.WorkflowOwners {
			in.BillingService.WorkflowOwners[idx] = owner.Hex()
		}
	}

	out, startErr := billingplatformservice.New(in.BillingService)
	if startErr != nil {
		return errors.Wrap(startErr, "failed to create Billing Platform Service")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Started Billing Service stack in %.2f seconds", stageGen.Elapsed().Seconds())))

	fmt.Println()
	framework.L.Info().Msgf("Billing GRPC Service External URL: %s", out.BillingPlatformService.BillingGRPCExternalURL)
	framework.L.Info().Msgf("Credit Reservation GRPC Service External URL: %s", out.BillingPlatformService.CreditGRPCExternalURL)

	fmt.Println()
	fmt.Print("To terminate Billing stack execute: `go run . env billing stop`\n\n")

	in.BillingService.Output = out

	return in.Store(envconfig.MustBillingStateFileAbsPath(relativePathToRepoRoot))
}
