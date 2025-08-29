package environment

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/deploy"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/trigger"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/verify"
	cronbasedtypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
	webapitriggerbasedtypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/web-trigger-based/types"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

func deployAndVerifyExampleWorkflowCmd() *cobra.Command {
	var (
		rpcURLFlag                  string
		gatewayURLFlag              string
		donIDFlag                   string
		exampleWorkflowTriggerFlag  string
		exampleWorkflowTimeoutFlag  string
		workflowRegistryAddressFlag string
	)
	cmd := &cobra.Command{
		Use:   "run-por-example",
		Short: "Runs v1 Proof-of-Reserve example workflow",
		Long:  `Deploys a simple Proof-of-Reserve workflow and, optionally, wait until it succeeds`,
		RunE: func(cmd *cobra.Command, args []string) error {
			timeout, timeoutErr := time.ParseDuration(exampleWorkflowTimeoutFlag)
			if timeoutErr != nil {
				return errors.Wrapf(timeoutErr, "failed to parse %s to time.Duration", exampleWorkflowTimeoutFlag)
			}

			return deployAndVerifyExampleWorkflow(cmd.Context(), rpcURLFlag, gatewayURLFlag, donIDFlag, timeout, exampleWorkflowTriggerFlag, workflowRegistryAddressFlag)
		},
	}

	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&exampleWorkflowTriggerFlag, "example-workflow-trigger", "y", "web-trigger", "Trigger for example workflow to deploy (web-trigger or cron)")
	cmd.Flags().StringVarP(&exampleWorkflowTimeoutFlag, "example-workflow-timeout", "u", "5m", "Time to wait until example workflow succeeds")
	cmd.Flags().StringVarP(&gatewayURLFlag, "gateway-url", "g", "http://localhost:5002", "Gateway URL (only for web API trigger-based workflow)")
	cmd.Flags().StringVarP(&donIDFlag, "don-id", "d", "vault", "DON ID (only for web API trigger-based workflow)")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "w", DefaultWorkflowRegistryAddress, "Workflow registry address")

	return cmd
}

type executableWorkflowFn = func(cmdContext context.Context, rpcURL, gatewayURL, donID, privateKey string, consumerContractAddress common.Address, feedID string, waitTime time.Duration, startTime time.Time) error

func executeWebTriggerBasedWorkflow(cmdContext context.Context, rpcURL, gatewayURL, donID, privateKey string, consumerContractAddress common.Address, feedID string, waitTime time.Duration, startTime time.Time) error {
	ticker := 5 * time.Second
	for {
		select {
		case <-time.After(waitTime):
			fmt.Print(libformat.PurpleText("\n[Stage 3/3] Example workflow failed to execute successfully in %.2f seconds\n", time.Since(startTime).Seconds()))

			return fmt.Errorf("example workflow failed to execute successfully within %s", waitTime)
		case <-time.Tick(ticker):
			triggerErr := trigger.WebAPITriggerValue(
				gatewayURL,
				donID,
				privateKey,
				5*time.Minute,
			)
			if triggerErr == nil {
				verifyTime := 25 * time.Second
				verifyErr := verify.ProofOfReserve(rpcURL, consumerContractAddress.Hex(), feedID, true, verifyTime)
				if verifyErr == nil {
					if isBlockscoutRunning(cmdContext) {
						fmt.Print(libformat.PurpleText("Open http://localhost/address/%s?tab=internal_txns to check consumer contract's transaction history\n", consumerContractAddress.Hex()))
					}

					return nil
				}

				fmt.Printf("\nTrying to verify workflow again in %.2f seconds...\n\n", ticker.Seconds())
			} else {
				framework.L.Debug().Msgf("failed to trigger web API trigger: %s", triggerErr)
			}
		}
	}
}

func executeCronBasedWorkflow(cmdContext context.Context, rpcURL, _, _, privateKey string, consumerContractAddress common.Address, feedID string, waitTime time.Duration, startTime time.Time) error {
	// we ignore return as if verification failed it will print that info
	verifyErr := verify.ProofOfReserve(rpcURL, consumerContractAddress.Hex(), feedID, true, waitTime)
	if verifyErr != nil {
		fmt.Print(libformat.PurpleText("\n[Stage 3/3] Example workflow failed to execute successfully in %.2f seconds\n", time.Since(startTime).Seconds()))
		return errors.Wrap(verifyErr, "failed to verify example workflow")
	}

	if isBlockscoutRunning(cmdContext) {
		fmt.Print(libformat.PurpleText("Open http://localhost/address/%s?tab=internal_txns to check consumer contract's transaction history\n", consumerContractAddress.Hex()))
	}

	return nil
}

