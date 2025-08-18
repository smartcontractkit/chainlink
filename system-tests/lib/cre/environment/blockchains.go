package environment

import (
	"context"
	"maps"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	tronprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	libnix "github.com/smartcontractkit/chainlink/system-tests/lib/nix"
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
	ctx context.Context,
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
		} else if bi.Type == "tron" {
			chainID, err := strconv.ParseUint(bi.ChainID, 10, 64)
			if err != nil {
				return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bi.ChainID)
			}
			selector, err := chainselectors.SelectorFromChainId(chainID)
			if err != nil {
				return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %d", bi.ChainID)
			}
			signerGen, err := tronprovider.SignerGenCTFDefault()
			if err != nil {
				return nil, pkgerrors.Wrap(err, "failed to create signer generator")
			}
			config := tronprovider.CTFChainProviderConfig{
				DeployerSignerGen: signerGen,
				Once:              &sync.Once{},
			}
			ctfProvider := tronprovider.NewCTFChainProvider(selector, config)
			_, err = ctfProvider.Initialize(ctx)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "failed to initialize blockchain")
			}
			chain := ctfProvider.BlockChain()
			tronChain, ok := chain.(tron.Chain)
			if !ok {
				return nil, pkgerrors.New("failed to cast chain to tron.Chain")
			}
			// For Tron chains, we need to use host.docker.internal to access the host's localhost
			// The external URL is typically like "http://localhost:50555/wallet"
			// We need to convert it to use host.docker.internal for internal connections
			externalHTTPURL := strings.Replace(tronChain.URL, "wallet", "jsonrpc", 1)
			internalHTTPURL := ""

			// Extract the port from the external URL dynamically
			if strings.Contains(externalHTTPURL, "localhost:") {
				// Use host.docker.internal to access the host's localhost from within the container
				// This works regardless of the port number
				internalHTTPURL = strings.Replace(externalHTTPURL, "localhost", "host.docker.internal", 1)
			} else {
				// Fallback to external URL if we can't determine the internal URL
				internalHTTPURL = externalHTTPURL
			}

			bcOut = &blockchain.Output{
				ChainID: bi.ChainID,
				Family:  "tron",
				Nodes: []*blockchain.Node{
					{
						InternalHTTPUrl: internalHTTPURL,
						ExternalHTTPUrl: externalHTTPURL,
					},
				},
			}
			blockchainOutput = append(blockchainOutput, &cre.WrappedBlockchainOutput{
				ChainSelector:      selector,
				ChainID:            chainID,
				BlockchainOutput:   bcOut,
				SethClient:         nil,
				DeployerPrivateKey: blockchain.TRONAccounts.PrivateKeys[0],
				TronChain:          &tronChain,
			})
			continue
		} else {
			bcOut, bcErr = blockchain.NewBlockchainNetwork(&bi)
			if bcErr != nil {
				return nil, pkgerrors.Wrap(bcErr, "failed to deploy blockchain")
			}
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

		chainID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bcOut.ChainID)
		}

		chainSelector, err := chainselectors.SelectorFromChainId(chainID)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %d", chainID)
		}

		sethClient, err = seth.NewClientBuilder().
			WithRpcUrl(bcOut.Nodes[0].ExternalWSUrl).
			WithPrivateKeys([]string{privateKey}).
			// do not check if there's a pending nonce nor check node's health
			WithProtections(false, false, seth.MustMakeDuration(time.Second)).
			Build()
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to create seth client")
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

func StartBlockchains(ctx context.Context, loggers BlockchainLoggers, input BlockchainsInput) (StartBlockchainsOutput, error) {
	blockchainsOutput, err := CreateBlockchains(ctx, loggers.lggr, input)
	if err != nil {
		return StartBlockchainsOutput{}, pkgerrors.Wrap(err, "failed to create blockchains")
	}

	chainsConfigs := make([]devenv.ChainConfig, 0)

	for _, bcOut := range blockchainsOutput {
		if bcOut.TronChain != nil {
			privateKey, err := crypto.HexToECDSA(bcOut.DeployerPrivateKey)
			if err != nil {
				return StartBlockchainsOutput{}, pkgerrors.Wrap(err, "failed to parse private key for Tron")
			}

			chainID, err := strconv.ParseInt(bcOut.BlockchainOutput.ChainID, 10, 64)
			if err != nil {
				return StartBlockchainsOutput{}, pkgerrors.Wrapf(err, "failed to parse Tron chain ID %s", bcOut.BlockchainOutput.ChainID)
			}

			deployerKey, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
			if err != nil {
				return StartBlockchainsOutput{}, pkgerrors.Wrap(err, "failed to create transactor for Tron")
			}

			chainsConfigs = append(chainsConfigs, devenv.ChainConfig{
				ChainID:   bcOut.BlockchainOutput.ChainID,
				ChainName: "Tron",
				ChainType: "EVM",
				WSRPCs:    []devenv.CribRPCs{{}},
				HTTPRPCs: []devenv.CribRPCs{{
					External: bcOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
					Internal: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
				}},
				DeployerKey:        deployerKey,
				PreferredURLScheme: deployment.URLSchemePreferenceHTTP,
			})
			continue
		}

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
