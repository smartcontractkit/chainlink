package environment

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/deploy"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/trigger"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/verify"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

func workflowCmds() *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Workflow management commands",
		Long:  `Commands to manage workflows`,
	}

	workflowCmd.AddCommand(deployAndVerifyExampleWorkflowCmd())
	workflowCmd.AddCommand(deployWorkflowCmd())

	return workflowCmd
}

func deployAndVerifyExampleWorkflowCmd() *cobra.Command {
	var (
		rpcURLFlag                 string
		gatewayURLFlag             string
		chainIDFlag                uint64
		exampleWorkflowTriggerFlag string
		exampleWorkflowTimeoutFlag string
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

			return deployAndVerifyExampleWorkflow(cmd.Context(), rpcURLFlag, gatewayURLFlag, chainIDFlag, timeout, exampleWorkflowTriggerFlag)
		},
	}

	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "c", 1337, "Chain ID")
	cmd.Flags().StringVarP(&exampleWorkflowTriggerFlag, "example-workflow-trigger", "y", "web-trigger", "Trigger for example workflow to deploy (web-trigger or cron)")
	cmd.Flags().StringVarP(&exampleWorkflowTimeoutFlag, "example-workflow-timeout", "u", "5m", "Time to wait until example workflow succeeds")
	cmd.Flags().StringVarP(&gatewayURLFlag, "gateway-url", "g", "http://localhost:5002", "Gateway URL (only for web API trigger-based workflow)")

	return cmd
}

type executableWorkflowFn = func(cmdContext context.Context, rpcURL, gatewayURL, privateKey string, consumerContractAddress common.Address, workflowData *workflowData, waitTime time.Duration, startTime time.Time) error

