package environment

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
		var bcOut *blockchain.Output
		var bcErr error
		var solPrivateKey solana.PrivateKey // lazy private key, will be initialized only if Solana is present
		isSolana := bi.Type == blockchain.FamilySolana

		if isSolana {
			solPrivateKey, _ = solana.NewRandomPrivateKey()
			bi.PublicKey = solPrivateKey.PublicKey().String()
			bi.ContractsDir = getSolProgramsPath(bi.ContractsDir)
		}

		if input.infra.Type == infra.CRIB {
			if input.nixShell == nil {
				return nil, pkgerrors.New("nix shell is nil")
			}

			deployCribBlockchainInput := &cre.DeployCribBlockchainInput{
				BlockchainInput: &bi,
				NixShell:        input.nixShell,
				CribConfigsDir:  cribConfigsDir,
				Namespace:       input.infra.CRIB.Namespace,
			}
			bcOut, bcErr = crib.DeployBlockchain(deployCribBlockchainInput)
			if bcErr != nil {
				return nil, pkgerrors.Wrap(bcErr, "failed to deploy blockchain")
			}
			err := infra.WaitForRPCEndpoint(testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
			}
		} else {
			bcOut, bcErr = blockchain.NewBlockchainNetwork(&bi)
			if bcErr != nil {
				return nil, pkgerrors.Wrap(bcErr, "failed to deploy blockchain")
			}
		}
		// handle solana here
		if isSolana {
			solClient := rpc.New(bcOut.Nodes[0].ExternalHTTPUrl)

			// we pass selector from input, because local solana chainID is unpredictable
			selector, ok := chainselectors.SolanaChainIdToChainSelector()[bi.ChainID]
			if !ok {
				return nil, pkgerrors.Errorf("selector not found for solana chainID '%s'", bi.ChainID)
			}

			err := cldf_solana_provider.WritePrivateKeyToPath(filepath.Join(bi.ContractsDir, "deploy-keypair.json"), solPrivateKey)
			if err != nil {
				return nil, fmt.Errorf("failed to save private key for solana: %w", err)
			}

			blockchainOutput = append(blockchainOutput, &cre.WrappedBlockchainOutput{
				BlockchainOutput: bcOut,
				SolClient:        solClient,
				SolChain: &cre.SolChain{
					ChainSelector: selector,
					ChainID:       bi.ChainID,
					PrivateKey:    solPrivateKey,
					ArtifactsDir:  bi.ContractsDir,
				},
			})

			continue
		}

		if pkErr := SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
			return nil, pkErr
		}

		privateKey := os.Getenv("PRIVATE_KEY")
		sethClient, err := seth.NewClientBuilder().
			WithRpcUrl(bcOut.Nodes[0].ExternalWSUrl).
			WithPrivateKeys([]string{privateKey}).
			// do not check if there's a pending nonce nor check node's health
			WithProtections(false, false, seth.MustMakeDuration(time.Second)).
			Build()
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to create seth client")
		}

		chainSelector, err := chainselectors.SelectorFromChainId(sethClient.Cfg.Network.ChainID)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %d", sethClient.Cfg.Network.ChainID)
		}
		chainID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bcOut.ChainID)
		}

		blockchainOutput = append(blockchainOutput, &cre.WrappedBlockchainOutput{
			ChainSelector:      chainSelector,
			ChainID:            chainID,
			BlockchainOutput:   bcOut,
			SethClient:         sethClient,
			DeployerPrivateKey: privateKey,
		})
	}
	return blockchainOutput, nil
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
		switch bcOut.BlockchainOutput.Family {
		case chainselectors.FamilyEVM:
			chainsConfigs = append(chainsConfigs, devenv.ChainConfig{
				ChainID:   strconv.FormatUint(bcOut.SethClient.Cfg.Network.ChainID, 10),
				ChainName: bcOut.SethClient.Cfg.Network.Name,
				ChainType: strings.ToUpper(bcOut.BlockchainOutput.Family),
				WSRPCs: []devenv.CribRPCs{{
					External: bcOut.BlockchainOutput.Nodes[0].ExternalWSUrl,
					Internal: bcOut.BlockchainOutput.Nodes[0].InternalWSUrl,
				}},
				HTTPRPCs: []devenv.CribRPCs{{
					External: bcOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
					Internal: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
				}},
				DeployerKey: bcOut.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the RPC node
			})
		case chainselectors.FamilySolana:
			chainsConfigs = append(chainsConfigs, devenv.ChainConfig{
				ChainID:   bcOut.SolChain.ChainID,
				ChainName: bcOut.SolChain.ChainName,
				ChainType: strings.ToUpper(bcOut.BlockchainOutput.Family),
				WSRPCs: []devenv.CribRPCs{{
					External: bcOut.BlockchainOutput.Nodes[0].ExternalWSUrl,
					Internal: bcOut.BlockchainOutput.Nodes[0].InternalWSUrl,
				}},
				HTTPRPCs: []devenv.CribRPCs{{
					External: bcOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
					Internal: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
				}},
				SolDeployerKey: bcOut.SolChain.PrivateKey,
				SolArtifactDir: bcOut.SolChain.ArtifactsDir,
			})
		}
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
