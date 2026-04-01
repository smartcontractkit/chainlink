package environment

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/deploy"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/verify"
	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/proof-of-reserve/cron-based/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func deployAndVerifyExampleWorkflowCmd() *cobra.Command {
	var (
		rpcURLFlag                  string
		workflowDonIDFlag           uint32
		exampleWorkflowTimeoutFlag  string
		workflowRegistryAddressFlag string
		contractsVersionFlag        string
	)
	cmd := &cobra.Command{
		Use:              "run-por-example",
		Short:            "Runs the PoR v2 cron example workflow",
		Long:             `Deploys the PoR v2 cron example workflow and waits until it succeeds`,
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			timeout, timeoutErr := time.ParseDuration(exampleWorkflowTimeoutFlag)
			if timeoutErr != nil {
				return errors.Wrapf(timeoutErr, "failed to parse %s to time.Duration", exampleWorkflowTimeoutFlag)
			}

			resolver, resolverErr := TryLoadLocalCREStateResolver()
			if resolverErr != nil {
				return errors.Wrap(resolverErr, "failed to load local CRE state")
			}

			rpcURL := rpcURLFlag
			if !cmd.Flags().Changed("rpc-url") && resolver != nil {
				if stateRPC, err := resolver.RegistryRPC(); err == nil {
					rpcURL = stateRPC
				}
			}

			workflowDONID := workflowDonIDFlag
			if !cmd.Flags().Changed("workflow-don-id") && resolver != nil {
				if stateDONID, err := resolver.WorkflowDONID(); err == nil {
					workflowDONID = stateDONID
				}
			}

			workflowRegistryAddress, contractsVersion, err := resolveContractAddressAndVersion(cmd, resolver, keystone_changeset.WorkflowRegistry, workflowRegistryAddressFlag, contractsVersionFlag, "workflow-registry-address")
			if err != nil {
				return errors.Wrap(err, "❌ failed to resolve workflow registry")
			}

			return deployAndVerifyExampleWorkflow(cmd.Context(), rpcURL, workflowDONID, timeout, workflowRegistryAddress, contractsVersion)
		},
	}

	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&exampleWorkflowTimeoutFlag, "example-workflow-timeout", "u", "5m", "Time to wait until example workflow succeeds (e.g. 10s, 1m, 1h)")
	cmd.Flags().Uint32VarP(&workflowDonIDFlag, "workflow-don-id", "d", 1, "DonID used in the workflow registry contract (integer starting with 1)")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "w", "", "Workflow registry address (if not provided, address from the state file will be used)")
	cmd.Flags().StringVar(&contractsVersionFlag, "with-contracts-version", "v2", "Version of workflow registry contract to use (v1 or v2)")

	return cmd
}

func executeCronBasedWorkflow(cmdContext context.Context, rpcURL string, consumerContractAddress common.Address, feedID string, waitTime time.Duration, startTime time.Time) error {
	// we ignore return as if verification failed it will print that info
	verifyErr := verify.ProofOfReserve(rpcURL, consumerContractAddress.Hex(), feedID, true, waitTime)
	if verifyErr != nil {
		fmt.Print(libformat.PurpleText("\n[Stage 4/4] Example workflow failed to execute successfully in %.2f seconds\n", time.Since(startTime).Seconds()))
		return errors.Wrap(verifyErr, "failed to verify example workflow")
	}

	if isBlockscoutRunning(cmdContext) {
		fmt.Print(libformat.PurpleText("Open http://localhost/address/%s?tab=internal_txns to check consumer contract's transaction history\n", consumerContractAddress.Hex()))
	}

	return nil
}

