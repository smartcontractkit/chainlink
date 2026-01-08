package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/shard_config/v1/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[TransferOwnershipInput] = TransferOwnership{}

// TransferOwnershipInput contains the input parameters for the TransferOwnership changeset.
type TransferOwnershipInput struct {
	ChainSelector  uint64                  `json:"chainSelector"`
	NewOwner       common.Address          `json:"newOwner"`
	ShardConfigRef datastore.AddressRefKey `json:"shardConfigRef"`
}

// TransferOwnership is a changeset that transfers ownership of a ShardConfig contract.
type TransferOwnership struct{}

// VerifyPreconditions validates the input parameters before applying the changeset.
func (t TransferOwnership) VerifyPreconditions(e cldf.Environment, config TransferOwnershipInput) error {
	if config.ChainSelector == 0 {
		return errors.New("chain selector must be provided")
	}
	_, err := chain_selectors.GetChainIDFromSelector(config.ChainSelector)
	if err != nil {
		return fmt.Errorf("invalid chain selector %d: %w", config.ChainSelector, err)
	}
	if config.NewOwner == (common.Address{}) {
		return errors.New("new owner address must be provided")
	}

	// Verify contract exists in datastore
	addr, err := e.DataStore.Addresses().Get(config.ShardConfigRef)
	if err != nil {
		return fmt.Errorf("shard config not found: %w", err)
	}
	if addr.ChainSelector != config.ChainSelector {
		return errors.New("chain selector mismatch with contract reference")
	}
	return nil
}

// Apply executes the changeset to transfer ownership.
func (t TransferOwnership) Apply(e cldf.Environment, config TransferOwnershipInput) (cldf.ChangesetOutput, error) {
	addr, err := e.DataStore.Addresses().Get(config.ShardConfigRef)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	deps := contracts.TransferOwnershipOpDeps{Env: &e}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.TransferOwnershipOp,
		deps,
		contracts.TransferOwnershipOpInput{
			ChainSelector: config.ChainSelector,
			ContractAddr:  addr.Address,
			NewOwner:      config.NewOwner,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
