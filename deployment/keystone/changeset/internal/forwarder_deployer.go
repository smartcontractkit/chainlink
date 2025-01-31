package internal

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	forwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
)

type KeystoneForwarderDeployer struct {
	lggr     logger.Logger
	contract *forwarder.KeystoneForwarder
}

func NewKeystoneForwarderDeployer() (*KeystoneForwarderDeployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderDeployer{lggr: lggr}, nil
}
func (c *KeystoneForwarderDeployer) deploy(req DeployRequest) (*DeployResponse, error) {
	est, err := estimateDeploymentGas(req.Chain.Client, forwarder.KeystoneForwarderABI)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}
	c.lggr.Debugf("Forwarder estimated gas: %d", est)

	forwarderAddr, tx, forwarder, err := forwarder.DeployKeystoneForwarder(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm and save KeystoneForwarder: %w", err)
	}
	tvStr, err := forwarder.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}
	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}
	resp := &DeployResponse{
		Address: forwarderAddr,
		Tx:      tx.Hash(),
		Tv:      tv,
	}
	c.contract = forwarder
	return resp, nil
}

type ConfigureForwarderContractsRequest struct {
	Dons []RegisteredDon

	ForwarderContracts map[uint64]*forwarder.KeystoneForwarder
	UseMCMS            bool
}

type ConfigureForwarderContractsResponse struct {
	OpsPerChain map[uint64]timelock.BatchChainOperation
}

// Depreciated: use [changeset.ConfigureForwarders] instead
// ConfigureForwarderContracts configures the forwarder contracts on all chains for the given DONS
// the address book is required to contain the an address of the deployed forwarder contract for every chain in the environment
func ConfigureForwarderContracts(env *deployment.Environment, req ConfigureForwarderContractsRequest) (*ConfigureForwarderContractsResponse, error) {
	opPerChain := make(map[uint64]timelock.BatchChainOperation)
	// configure forwarders on all chains
	for sel, contract := range req.ForwarderContracts {
		chain, ok := env.Chains[sel]
		if !ok {
			return nil, fmt.Errorf("failed to get chain for selector %d", sel)
		}
		ops, err := configureForwarder(env.Logger, chain, contract, req.Dons, req.UseMCMS)
		if err != nil {
			return nil, fmt.Errorf("failed to configure forwarder for chain selector %d: %w", chain.Selector, err)
		}
		for k, op := range ops {
			opPerChain[k] = op
		}
	}

	return &ConfigureForwarderContractsResponse{
		OpsPerChain: opPerChain,
	}, nil
}
