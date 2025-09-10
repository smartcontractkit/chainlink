package environment

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/http/common"
	"github.com/fbsobreira/gotron-sdk/pkg/http/soliditynode"
	"google.golang.org/grpc/credentials/insecure"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	tronprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/tron/provider/rpcclient"
	tron_sdk "github.com/smartcontractkit/chainlink-tron/relayer/sdk"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	deployment_devenv "github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

// BuildFromSavedState rebuilds the CLDF environment and per‑chain clients from
// artifacts produced by a previous local CRE run.
// Inputs:
//   - cachedInput: outputs from starting the environment via CTFv2 configs
//     (node sets, Job Distributor, blockchain nodes).
//   - envArtifact: CLDF deployment output including JD config and DON
//     topology/metadata.
//
// Artifact paths are recorded in `artifact_paths.json` in the environment
// directory (typically `core/scripts/cre/environment`).
// Returns the reconstructed CLDF environment, wrapped blockchain outputs, and an error.
func BuildFromSavedState(ctx context.Context, cldLogger logger.Logger, cachedInput *envconfig.Config, envArtifact *EnvArtifact) (*cre.FullCLDEnvironmentOutput, []*cre.WrappedBlockchainOutput, error) {
	if cachedInput == nil {
		return nil, nil, errors.New("cached input cannot be nil")
	}

	if envArtifact == nil {
		return nil, nil, errors.New("environment artifact cannot be nil")
	}

	if pkErr := SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return nil, nil, pkErr
	}
	// just in case
	if pkErr := SetDefaultSolanaPrivateKeyIfEmpty(defaultSolanaPrivateKey); pkErr != nil {
		return nil, nil, pkErr
	}

	wrappedBlockchainOutputs := make([]*cre.WrappedBlockchainOutput, 0)

	for _, bc := range cachedInput.Blockchains {
		chainID, chainIDErr := strconv.ParseUint(bc.ChainID, 10, 64)
		if chainIDErr != nil {
			return nil, nil, errors.Wrapf(chainIDErr, "failed to parse chain id %s", bc.ChainID)
		}

		chainSelector, chainSelectorErr := chainselectors.SelectorFromChainId(chainID)
		if chainSelectorErr != nil {
			return nil, nil, errors.Wrapf(chainSelectorErr, "failed to get chain selector for chain id %d", chainID)
		}

		// Handle Tron chains differently - they don't use Seth clients
		if bc.Type == blockchain.FamilyTron {
			// For Tron chains, reconstruct the TronChain from cached configuration
			tronChain, tronErr := reconstructTronChain(ctx, bc, chainSelector, chainID)
			if tronErr != nil {
				return nil, nil, errors.Wrapf(tronErr, "failed to reconstruct Tron chain for chain ID %d", chainID)
			}

			wrappedBlockchainOutputs = append(wrappedBlockchainOutputs, &cre.WrappedBlockchainOutput{
				BlockchainOutput:   bc.Out,
				SethClient:         nil, // Tron chains don't use Seth clients
				ChainSelector:      chainSelector,
				ChainID:            chainID,
				DeployerPrivateKey: blockchain.TRONAccounts.PrivateKeys[0],
				TronChain:          tronChain,
			})
			continue
		}

		// Handle EVM chains with Seth clients
		sethClient, sethErr := seth.NewClientBuilder().
			WithRpcUrl(bc.Out.Nodes[0].ExternalWSUrl).
			WithPrivateKeys([]string{os.Getenv("PRIVATE_KEY")}).
			// do not check if there's a pending nonce nor check node's health
			WithProtections(false, false, seth.MustMakeDuration(time.Second)).
			Build()
		if sethErr != nil {
			return nil, nil, errors.Wrap(sethErr, "failed to create seth client")
		}

		wrappedBlockchainOutputs = append(wrappedBlockchainOutputs, &cre.WrappedBlockchainOutput{
			BlockchainOutput:   bc.Out,
			SethClient:         sethClient,
			ChainSelector:      chainSelector,
			ChainID:            chainID,
			DeployerPrivateKey: sethClient.Cfg.Network.PrivateKeys[0],
		})
		if bc.Type == blockchain.FamilySolana {
			initErr := initSolanaInput(&bc)
			if initErr != nil {
				return nil, nil, errors.Wrap(initErr, "failed to init solana")
			}
			w, err := wrapSolana(&bc, bc.Out)
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to wrap solana")
			}
			wrappedBlockchainOutputs = append(wrappedBlockchainOutputs, w)
			continue
		}
		w, err := wrapEVM(bc.Out)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to wrap evm")
		}
		wrappedBlockchainOutputs = append(wrappedBlockchainOutputs, w)
	}

	addressBook := cldf.NewMemoryAddressBookFromMap(envArtifact.AddressBook)
	datastore := datastore.NewMemoryDataStore()
	for _, addrRef := range envArtifact.AddressRefs {
		addErr := datastore.AddressRefStore.Add(addrRef)
		if addErr != nil {
			return nil, nil, errors.Wrapf(addErr, "failed to add address ref to datastore %v", addrRef)
		}
	}

	allNodeInfo := make([]deployment_devenv.NodeInfo, 0)
	allNodeIDs := make([]string, 0)

	for idx, don := range envArtifact.DONs {
		_, ok := envArtifact.Nodes[don.DonName]
		if !ok {
			return nil, nil, errors.Errorf("no nodes found for don %s", don.DonName)
		}

		for id := range envArtifact.Nodes[don.DonName].Nodes {
			allNodeIDs = append(allNodeIDs, id)
		}

		bootstrapNodes, err := crenode.FindManyWithLabel(envArtifact.Topology.DonsWithMetadata[idx].NodesMetadata, &cre.Label{Key: crenode.NodeTypeKey, Value: cre.BootstrapNode}, crenode.EqualLabels)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to find bootstrap nodes")
		}

		nodeInfo, err := crenode.GetNodeInfo(cachedInput.NodeSets[idx].Out, cachedInput.NodeSets[idx].Name, don.DonID, len(bootstrapNodes))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get node info for don %s", don.DonName)
		}

		offChain, offChainErr := deployment_devenv.NewJDClient(ctx, deployment_devenv.JDConfig{
			WSRPC:    envArtifact.JdConfig.ExternalGRPCUrl,
			GRPC:     envArtifact.JdConfig.ExternalGRPCUrl,
			Creds:    insecure.NewCredentials(),
			NodeInfo: nodeInfo,
		})
		if offChainErr != nil {
			return nil, nil, errors.Wrapf(offChainErr, "failed to create offchain client for don %s", don.DonName)
		}

		jd, ok := offChain.(*deployment_devenv.JobDistributor)
		if !ok {
			return nil, nil, errors.Errorf("offchain client is not a JobDistributor for don %s", don.DonName)
		}

		registeredDon, donErr := deployment_devenv.NewRegisteredDON(ctx, nodeInfo, *jd)
		if donErr != nil {
			return nil, nil, errors.Wrapf(donErr, "failed to create DON for don %s", don.DonName)
		}

		envArtifact.Topology.DonsWithMetadata[idx].DON = registeredDon
		allNodeInfo = append(allNodeInfo, nodeInfo...)
	}

	offChain, offChainErr := deployment_devenv.NewJDClient(ctx, deployment_devenv.JDConfig{
		WSRPC:    envArtifact.JdConfig.ExternalGRPCUrl,
		GRPC:     envArtifact.JdConfig.ExternalGRPCUrl,
		Creds:    insecure.NewCredentials(),
		NodeInfo: allNodeInfo,
	})
	if offChainErr != nil {
		return nil, nil, errors.Wrapf(offChainErr, "failed to create offchain client")
	}

	// Create a map of chain.BlockChain objects for both EVM and Tron chains
	// This matches the same pattern used in the original environment setup
	blockChainsMap := make(map[uint64]chain.BlockChain)
	chainConfigs := make([]deployment_devenv.ChainConfig, 0, len(wrappedBlockchainOutputs))

	for _, output := range wrappedBlockchainOutputs {
		if output.TronChain != nil {
			// Add Tron chains directly to the map
			blockChainsMap[output.ChainSelector] = *output.TronChain
		}
		cfg, cfgErr := cre.ChainConfigFromWrapped(output)
		if cfgErr != nil {
			return nil, nil, errors.Wrapf(cfgErr, "failed to build chain config from write for blockchain %s", output.BlockchainOutput.Family)
		}
		chainConfigs = append(chainConfigs, cfg)
	}

	// Create BlockChains object using the same pattern as original environment setup
	blockChains := chain.NewBlockChains(blockChainsMap)

	cldEnv := cldf.NewEnvironment(
		"cre",
		cldLogger,
		addressBook,
		datastore.Seal(),
		allNodeIDs,
		offChain,
		func() context.Context {
			return ctx
		},
		cldf.XXXGenerateTestOCRSecrets(),
		blockChains,
	)

	return &cre.FullCLDEnvironmentOutput{
		Environment: cldEnv,
		DonTopology: &envArtifact.Topology,
	}, wrappedBlockchainOutputs, nil
}

