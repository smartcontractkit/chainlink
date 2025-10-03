package kubernetes

import (
	"maps"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type EVMDeployer struct {
	namespace      string
	cribConfigsDir string
	testLogger     zerolog.Logger
	commonLogger   logger.Logger
}

func NewEVMDeployer(testLogger zerolog.Logger, commonLogger logger.Logger, namespace string, cribConfigsDir string) *EVMDeployer {
	return &EVMDeployer{
		namespace:      namespace,
		cribConfigsDir: cribConfigsDir,
		testLogger:     testLogger,
		commonLogger:   commonLogger,
	}
}

func (e *EVMDeployer) Deploy(input *blockchain.Input) (*creblockchains.DeployedBLockchain, error) {
	deployCribBlockchainInput := &cre.DeployCribBlockchainInput{
		BlockchainInput: input,
		CribConfigsDir:  e.cribConfigsDir,
		Namespace:       e.namespace,
	}

	bcOut, err := crib.DeployBlockchain(deployCribBlockchainInput)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to deploy blockchain")
	}

	err = infra.WaitForRPCEndpoint(e.testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
	}

	w, wrapErr := creblockchains.WrapEVM(bcOut)
	if wrapErr != nil {
		return nil, pkgerrors.Wrap(wrapErr, "failed to wrap EVM")
	}

	chainsConfigs := make([]devenv.ChainConfig, 0)
	cfg, cfgErr := cre.ChainConfigFromWrapped(w)
	if cfgErr != nil {
		return nil, pkgerrors.Wrap(cfgErr, "failed to wrap blockchain output to chain config")
	}
	chainsConfigs = append(chainsConfigs, cfg)

	cldfBlockchain, err := devenv.NewChains(e.commonLogger, chainsConfigs)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create chains")
	}

	return &creblockchains.DeployedBLockchain{
		CldfBlockchain: maps.Collect(cldfBlockchain.All())[w.ChainSelector],
		Blockchain:     w,
	}, nil
}
