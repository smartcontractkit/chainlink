package internal

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/balance_reader"
	"github.com/smartcontractkit/chainlink/deployment"
)

type BalanceReaderDeployer struct {
	lggr     logger.Logger
	contract *balance_reader.BalanceReader
}

func NewBalanceReaderDeployer() (*BalanceReaderDeployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &BalanceReaderDeployer{lggr: lggr}, nil
}
func (c *BalanceReaderDeployer) deploy(req DeployRequest) (*DeployResponse, error) {
	est, err := estimateDeploymentGas(req.Chain.Client, balance_reader.BalanceReaderABI)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}
	c.lggr.Debugf("BalanceReader estimated gas: %d", est)

	balReaderAddr, tx, balReader, err := balance_reader.DeployBalanceReader(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy BalanceReader: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm and save BalanceReader: %w", err)
	}
	tvStr, err := balReader.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}
	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}
	resp := &DeployResponse{
		Address: balReaderAddr,
		Tx:      tx.Hash(),
		Tv:      tv,
	}
	c.contract = balReader
	return resp, nil
}
