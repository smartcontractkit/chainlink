package environment

import (
	"os"
	"strconv"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libinfra "github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	libnix "github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

type BlockchainsInput struct {
	blockchainsInput []*cretypes.WrappedBlockchainInput
	infra            *libtypes.InfraInput
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
) ([]*cretypes.WrappedBlockchainOutput, error) {
	if len(input.blockchainsInput) == 0 {
		return nil, pkgerrors.New("blockchain input is nil")
	}

	blockchainOutput := make([]*cretypes.WrappedBlockchainOutput, 0)
	for _, bi := range input.blockchainsInput {
		var bcOut *blockchain.Output
		var bcErr error
		if input.infra.InfraType == libtypes.CRIB {
			if input.nixShell == nil {
				return nil, pkgerrors.New("nix shell is nil")
			}

			deployCribBlockchainInput := &cretypes.DeployCribBlockchainInput{
				BlockchainInput: &bi.Input,
				NixShell:        input.nixShell,
				CribConfigsDir:  cribConfigsDir,
			}
			bcOut, bcErr = crib.DeployBlockchain(deployCribBlockchainInput)
			if bcErr != nil {
				return nil, pkgerrors.Wrap(bcErr, "failed to deploy blockchain")
			}
			err := libinfra.WaitForRPCEndpoint(testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
			}
		} else {
			bcOut, bcErr = blockchain.NewBlockchainNetwork(&bi.Input)
			if bcErr != nil {
				return nil, pkgerrors.Wrap(bcErr, "failed to deploy blockchain")
			}
		}

		pkey := os.Getenv("PRIVATE_KEY")
		if pkey == "" {
			return nil, pkgerrors.New("PRIVATE_KEY env var must be set")
		}

		sethClient, err := seth.NewClientBuilder().
			WithRpcUrl(bcOut.Nodes[0].ExternalWSUrl).
			WithPrivateKeys([]string{pkey}).
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

		blockchainOutput = append(blockchainOutput, &cretypes.WrappedBlockchainOutput{
			ChainSelector:      chainSelector,
			ChainID:            chainID,
			BlockchainOutput:   bcOut,
			SethClient:         sethClient,
			DeployerPrivateKey: pkey,
			ReadOnly:           bi.ReadOnly,
		})
	}

	return blockchainOutput, nil
}

type BlockchainLoggers struct {
	lggr       zerolog.Logger
	singleFile logger.Logger
}

type StartBlockchainsOutput struct {
	BlockChainOutputs []*cretypes.WrappedBlockchainOutput
	BlockChains       map[uint64]chain.BlockChain
}

func StartBlockchains(loggers BlockchainLoggers, input BlockchainsInput) (StartBlockchainsOutput, error) {
	blockchainsOutput, err := CreateBlockchains(loggers.lggr, input)
	if err != nil {
		return StartBlockchainsOutput{}, pkgerrors.Wrap(err, "failed to create blockchains")
	}

	chainsConfigs := make([]devenv.ChainConfig, 0)

	for _, bcOut := range blockchainsOutput {
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

	allChains, _, err := devenv.NewChains(loggers.singleFile, chainsConfigs)
	if err != nil {
		return StartBlockchainsOutput{}, pkgerrors.Wrap(err, "failed to create chains")
	}

	blockChains := make(map[uint64]chain.BlockChain, len(allChains))
	for selector, ch := range allChains {
		blockChains[selector] = ch
	}

	return StartBlockchainsOutput{
		BlockChainOutputs: blockchainsOutput,
		BlockChains:       blockChains,
	}, nil
}
