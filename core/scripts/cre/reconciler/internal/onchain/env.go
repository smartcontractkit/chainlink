package onchain

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider/rpcclient"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	cldfjd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cldflogger "github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	griddleinfra "github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	creaptos "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/aptos"
	creevm "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	cresolana "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

// SolanaPrivateKeyEnv is the ONLY place the Solana deployer/signer key may come
// from — required by cresolana's Deploy() even for an already-running,
// externally-provided (Kubernetes) chain; there is no default.
const SolanaPrivateKeyEnv = "SOLANA_PRIVATE_KEY"

// buildCldfEnv creates a cldf.Environment with an EVM chain provider for the
// desired state's registry chain and an optional JD client.
func (d *Deployer) buildCldfEnv(ctx context.Context, desired *domain.DesiredState) (*cldf.Environment, uint64, error) {
	registryChain, ok := desired.RegistryChain()
	if !ok {
		return nil, 0, errors.New("no registry chain declared in desired state (exactly one [[chains]] entry must set registry = true)")
	}

	chainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(
		strconv.FormatUint(registryChain.ChainID, 10),
		chainselectors.FamilyEVM,
	)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "no chain details for chain ID %d", registryChain.ChainID)
	}

	provider, err := cldf_evm_provider.NewRPCChainProvider(
		chainDetails.ChainSelector,
		cldf_evm_provider.RPCChainProviderConfig{
			DeployerTransactorGen: cldf_evm_provider.TransactorFromRaw(d.deployerKey),
			RPCs: []rpcclient.RPC{
				{
					Name:               "default",
					WSURL:              registryChain.WSURL,
					HTTPURL:            registryChain.HTTPURL,
					PreferredURLScheme: rpcclient.URLSchemePreferenceHTTP,
				},
			},
			ConfirmFunctor: cldf_evm_provider.ConfirmFuncGeth(
				3*time.Minute,
				cldf_evm_provider.WithTickInterval(5*time.Millisecond),
			),
		},
	).Initialize(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create EVM chain provider")
	}

	cldfLogger, err := cldflogger.New()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create cldf logger")
	}

	offchainClient, err := buildOffchainClient(desired.JD)
	if err != nil {
		return nil, 0, err
	}

	env := &cldf.Environment{
		Name:              desired.JD.Environment,
		Logger:            cldfLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         datastore.NewMemoryDataStore().Seal(),
		Offchain:          offchainClient,
		GetContext: func() context.Context {
			return ctx
		},
		BlockChains: cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{provider}),
		OCRSecrets:  focr.XXXGenerateTestOCRSecrets(),
		OperationsBundle: operations.NewBundle(
			func() context.Context { return ctx },
			cldfLogger,
			operations.NewMemoryReporter(),
		),
	}

	return env, chainDetails.ChainSelector, nil
}

// buildOffchainClient constructs a JD offchain.Client from JD config alone —
// no EVM/chain provider setup required, unlike the rest of buildCldfEnv. This
// is what a JD-only preflight check (e.g. node-label validation) can call
// directly, without paying for full CLDF environment construction. Returns a
// nil client (not an error) if JD isn't configured or the access token isn't
// set, matching buildCldfEnv's existing behavior.
func buildOffchainClient(jdCfg domain.JDConfig) (offchain.Client, error) {
	token := griddleinfra.JDAccessToken()
	if jdCfg.GRPC == "" || token == "" {
		return nil, nil
	}

	useTLS := jdCfg.UseTLS
	if !useTLS && strings.Contains(jdCfg.GRPC, ":443") {
		useTLS = true
	}

	var creds credentials.TransportCredentials
	if useTLS {
		creds = credentials.NewTLS(nil)
	} else {
		creds = insecure.NewCredentials()
	}

	jdClient, err := cldfjd.NewJDClient(cldfjd.JDConfig{
		GRPC:  jdCfg.GRPC,
		Creds: creds,
		Auth: oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
			TokenType:   "Bearer",
		}),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create JD client")
	}
	return jdClient, nil
}

// syncAddressBook copies every address ref currently in env.DataStore into
// state.Addresses, so state.toml always reflects everything deployed so far —
// not just CapabilitiesRegistry/WorkflowRegistry. Without this, contracts
// deployed by later phases (e.g. Feature.PreEnvStartup deploying a Forwarder)
// are lost the moment the process exits, since only state.Addresses is persisted.
func (d *Deployer) syncAddressBook(env *cldf.Environment, state *domain.StateFile) error {
	refs, err := env.DataStore.Addresses().Fetch()
	if err != nil {
		return errors.Wrap(err, "failed to fetch address refs from datastore")
	}
	for _, ref := range refs {
		version := ""
		if ref.Version != nil {
			version = ref.Version.String()
		}
		state.SetAddress(domain.AddressRef{
			ChainSelector: domain.ChainSelector(ref.ChainSelector),
			Address:       ref.Address,
			Type:          ref.Type.String(),
			Version:       version,
			Qualifier:     ref.Qualifier,
		})
	}
	return nil
}

