package aptos

import (
	"context"
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	aptossdk "github.com/aptos-labs/aptos-go-sdk"
	"github.com/pelletier/go-toml/v2"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-retry"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	aptosplatform "github.com/smartcontractkit/chainlink-aptos/bindings/platform"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	aptoschangeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/aptos"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	aptoschain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/aptos"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func forwarderAddress(ds datastore.DataStore, chainSelector uint64) (string, bool) {
	key := datastore.NewAddressRefKey(
		chainSelector,
		datastore.ContractType(forwarderContractType),
		forwarderContractVersion,
		forwarderQualifier,
	)
	ref, err := ds.Addresses().Get(key)
	if err != nil {
		return "", false
	}
	return ref.Address, true
}

func mustForwarderAddress(ds datastore.DataStore, chainSelector uint64) string {
	addr, ok := forwarderAddress(ds, chainSelector)
	if !ok {
		panic(fmt.Sprintf("missing Aptos forwarder address for chain selector %d", chainSelector))
	}
	return addr
}

func ensureForwardersForChains(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	chainIDs []uint64,
) (map[uint64]string, error) {
	forwardersByChainID := make(map[uint64]string, len(chainIDs))
	for _, chainID := range chainIDs {
		aptosChain, err := findAptosChainByChainID(creEnv.Blockchains, chainID)
		if err != nil {
			return nil, err
		}

		forwarderAddress, err := ensureForwarder(ctx, testLogger, creEnv, aptosChain)
		if err != nil {
			return nil, err
		}
		forwardersByChainID[chainID] = forwarderAddress
	}
	return forwardersByChainID, nil
}

func patchNodeTOML(don *cre.DonMetadata, forwardersByChainID map[uint64]string) error {
	for nodeIndex := range don.MustNodeSet().NodeSpecs {
		currentConfig := don.MustNodeSet().NodeSpecs[nodeIndex].Node.TestConfigOverrides
		if strings.TrimSpace(currentConfig) == "" {
			return fmt.Errorf("missing node config for node index %d in DON %q", nodeIndex, don.Name)
		}

		var typedConfig corechainlink.Config
		if err := toml.Unmarshal([]byte(currentConfig), &typedConfig); err != nil {
			return fmt.Errorf("failed to unmarshal config for node index %d: %w", nodeIndex, err)
		}

		for chainID, forwarderAddress := range forwardersByChainID {
			if err := setForwarderAddress(&typedConfig, strconv.FormatUint(chainID, 10), forwarderAddress); err != nil {
				return fmt.Errorf("failed to patch Aptos forwarder address for node index %d: %w", nodeIndex, err)
			}
		}

		stringifiedConfig, err := toml.Marshal(typedConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal patched config for node index %d: %w", nodeIndex, err)
		}
		don.MustNodeSet().NodeSpecs[nodeIndex].Node.TestConfigOverrides = string(stringifiedConfig)
	}

	return nil
}

func setForwarderAddress(cfg *corechainlink.Config, chainID, forwarderAddress string) error {
	for i := range cfg.Aptos {
		raw := map[string]any(cfg.Aptos[i])
		if fmt.Sprint(raw["ChainID"]) != chainID {
			continue
		}

		workflow := make(map[string]any)
		switch existing := raw["Workflow"].(type) {
		case map[string]any:
			for k, v := range existing {
				workflow[k] = v
			}
		case corechainlink.RawConfig:
			for k, v := range existing {
				workflow[k] = v
			}
		case nil:
		default:
			return fmt.Errorf("unexpected Aptos workflow config type %T", existing)
		}
		workflow["ForwarderAddress"] = forwarderAddress
		raw["Workflow"] = workflow
		cfg.Aptos[i] = corechainlink.RawConfig(raw)
		return nil
	}

	return fmt.Errorf("Aptos chain %s not found in node config", chainID)
}

