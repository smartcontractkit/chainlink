package changeset

import (
	"errors"
	"fmt"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/shard_config/v1/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[UpdateShardCountInput] = UpdateShardCount{}

// UpdateShardCountInput contains the input parameters for the UpdateShardCount changeset.
type UpdateShardCountInput struct {
	ChainSelector  uint64                  `json:"chainSelector"`
	NewShardCount  uint64                  `json:"newShardCount"`
	ShardConfigRef datastore.AddressRefKey `json:"shardConfigRef"`
}

// UpdateShardCount is a changeset that updates the desired shard count on a ShardConfig contract.
type UpdateShardCount struct{}

// VerifyPreconditions validates the input parameters before applying the changeset.
func (u UpdateShardCount) VerifyPreconditions(e cldf.Environment, config UpdateShardCountInput) error {
	if config.ChainSelector == 0 {
		return errors.New("chain selector must be provided")
	}
	_, err := chain_selectors.GetChainIDFromSelector(config.ChainSelector)
	if err != nil {
		return fmt.Errorf("invalid chain selector %d: %w", config.ChainSelector, err)
	}
	if config.NewShardCount == 0 {
		return errors.New("new shard count must be greater than 0")
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

// Apply executes the changeset to update the shard count.
func (u UpdateShardCount) Apply(e cldf.Environment, config UpdateShardCountInput) (cldf.ChangesetOutput, error) {
	addr, err := e.DataStore.Addresses().Get(config.ShardConfigRef)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	deps := contracts.UpdateShardCountOpDeps{Env: &e}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.UpdateShardCountOp,
		deps,
		contracts.UpdateShardCountOpInput{
			ChainSelector: config.ChainSelector,
			ContractAddr:  addr.Address,
			NewShardCount: config.NewShardCount,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