func deployAndVerifyExampleWorkflow(cmdContext context.Context, rpcURL string, workflowDonID uint32, timeout time.Duration, workflowRegistryAddress string, contractsVersion *semver.Version) error {
	totalStart := time.Now()
	start := time.Now()

	if pkErr := environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return pkErr
	}

	fmt.Print(libformat.PurpleText("[Stage 1/4] Deploying Permissionless Feeds Consumer\n\n"))
	consumerContractAddress, consumerErr := deploy.PermissionlessFeedsConsumer(rpcURL)
	if consumerErr != nil {
		return errors.Wrap(consumerErr, "failed to deploy Permissionless Feeds Consumer contract")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 1/4] Deployed Permissionless Feeds Consumer in %.2f seconds\n", time.Since(start).Seconds()))

	fmt.Print(libformat.PurpleText("[Stage 2/4] Deploying Balance Reader\n\n"))
	balanceReaderContractAddress, balanceReaderErr := deploy.BalanceReader(rpcURL)
	if balanceReaderErr != nil {
		return errors.Wrap(balanceReaderErr, "failed to deploy Balance Reader contract")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 2/4] Deployed Balance Reader in %.2f seconds\n", time.Since(start).Seconds()))

	start = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 3/4] Registering PoR v2 cron example workflow\n\n"))

	workflowName := "por-v2-cron-example"
	workflowFilePath := "examples/workflows/v2/proof-of-reserve/cron-based/main.go"
	feedID := "0x018e16c39e0003200000000000000000"
	chainID, chainSelector, chainErr := deploy.ChainMetadata(rpcURL)
	if chainErr != nil {
		return errors.Wrap(chainErr, "failed to resolve chain metadata for PoR config")
	}

	addressesToRead, addressesErr := deploy.CreateAndFundAddresses(rpcURL, 2, big.NewInt(10))
	if addressesErr != nil {
		return errors.Wrap(addressesErr, "failed to create and fund addresses for PoR config")
	}

	configFilePath, configErr := buildAndSavePoRV2CronConfig(
		consumerContractAddress.Hex(),
		balanceReaderContractAddress.Hex(),
		feedID,
		chainSelector,
		corevm.GenerateWriteTargetName(chainID),
		addressesToRead,
		filepath.Dir(workflowFilePath),
	)
	if configErr != nil {
		return errors.Wrap(configErr, "failed to build and save PoR config")
	}

	defer func() {
		_ = os.Remove(configFilePath)
	}()

	deployErr := compileCopyAndRegisterWorkflow(cmdContext, workflowFilePath, workflowName, "", workflowRegistryAddress, "", creworkflow.DefaultWorkflowNodePattern, creworkflow.DefaultWorkflowTargetDir, configFilePath, "", "", rpcURL, contractsVersion, contractsVersion, workflowDonID)
	if deployErr != nil {
		return errors.Wrap(deployErr, "failed to deploy example workflow")
	}

	// Print workflow owner and name for debugging purposes
	workflowOwner := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266") // Default Anvil address
	fmt.Printf("Workflow Owner: %s\n", workflowOwner.Hex())
	fmt.Printf("Workflow Name: %s\n", workflowName)

	fmt.Print(libformat.PurpleText("\n[Stage 3/4] Registered workflow in %.2f seconds\n", time.Since(start).Seconds()))
	fmt.Print(libformat.PurpleText("[Stage 4/4] Waiting for %.2f seconds for workflow to execute successfully\n\n", timeout.Seconds()))

	var pauseWorkflow = func() {
		fmt.Print(libformat.PurpleText("\n[Stage 4/4] Example workflow executed in %.2f seconds\n", time.Since(totalStart).Seconds()))
		start = time.Now()
		fmt.Print(libformat.PurpleText("\n[CLEANUP] Deleting example workflow\n\n"))
		deleteErr := deleteAllWorkflows(cmdContext, rpcURL, workflowRegistryAddress, contractsVersion)
		if deleteErr != nil {
			fmt.Printf("Failed to delete example workflow: %s\nPlease delete it manually\n", deleteErr)
		}

		fmt.Print(libformat.PurpleText("\n[CLEANUP] Deleted example workflow in %.2f seconds\n\n", time.Since(start).Seconds()))
	}
	defer pauseWorkflow()

	if pkErr := environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return pkErr
	}

	return executeCronBasedWorkflow(cmdContext, rpcURL, *consumerContractAddress, feedID, timeout, totalStart)
}

func buildAndSavePoRV2CronConfig(dataFeedsCacheAddress, balanceReaderAddress, feedID string, chainSelector uint64, writeTargetName string, addressesToRead []common.Address, folder string) (string, error) {
	if feedID == "" {
		return "", errors.New("feedID is empty")
	}
	if len(addressesToRead) < 2 {
		return "", errors.New("at least two addresses are required for the PoR v2 example")
	}

	cfg := portypes.WorkflowConfig{
		ChainSelector: chainSelector,
		ComputeConfig: portypes.ComputeConfig{
			DataFeedsCacheAddress: dataFeedsCacheAddress,
			URL:                   "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
			FeedID:                feedID,
			WriteTargetName:       writeTargetName,
		},
		BalanceReaderConfig: portypes.BalanceReaderConfig{
			BalanceReaderAddress: balanceReaderAddress,
			AddressesToRead:      addressesToRead,
		},
	}

	yaml, yamlErr := yaml.Marshal(cfg)
	if yamlErr != nil {
		return "", errors.Wrap(yamlErr, "failed to marshal config to YAML")
	}

	filePath := filepath.Join(folder, "config.yaml")
	writeErr := os.WriteFile(filePath, yaml, 0644) //nolint:gosec // G306: we want it to be readable by everyone
	if writeErr != nil {
		return "", errors.Wrap(writeErr, "failed to write config to file")
	}

	return filePath, nil
}