// ensureForwarder makes sure a forwarder exists for the Aptos chain selector and
// returns its address. In local Docker environments it will deploy the forwarder
// once and cache the resulting address in the CRE datastore; in non-Docker
// environments it only reuses an address that has already been injected.
func ensureForwarder(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	chain *aptoschain.Blockchain,
) (string, error) {
	if addr, ok := forwarderAddress(creEnv.CldfEnvironment.DataStore, chain.ChainSelector()); ok {
		return addr, nil
	}
	if !creEnv.Provider.IsDocker() {
		return "", fmt.Errorf("missing Aptos forwarder address for chain selector %d", chain.ChainSelector())
	}

	nodeURL, err := chain.NodeURL()
	if err != nil {
		return "", fmt.Errorf("invalid Aptos node URL for chain selector %d: %w", chain.ChainSelector(), err)
	}
	client, err := chain.NodeClient()
	if err != nil {
		return "", fmt.Errorf("failed to create Aptos client for chain selector %d (%s): %w", chain.ChainSelector(), nodeURL, err)
	}
	deployerAccount, err := chain.LocalDeployerAccount()
	if err != nil {
		return "", fmt.Errorf("failed to create Aptos deployer signer: %w", err)
	}
	deploymentChain, err := chain.LocalDeploymentChain()
	if err != nil {
		return "", fmt.Errorf("failed to build Aptos deployment chain for chain selector %d: %w", chain.ChainSelector(), err)
	}

	owner := deployerAccount.AccountAddress()
	if _, accountErr := client.Account(owner); accountErr != nil {
		if fundErr := chain.Fund(ctx, owner.StringLong(), 100_000_000); fundErr != nil {
			testLogger.Warn().
				Uint64("chainSelector", chain.ChainSelector()).
				Str("nodeURL", nodeURL).
				Err(fundErr).
				Msg("Aptos deployer account not confirmed visible yet; proceeding with deploy retries")
		}
	}

	var deployedAddress string
	var pendingTxHash string
	var lastDeployErr error
	if retryErr := retry.Do(ctx, retry.WithMaxDuration(3*time.Minute, retry.NewFibonacci(500*time.Millisecond)), func(ctx context.Context) error {
		deploymentResp, deployErr := aptoschangeset.DeployPlatform(deploymentChain, owner, nil)
		if deployErr != nil {
			lastDeployErr = deployErr
			if fundErr := chain.Fund(ctx, owner.StringLong(), 1_000_000_000_000); fundErr != nil {
				testLogger.Warn().
					Uint64("chainSelector", chain.ChainSelector()).
					Err(fundErr).
					Msg("failed to re-fund Aptos deployer account during deploy retry")
			}
			return retry.RetryableError(fmt.Errorf("deploy-to-object failed: %w", deployErr))
		}
		if deploymentResp == nil {
			lastDeployErr = pkgerrors.New("nil deployment response")
			return retry.RetryableError(pkgerrors.New("DeployPlatform returned nil response"))
		}
		deployedAddress = deploymentResp.Address.StringLong()
		pendingTxHash = deploymentResp.Tx
		return nil
	}); retryErr != nil {
		if lastDeployErr != nil {
			return "", fmt.Errorf("failed to deploy Aptos platform forwarder for chain selector %d after retries: %w", chain.ChainSelector(), stderrors.Join(lastDeployErr, retryErr))
		}
		return "", fmt.Errorf("failed to deploy Aptos platform forwarder for chain selector %d after retries: %w", chain.ChainSelector(), retryErr)
	}

	addr, err := normalizeForwarderAddress(deployedAddress)
	if err != nil {
		return "", fmt.Errorf("invalid Aptos forwarder address parsed from deployment output for chain selector %d: %w", chain.ChainSelector(), err)
	}

	if err := addForwarderToDataStore(creEnv, chain.ChainSelector(), addr); err != nil {
		return "", err
	}

	testLogger.Info().
		Uint64("chainSelector", chain.ChainSelector()).
		Str("nodeURL", nodeURL).
		Str("txHash", pendingTxHash).
		Str("forwarderAddress", addr).
		Msg("Aptos platform forwarder deployed")

	return addr, nil
}

// addForwarderToDataStore seals a new datastore snapshot with the Aptos
// forwarder address so later setup phases can reuse it without redeploying.
func addForwarderToDataStore(creEnv *cre.Environment, chainSelector uint64, address string) error {
	memoryDatastore, err := crecontracts.NewDataStoreFromExisting(creEnv.CldfEnvironment.DataStore)
	if err != nil {
		return fmt.Errorf("failed to create memory datastore: %w", err)
	}

	err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
		Address:       address,
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(forwarderContractType),
		Version:       forwarderContractVersion,
		Qualifier:     forwarderQualifier,
	})
	if err != nil && !stderrors.Is(err, datastore.ErrAddressRefExists) {
		return fmt.Errorf("failed to add Aptos forwarder address to datastore: %w", err)
	}

	creEnv.CldfEnvironment.DataStore = memoryDatastore.Seal()
	return nil
}

