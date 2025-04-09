package internal

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	DeploymentBlockLabel = "deployment-block"
	DeploymentHashLabel  = "deployment-hash"
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
func (c *KeystoneForwarderDeployer) deploy(ctx context.Context, req DeployRequest) (*DeployResponse, error) {
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
	txHash := tx.Hash()
	txReceipt, err := req.Chain.Client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}
	tv.Labels.Add(fmt.Sprintf("%s: %s", DeploymentHashLabel, txHash.Hex()))
	tv.Labels.Add(fmt.Sprintf("%s: %s", DeploymentBlockLabel, txReceipt.BlockNumber.String()))
	resp := &DeployResponse{
		Address: forwarderAddr,
		Tx:      txHash,
		Tv:      tv,
	}
	c.contract = forwarder
	return resp, nil
}

type ConfigureForwarderContractsRequest struct {
	Dons []RegisteredDon

	Chains  map[uint64]struct{} // list of chains for which request will be executed. If empty, request is applied to all chains
	UseMCMS bool
}
type ConfigureForwarderContractsResponse struct {
	OpsPerChain map[uint64]mcmstypes.BatchOperation
}

// Depreciated: use [changeset.ConfigureForwardContracts] instead
// ConfigureForwardContracts configures the forwarder contracts on all chains for the given DONS
// the address book is required to contain the an address of the deployed forwarder contract for every chain in the environment
func ConfigureForwardContracts(env *deployment.Environment, req ConfigureForwarderContractsRequest) (*ConfigureForwarderContractsResponse, error) {
	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

	opPerChain := make(map[uint64]mcmstypes.BatchOperation)
	// configure forwarders on all chains
	for _, chain := range env.Chains {
		if _, shouldInclude := req.Chains[chain.Selector]; len(req.Chains) > 0 && !shouldInclude {
			continue
		}
		// get the forwarder contract for the chain
		contracts, ok := contractSetsResp.ContractSets[chain.Selector]
		if !ok {
			return nil, fmt.Errorf("failed to get contract set for chain %d", chain.Selector)
		}
		ops, err := configureForwarder(env.Logger, chain, contracts.Forwarder, req.Dons, req.UseMCMS)
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
