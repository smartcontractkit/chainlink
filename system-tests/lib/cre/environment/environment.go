package environment

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	aptos "github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/api"
	aptoscrypto "github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-retry"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	aptosplatform "github.com/smartcontractkit/chainlink-aptos/bindings/platform"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	donconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	gateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/sharding"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/worker"
)

type SetupOutput struct {
	WorkflowRegistryConfigurationOutput *cre.WorkflowRegistryOutput
	CreEnvironment                      *cre.Environment
	Dons                                *cre.Dons
	NodeOutput                          []*cre.NodeSetOutput
	S3ProviderOutput                    *s3provider.Output
	GatewayConnectors                   *cre.GatewayConnectors
}

type SetupInput struct {
	NodeSets               []*cre.NodeSet
	BlockchainsInput       []*blockchain.Input
	JdInput                *jd.Input
	Provider               infra.Provider
	ContractVersions       map[cre.ContractType]*semver.Version
	WithV2Registries       bool
	OCR3Config             *keystone_changeset.OracleConfig
	DONTimeConfig          *keystone_changeset.OracleConfig
	VaultOCR3Config        *keystone_changeset.OracleConfig
	S3ProviderInput        *s3provider.Input
	CapabilityConfigs      cre.CapabilityConfigs
	CopyCapabilityBinaries bool // if true, copy capability binaries to the containers (if false, we assume that the plugins image already has them)
	Capabilities           []cre.InstallableCapability
	Features               cre.Features
	GatewayWhitelistConfig gateway.WhitelistConfig
	BlockchainDeployers    map[blockchain.ChainFamily]blockchains.Deployer

	// allow to pass custom transformers for extensibility
	ConfigFactoryFunctions               []cre.NodeConfigTransformerFn
	JobSpecFactoryFunctions              []cre.JobSpecFn
	CapabilitiesContractFactoryFunctions []cre.CapabilityRegistryConfigFn

	StageGen *stagegen.StageGen
	// Optional map of Aptos chain selector -> forwarder address to inject into Aptos node config and write-aptos job config.
	AptosForwarderAddresses map[uint64]string
}

func (s *SetupInput) Validate() error {
	if s == nil {
		return pkgerrors.New("input is nil")
	}

	if len(s.NodeSets) == 0 {
		return pkgerrors.New("at least one nodeSet is required")
	}

	if len(s.BlockchainsInput) == 0 {
		return pkgerrors.New("at least one blockchain is required")
	}

	if s.JdInput == nil {
		return pkgerrors.New("jd input is nil")
	}

	return nil
}

const (
	aptosForwarderAddressEnvVar   = "CRE_APTOS_FORWARDER_ADDRESS"
	aptosForwarderAddressesEnvVar = "CRE_APTOS_FORWARDER_ADDRESSES"
	aptosContractsPathEnvVar      = "CRE_APTOS_CONTRACTS_PATH"
	aptosAddressHexLen            = 64
	aptosForwarderConfigVersion   = 1
)

