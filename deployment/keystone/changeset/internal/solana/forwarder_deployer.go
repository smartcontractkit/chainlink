package solana

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/wsrpc/logger"
)

type KeystoneForwarderDeployer struct {
	lggr        logger.Logger
	programPath string
}

func NewKeystoneForwarderDeployer() (*KeystoneForwarderDeployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderDeployer{lggr: lggr}, nil
}

func (c *KeystoneForwarderDeployer) deploy(ctx context.Context, req DeployRequest) (*DeployResponse, error) {
	pubKey, err := req.Chain.DeployProgram(c.lggr, req.SolProgramInfo, false, true)
	if err != nil {
		return nil, err
	}

	key, err := solana.PublicKeyFromBase58(pubKey)
	if err != nil {
		return nil, err
	}

	// TODO PLEX-354 mock for now, implement deploy logic here
	// TODO add labels with deployment block/hash
	return &DeployResponse{
		Address: key,
		Tx:      solana.Signature{},
		Tv:      deployment.MustTypeAndVersionFromString("KeystoneForwarder 1.0.0"),
	}, nil
}
