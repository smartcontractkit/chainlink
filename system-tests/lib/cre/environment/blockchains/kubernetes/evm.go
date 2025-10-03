package kubernetes

import (
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type EVMDeployer struct {
	namespace      string
	cribConfigsDir string
	testLogger     zerolog.Logger
}

func NewEVMDeployer(testLogger zerolog.Logger, namespace string, cribConfigsDir string) *EVMDeployer {
	return &EVMDeployer{
		namespace:      namespace,
		cribConfigsDir: cribConfigsDir,
		testLogger:     testLogger,
	}
}

func (e *EVMDeployer) Deploy(input *blockchain.Input) (*cre.Blockchain, error) {
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

	return creblockchains.WrapEVM(bcOut)
}