func executeWebTriggerBasedWorkflow(cmdContext context.Context, rpcURL, gatewayURL, privateKey string, consumerContractAddress common.Address, workflowData *workflowData, waitTime time.Duration, startTime time.Time) error {
	ticker := 5 * time.Second
	for {
		select {
		case <-time.After(waitTime):
			fmt.Print(libformat.PurpleText("\n[Stage 3/3] Example workflow failed to execute successfully in %.2f seconds\n", time.Since(startTime).Seconds()))

			return fmt.Errorf("example workflow failed to execute successfully within %s", waitTime)
		case <-time.Tick(ticker):
			triggerErr := trigger.WebAPITriggerValue(gatewayURL, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "0x9A99f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE", privateKey, 5*time.Minute)
			if triggerErr == nil {
				verifyTime := 25 * time.Second
				verifyErr := verify.ProofOfReserve(rpcURL, consumerContractAddress.Hex(), workflowData.FeedID, true, verifyTime)
				if verifyErr == nil {
					if isBlockscoutRunning(cmdContext) {
						fmt.Print(libformat.PurpleText("Open http://localhost/address/0x9A9f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE?tab=internal_txns to check consumer contract's transaction history\n"))
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

func executeCronBasedWorkflow(cmdContext context.Context, rpcURL, _, privateKey string, consumerContractAddress common.Address, workflowData *workflowData, waitTime time.Duration, startTime time.Time) error {
	// we ignore return as if verification failed it will print that info
	verifyErr := verify.ProofOfReserve(rpcURL, consumerContractAddress.Hex(), workflowData.FeedID, true, waitTime)
	if verifyErr != nil {
		fmt.Print(libformat.PurpleText("\n[Stage 3/3] Example workflow failed to execute successfully in %.2f seconds\n", time.Since(startTime).Seconds()))
		return errors.Wrap(verifyErr, "failed to verify example workflow")
	}

	if isBlockscoutRunning(cmdContext) {
		fmt.Print(libformat.PurpleText("Open http://localhost/address/0x9A9f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE?tab=internal_txns to check consumer contract's transaction history\n"))
	}

	return nil
}

func deployAndVerifyExampleWorkflow(cmdContext context.Context, rpcURL, gatewayURL string, chainID uint64, timeout time.Duration, exampleWorkflowTriggerFlag string) error {
	totalStart := time.Now()
	start := time.Now()

	var executableWorkflowFunction executableWorkflowFn

	var workflowData *workflowData
	var workflowDataErr error
	if strings.EqualFold(exampleWorkflowTriggerFlag, WorkflowTriggerCron) {
		workflowData, workflowDataErr = readWorkflowData(WorkflowTriggerCron)
		executableWorkflowFunction = executeCronBasedWorkflow
	} else {
		workflowData, workflowDataErr = readWorkflowData(WorkflowTriggerWebTrigger)
		executableWorkflowFunction = executeWebTriggerBasedWorkflow
	}

	if workflowDataErr != nil {
		return errors.Wrap(workflowDataErr, "failed to read workflow data")
	}

	fmt.Print(libformat.PurpleText("[Stage 1/3] Deploying Permissionless Feeds Consumer\n\n"))
	consumerContractAddress, consumerErr := deploy.PermissionlessFeedsConsumer(rpcURL)
	if consumerErr != nil {
		return errors.Wrap(consumerErr, "failed to deploy Permissionless Feeds Consumer contract")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 1/3] Deployed Permissionless Feeds Consumer in %.2f seconds\n", time.Since(start).Seconds()))

	start = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 2/3] Registering example Proof-of-Reserve workflow\n\n"))

	deployErr := deployExampleWorkflow(chainID, *workflowData)
	if deployErr != nil {
		return errors.Wrap(deployErr, "failed to deploy example workflow")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 2/3] Registered workflow in %.2f seconds\n", time.Since(start).Seconds()))
	fmt.Print(libformat.PurpleText("[Stage 3/3] Waiting for %.2f seconds for workflow to execute successfully\n\n", timeout.Seconds()))

	var pauseWorkflow = func() {
		fmt.Print(libformat.PurpleText("\n[Stage 3/3] Example workflow executed in %.2f seconds\n", time.Since(totalStart).Seconds()))
		start = time.Now()
		fmt.Print(libformat.PurpleText("\n[CLEANUP] Pausing example workflow\n\n"))
		pauseErr := pauseExampleWorkflow(chainID)
		if pauseErr != nil {
			fmt.Printf("Failed to pause example workflow: %s\nPlease pause it manually\n", pauseErr)
		}

		fmt.Print(libformat.PurpleText("\n[CLEANUP] Paused example workflow in %.2f seconds\n\n", time.Since(start).Seconds()))
	}
	defer pauseWorkflow()

	return executableWorkflowFunction(cmdContext, rpcURL, gatewayURL, os.Getenv("PRIVATE_KEY"), *consumerContractAddress, workflowData, timeout, totalStart)
}

var creCLI = "cre_v0.2.0_darwin_arm64"
var exampleWorkflowName = "exampleworkflow"

func prepareCLIInput(chainID uint64) (*cretypes.ManageWorkflowWithCRECLIInput, error) {
	if !isCRECLIIsAvailable() {
		if downloadErr := tryToDownloadCRECLI(); downloadErr != nil {
			return nil, errors.Wrapf(downloadErr, "failed to download %s", creCLI)
		}
	}

	if os.Getenv("CRE_GITHUB_API_TOKEN") == "" {
		// set fake token to satisfy CRE CLI
		_ = os.Setenv("CRE_GITHUB_API_TOKEN", "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")
	}

	chainSelector, chainSelectorErr := chainselectors.SelectorFromChainId(chainID)
	if chainSelectorErr != nil {
		return nil, errors.Wrapf(chainSelectorErr, "failed to find chain selector for chainID %d", chainID)
	}

	CRECLIAbsPath, CRECLIAbsPathErr := creCLIAbsPath()
	if CRECLIAbsPathErr != nil {
		return nil, errors.Wrapf(CRECLIAbsPathErr, "failed to get absolute path of the %s binary", creCLI)
	}

	deployerPrivateKey := os.Getenv("PRIVATE_KEY")
	if deployerPrivateKey == "" {
		deployerPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	}

	privateKey, pkErr := crypto.HexToECDSA(deployerPrivateKey)
	if pkErr != nil {
		return nil, errors.Wrap(pkErr, "failed to parse the private key")
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	deployerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	cliSettingsFileName := "cre.yaml"
	if _, cliFileErr := os.Stat(cliSettingsFileName); os.IsNotExist(cliFileErr) {
		return nil, errors.Wrap(cliFileErr, "CRE CLI settings file not found")
	}

	cliSettingsFile, cliSettingsFilhErr := os.OpenFile(cliSettingsFileName, os.O_RDONLY, 0600)
	if cliSettingsFilhErr != nil {
		return nil, errors.Wrap(cliSettingsFilhErr, "failed to open the CRE CLI settings file")
	}

	return &cretypes.ManageWorkflowWithCRECLIInput{
		ChainSelector:            chainSelector,
		WorkflowDonID:            1,
		WorkflowOwnerAddress:     deployerAddress,
		CRECLIPrivateKey:         deployerPrivateKey,
		CRECLIAbsPath:            CRECLIAbsPath,
		CRESettingsFile:          cliSettingsFile,
		WorkflowName:             exampleWorkflowName,
		ShouldCompileNewWorkflow: false,
		CRECLIProfile:            "test",
	}, nil
}

func deployExampleWorkflow(chainID uint64, workflowData workflowData) error {
	registerWorkflowInput, registerWorkflowInputErr := prepareCLIInput(chainID)
	if registerWorkflowInputErr != nil {
		return errors.Wrap(registerWorkflowInputErr, "failed to prepare CLI input")
	}

	registerWorkflowInput.ExistingWorkflow = &cretypes.ExistingWorkflow{
		BinaryURL: workflowData.BinaryURL,
		ConfigURL: &workflowData.ConfigURL,
	}

	registerErr := creworkflow.RegisterWithCRECLI(*registerWorkflowInput)
	if registerErr != nil {
		return errors.Wrap(registerErr, "failed to register workflow")
	}

	return nil
}

func pauseExampleWorkflow(chainID uint64) error {
	pauseWorkflowInput, pauseWorkflowInputErr := prepareCLIInput(chainID)
	if pauseWorkflowInputErr != nil {
		return errors.Wrap(pauseWorkflowInputErr, "failed to prepare CLI input")
	}

	pauseErr := creworkflow.PauseWithCRECLI(*pauseWorkflowInput)
	if pauseErr != nil {
		return errors.Wrap(pauseErr, "failed to pause workflow")
	}

	return nil
}

type workflowData struct {
	BinaryURL string `json:"binary_url"`
	ConfigURL string `json:"config_url"`
	FeedID    string `json:"feed_id"`
}

func readWorkflowData(workflowTrigger string) (*workflowData, error) {
	var path string
	if strings.EqualFold(workflowTrigger, WorkflowTriggerCron) {
		path = "./examples/workflows/proof-of-reserve/cron-based/workflow_data.json"
	} else {
		path = "./examples/workflows/proof-of-reserve/web-trigger-based/workflow_data.json"
	}

	wdFileContent, wdFileErr := os.ReadFile(path)
	if wdFileErr != nil {
		return nil, errors.Wrap(wdFileErr, "failed to open workflow_data.json file")
	}

	wdData := &workflowData{}
	unmarshallErr := json.Unmarshal(wdFileContent, wdData)
	if unmarshallErr != nil {
		return nil, errors.Wrap(unmarshallErr, "failed to unmarshall workflow data")
	}

	return wdData, nil
}

func deployWorkflowCmd() *cobra.Command {
	var (
		workflowFilePathFlag        string
		configFilePathFlag          string
		containerTargetDirFlag      string
		containerNamePatternFlag    string
		workflowNameFlag            string
		workflowOwnerAddressFlag    string
		workflowRegistryAddressFlag string
		chainIDFlag                 uint64
		rpcURLFlag                  string
	)

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Compiles and uploads a workflow to the environment",
		Long:  `Compiles and uploads a workflow to the environment by copying it to workflow nodes and registering with the workflow registry`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("\n⚙️ Compiling workflow from %s\n", workflowFilePathFlag)

			compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFilePathFlag, workflowNameFlag)
			if compileErr != nil {
				return errors.Wrap(compileErr, "❌ failed to compile workflow")
			}

			fmt.Printf("\n✅ Workflow compiled and compressed successfully\n\n")

			copyErr := creworkflow.CopyWorkflowToDockerContainers(compressedWorkflowWasmPath, containerNamePatternFlag, containerTargetDirFlag)
			if copyErr != nil {
				return errors.Wrap(copyErr, "❌ failed to copy workflow to Docker container")
			}

			fmt.Printf("\n✅ Workflow copied to Docker containers\n")
			fmt.Printf("\n⚙️ Creating Seth client\n\n")

			var privateKey string
			if os.Getenv("PRIVATE_KEY") != "" {
				privateKey = os.Getenv("PRIVATE_KEY")
			} else {
				privateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
			}

			sethClient, scErr := seth.NewClientBuilder().
				WithRpcUrl(rpcURLFlag).
				WithPrivateKeys([]string{privateKey}).
				WithProtections(false, false, seth.MustMakeDuration(time.Minute)).
				Build()
			if scErr != nil {
				return errors.Wrap(scErr, "failed to create Seth client")
			}

			var configURL *string
			if configFilePathFlag != "" {
				configURL = &configFilePathFlag
			}

			fmt.Printf("\n⚙️ Deleting all workflows from the workflow registry\n\n")

			deleteErr := creworkflow.DeleteAllWithContract(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddressFlag))
			if deleteErr != nil {
				return errors.Wrapf(deleteErr, "❌ failed to delete all workflows from the registry %s", workflowRegistryAddressFlag)
			}

			fmt.Printf("\n⚙️ Registering workflow %s with the workflow registry\n\n", workflowNameFlag)

			registerErr := creworkflow.RegisterWithContract(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddressFlag), 1, workflowNameFlag, "file://"+compressedWorkflowWasmPath, configURL, nil, &containerTargetDirFlag)
			if registerErr != nil {
				return errors.Wrapf(registerErr, "❌ failed to register workflow %s", workflowNameFlag)
			}

			defer func() {
				_ = os.Remove(compressedWorkflowWasmPath)
			}()

			fmt.Printf("\n✅ Workflow registered successfully\n\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "workflow-file-path", "w", "./examples/workflows/v2/cron/main.go", "Path to the workflow file")
	cmd.Flags().StringVarP(&configFilePathFlag, "config-file-path", "c", "", "Path to the config file")
	cmd.Flags().StringVarP(&containerTargetDirFlag, "container-target-dir", "t", "/home/chainlink/workflows", "Path to the target directory in the Docker container")
	cmd.Flags().StringVarP(&containerNamePatternFlag, "container-name-pattern", "n", "workflow-node", "Pattern to match the container name")
	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "o", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0", "Workflow registry address")
	cmd.Flags().StringVarP(&workflowNameFlag, "workflow-name", "m", "exampleworkflow", "Workflow name")

	return cmd
}
