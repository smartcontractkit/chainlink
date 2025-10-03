package blockchains

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	pkgerrors "github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana/provider"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

type Deployer interface {
	Deploy(input *blockchain.Input) (*DeployedBLockchain, error)
}

type DeployedBLockchain struct {
	CldfBlockchain chain.BlockChain
	Blockchain     *cre.WrappedBlockchainOutput
}

type ChainFamily string

// TODO move to the CTF
func typeToFamily(t string) (ChainFamily, error) {
	switch t {
	case blockchain.TypeAnvil, blockchain.TypeGeth, blockchain.TypeBesu, blockchain.TypeAnvilZKSync:
		return ChainFamily(blockchain.FamilyEVM), nil
	case blockchain.TypeSolana:
		return ChainFamily(blockchain.FamilySolana), nil
	case blockchain.TypeAptos:
		return ChainFamily(blockchain.FamilyAptos), nil
	case blockchain.TypeSui:
		return ChainFamily(blockchain.FamilySui), nil
	case blockchain.TypeTron:
		return ChainFamily(blockchain.FamilyTron), nil
	case blockchain.TypeTon:
		return ChainFamily(blockchain.FamilyTon), nil
	default:
		return "", fmt.Errorf("blockchain type is not supported or empty: %s", t)
	}
}

type DeployedBlockchains struct {
	Outputs         []*cre.WrappedBlockchainOutput
	CldfBlockChains map[uint64]chain.BlockChain
}

func (s *DeployedBlockchains) RegistryChain() *cre.WrappedBlockchainOutput {
	return s.Outputs[0]
}

func Start(
	inputs []*blockchain.Input,
	deployers map[ChainFamily]Deployer,
) (*DeployedBlockchains, error) {
	deployedBlockchains := make([]*DeployedBLockchain, 0, len(inputs))

	for _, input := range inputs {
		chainFamily, chErr := typeToFamily(input.Type)
		if chErr != nil {
			return nil, chErr
		}

		deployer, ok := deployers[chainFamily]
		if !ok {
			return nil, fmt.Errorf("no deployer found for blockchain type %s", input.Type)
		}

		deployedBlockchain, deployErr := deployer.Deploy(input)
		if deployErr != nil {
			return nil, pkgerrors.Wrapf(deployErr, "failed to deploy blockchain of type %s", input.Type)
		}

		deployedBlockchains = append(deployedBlockchains, deployedBlockchain)
	}

	outputs := make([]*cre.WrappedBlockchainOutput, 0, len(deployedBlockchains))
	for _, db := range deployedBlockchains {
		outputs = append(outputs, db.Blockchain)
	}
	cldfBlockchains := make(map[uint64]chain.BlockChain)
	for _, db := range deployedBlockchains {
		cldfBlockchains[db.Blockchain.ChainSelector] = db.CldfBlockchain
	}

	return &DeployedBlockchains{
		Outputs:         outputs,
		CldfBlockChains: cldfBlockchains,
	}, nil
}

func WrapEVM(bcOut *blockchain.Output) (*cre.WrappedBlockchainOutput, error) {
	if err := SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); err != nil {
		return nil, err
	}

	priv := os.Getenv("PRIVATE_KEY")
	sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(bcOut.Nodes[0].ExternalWSUrl).
		WithPrivateKeys([]string{priv}).
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create seth client")
	}

	selector, err := chainselectors.SelectorFromChainId(sethClient.Cfg.Network.ChainID)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %d", sethClient.Cfg.Network.ChainID)
	}

	chainID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bcOut.ChainID)
	}

	return &cre.WrappedBlockchainOutput{
		ChainSelector:      selector,
		ChainID:            chainID,
		BlockchainOutput:   bcOut,
		SethClient:         sethClient,
		DeployerPrivateKey: priv,
	}, nil
}

func WrapTron(bi *blockchain.Input, bcOut *blockchain.Output) (*cre.WrappedBlockchainOutput, error) {
	chainID, err := strconv.ParseUint(bi.ChainID, 10, 64)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bi.ChainID)
	}
	selector, err := chainselectors.SelectorFromChainId(chainID)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %s", bi.ChainID)
	}

	// if jsonrpc is not present, add it
	if !strings.HasSuffix(bcOut.Nodes[0].ExternalHTTPUrl, "/jsonrpc") {
		bcOut.Nodes[0].ExternalHTTPUrl += "/jsonrpc"
	}
	if !strings.HasSuffix(bcOut.Nodes[0].InternalHTTPUrl, "/jsonrpc") {
		bcOut.Nodes[0].InternalHTTPUrl += "/jsonrpc"
	}

	externalHTTPURL := bcOut.Nodes[0].ExternalHTTPUrl
	internalHTTPURL := bcOut.Nodes[0].InternalHTTPUrl

	return &cre.WrappedBlockchainOutput{
		ChainSelector: selector,
		ChainID:       chainID,
		BlockchainOutput: &blockchain.Output{
			ChainID: bi.ChainID,
			Family:  blockchain.FamilyTron,
			Nodes: []*blockchain.Node{
				{
					InternalHTTPUrl: internalHTTPURL,
					ExternalHTTPUrl: externalHTTPURL,
				},
			},
		},
		SethClient:         nil,
		DeployerPrivateKey: blockchain.TRONAccounts.PrivateKeys[0],
	}, nil
}

// Will be set as --mint when spin up local solana validator, unless env variable with a different key provided
var DefaultSolanaPrivateKey = solana.MustPrivateKeyFromBase58("4u2itaM9r5kxsmoti3GMSDZrQEFpX14o6qPWY9ZrrYTR6kduDBr4YAZJsjawKzGP3wDzyXqterFmfcLUmSBro5AT")

func WrapSolana(bi *blockchain.Input, bcOut *blockchain.Output) (*cre.WrappedBlockchainOutput, error) {
	sel, ok := chainselectors.SolanaChainIdToChainSelector()[bi.ChainID]
	if !ok {
		return nil, fmt.Errorf("selector not found for solana chainID '%s'", bi.ChainID)
	}
	// shouldn't be empty, since we call initSolana before wrap, but just in case
	setErr := SetDefaultSolanaPrivateKeyIfEmpty(DefaultSolanaPrivateKey)
	if setErr != nil {
		return nil, fmt.Errorf("set default private key solana failed: %w", setErr)
	}

	envp := os.Getenv("SOLANA_PRIVATE_KEY")
	pk, err := solana.PrivateKeyFromBase58(envp)
	if err != nil {
		return nil, errors.New("failed to decode private key for solana")
	}

	if err := cldf_solana_provider.WritePrivateKeyToPath(filepath.Join(bi.ContractsDir, "deploy-keypair.json"), pk); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to save private key for solana")
	}

	return &cre.WrappedBlockchainOutput{
		BlockchainOutput: bcOut,
		SolClient:        rpc.New(bcOut.Nodes[0].ExternalHTTPUrl),
		SolChain: &cre.SolChain{
			ChainSelector: sel,
			ChainID:       bi.ChainID,
			PrivateKey:    pk,
			ArtifactsDir:  bi.ContractsDir,
		},
	}, nil
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
