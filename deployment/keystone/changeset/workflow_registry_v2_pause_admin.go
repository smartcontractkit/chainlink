package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// AdminPauseWorkflowRequest pauses a specific workflow
type AdminPauseWorkflowRequest struct {
	ChainSelector uint64
	WorkflowID    [32]byte // Workflow ID to pause
}

// AdminPauseWorkflow changeset to pause a specific workflow
var _ cldf.ChangeSet[AdminPauseWorkflowRequest] = AdminPauseWorkflow

func AdminPauseWorkflow(env cldf.Environment, req AdminPauseWorkflowRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.AdminPauseWorkflow(chain.DeployerKey, req.WorkflowID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to admin pause workflow: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm admin pause workflow transaction: %w", err)
	}

	env.Logger.Infof("Successfully admin paused workflow %x on chain %d, tx: %s", req.WorkflowID, req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// AdminBatchPauseWorkflowsRequest pauses multiple workflows in a batch
type AdminBatchPauseWorkflowsRequest struct {
	ChainSelector uint64
	WorkflowIDs   [][32]byte // List of workflow IDs to pause
}

// AdminBatchPauseWorkflows changeset to pause multiple workflows in a batch
var _ cldf.ChangeSet[AdminBatchPauseWorkflowsRequest] = AdminBatchPauseWorkflows

func AdminBatchPauseWorkflows(env cldf.Environment, req AdminBatchPauseWorkflowsRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.AdminBatchPauseWorkflows(chain.DeployerKey, req.WorkflowIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to admin batch pause workflows: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm admin batch pause workflows transaction: %w", err)
	}

	env.Logger.Infof("Successfully admin batch paused %d workflows on chain %d, tx: %s", len(req.WorkflowIDs), req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// AdminPauseAllByOwnerRequest pauses all workflows for a specific owner
type AdminPauseAllByOwnerRequest struct {
	ChainSelector uint64
	Owner         common.Address // Owner whose workflows should be paused
}

// AdminPauseAllByOwner changeset to pause all workflows for a specific owner
var _ cldf.ChangeSet[AdminPauseAllByOwnerRequest] = AdminPauseAllByOwner

func AdminPauseAllByOwner(env cldf.Environment, req AdminPauseAllByOwnerRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.AdminPauseAllByOwner(chain.DeployerKey, req.Owner)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to admin pause all by owner: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm admin pause all by owner transaction: %w", err)
	}

	env.Logger.Infof("Successfully admin paused all workflows for owner %s on chain %d, tx: %s", req.Owner.String(), req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// AdminPauseAllByDONRequest pauses all workflows for a specific DON family
type AdminPauseAllByDONRequest struct {
	ChainSelector uint64
	DONFamily     string // DON family whose workflows should be paused
}

// AdminPauseAllByDON changeset to pause all workflows for a specific DON family
var _ cldf.ChangeSet[AdminPauseAllByDONRequest] = AdminPauseAllByDON

func AdminPauseAllByDON(env cldf.Environment, req AdminPauseAllByDONRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.AdminPauseAllByDON(chain.DeployerKey, req.DONFamily)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to admin pause all by DON: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm admin pause all by DON transaction: %w", err)
	}

	env.Logger.Infof("Successfully admin paused all workflows for DON family %s on chain %d, tx: %s", req.DONFamily, req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}