// configureForwarders writes the final DON membership and signer set to each
// Aptos forwarder after the DON has started and contract DON IDs are known.
func configureForwarders(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	creEnv *cre.Environment,
	chainIDs []uint64,
) error {
	workers, err := don.Workers()
	if err != nil {
		return fmt.Errorf("failed to get worker nodes for DON %q: %w", don.Name, err)
	}
	f := (len(workers) - 1) / 3
	if f <= 0 {
		return fmt.Errorf("invalid Aptos DON %q fault tolerance F=%d (workers=%d)", don.Name, f, len(workers))
	}
	if f > 255 {
		return fmt.Errorf("aptos DON %q fault tolerance F=%d exceeds u8", don.Name, f)
	}
	forwarderF := uint8(f)

	donIDUint32, err := aptosDonIDUint32(don.ID)
	if err != nil {
		return fmt.Errorf("invalid DON id for Aptos forwarder config: %w", err)
	}

	oracles, err := donOraclePublicKeys(ctx, don)
	if err != nil {
		return err
	}

	for _, chainID := range chainIDs {
		aptosChain, err := findAptosChainByChainID(creEnv.Blockchains, chainID)
		if err != nil {
			return err
		}

		nodeURL, err := aptosChain.NodeURL()
		if err != nil {
			return fmt.Errorf("invalid Aptos node URL for chain selector %d: %w", aptosChain.ChainSelector(), err)
		}
		client, err := aptosChain.NodeClient()
		if err != nil {
			return fmt.Errorf("failed to create Aptos client for chain selector %d (%s): %w", aptosChain.ChainSelector(), nodeURL, err)
		}
		deployerAccount, err := aptosChain.LocalDeployerAccount()
		if err != nil {
			return fmt.Errorf("failed to create Aptos deployer signer for forwarder config: %w", err)
		}
		deployerAddress := deployerAccount.AccountAddress()

		if _, accountErr := client.Account(deployerAddress); accountErr != nil {
			if fundErr := aptosChain.Fund(ctx, deployerAddress.StringLong(), 100_000_000); fundErr != nil {
				testLogger.Warn().
					Uint64("chainSelector", aptosChain.ChainSelector()).
					Str("nodeURL", nodeURL).
					Err(fundErr).
					Msg("Aptos deployer account not confirmed visible yet; proceeding with forwarder set_config retries")
			}
		}

		forwarderHex := mustForwarderAddress(creEnv.CldfEnvironment.DataStore, aptosChain.ChainSelector())
		var forwarderAddr aptossdk.AccountAddress
		if err := forwarderAddr.ParseStringRelaxed(forwarderHex); err != nil {
			return fmt.Errorf("invalid Aptos forwarder address for chain selector %d: %w", aptosChain.ChainSelector(), err)
		}
		forwarderContract := aptosplatform.Bind(forwarderAddr, client).Forwarder()

		var pendingTxHash string
		var lastSetConfigErr error
		if err := retry.Do(ctx, retry.WithMaxDuration(2*time.Minute, retry.NewFibonacci(500*time.Millisecond)), func(ctx context.Context) error {
			pendingTx, err := forwarderContract.SetConfig(&bind.TransactOpts{Signer: deployerAccount}, donIDUint32, forwarderConfigVersion, forwarderF, oracles)
			if err != nil {
				lastSetConfigErr = err
				if fundErr := aptosChain.Fund(ctx, deployerAddress.StringLong(), 1_000_000_000_000); fundErr != nil {
					testLogger.Warn().
						Uint64("chainSelector", aptosChain.ChainSelector()).
						Err(fundErr).
						Msg("failed to fund Aptos deployer account during set_config retry")
				}
				return retry.RetryableError(fmt.Errorf("set_config transaction submit failed: %w", err))
			}
			pendingTxHash = pendingTx.Hash
			receipt, err := client.WaitForTransaction(pendingTxHash)
			if err != nil {
				lastSetConfigErr = err
				return retry.RetryableError(fmt.Errorf("waiting for set_config transaction failed: %w", err))
			}
			if !receipt.Success {
				lastSetConfigErr = fmt.Errorf("vm status: %s", receipt.VmStatus)
				return retry.RetryableError(fmt.Errorf("set_config transaction failed: %s", receipt.VmStatus))
			}
			return nil
		}); err != nil {
			if lastSetConfigErr != nil {
				return fmt.Errorf("failed to configure Aptos forwarder %s for DON %q on chain selector %d: %w", forwarderHex, don.Name, aptosChain.ChainSelector(), stderrors.Join(lastSetConfigErr, err))
			}
			return fmt.Errorf("failed to configure Aptos forwarder %s for DON %q on chain selector %d: %w", forwarderHex, don.Name, aptosChain.ChainSelector(), err)
		}

		testLogger.Info().
			Str("donName", don.Name).
			Uint64("donID", don.ID).
			Uint64("chainSelector", aptosChain.ChainSelector()).
			Str("txHash", pendingTxHash).
			Str("forwarderAddress", forwarderHex).
			Msg("configured Aptos forwarder set_config")
	}

	return nil
}