func normalizeAptosForwarderAddress(raw string) (string, error) {
	addr := strings.TrimSpace(raw)
	if addr == "" {
		return "", errors.New("empty address")
	}
	addr = strings.TrimPrefix(strings.TrimPrefix(addr, "0x"), "0X")
	if len(addr) != aptosAddressHexLen {
		return "", fmt.Errorf("expected %d hex chars, got %d", aptosAddressHexLen, len(addr))
	}
	for _, ch := range addr {
		isHex := (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
		if !isHex {
			return "", fmt.Errorf("address contains non-hex char %q", ch)
		}
	}
	return "0x" + addr, nil
}

func resolveAptosForwarderAddresses(
	testLogger zerolog.Logger,
	deployedBlockchains []blockchains.Blockchain,
	configured map[uint64]string,
) (map[uint64]string, error) {
	aptosSelectors := make([]uint64, 0)
	aptosSelectorSet := make(map[uint64]struct{})
	for _, bc := range deployedBlockchains {
		if bc.IsFamily(chainselectors.FamilyAptos) {
			aptosSelectors = append(aptosSelectors, bc.ChainSelector())
			aptosSelectorSet[bc.ChainSelector()] = struct{}{}
		}
	}
	if len(aptosSelectors) == 0 {
		return nil, nil
	}

	out := make(map[uint64]string)
	for selector, addrRaw := range configured {
		if _, ok := aptosSelectorSet[selector]; !ok {
			continue
		}
		addr, err := normalizeAptosForwarderAddress(addrRaw)
		if err != nil {
			return nil, fmt.Errorf("invalid configured Aptos forwarder address for selector %d: %w", selector, err)
		}
		out[selector] = addr
	}

	// Optional: CRE_APTOS_FORWARDER_ADDRESSES='{"4457093679053095497":"0x..."}'
	if rawMap := strings.TrimSpace(os.Getenv(aptosForwarderAddressesEnvVar)); rawMap != "" {
		var decoded map[string]string
		if err := json.Unmarshal([]byte(rawMap), &decoded); err != nil {
			return nil, fmt.Errorf("invalid %s JSON: %w", aptosForwarderAddressesEnvVar, err)
		}
		for selectorRaw, addrRaw := range decoded {
			selector, err := strconv.ParseUint(selectorRaw, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid chain selector key %q in %s: %w", selectorRaw, aptosForwarderAddressesEnvVar, err)
			}
			addr, err := normalizeAptosForwarderAddress(addrRaw)
			if err != nil {
				return nil, fmt.Errorf("invalid forwarder address for selector %d in %s: %w", selector, aptosForwarderAddressesEnvVar, err)
			}
			out[selector] = addr
		}
	}

	// Optional: CRE_APTOS_FORWARDER_ADDRESS='0x...'
	if rawSingle := strings.TrimSpace(os.Getenv(aptosForwarderAddressEnvVar)); rawSingle != "" {
		addr, err := normalizeAptosForwarderAddress(rawSingle)
		if err != nil {
			return nil, fmt.Errorf("invalid %s: %w", aptosForwarderAddressEnvVar, err)
		}
		for _, selector := range aptosSelectors {
			if _, exists := out[selector]; !exists {
				out[selector] = addr
			}
		}
	}

	if len(out) == 0 {
		return nil, nil
	}
	testLogger.Info().Interface("aptosForwarderAddresses", out).Msg("resolved Aptos forwarder addresses for node/job config")
	return out, nil
}

func configureAptosContractsPathFromEnv(testLogger zerolog.Logger, chainInputs []*blockchain.Input) error {
	contractsPath := strings.TrimSpace(os.Getenv(aptosContractsPathEnvVar))
	if contractsPath == "" {
		return nil
	}
	stat, err := os.Stat(contractsPath)
	if err != nil {
		return fmt.Errorf("%s path is not accessible (%s): %w", aptosContractsPathEnvVar, contractsPath, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s must point to a directory, got %s", aptosContractsPathEnvVar, contractsPath)
	}

	updated := 0
	for _, in := range chainInputs {
		if in.Type != blockchain.TypeAptos {
			continue
		}
		if strings.TrimSpace(in.ContractsDir) == "" {
			in.ContractsDir = contractsPath
			updated++
		}
	}
	if updated > 0 {
		testLogger.Info().Str("contractsDir", contractsPath).Int("aptosChainsUpdated", updated).Msg("configured Aptos contracts dir from env")
	}
	return nil
}

func deployMissingAptosForwarders(
	ctx context.Context,
	testLogger zerolog.Logger,
	provider infra.Provider,
	deployedBlockchains []blockchains.Blockchain,
	current map[uint64]string,
) (map[uint64]string, error) {
	missing := make([]blockchains.Blockchain, 0)
	for _, bc := range deployedBlockchains {
		if !bc.IsFamily(chainselectors.FamilyAptos) {
			continue
		}
		addr := ""
		if current != nil {
			addr = current[bc.ChainSelector()]
		}
		if addr != "" {
			continue
		}
		missing = append(missing, bc)
	}
	if len(missing) == 0 {
		return current, nil
	}
	if !provider.IsDocker() {
		missingSelectors := make([]uint64, 0, len(missing))
		for _, bc := range missing {
			missingSelectors = append(missingSelectors, bc.ChainSelector())
		}
		return current, fmt.Errorf(
			"missing Aptos forwarder address for chain selectors %v (set aptos_forwarder_addresses in config or %s/%s env vars)",
			missingSelectors,
			aptosForwarderAddressesEnvVar,
			aptosForwarderAddressEnvVar,
		)
	}

	var deployerPrivateKey aptoscrypto.Ed25519PrivateKey
	if err := deployerPrivateKey.FromHex(blockchain.DefaultAptosPrivateKey); err != nil {
		return current, fmt.Errorf("failed to parse Aptos deployer private key: %w", err)
	}
	deployerAccount, err := aptos.NewAccountFromSigner(&deployerPrivateKey)
	if err != nil {
		return current, fmt.Errorf("failed to create Aptos deployer signer: %w", err)
	}

	if current == nil {
		current = make(map[uint64]string)
	}

	for _, bc := range missing {
		output := bc.CtfOutput()
		if output == nil || len(output.Nodes) == 0 {
			return current, fmt.Errorf("missing Aptos node output for chain selector %d", bc.ChainSelector())
		}

		nodeURL := strings.TrimSpace(output.Nodes[0].ExternalHTTPUrl)
		if nodeURL == "" {
			return current, fmt.Errorf("missing Aptos external node URL for chain selector %d", bc.ChainSelector())
		}
		nodeURL, err = normalizeAptosNodeURL(nodeURL)
		if err != nil {
			return current, fmt.Errorf("invalid Aptos node URL for chain selector %d: %w", bc.ChainSelector(), err)
		}

		chainID := bc.ChainID()
		if chainID > 255 {
			return current, fmt.Errorf("Aptos chain id %d does not fit in uint8 for chain selector %d", chainID, bc.ChainSelector())
		}

		client, clientErr := aptos.NewNodeClient(nodeURL, uint8(chainID))
		if clientErr != nil {
			return current, fmt.Errorf("failed to create Aptos client for chain selector %d (%s): %w", bc.ChainSelector(), nodeURL, clientErr)
		}

		owner := deployerAccount.AccountAddress()
		containerName := ""
		if output != nil {
			containerName = output.ContainerName
		}
		if ensureErr := ensureAptosAccountVisible(ctx, testLogger, client, nodeURL, owner, bc.ChainSelector(), containerName); ensureErr != nil {
			testLogger.Warn().
				Uint64("chainSelector", bc.ChainSelector()).
				Str("nodeURL", nodeURL).
				Err(ensureErr).
				Msg("Aptos deployer account not confirmed visible yet; proceeding with deploy retries")
		}

		var objectAddress aptos.AccountAddress
		var pendingTxHash string
		var lastDeployErr error
		deployCtx, deployCancel := context.WithTimeout(ctx, 3*time.Minute)
		defer deployCancel()
		deployRetryErr := retry.Do(deployCtx, retry.WithMaxDuration(3*time.Minute, retry.NewFibonacci(500*time.Millisecond)), func(ctx context.Context) error {
			var deployErr error
			var pendingTx *api.PendingTransaction
			objectAddress, pendingTx, _, deployErr = aptosplatform.DeployToObject(deployerAccount, client, owner)
			if deployErr != nil {
				lastDeployErr = deployErr
				if containerName != "" {
					if fundErr := fundAptosAccountInContainer(ctx, containerName, owner.StringLong()); fundErr != nil {
						testLogger.Warn().
							Uint64("chainSelector", bc.ChainSelector()).
							Str("containerName", containerName).
							Err(fundErr).
							Msg("failed to re-fund Aptos deployer account during deploy retry")
					}
				}
				return retry.RetryableError(fmt.Errorf("deploy-to-object failed: %w", deployErr))
			}
			if pendingTx == nil {
				lastDeployErr = errors.New("nil pending transaction")
				return retry.RetryableError(fmt.Errorf("deploy-to-object returned nil pending transaction"))
			}
			pendingTxHash = pendingTx.Hash
			receipt, waitErr := client.WaitForTransaction(pendingTxHash)
			if waitErr != nil {
				lastDeployErr = waitErr
				return retry.RetryableError(fmt.Errorf("wait for deployment tx failed: %w", waitErr))
			}
			if !receipt.Success {
				return fmt.Errorf("Aptos forwarder deployment tx %s failed on chain selector %d: %s", pendingTx.Hash, bc.ChainSelector(), receipt.VmStatus)
			}
			return nil
		})
		if deployRetryErr != nil {
			if lastDeployErr != nil {
				return current, fmt.Errorf("failed to deploy Aptos platform forwarder for chain selector %d after retries (last error: %v): %w", bc.ChainSelector(), lastDeployErr, deployRetryErr)
			}
			return current, fmt.Errorf("failed to deploy Aptos platform forwarder for chain selector %d after retries: %w", bc.ChainSelector(), deployRetryErr)
		}

		addr, normErr := normalizeAptosForwarderAddress(objectAddress.StringLong())
		if normErr != nil {
			return current, fmt.Errorf("invalid Aptos forwarder address parsed from deployment output for chain selector %d: %w", bc.ChainSelector(), normErr)
		}
		current[bc.ChainSelector()] = addr
		testLogger.Info().
			Uint64("chainSelector", bc.ChainSelector()).
			Str("nodeURL", nodeURL).
			Str("txHash", pendingTxHash).
			Str("forwarderAddress", addr).
			Msg("Aptos platform forwarder deployed")
	}

	return current, nil
}

func ensureAptosAccountVisible(
	ctx context.Context,
	testLogger zerolog.Logger,
	client *aptos.NodeClient,
	nodeURL string,
	address aptos.AccountAddress,
	chainSelector uint64,
	containerName string,
) error {
	ensureCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	// Fast path: account already visible.
	if _, err := client.Account(address); err == nil {
		return nil
	}

	faucetURL, faucetErr := aptosFaucetURLFromNodeURL(nodeURL)
	if faucetErr != nil {
		return faucetErr
	}

	faucetClient, faucetClientErr := aptos.NewFaucetClient(client, faucetURL)
	if faucetClientErr != nil {
		testLogger.Warn().
			Uint64("chainSelector", chainSelector).
			Str("faucetURL", faucetURL).
			Err(faucetClientErr).
			Msg("failed to create Aptos faucet client; will try container fallback")
	}

	var lastErr error
	fundedViaAPI := false
	fundedViaContainer := false
	bo := retry.WithMaxRetries(20, retry.NewFibonacci(300*time.Millisecond))
	err := retry.Do(ensureCtx, bo, func(ctx context.Context) error {
		if _, accountErr := client.Account(address); accountErr == nil {
			return nil
		} else {
			lastErr = accountErr
		}

		if !fundedViaAPI && faucetClientErr == nil {
			// Local Aptos testnet faucet can race with initial node startup.
			// Try funding once, then wait for account visibility.
			const fundAmount = uint64(1_000_000_000)
			if fundErr := faucetClient.Fund(address, fundAmount); fundErr != nil {
				lastErr = fundErr
				return retry.RetryableError(fmt.Errorf("failed to fund Aptos deployer account via faucet (%s): %w", faucetURL, fundErr))
			}
			fundedViaAPI = true
			testLogger.Info().
				Uint64("chainSelector", chainSelector).
				Str("nodeURL", nodeURL).
				Str("faucetURL", faucetURL).
				Str("account", address.StringLong()).
				Uint64("amount", fundAmount).
				Msg("Funded Aptos deployer account while waiting for account visibility")
		}

		if !fundedViaContainer && containerName != "" {
			if fundErr := fundAptosAccountInContainer(ctx, containerName, address.StringLong()); fundErr != nil {
				lastErr = fundErr
				testLogger.Warn().
					Uint64("chainSelector", chainSelector).
					Str("containerName", containerName).
					Err(fundErr).
					Msg("failed to fund Aptos deployer account via container CLI fallback")
			} else {
				fundedViaContainer = true
				testLogger.Info().
					Uint64("chainSelector", chainSelector).
					Str("containerName", containerName).
					Str("account", address.StringLong()).
					Msg("Funded Aptos deployer account via container CLI fallback")
			}
		}

		if lastErr != nil {
			return retry.RetryableError(fmt.Errorf("Aptos account not visible yet on %s: %w", nodeURL, lastErr))
		}
		return retry.RetryableError(fmt.Errorf("Aptos account not visible yet on %s", nodeURL))
	})
	if err != nil {
		if lastErr != nil {
			return fmt.Errorf("account %s not visible on %s within timeout (last error: %v): %w", address.StringLong(), nodeURL, lastErr, err)
		}
		return fmt.Errorf("account %s not visible on %s within timeout: %w", address.StringLong(), nodeURL, err)
	}

	return nil
}

func fundAptosAccountInContainer(ctx context.Context, containerName string, account string) error {
	dc, err := framework.NewDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	cmd := []string{
		"aptos", "account", "fund-with-faucet",
		"--account", account,
		"--amount", "1000000000000",
	}
	if _, err = dc.ExecContainerWithContext(ctx, containerName, cmd); err != nil {
		return fmt.Errorf("failed to execute aptos faucet funding command in container %s: %w", containerName, err)
	}
	return nil
}

func aptosFaucetURLFromNodeURL(nodeURL string) (string, error) {
	parsed, err := url.Parse(nodeURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse Aptos node URL %q: %w", nodeURL, err)
	}
	host := parsed.Hostname()
	if host == "" {
		return "", fmt.Errorf("Aptos node URL %q has empty host", nodeURL)
	}

	// Aptos local testnet faucet is exposed on 8081.
	parsed.Host = fmt.Sprintf("%s:%s", host, blockchain.DefaultAptosFaucetPort)
	parsed.Path = ""
	parsed.RawPath = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func normalizeAptosNodeURL(nodeURL string) (string, error) {
	parsed, err := url.Parse(nodeURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse Aptos node URL %q: %w", nodeURL, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("Aptos node URL %q must include scheme and host", nodeURL)
	}
	trimmedPath := strings.TrimRight(parsed.Path, "/")
	if trimmedPath == "" {
		parsed.Path = "/v1"
	} else if trimmedPath != "/v1" {
		parsed.Path = trimmedPath + "/v1"
	}
	parsed.RawPath = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func parseAptosOCR2OnchainPublicKey(raw string) ([]byte, error) {
	key := strings.TrimSpace(raw)
	if key == "" {
		return nil, errors.New("empty OCR2 onchain public key")
	}
	// Some OCR2 exports include prefixes; keep only the final token when present.
	if strings.Contains(key, "_") {
		parts := strings.Split(key, "_")
		key = parts[len(parts)-1]
	}
	key = strings.TrimPrefix(strings.TrimPrefix(key, "0x"), "0X")
	if key == "" {
		return nil, errors.New("empty OCR2 onchain public key after normalization")
	}
	decoded, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OCR2 onchain public key %q: %w", raw, err)
	}
	if len(decoded) == 0 {
		return nil, fmt.Errorf("decoded OCR2 onchain public key %q is empty", raw)
	}
	return decoded, nil
}

func aptosDonOraclePublicKeys(don *cre.Don) ([][]byte, error) {
	workers, err := don.Workers()
	if err != nil {
		return nil, fmt.Errorf("failed to list worker nodes for DON %q: %w", don.Name, err)
	}

	oracles := make([][]byte, 0, len(workers))
	for _, worker := range workers {
		ocr2ID := ""
		if worker.Keys != nil && worker.Keys.OCR2BundleIDs != nil {
			ocr2ID = worker.Keys.OCR2BundleIDs[chainselectors.FamilyAptos]
		}
		if ocr2ID == "" {
			// Fallback: fetch directly from node in case cached key IDs were not populated.
			fetchedID, fetchErr := worker.Clients.GQLClient.FetchOCR2KeyBundleID(context.Background(), strings.ToUpper(chainselectors.FamilyAptos))
			if fetchErr != nil {
				return nil, fmt.Errorf("missing Aptos OCR2 bundle id for worker %q in DON %q and fallback fetch failed: %w", worker.Name, don.Name, fetchErr)
			}
			if fetchedID == "" {
				return nil, fmt.Errorf("missing Aptos OCR2 bundle id for worker %q in DON %q", worker.Name, don.Name)
			}
			ocr2ID = fetchedID
			if worker.Keys != nil {
				if worker.Keys.OCR2BundleIDs == nil {
					worker.Keys.OCR2BundleIDs = make(map[string]string)
				}
				worker.Keys.OCR2BundleIDs[chainselectors.FamilyAptos] = ocr2ID
			}
		}

		exported, expErr := worker.ExportOCR2Keys(ocr2ID)
		if expErr != nil {
			return nil, fmt.Errorf("failed to export Aptos OCR2 key for worker %q (bundle %s): %w", worker.Name, ocr2ID, expErr)
		}

		pubkey, parseErr := parseAptosOCR2OnchainPublicKey(exported.OnchainPublicKey)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid Aptos OCR2 onchain public key for worker %q: %w", worker.Name, parseErr)
		}
		oracles = append(oracles, pubkey)
	}

	return oracles, nil
}

func configureAptosForwarderContracts(
	ctx context.Context,
	testLogger zerolog.Logger,
	provider infra.Provider,
	deployedBlockchains []blockchains.Blockchain,
	dons *cre.Dons,
	forwarderAddresses map[uint64]string,
) error {
	if dons == nil || len(forwarderAddresses) == 0 {
		return nil
	}

	var deployerPrivateKey aptoscrypto.Ed25519PrivateKey
	if err := deployerPrivateKey.FromHex(blockchain.DefaultAptosPrivateKey); err != nil {
		return fmt.Errorf("failed to parse Aptos deployer private key for forwarder config: %w", err)
	}
	deployerAccount, err := aptos.NewAccountFromSigner(&deployerPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to create Aptos deployer signer for forwarder config: %w", err)
	}
	deployerAddress := deployerAccount.AccountAddress()

	aptosChainsByChainID := make(map[uint64]blockchains.Blockchain)
	for _, bc := range deployedBlockchains {
		if bc.IsFamily(chainselectors.FamilyAptos) {
			aptosChainsByChainID[bc.ChainID()] = bc
		}
	}

	for _, don := range dons.DonsWithFlag(cre.WriteAptosCapability) {
		aptosChainIDs, chainErr := don.GetEnabledChainIDsForCapability(cre.WriteAptosCapability)
		if chainErr != nil {
			return fmt.Errorf("failed to get Aptos chain IDs for DON %q: %w", don.Name, chainErr)
		}
		if len(aptosChainIDs) == 0 {
			continue
		}

		workers, workerErr := don.Workers()
		if workerErr != nil {
			return fmt.Errorf("failed to get worker nodes for DON %q: %w", don.Name, workerErr)
		}
		f := (len(workers) - 1) / 3
		if f <= 0 {
			return fmt.Errorf("invalid Aptos DON %q fault tolerance F=%d (workers=%d)", don.Name, f, len(workers))
		}
		if f > 255 {
			return fmt.Errorf("Aptos DON %q fault tolerance F=%d exceeds u8", don.Name, f)
		}
		if don.ID > uint64(^uint32(0)) {
			return fmt.Errorf("DON %q id %d exceeds u32 for Aptos forwarder config", don.Name, don.ID)
		}

		oracles, oracleErr := aptosDonOraclePublicKeys(don)
		if oracleErr != nil {
			return oracleErr
		}

		for _, chainID := range aptosChainIDs {
			bc, ok := aptosChainsByChainID[chainID]
			if !ok {
				return fmt.Errorf("Aptos chain id %d enabled for DON %q but no Aptos blockchain is configured", chainID, don.Name)
			}

			forwarderHex := strings.TrimSpace(forwarderAddresses[bc.ChainSelector()])
			if forwarderHex == "" {
				return fmt.Errorf("missing Aptos forwarder address for chain selector %d (chain id %d)", bc.ChainSelector(), chainID)
			}

			output := bc.CtfOutput()
			if output == nil || len(output.Nodes) == 0 {
				return fmt.Errorf("missing Aptos node output for chain selector %d", bc.ChainSelector())
			}
			nodeURL, normErr := normalizeAptosNodeURL(output.Nodes[0].ExternalHTTPUrl)
			if normErr != nil {
				return fmt.Errorf("invalid Aptos node URL for chain selector %d: %w", bc.ChainSelector(), normErr)
			}
			if bc.ChainID() > 255 {
				return fmt.Errorf("Aptos chain id %d does not fit in uint8 for chain selector %d", bc.ChainID(), bc.ChainSelector())
			}

			client, clientErr := aptos.NewNodeClient(nodeURL, uint8(bc.ChainID()))
			if clientErr != nil {
				return fmt.Errorf("failed to create Aptos client for chain selector %d (%s): %w", bc.ChainSelector(), nodeURL, clientErr)
			}

			if ensureErr := ensureAptosAccountVisible(ctx, testLogger, client, nodeURL, deployerAddress, bc.ChainSelector(), output.ContainerName); ensureErr != nil {
				if !provider.IsDocker() {
					return fmt.Errorf("failed to ensure Aptos deployer account visibility for chain selector %d: %w", bc.ChainSelector(), ensureErr)
				}
				testLogger.Warn().
					Uint64("chainSelector", bc.ChainSelector()).
					Str("nodeURL", nodeURL).
					Err(ensureErr).
					Msg("Aptos deployer account not confirmed visible yet; proceeding with forwarder set_config retries")
			}

			var forwarderAddr aptos.AccountAddress
			if parseErr := forwarderAddr.ParseStringRelaxed("0x" + strings.TrimPrefix(forwarderHex, "0x")); parseErr != nil {
				return fmt.Errorf("invalid Aptos forwarder address for chain selector %d: %w", bc.ChainSelector(), parseErr)
			}
			forwarderContract := aptosplatform.Bind(forwarderAddr, client).Forwarder()
			testLogger.Info().
				Str("donName", don.Name).
				Uint64("donID", don.ID).
				Uint64("chainSelector", bc.ChainSelector()).
				Uint64("chainID", chainID).
				Str("forwarderAddress", "0x"+strings.TrimPrefix(forwarderHex, "0x")).
				Uint32("configVersion", aptosForwarderConfigVersion).
				Int("workerCount", len(workers)).
				Int("f", f).
				Int("oraclesCount", len(oracles)).
				Msg("configuring Aptos forwarder set_config")

			var pendingTxHash string
			var lastSetConfigErr error
			setConfigErr := retry.Do(
				ctx,
				retry.WithMaxDuration(2*time.Minute, retry.NewFibonacci(500*time.Millisecond)),
				func(ctx context.Context) error {
					pendingTx, txErr := forwarderContract.SetConfig(&bind.TransactOpts{
						Signer: deployerAccount,
					}, uint32(don.ID), aptosForwarderConfigVersion, byte(f), oracles)
					if txErr != nil {
						lastSetConfigErr = txErr
						if output.ContainerName != "" {
							if fundErr := fundAptosAccountInContainer(ctx, output.ContainerName, deployerAddress.StringLong()); fundErr != nil {
								testLogger.Warn().
									Uint64("chainSelector", bc.ChainSelector()).
									Str("containerName", output.ContainerName).
									Err(fundErr).
									Msg("failed to fund Aptos deployer account during set_config retry")
							}
						}
						return retry.RetryableError(fmt.Errorf("set_config transaction submit failed: %w", txErr))
					}
					pendingTxHash = pendingTx.Hash
					receipt, waitErr := client.WaitForTransaction(pendingTxHash)
					if waitErr != nil {
						lastSetConfigErr = waitErr
						return retry.RetryableError(fmt.Errorf("waiting for set_config transaction failed: %w", waitErr))
					}
					if !receipt.Success {
						lastSetConfigErr = fmt.Errorf("vm status: %s", receipt.VmStatus)
						return retry.RetryableError(fmt.Errorf("set_config transaction failed: %s", receipt.VmStatus))
					}
					return nil
				},
			)
			if setConfigErr != nil {
				if lastSetConfigErr != nil {
					return fmt.Errorf("failed to configure Aptos forwarder %s for DON %q on chain selector %d (last error: %v): %w", forwarderHex, don.Name, bc.ChainSelector(), lastSetConfigErr, setConfigErr)
				}
				return fmt.Errorf("failed to configure Aptos forwarder %s for DON %q on chain selector %d: %w", forwarderHex, don.Name, bc.ChainSelector(), setConfigErr)
			}

			testLogger.Info().
				Str("donName", don.Name).
				Uint64("donID", don.ID).
				Uint64("chainSelector", bc.ChainSelector()).
				Uint64("chainID", chainID).
				Str("forwarderAddress", "0x"+strings.TrimPrefix(forwarderHex, "0x")).
				Int("workerCount", len(workers)).
				Int("f", f).
				Int("oraclesCount", len(oracles)).
				Str("txHash", pendingTxHash).
				Msg("configured Aptos forwarder set_config")
		}
	}

	return nil
}

// waitForWorkflowNodesEVMChainConfig waits until all workflow DON nodes have registered
// an EVM chain config with the Job Distributor. This is required before configuring
// the Capability Registry (which uses deployment.NodeInfo and requires EVM chain/OCR2
// config for the registry chain, e.g. in Aptos topology where workflow nodes have
// both EVM and Aptos chains).
func waitForWorkflowNodesEVMChainConfig(ctx context.Context, oc deployment.NodeChainConfigsLister, dons *cre.Dons, testLogger zerolog.Logger) error {
	var workflowDon *cre.Don
	for _, don := range dons.List() {
		if don.HasFlag(cre.WorkflowDON) {
			workflowDon = don
			break
		}
	}
	if workflowDon == nil {
		return nil
	}
	nodeIDs := workflowDon.JDNodeIDs()
	if len(nodeIDs) == 0 {
		return nil
	}
	err := retry.Do(ctx, retry.WithMaxDuration(2*time.Minute, retry.NewFibonacci(2*time.Second)), func(ctx context.Context) error {
		_, err := deployment.NodeInfo(nodeIDs, oc)
		if err == nil {
			return nil
		}
		if errors.Is(err, deployment.ErrMissingNodeMetadata) {
			testLogger.Debug().Err(err).Msg("Workflow nodes have not yet registered EVM chain config with JD, retrying")
			return retry.RetryableError(err)
		}
		return err
	})
	return err
}

func SetupTestEnvironment(
	ctx context.Context,
	testLogger zerolog.Logger,
	singleFileLogger logger.Logger,
	input *SetupInput,
	relativePathToRepoRoot string,
) (*SetupOutput, error) {
	if input == nil {
		return nil, pkgerrors.New("input is nil")
	}

	//TODO: remove these checks in December 2025, when everyone has migrated
	if val := os.Getenv("E2E_JD_IMAGE"); val != "" {
		return nil, errors.New("E2E_JD_IMAGE and E2E_JD_VERSION are deprecated, please use CTF_JD_IMAGE instead to specify the Job Distributor image with tag")
	}

	if val := os.Getenv("E2E_TEST_CHAINLINK_IMAGE"); val != "" {
		return nil, errors.New("E2E_TEST_CHAINLINK_IMAGE and E2E_TEST_CHAINLINK_VERSION are deprecated, please use CTF_CHAINLINK_IMAGE instead to specify the Chainlink Node image with tag")
	}

	if err := input.Validate(); err != nil {
		return nil, pkgerrors.Wrap(err, "input validation failed")
	}
	if cfgErr := configureAptosContractsPathFromEnv(testLogger, input.BlockchainsInput); cfgErr != nil {
		return nil, cfgErr
	}

	s3Output, s3Err := workflow.StartS3(testLogger, input.S3ProviderInput, input.StageGen)
	if s3Err != nil {
		return nil, pkgerrors.Wrap(s3Err, "failed to start S3 provider")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting %d blockchain(s)", len(input.BlockchainsInput))))

	deployedBlockchains, startErr := blockchains.Start(
		ctx,
		testLogger,
		singleFileLogger,
		input.BlockchainsInput,
		input.BlockchainDeployers,
	)
	if startErr != nil {
		return nil, pkgerrors.Wrap(startErr, "failed to start blockchains")
	}

	creEnvironment := &cre.Environment{
		Blockchains:           deployedBlockchains.Outputs,
		ContractVersions:      input.ContractVersions,
		Provider:              input.Provider,
		RegistryChainSelector: deployedBlockchains.RegistryChain().ChainSelector(),
	}

	aptosForwarderAddresses, addrErr := resolveAptosForwarderAddresses(testLogger, deployedBlockchains.Outputs, input.AptosForwarderAddresses)
	if addrErr != nil {
		return nil, fmt.Errorf("failed to resolve Aptos forwarder addresses: %w", addrErr)
	}
	aptosForwarderAddresses, addrErr = deployMissingAptosForwarders(ctx, testLogger, input.Provider, deployedBlockchains.Outputs, aptosForwarderAddresses)
	if addrErr != nil {
		return nil, fmt.Errorf("failed to auto-deploy missing Aptos forwarders: %w", addrErr)
	}
	creEnvironment.AptosForwarderAddresses = aptosForwarderAddresses

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Blockchains started in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Deploying Workflow and Capability Registry contracts")))

	deployKeystoneContractsOutput, deployErr := crecontracts.DeployKeystoneContracts(
		ctx,
		testLogger,
		singleFileLogger,
		crecontracts.DeployKeystoneContractsInput{
			CldfEnvironment:  newCldfEnvironment(ctx, singleFileLogger, deployedBlockchains.CldfBlockChains),
			CtfBlockchains:   deployedBlockchains.Outputs,
			ContractVersions: input.ContractVersions,
			WithV2Registries: input.WithV2Registries,
		},
	)
	if deployErr != nil {
		return nil, pkgerrors.Wrap(deployErr, "failed to deploy Keystone contracts")
	}
	creEnvironment.CldfEnvironment = deployKeystoneContractsOutput.Env

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Workflow and Capability Registry contracts deployed in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Preparing DONs configuration")))

	topology, tErr := cre.NewTopology(input.NodeSets, creEnvironment.Provider, input.CapabilityConfigs)
	if tErr != nil {
		return nil, pkgerrors.Wrap(tErr, "failed to create topology")
	}

	updatedNodeSets, topoErr := donconfig.PrepareNodeTOMLs(
		ctx,
		topology,
		creEnvironment,
		input.NodeSets,
		input.Capabilities,
		input.ConfigFactoryFunctions,
	)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs configuration prepared in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Applying Features before environment startup")))
	var donsCapabilities = make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	var capabilityToOCR3Config = make(map[string]*ocr3.OracleConfig)
	for _, feature := range input.Features.List() {
		for _, donMetadata := range topology.DonsMetadataWithFlag(feature.Flag()) {
			testLogger.Info().Msgf("Executing PreEnvStartup for feature %s for don '%s'", feature.Flag(), donMetadata.Name)
			output, preErr := feature.PreEnvStartup(
				ctx,
				testLogger,
				donMetadata,
				topology,
				creEnvironment,
			)
			if preErr != nil {
				return nil, fmt.Errorf("failed to execute PreEnvStartup for feature %s: %w", feature.Flag(), preErr)
			}
			if output != nil {
				if donsCapabilities[donMetadata.ID] == nil {
					donsCapabilities[donMetadata.ID] = []keystone_changeset.DONCapabilityWithConfig{}
				}
				donsCapabilities[donMetadata.ID] = append(donsCapabilities[donMetadata.ID], output.DONCapabilityWithConfig...)
				maps.Copy(capabilityToOCR3Config, output.CapabilityToOCR3Config)
			}
			testLogger.Info().Msgf("PreEnvStartup for feature %s executed successfully", feature.Flag())
		}
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Applied Features in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	queue := worker.New(ctx, 10)
	defer queue.StopAndWait() // Ensure cleanup on any exit path

	jdStartedFuture := queue.SubmitAny(func(ctx context.Context) (any, error) {
		// TODO: pass context after we update the CTF to accept context, when creating new JD instance
		jdOutput, startJDErr := StartJD(ctx, testLogger, *input.JdInput, input.Provider)
		if startJDErr != nil {
			return nil, pkgerrors.Wrap(startJDErr, "failed to start Job Distributor")
		}
		return jdOutput, nil
	})

	donsStartedFuture := queue.SubmitAny(func(ctx context.Context) (any, error) {
		nodeSetOutput, startDonsErr := StartDONs(ctx, testLogger, topology, input.Provider, deployedBlockchains.RegistryChain().CtfOutput(), input.CapabilityConfigs, input.CopyCapabilityBinaries, updatedNodeSets)
		if startDonsErr != nil {
			return nil, pkgerrors.Wrap(startDonsErr, "failed to start DONs")
		}

		return nodeSetOutput, nil
	})

	// Await both futures to ensure proper cleanup even if one fails
	startedJD, jdStartErr := worker.AwaitAs[*StartedJD](ctx, jdStartedFuture)
	startedDONs, donStartErr := worker.AwaitAs[*StartedDONs](ctx, donsStartedFuture)

	// Check errors after both awaits complete
	// If both failed, prefer the non-context-cancelled error as it's likely the root cause
	if jdStartErr != nil && donStartErr != nil {
		// If one is context.Canceled, it was likely caused by the other task's error
		if pkgerrors.Is(jdStartErr, context.Canceled) && !pkgerrors.Is(donStartErr, context.Canceled) {
			return nil, pkgerrors.Wrap(donStartErr, "failed to start DONs")
		}
		if pkgerrors.Is(donStartErr, context.Canceled) && !pkgerrors.Is(jdStartErr, context.Canceled) {
			return nil, pkgerrors.Wrap(jdStartErr, "failed to start Job Distributor")
		}
		// Both real errors
		return nil, pkgerrors.Wrap(errors.Join(fmt.Errorf("JD failed to start: %w", jdStartErr), fmt.Errorf("DONs failed to start: %w", donStartErr)), "failed to start Job Distributor AND Dons")
	}
	if jdStartErr != nil {
		return nil, pkgerrors.Wrap(jdStartErr, "failed to start Job Distributor")
	}
	if donStartErr != nil {
		return nil, pkgerrors.Wrap(donStartErr, "failed to start DONs")
	}
	dons := cre.NewDons(startedDONs.DONs(), topology.GatewayConnectors)
	deployKeystoneContractsOutput.Env.Offchain = startedJD.Client

	linkDonsToJDInput := &cre.LinkDonsToJDInput{
		Blockchains:     deployedBlockchains.Outputs,
		CldfEnvironment: deployKeystoneContractsOutput.Env,
		Topology:        topology,
		Dons:            dons,
	}

	cldErr := cre.LinkToJobDistributor(ctx, linkDonsToJDInput)
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to link DONs to Job Distributor")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs and Job Distributor started and linked in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Creating Jobs with Job Distributor")))

	gJobErr := gateway.CreateJobs(ctx, creEnvironment, dons, topology.GatewayConfigs, input.GatewayWhitelistConfig)
	if gJobErr != nil {
		return nil, pkgerrors.Wrap(gJobErr, "failed to create gateway jobs with Job Distributor")
	}

	// Deprecated: use Features instead. Support for InstallableCapability will be removed in the future.
	jobSpecFactoryFunctions := make([]cre.JobSpecFn, 0)
	for _, capability := range input.Capabilities {
		jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, capability.JobSpecFn())
	}

	// allow to pass custom job spec factories for extensibility
	jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, input.JobSpecFactoryFunctions...)

	createJobsDeps := CreateJobsWithJdOpDeps{
		Logger:                        testLogger,
		SingleFileLogger:              singleFileLogger,
		RegistryChainBlockchainOutput: deployedBlockchains.RegistryChain().CtfOutput(),
		JobSpecFactoryFunctions:       jobSpecFactoryFunctions,
		CreEnvironment:                creEnvironment,
		Dons:                          dons,
		NodeSets:                      input.NodeSets,
		Capabilities:                  input.Capabilities,
	}
	_, createJobsErr := operations.ExecuteOperation(deployKeystoneContractsOutput.Env.OperationsBundle, CreateJobsWithJdOp, createJobsDeps, CreateJobsWithJdOpInput{})
	if createJobsErr != nil {
		return nil, pkgerrors.Wrap(createJobsErr, "failed to create jobs with Job Distributor")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Jobs created in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Funding Chainlink nodes")))

	fundingPerChainFamilyForEachNode := map[string]uint64{
		chainselectors.FamilyEVM:    10000000000000000, // 0.01 ETH
		chainselectors.FamilySolana: 50_000_000_000,    // 50 SOL
		chainselectors.FamilyTron:   100_000_000,       // 100 TRX in SUN
		chainselectors.FamilyAptos:  1_000_000_000_000, // 1,000 APT (octas) for local devnet sender accounts
	}

	fErr := FundNodes(
		ctx,
		testLogger,
		dons,
		deployedBlockchains.Outputs,
		fundingPerChainFamilyForEachNode,
	)
	if fErr != nil {
		return nil, pkgerrors.Wrap(fErr, "failed to fund chainlink nodes")
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Chainlink nodes funded in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Configuring Workflow and Capability Registry contracts")))

	// Wait for workflow DON nodes to register EVM chain config with the JD before configuring
	// the Capability Registry (required for Aptos and other topologies where nodes have multiple chains).
	if waitErr := waitForWorkflowNodesEVMChainConfig(ctx, deployKeystoneContractsOutput.Env.Offchain, dons, testLogger); waitErr != nil {
		return nil, pkgerrors.Wrap(waitErr, "failed waiting for workflow nodes to register EVM chain config with Job Distributor")
	}

	wfRegVersion := input.ContractVersions[keystone_changeset.WorkflowRegistry.String()]
	workflowRegistryConfigurationOutput, wfErr := workflow.ConfigureWorkflowRegistry(
		ctx,
		testLogger,
		singleFileLogger,
		&cre.WorkflowRegistryInput{
			ContractAddress: common.HexToAddress(crecontracts.MustGetAddressFromDataStore(deployKeystoneContractsOutput.Env.DataStore, deployedBlockchains.RegistryChain().ChainSelector(), keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")),
			ContractVersion: cldf.TypeAndVersion{Version: *wfRegVersion},
			ChainSelector:   deployedBlockchains.RegistryChain().ChainSelector(),
			CldEnv:          deployKeystoneContractsOutput.Env,
			AllowedDonIDs:   topology.WorkflowDONIDs,
			WorkflowOwners:  []common.Address{deployedBlockchains.RegistryChain().(*evm.Blockchain).SethClient.MustGetRootKeyAddress()}, // registry chain is always EVM
		},
	)
	if wfErr != nil {
		return nil, pkgerrors.Wrap(wfErr, "failed to configure workflow registry")
	}

	wfFiltersFuture := queue.SubmitErr(func(ctx context.Context) error {
		// we currently have no way of checking if filters were registered in Kubernetes mode
		// as we don't have a way to get its database connection string
		if !input.Provider.IsDocker() {
			return nil
		}

		fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Waiting for Workflow Registry filters registration\n\n"))
		defer fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Finished waiting for Workflow Registry filters registration\n\n"))

		// this operation can always safely run in the background, since it doesn't change on-chain state, it only reads data from databases
		switch wfRegVersion.Major() {
		case 2:
			// There are no filters registered with the V2 WF Registry Syncer
			return nil
		default:
			return workflow.WaitForAllNodesToHaveExpectedFiltersRegistered(ctx, singleFileLogger, testLogger, deployedBlockchains.RegistryChain().ChainID(), dons, updatedNodeSets)
		}
	})

	capRegInput := cre.ConfigureCapabilityRegistryInput{
		ChainSelector: deployedBlockchains.RegistryChain().ChainSelector(),
		CldEnv:        creEnvironment.CldfEnvironment,
		Blockchains:   deployedBlockchains.Outputs,
		Topology:      topology,
		CapabilitiesRegistryAddress: ptr.Ptr(crecontracts.MustGetAddressFromMemoryDataStore(
			deployKeystoneContractsOutput.MemoryDataStore,
			deployedBlockchains.RegistryChain().ChainSelector(),
			keystone_changeset.CapabilitiesRegistry.String(),
			input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()],
			""),
		),
		NodeSets:                 input.NodeSets,
		WithV2Registries:         input.WithV2Registries,
		DONCapabilityWithConfigs: make(map[uint64][]keystone_changeset.DONCapabilityWithConfig),
		CapabilityToOCR3Config:   capabilityToOCR3Config,
	}

	for _, capability := range input.Capabilities {
		configFn := capability.CapabilityRegistryV1ConfigFn()
		capRegInput.CapabilityRegistryConfigFns = append(capRegInput.CapabilityRegistryConfigFns, configFn)
	}
	capRegInput.CapabilityRegistryConfigFns = append(capRegInput.CapabilityRegistryConfigFns, input.CapabilitiesContractFactoryFunctions...)
	maps.Copy(capRegInput.DONCapabilityWithConfigs, donsCapabilities)

	_, capRegErr := crecontracts.ConfigureCapabilityRegistry(capRegInput)
	if capRegErr != nil {
		return nil, pkgerrors.Wrap(capRegErr, "failed to configure Capability Registry contracts")
	}
	// Match Aptos data-feeds setup flow: deploy forwarder, then set forwarder config
	// before write workflows execute.
	if cfgErr := configureAptosForwarderContracts(
		ctx,
		testLogger,
		input.Provider,
		deployedBlockchains.Outputs,
		dons,
		aptosForwarderAddresses,
	); cfgErr != nil {
		return nil, fmt.Errorf("failed to configure Aptos forwarders: %w", cfgErr)
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Workflow and Capability Registry contracts configured in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Applying Features after environment startup")))

	for _, feature := range input.Features.List() {
		for _, don := range dons.DonsWithFlag(feature.Flag()) {
			testLogger.Info().Msgf("Executing PostEnvStartup for feature %s for don '%s'", feature.Flag(), don.Name)
			if pErr := feature.PostEnvStartup(
				ctx,
				testLogger,
				don,
				dons,
				creEnvironment,
			); pErr != nil {
				return nil, fmt.Errorf("failed to execute PostEnvStartup for feature %s: %w", feature.Flag(), pErr)
			}
			testLogger.Info().Msgf("PostEnvStartup for feature %s executed successfully", feature.Flag())
		}
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Features applied in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	// Sharding setup moved AFTER PostEnvStartup to ensure OCR3 configs work properly
	if topology.DonsMetadata.ShardingEnabled() {
		fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Setting up Sharding")))
		err := sharding.SetupSharding(ctx, sharding.SetupShardingInput{
			Logger:   testLogger,
			CreEnv:   creEnvironment,
			Topology: topology,
			Dons:     dons,
		})
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to setup Sharding")
		}
		fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Sharding setup in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	}

	if err := worker.AwaitErr(ctx, wfFiltersFuture); err != nil {
		return nil, pkgerrors.Wrap(err, "failed while waiting for workflow registry filters registration")
	}

	appendOutputsToInput(input, startedDONs.NodeOutputs(), deployedBlockchains.Outputs, startedJD.JDOutput)

	if err := workflowRegistryConfigurationOutput.Store(config.MustWorkflowRegistryStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to store workflow registry configuration output")
	}

	return &SetupOutput{
		WorkflowRegistryConfigurationOutput: workflowRegistryConfigurationOutput, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		Dons:                                dons,
		NodeOutput:                          startedDONs.NodeOutputs(),
		CreEnvironment:                      creEnvironment,
		S3ProviderOutput:                    s3Output,
		GatewayConnectors:                   topology.GatewayConnectors,
	}, nil
}

func appendOutputsToInput(input *SetupInput, nodeSetOutput []*cre.NodeSetOutput, blockchains []blockchains.Blockchain, jdOutput *jd.Output) {
	// append the nodeset output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	for idx, nsOut := range nodeSetOutput {
		input.NodeSets[idx].Out = nsOut.Output
	}

	for idx, blockchain := range blockchains {
		input.BlockchainsInput[idx].Out = blockchain.CtfOutput()
	}

	// append the jd output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	input.JdInput.Out = jdOutput
}

func newCldfEnvironment(ctx context.Context, singleFileLogger logger.Logger, cldfBlockchains cldf_chain.BlockChains) *cldf.Environment {
	allChainsCLDEnvironment := &cldf.Environment{
		Name:              cre.EnvironmentName,
		Logger:            singleFileLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(), // can't set it to nil, because some changesets save addresses both to the address book and datastore
		DataStore:         datastore.NewMemoryDataStore().Seal(),
		GetContext: func() context.Context {
			return ctx
		},
		BlockChains: cldfBlockchains,
		OCRSecrets:  focr.XXXGenerateTestOCRSecrets(),
		OperationsBundle: operations.NewBundle(
			func() context.Context { return ctx },
			singleFileLogger, operations.NewMemoryReporter()),
	}

	return allChainsCLDEnvironment
}