func hydrateMemoryDataStoreFromState(state *domain.StateFile, chainSelector uint64) (*datastore.MemoryDataStore, error) {
	memDs := datastore.NewMemoryDataStore()
	for _, ref := range state.Addresses {
		if uint64(ref.ChainSelector) != chainSelector {
			continue
		}
		version, err := semver.NewVersion(ref.Version)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid version %q for contract %s", ref.Version, ref.Type)
		}
		if err := memDs.Addresses().Add(datastore.AddressRef{
			ChainSelector: uint64(ref.ChainSelector),
			Address:       ref.Address,
			Type:          datastore.ContractType(ref.Type),
			Version:       version,
			Qualifier:     ref.Qualifier,
		}); err != nil {
			return nil, errors.Wrapf(err, "failed to add %s address to datastore", ref.Type)
		}
	}
	return memDs, nil
}

// buildCreEnvironment wraps the cldf environment with CRE-specific context
// used by Features and workflow configuration.
func (d *Deployer) buildCreEnvironment(
	ctx context.Context,
	desired *domain.DesiredState,
	cldfEnv *cldf.Environment,
	registryChainSelector uint64,
) (*cre.Environment, error) {
	provider := infra.Provider{Type: infra.Kubernetes}
	evmDeployer := creevm.NewDeployer(d.log, &provider)
	solanaDeployer := cresolana.NewDeployer(d.log, &provider)
	aptosDeployer := creaptos.NewDeployer(d.log, &provider)

	var evmChains int
	blockchainsOut := make([]creblockchains.Blockchain, 0, len(desired.Chains))
	for _, chain := range desired.Chains {
		chainIDStr := strconv.FormatUint(chain.ChainID, 10)

		var bc creblockchains.Blockchain
		var err error
		switch chain.Family {
		case cre.EVMCapability:
			evmChains++
			bc, err = evmDeployer.Deploy(ctx, &blockchain.Input{
				Type:    blockchain.TypeAnvil,
				ChainID: chainIDStr,
				Out: &blockchain.Output{
					Type:    blockchain.TypeAnvil,
					Family:  chainselectors.FamilyEVM,
					ChainID: chainIDStr,
					Nodes: []*blockchain.Node{{
						ExternalWSUrl:   chain.WSURL,
						ExternalHTTPUrl: chain.HTTPURL,
					}},
				},
			})
		case cre.SolanaCapability:
			if os.Getenv(SolanaPrivateKeyEnv) == "" {
				return nil, errors.Errorf("chain_id %d (family solana): %s env var is required to build the Solana blockchain provider", chain.ChainID, SolanaPrivateKeyEnv)
			}
			genesisHash := chain.SolanaGenesisHash()
			contractsDir, dirErr := os.MkdirTemp("", "reconciler-solana-")
			if dirErr != nil {
				return nil, errors.Wrap(dirErr, "failed to create solana contracts dir")
			}
			bc, err = solanaDeployer.Deploy(ctx, &blockchain.Input{
				Type:         blockchain.TypeSolana,
				ChainID:      genesisHash,
				ContractsDir: contractsDir,
				Out: &blockchain.Output{
					Type:    blockchain.TypeSolana,
					Family:  chainselectors.FamilySolana,
					ChainID: genesisHash,
					Nodes: []*blockchain.Node{{
						ExternalWSUrl:   chain.WSURL,
						ExternalHTTPUrl: chain.HTTPURL,
					}},
				},
			})
		case cre.AptosCapability:
			bc, err = aptosDeployer.Deploy(ctx, &blockchain.Input{
				Type:    blockchain.TypeAptos,
				ChainID: chainIDStr,
				Out: &blockchain.Output{
					Type:    blockchain.TypeAptos,
					Family:  chainselectors.FamilyAptos,
					ChainID: chainIDStr,
					Nodes: []*blockchain.Node{{
						ExternalHTTPUrl: chain.HTTPURL,
					}},
				},
			})
		default:
			return nil, errors.Errorf("chain_id %d: unsupported family %q", chain.ChainID, chain.Family)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build %s blockchain for chain %s", chain.Family, chainIDStr)
		}
		blockchainsOut = append(blockchainsOut, bc)
	}

	if evmChains == 0 {
		return nil, errors.New("no EVM chains declared in desired state")
	}

	contractVersions, err := cre.ContractVersionsProviderFromDataStore(cldfEnv.DataStore)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve contract versions")
	}

	return &cre.Environment{
		CldfEnvironment:       cldfEnv,
		Blockchains:           blockchainsOut,
		RegistryChainSelector: registryChainSelector,
		Provider:              provider,
		ContractVersions:      contractVersions.ContractVersions(),
	}, nil
}

func deployerAddress(deployerKey string) (common.Address, error) {
	privateKeyHex := strings.TrimPrefix(deployerKey, "0x")
	privateKey, err := ethcrypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to parse deployer private key")
	}

	return ethcrypto.PubkeyToAddress(privateKey.PublicKey), nil
}