func deployAndVerifyExampleWorkflow(cmdContext context.Context, rpcURL, gatewayURL, donID string, timeout time.Duration, exampleWorkflowTrigger, workflowRegistryAddress string) error {
	totalStart := time.Now()
	start := time.Now()

	if pkErr := creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return pkErr
	}

	fmt.Print(libformat.PurpleText("[Stage 1/3] Deploying Permissionless Feeds Consumer\n\n"))
	consumerContractAddress, consumerErr := deploy.PermissionlessFeedsConsumer(rpcURL)
	if consumerErr != nil {
		return errors.Wrap(consumerErr, "failed to deploy Permissionless Feeds Consumer contract")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 1/3] Deployed Permissionless Feeds Consumer in %.2f seconds\n", time.Since(start).Seconds()))

	start = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 2/3] Registering example Proof-of-Reserve workflow\n\n"))

	var executableWorkflowFunction executableWorkflowFn

	var workflowName string
	var workflowFilePath string
	var configFilePath string
	var configErr error
	feedID := "0x018e16c39e0003200000000000000000"

	if strings.EqualFold(exampleWorkflowTrigger, WorkflowTriggerCron) {
		workflowName = "cron-based-proof-of-reserve"
		workflowFilePath = "examples/workflows/v1/proof-of-reserve/cron-based/main.go"
		configFilePath, configErr = builAndSavePoRCronConfig(consumerContractAddress.Hex(), feedID, filepath.Dir(workflowFilePath))
		if configErr != nil {
			return errors.Wrap(configErr, "failed to build and save PoR config")
		}
		executableWorkflowFunction = executeCronBasedWorkflow
	} else {
		workflowName = "web-trigger-based-proof-of-reserve"
		workflowFilePath = "examples/workflows/v1/proof-of-reserve/web-trigger-based/main.go"
		configFilePath, configErr = builAndSavePoRWebTriggerConfig(consumerContractAddress.Hex(), feedID, filepath.Dir(workflowFilePath))
		if configErr != nil {
			return errors.Wrap(configErr, "failed to build and save PoR config")
		}
		executableWorkflowFunction = executeWebTriggerBasedWorkflow
	}

	defer func() {
		_ = os.Remove(configFilePath)
	}()

	parsed, err := strconv.ParseUint(donID, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse DON ID %s: %w", donID, err)
	}
	did := uint32(parsed)

	deployErr := compileCopyAndRegisterWorkflow(cmdContext, workflowFilePath, workflowName, "", workflowRegistryAddress, "", creworkflow.DefaultWorkflowNodePattern, creworkflow.DefaultWorkflowTargetDir, configFilePath, "", rpcURL, did)
	if deployErr != nil {
		return errors.Wrap(deployErr, "failed to deploy example workflow")
	}

	// Print workflow owner and name for debugging purposes
	workflowOwner := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266") // Default Anvil address
	fmt.Printf("Workflow Owner: %s\n", workflowOwner.Hex())
	fmt.Printf("Workflow Name: %s\n", workflowName)

	fmt.Print(libformat.PurpleText("\n[Stage 2/3] Registered workflow in %.2f seconds\n", time.Since(start).Seconds()))
	fmt.Print(libformat.PurpleText("[Stage 3/3] Waiting for %.2f seconds for workflow to execute successfully\n\n", timeout.Seconds()))

	var pauseWorkflow = func() {
		fmt.Print(libformat.PurpleText("\n[Stage 3/3] Example workflow executed in %.2f seconds\n", time.Since(totalStart).Seconds()))
		start = time.Now()
		fmt.Print(libformat.PurpleText("\n[CLEANUP] Deleting example workflow\n\n"))
		deleteErr := deleteAllWorkflows(cmdContext, rpcURL, workflowRegistryAddress)
		if deleteErr != nil {
			fmt.Printf("Failed to delete example workflow: %s\nPlease delete it manually\n", deleteErr)
		}

		fmt.Print(libformat.PurpleText("\n[CLEANUP] Deleted example workflow in %.2f seconds\n\n", time.Since(start).Seconds()))
	}
	defer pauseWorkflow()

	if pkErr := creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return pkErr
	}

	return executableWorkflowFunction(cmdContext, rpcURL, gatewayURL, donID, os.Getenv("PRIVATE_KEY"), *consumerContractAddress, feedID, timeout, totalStart)
}

func builAndSavePoRWebTriggerConfig(dataFeedsCacheAddress, feedID, folder string) (string, error) {
	cfg := webapitriggerbasedtypes.WorkflowConfig{
		DataFeedsCacheAddress: dataFeedsCacheAddress,
		AllowedTriggerSender:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		AllowedTriggerTopic:   "sendValue",
		FeedID:                feedID,
		WriteTargetName:       "write_geth-testnet@1.0.0",
	}

	yaml, yamlErr := yaml.Marshal(cfg)
	if yamlErr != nil {
		return "", errors.Wrap(yamlErr, "failed to marshal config to YAML")
	}

	filePath := filepath.Join(folder, "web_trigger_config.yaml")
	writeErr := os.WriteFile(filePath, yaml, 0644) //nolint:gosec // G306: we want it to be readable by everyone
	if writeErr != nil {
		return "", errors.Wrap(writeErr, "failed to write config to file")
	}

	return filePath, nil
}

func builAndSavePoRCronConfig(dataFeedsCacheAddress, feedID, folder string) (string, error) {
	if feedID == "" {
		return "", errors.New("feedID is empty")
	}

	cfg := cronbasedtypes.WorkflowConfig{
		ComputeConfig: cronbasedtypes.ComputeConfig{
			DataFeedsCacheAddress: dataFeedsCacheAddress,
			URL:                   "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
			FeedID:                feedID,
			WriteTargetName:       "write_geth-testnet@1.0.0",
		},
	}

	yaml, yamlErr := yaml.Marshal(cfg)
	if yamlErr != nil {
		return "", errors.Wrap(yamlErr, "failed to marshal config to YAML")
	}

	filePath := filepath.Join(folder, "cron_config.yaml")
	writeErr := os.WriteFile(filePath, yaml, 0644) //nolint:gosec // G306: we want it to be readable by everyone
	if writeErr != nil {
		return "", errors.Wrap(writeErr, "failed to write config to file")
	}

	return filePath, nil
}
