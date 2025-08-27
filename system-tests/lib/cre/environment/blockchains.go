package environment

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	libnix "github.com/smartcontractkit/chainlink/system-tests/lib/nix"

	cldf_solana_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana/provider"
)

type BlockchainsInput struct {
	blockchainsInput []blockchain.Input
	infra            *infra.Input
	nixShell         *libnix.Shell
}

type BlockchainOutput struct {
	ChainSelector      uint64
	ChainID            uint64
	BlockchainOutput   *blockchain.Output
	SethClient         *seth.Client
	DeployerPrivateKey string
}

func CreateBlockchains(
	testLogger zerolog.Logger,
	input BlockchainsInput,
) ([]*cre.WrappedBlockchainOutput, error) {
	if len(input.blockchainsInput) == 0 {
		return nil, pkgerrors.New("blockchain input is nil")
	}

	blockchainOutput := make([]*cre.WrappedBlockchainOutput, 0)

	for _, bi := range input.blockchainsInput {
		isSolana := bi.Type == blockchain.FamilySolana

		if isSolana {
			err := initSolanaInput(&bi)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "failed to init Solana input")
			}
		}

		bcOut, err := deployBlockchain(testLogger, input.infra, input.nixShell, bi)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s", bi.Type)
		}

		if isSolana {
			w, wrapErr := wrapSolana(&bi, bcOut)
			if wrapErr != nil {
				return nil, pkgerrors.Wrap(wrapErr, "failed to wrap Solana")
			}

			blockchainOutput = append(blockchainOutput, w)
			continue
		}

		w, wrapErr := wrapEVM(bcOut)
		if wrapErr != nil {
			return nil, pkgerrors.Wrap(wrapErr, "failed to wrap EVM")
		}

		blockchainOutput = append(blockchainOutput, w)
	}
	return blockchainOutput, nil
}

// Will be set as --mint when spin up local solana validator, unless env variable with a different key provided
var defaultSolanaPrivateKey = solana.MustPrivateKeyFromBase58("4u2itaM9r5kxsmoti3GMSDZrQEFpX14o6qPWY9ZrrYTR6kduDBr4YAZJsjawKzGP3wDzyXqterFmfcLUmSBro5AT")

func initSolanaInput(bi *blockchain.Input) error {
	err := SetDefaultSolanaPrivateKeyIfEmpty(defaultSolanaPrivateKey)
	if err != nil {
		return errors.New("failed to set default solana private key")
	}
	bi.PublicKey = defaultSolanaPrivateKey.PublicKey().String()
	bi.ContractsDir = getSolProgramsPath(bi.ContractsDir)
	return nil
}

func deployBlockchain(testLogger zerolog.Logger, infraIn *infra.Input, nixShell *libnix.Shell, bi blockchain.Input) (*blockchain.Output, error) {
	if infraIn.Type != infra.CRIB {
		bcOut, err := blockchain.NewBlockchainNetwork(&bi)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", bi.Type, bi.ChainID)
		}

		return bcOut, nil
	}

	if nixShell == nil {
		return nil, pkgerrors.New("nix shell is nil")
	}

	deployCribBlockchainInput := &cre.DeployCribBlockchainInput{
		BlockchainInput: &bi,
		NixShell:        nixShell,
		CribConfigsDir:  cribConfigsDir,
		Namespace:       infraIn.CRIB.Namespace,
	}
	bcOut, err := crib.DeployBlockchain(deployCribBlockchainInput)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to deploy blockchain")
	}

	err = infra.WaitForRPCEndpoint(testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
	}

	return bcOut, nil
}

func wrapSolana(bi *blockchain.Input, bcOut *blockchain.Output) (*cre.WrappedBlockchainOutput, error) {
	sel, ok := chainselectors.SolanaChainIdToChainSelector()[bi.ChainID]
	if !ok {
		return nil, fmt.Errorf("selector not found for solana chainID '%s'", bi.ChainID)
	}
	// shouldn't be empty, since we call initSolana before wrap, but just in case
	setErr := SetDefaultSolanaPrivateKeyIfEmpty(defaultSolanaPrivateKey)
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

func wrapEVM(bcOut *blockchain.Output) (*cre.WrappedBlockchainOutput, error) {
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

type BlockchainLoggers struct {
	lggr       zerolog.Logger
	singleFile logger.Logger
}

type StartBlockchainsOutput struct {
	BlockChainOutputs []*cre.WrappedBlockchainOutput
	BlockChains       map[uint64]chain.BlockChain
}

func StartBlockchains(loggers BlockchainLoggers, input BlockchainsInput) (StartBlockchainsOutput, error) {
	blockchainsOutput, err := CreateBlockchains(loggers.lggr, input)
	if err != nil {
		return StartBlockchainsOutput{}, pkgerrors.Wrap(err, "failed to create blockchains")
	}

	chainsConfigs := make([]devenv.ChainConfig, 0)
	for _, bcOut := range blockchainsOutput {
		cfg, cfgErr := cre.ChainConfigFromWrapped(bcOut)
		if cfgErr != nil {
			return StartBlockchainsOutput{}, pkgerrors.Wrap(err, "failed to wrap blockchain output to chain config")
		}
		chainsConfigs = append(chainsConfigs, cfg)
	}

	blockChains, err := devenv.NewChains(loggers.singleFile, chainsConfigs)
	if err != nil {
		return StartBlockchainsOutput{}, pkgerrors.Wrap(err, "failed to create chains")
	}

	return StartBlockchainsOutput{
		BlockChainOutputs: blockchainsOutput,
		BlockChains:       maps.Collect(blockChains.All()),
	}, nil
}

func getSolProgramsPath(path string) string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, path)
}
