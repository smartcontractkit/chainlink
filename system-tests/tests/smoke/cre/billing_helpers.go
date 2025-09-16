package cre

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	libcre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func loadBillingStackCache(relativePathToRepoRoot string) (*config.BillingConfig, error) {
	c := &config.BillingConfig{}
	if loadErr := c.Load(config.MustBillingStateFileAbsPath(relativePathToRepoRoot)); loadErr != nil {
		return nil, errors.Wrap(loadErr, "failed to load billing stack cache")
	}

	return c, nil
}

func startBillingStackIfIsNotRunning(relativePathToRepoRoot, environmentDir, streamsAPIURL string, testEnv *TestEnvironment) error {
	if !config.BillingStateFileExists(relativePathToRepoRoot) {
		// set env vars for billing config
		cache, err := loadWorkflowRegistryCache(relativePathToRepoRoot)
		if err != nil {
			return errors.Wrap(err, "failed to load workflow registry cache")
		}

		if len(testEnv.WrappedBlockchainOutputs) == 0 {
			return errors.New("no blockchain outputs found in the test environment")
		}

		replaceHost := func(url string) string {
			return strings.Replace(url, "127.0.0.1", "host.docker.internal", 1)
		}

		for _, ref := range testEnv.EnvArtifact.AddressRefs {
			switch ref.Type {
			case "WorkflowRegistry":
				if cache.ChainSelector == ref.ChainSelector {
					os.Setenv("MAINNET_WORKFLOW_REGISTRY_CONTRACT_ADDRESS", ref.Address)
				}
			case "CapabilitiesRegistry":
				if cache.ChainSelector == ref.ChainSelector {
					os.Setenv("MAINNET_CAPABILITIES_REGISTRY_CONTRACT_ADDRESS", ref.Address)
				}
			default:
				continue
			}
		}

		os.Setenv("MAINNET_WORKFLOW_REGISTRY_CHAIN_SELECTOR", strconv.FormatUint(cache.ChainSelector, 10))
		os.Setenv("MAINNET_CAPABILITIES_REGISTRY_CHAIN_SELECTOR", strconv.FormatUint(cache.ChainSelector, 10))
		os.Setenv("STREAMS_API_URL", replaceHost(streamsAPIURL))
		os.Setenv("STREAMS_API_KEY", "cannot be empty")
		os.Setenv("STREAMS_API_SECRET", "cannot be empty")
		os.Setenv("TEST_OWNERS", strings.Join(cache.WorkflowOwnersStrings(), ","))

		// Select the appropriate chain for billing service from available chains in the environment.
		// otherwise, if RPCURL is defined, billing service can be used standalone
		if len(testEnv.WrappedBlockchainOutputs) != 0 {
			var selectedChain *blockchain.Output

			for _, chain := range testEnv.WrappedBlockchainOutputs {
				if chain.ChainSelector == cache.ChainSelector {
					selectedChain = chain.BlockchainOutput
				}
			}

			if selectedChain == nil || len(selectedChain.Nodes) == 0 {
				return errors.Wrap(err, fmt.Sprintf("configured chain selector does not exist in the current topology: %d", cache.ChainSelector))
			}

			rpcURL := replaceHost(selectedChain.Nodes[0].ExternalHTTPUrl)

			os.Setenv("MAINNET_WORKFLOW_REGISTRY_RPC_URL", rpcURL)
			os.Setenv("MAINNET_CAPABILITIES_REGISTRY_RPC_URL", rpcURL)
		}

		framework.L.Info().Str("state file", config.MustBillingStateFileAbsPath(relativePathToRepoRoot)).Msg("Billing state file was not found. Starting Billing...")
		cmd := exec.Command("go", "run", ".", "env", "billing", "start")
		cmd.Dir = environmentDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			return errors.Wrap(cmdErr, "failed to start Billing Platform Service")
		}
	}
	framework.L.Info().Msg("Billing Platform Service is running.")
	return nil
}

func loadWorkflowRegistryCache(relativePathToRepoRoot string) (*libcre.WorkflowRegistryOutput, error) {
	previousCTFconfigs := os.Getenv("CTF_CONFIGS")
	defer func() {
		setErr := os.Setenv("CTF_CONFIGS", previousCTFconfigs)
		if setErr != nil {
			framework.L.Warn().Err(setErr).Msg("failed to restore previous CTF_CONFIGS env var")
		}
	}()

	setErr := os.Setenv("CTF_CONFIGS", config.MustWorkflowRegistryStateFileAbsPath(relativePathToRepoRoot))
	if setErr != nil {
		return nil, errors.Wrap(setErr, "failed to set CTF_CONFIGS env var")
	}

	return framework.Load[libcre.WorkflowRegistryOutput](nil)
}
