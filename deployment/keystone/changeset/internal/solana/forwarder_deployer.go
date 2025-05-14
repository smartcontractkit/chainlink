package solana

import (
	"context"

	"github.com/gagliardetto/solana-go"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
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
	// TODO PLEX-354 mock for now, implement deploy logic here
	// TODO add labels with deployment block/hash
	tv, err := cldf.TypeAndVersionFromString("CallProxy 1.0.0 SA staging")
	if err != nil {
		return &DeployResponse{}, err
	}
	return &DeployResponse{
		Address: solana.PublicKey{},
		Tx:      solana.Signature{},
		Tv:      tv,
	}, nil
}
