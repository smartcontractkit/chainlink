package solana

import (
	"context"

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
	return nil, nil
}
