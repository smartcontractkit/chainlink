package changeset

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/shard_config/v1/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[DeployShardConfigInput] = DeployShardConfig{}

// DeployShardConfigInput contains the input parameters for the DeployShardConfig changeset.
type DeployShardConfigInput struct {
	ChainSelector     uint64 `json:"chainSelector"`
	InitialShardCount uint64 `json:"initialShardCount"`
	Qualifier         string `json:"qualifier,omitempty"`
}

// DeployShardConfig is a changeset that deploys a ShardConfig contract.
type DeployShardConfig struct{}

// VerifyPreconditions validates the input parameters before applying the changeset.
func (d DeployShardConfig) VerifyPreconditions(e cldf.Environment, config DeployShardConfigInput) error {
	if config.ChainSelector == 0 {
		return errors.New("chain selector must be provided")
	}
	_, err := chain_selectors.GetChainIDFromSelector(config.ChainSelector)
	if err != nil {
		return fmt.Errorf("invalid chain selector %d: %w", config.ChainSelector, err)
	}
	if config.InitialShardCount == 0 {
		return errors.New("initial shard count must be greater than 0")
	}
	return nil
}

// Apply executes the changeset to deploy a ShardConfig contract.
func (d DeployShardConfig) Apply(e cldf.Environment, config DeployShardConfigInput) (cldf.ChangesetOutput, error) {
	deps := contracts.DeployShardConfigOpDeps{Env: &e}

	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.DeployShardConfigOp,
		deps,
		contracts.DeployShardConfigOpInput{
			ChainSelector:     config.ChainSelector,
			InitialShardCount: config.InitialShardCount,
			Qualifier:         config.Qualifier,
		},
		operations.WithRetry[contracts.DeployShardConfigOpInput, contracts.DeployShardConfigOpDeps](),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	ds := datastore.NewMemoryDataStore()
	version, err := semver.NewVersion(report.Output.Version)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	labels := datastore.NewLabelSet()
	for _, label := range report.Output.Labels {
		labels.Add(label)
	}

	addressRef := datastore.AddressRef{
		ChainSelector: report.Output.ChainSelector,
		Address:       report.Output.Address,
		Type:          datastore.ContractType(report.Output.Type),
		Version:       version,
		Qualifier:     report.Output.Qualifier,
		Labels:        labels,
	}

	if err := ds.Addresses().Add(addressRef); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