// reconstructTronChain reconnects to an existing Tron node using cached configuration
// without spinning up new containers
func reconstructTronChain(ctx context.Context, bc blockchain.Input, chainSelector, chainID uint64) (*tron.Chain, error) {
	if bc.Out == nil || len(bc.Out.Nodes) == 0 {
		return nil, fmt.Errorf("blockchain output is nil or has no nodes")
	}

	// Extract URLs from the cached blockchain output
	// For Tron, ExternalHTTPUrl should contain the base URL
	baseURL := bc.Out.Nodes[0].ExternalHTTPUrl
	if baseURL == "" {
		return nil, fmt.Errorf("no external HTTP URL found for Tron chain")
	}

	// Construct the Tron node URLs (wallet and walletsolidity endpoints)
	fullNodeURL := strings.Replace(baseURL, "/jsonrpc", "/wallet", 1)
	solidityNodeURL := strings.Replace(baseURL, "/jsonrpc", "/walletsolidity", 1)

	// Parse URLs for client connections
	fullNodeUrlObj, err := url.Parse(fullNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse full node URL %s: %w", fullNodeURL, err)
	}
	solidityNodeUrlObj, err := url.Parse(solidityNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse solidity node URL %s: %w", solidityNodeURL, err)
	}

	// Create combined client to interact with both full node and solidity node
	combinedClient, err := tron_sdk.CreateCombinedClient(fullNodeUrlObj, solidityNodeUrlObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create combined client: %w", err)
	}

	// Create signer generator for the default Tron account
	signerGen, err := tronprovider.SignerGenCTFDefault()
	if err != nil {
		return nil, fmt.Errorf("failed to create signer generator: %w", err)
	}

	// Get deployer address from the signer generator
	deployerAddr, err := signerGen.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get deployer address: %w", err)
	}

	// Create RPC client wrapper that uses the signer generator's signing function
	client := rpcclient.New(combinedClient, signerGen.Sign)

	// Construct the Tron chain instance with all the helper methods
	tronChain := &tron.Chain{
		ChainMetadata: tron.ChainMetadata{
			Selector: chainSelector,
		},
		Client:   combinedClient,
		SignHash: signerGen.Sign,
		Address:  deployerAddr,
		URL:      fullNodeURL,
		// Helper for sending and confirming transactions
		SendAndConfirm: func(ctx context.Context, tx *common.Transaction, opts *tron.ConfirmRetryOptions) (*soliditynode.TransactionInfo, error) {
			options := tron.DefaultConfirmRetryOptions()
			if opts != nil {
				options = opts
			}
			return client.SendAndConfirmTx(ctx, tx, options)
		},
		// Helper for deploying a contract and waiting for confirmation
		DeployContractAndConfirm: func(
			ctx context.Context, contractName string, abi string, bytecode string, params []interface{}, opts *tron.DeployOptions,
		) (address.Address, *soliditynode.TransactionInfo, error) {
			options := tron.DefaultDeployOptions()
			if opts != nil {
				options = opts
			}

			deployResponse, err := combinedClient.DeployContract(
				deployerAddr, contractName, abi, bytecode, options.OeLimit, options.CurPercent, options.FeeLimit, params,
			)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create deploy contract transaction: %w", err)
			}

			txInfo, err := client.SendAndConfirmTx(ctx, &deployResponse.Transaction, options.ConfirmRetryOptions)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to confirm deploy contract transaction: %w", err)
			}

			contractAddr, err := address.StringToAddress(txInfo.ContractAddress)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse contract address: %w", err)
			}

			if err := client.CheckContractDeployed(contractAddr); err != nil {
				return nil, nil, fmt.Errorf("contract deployment check failed: %w", err)
			}

			return contractAddr, txInfo, nil
		},
		// Helper for triggering a contract method and waiting for confirmation
		TriggerContractAndConfirm: func(
			ctx context.Context, contractAddr address.Address, functionName string, params []interface{}, opts *tron.TriggerOptions,
		) (*soliditynode.TransactionInfo, error) {
			options := tron.DefaultTriggerOptions()
			if opts != nil {
				options = opts
			}

			if err := client.CheckContractDeployed(contractAddr); err != nil {
				return nil, fmt.Errorf("contract deployment check failed: %w", err)
			}

			contractResponse, err := combinedClient.TriggerSmartContract(
				deployerAddr, contractAddr, functionName, params, options.FeeLimit, options.TAmount,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create trigger contract transaction: %w", err)
			}

			return client.SendAndConfirmTx(ctx, contractResponse.Transaction, options.ConfirmRetryOptions)
		},
	}

	return tronChain, nil
}

func SetDefaultPrivateKeyIfEmpty(defaultPrivateKey string) error {
	if os.Getenv("PRIVATE_KEY") == "" {
		setErr := os.Setenv("PRIVATE_KEY", defaultPrivateKey)
		if setErr != nil {
			return fmt.Errorf("failed to set PRIVATE_KEY environment variable: %w", setErr)
		}
		framework.L.Info().Msgf("Set PRIVATE_KEY environment variable to default value: %s", os.Getenv("PRIVATE_KEY"))
	}

	return nil
}

func SetDefaultSolanaPrivateKeyIfEmpty(key solana.PrivateKey) error {
	if os.Getenv("SOLANA_PRIVATE_KEY") == "" {
		setErr := os.Setenv("SOLANA_PRIVATE_KEY", key.String())
		if setErr != nil {
			return fmt.Errorf("failed to set SOLANA_PRIVATE_KEY environment variable: %w", setErr)
		}
		framework.L.Info().Msgf("Set SOLANA_PRIVATE_KEY environment variable to default value: %s", os.Getenv("PRIVATE_KEY"))
	}

	return nil
}
