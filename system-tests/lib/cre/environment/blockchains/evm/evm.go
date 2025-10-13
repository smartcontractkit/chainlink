package evm

import (
	"fmt"
	"os"
	"strconv"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type Deployer struct {
	provider       infra.Provider
	testLogger     zerolog.Logger
	cribConfigsDir string
}

func NewDeployer(testLogger zerolog.Logger, provider *infra.Provider, cribConfigsDir string) *Deployer {
	return &Deployer{
		provider:       *provider,
		testLogger:     testLogger,
		cribConfigsDir: cribConfigsDir,
	}
}

func (e *Deployer) Deploy(input *blockchain.Input) (*cre.Blockchain, error) {
	var bcOut *blockchain.Output
	var err error

	if e.provider.IsCRIB() {
		deployCribBlockchainInput := &cre.DeployCribBlockchainInput{
			BlockchainInput: input,
			CribConfigsDir:  e.cribConfigsDir,
			Namespace:       e.provider.CRIB.Namespace,
		}

		bcOut, err := crib.DeployBlockchain(deployCribBlockchainInput)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to deploy blockchain")
		}

		err = infra.WaitForRPCEndpoint(e.testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
		}
	} else {
		bcOut, err = blockchain.NewBlockchainNetwork(input)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
		}
	}

	if err := setDefaultPrivateKeyIfEmpty(); err != nil {
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

	return &cre.Blockchain{
		ChainSelector:      selector,
		ChainID:            chainID,
		CtfOutput:          bcOut,
		SethClient:         sethClient,
		DeployerPrivateKey: priv,
	}, nil
}

func setDefaultPrivateKeyIfEmpty() error {
	if os.Getenv("PRIVATE_KEY") == "" {
		setErr := os.Setenv("PRIVATE_KEY", blockchain.DefaultAnvilPrivateKey)
		if setErr != nil {
			return fmt.Errorf("failed to set PRIVATE_KEY environment variable: %w", setErr)
		}
		framework.L.Info().Msgf("Set PRIVATE_KEY environment variable to default value: %s", os.Getenv("PRIVATE_KEY"))
	}

	return nil
}
