package contracts

import (
	"fmt"
	"math/big"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

const (
	UpsertWorkflowUserDescription     = "upsertWorkflow on workflow registry v2 as user"
	PauseWorkflowUserDescription      = "pauseWorkflow on workflow registry v2 as user"
	BatchPauseWorkflowUserDescription = "batchPauseWorkflow on workflow registry v2 as user"
	ActivateWorkflowUserDescription   = "activateWorkflow on workflow registry v2 as user"
	DeleteWorkflowUserDescription     = "deleteWorkflow on workflow registry v2 a user"
	LinkOwnerUserDescription          = "linkOwner on workflow registry v2 as user"
	UnlinkOwnerUserDescription        = "unlinkOwner on workflow registry v2 as user"
	AllowlistRequest                  = "allowlistRequest on workflow registry v2 as user"
)

type UserUpsertWorkflowOpInput struct {
	WorkflowID     [32]byte `json:"workflowID"`
	WorkflowName   string   `json:"workflowName"`
	WorkflowTag    string   `json:"workflowTag"`
	WorkflowStatus uint8    `json:"workflowStatus"`
	DonFamily      string   `json:"donFamily"`
	BinaryURL      string   `json:"binaryURL"`
	ConfigURL      string   `json:"configURL"`
	Attributes     []byte   `json:"attributes"`
	KeepAlive      bool     `json:"keepAlive"`

	ChainSelector uint64                `json:"chainSelector"`
	MCMSConfig    *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
	Qualifier     string                `json:"qualifier"`
}

type UserUpsertWorkflowOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var UserUpsertWorkflowOp = operations.NewOperation(
	"user-upsert-workflow-op",
	semver.MustParse("1.0.0"),
	"User Upsert Workflow in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UserUpsertWorkflowOpInput) (UserUpsertWorkflowOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.Registry.UpsertWorkflow(opts, input.WorkflowName, input.WorkflowTag, input.WorkflowID, input.WorkflowStatus, input.DonFamily, input.BinaryURL, input.ConfigURL, input.Attributes, input.KeepAlive)
			if err != nil {
				return nil, fmt.Errorf("failed to call UpsertWorkflow: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UserUpsertWorkflowOpOutput{}, fmt.Errorf("failed to execute UpsertWorkflow: %w", err)
		}

		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for UpsertWorkflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully user upserted workflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		}

		return UserUpsertWorkflowOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

type UserPauseWorkflowOpInput struct {
	WorkflowID [32]byte `json:"workflowID"`

	ChainSelector uint64                `json:"chainSelector"`
	MCMSConfig    *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
	Qualifier     string                `json:"qualifier"`
}

type UserPauseWorkflowOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var UserPauseWorkflowOp = operations.NewOperation(
	"user-pause-workflow-op",
	semver.MustParse("1.0.0"),
	"User Pause Workflow in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UserPauseWorkflowOpInput) (UserPauseWorkflowOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.Registry.PauseWorkflow(opts, input.WorkflowID)
			if err != nil {
				return nil, fmt.Errorf("failed to call PauseWorkflow: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UserPauseWorkflowOpOutput{}, fmt.Errorf("failed to execute PauseWorkflow: %w", err)
		}
		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for PauseWorkflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully user paused workflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		}
		return UserPauseWorkflowOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

type UserBatchPauseWorkflowOpInput struct {
	WorkflowIDs [][32]byte `json:"workflowID"`

	ChainSelector uint64                `json:"chainSelector"`
	MCMSConfig    *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
	Qualifier     string                `json:"qualifier"`
}

type UserBatchPauseWorkflowOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var UserBatchPauseWorkflowOp = operations.NewOperation(
	"user-batch-pause-workflow-op",
	semver.MustParse("1.0.0"),
	"User Batch Pause Workflow in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UserBatchPauseWorkflowOpInput) (UserBatchPauseWorkflowOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.Registry.BatchPauseWorkflows(opts, input.WorkflowIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to call PauseWorkflow: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UserBatchPauseWorkflowOpOutput{}, fmt.Errorf("failed to execute PauseWorkflow: %w", err)
		}
		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for PauseWorkflow %x on chain %d", input.WorkflowIDs, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully user paused workflow %x on chain %d", input.WorkflowIDs, input.ChainSelector)
		}
		return UserBatchPauseWorkflowOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

type UserActivateWorkflowOpInput struct {
	WorkflowID [32]byte `json:"workflowID"`
	DonFamily  string   `json:"donFamily"`

	ChainSelector uint64                `json:"chainSelector"`
	MCMSConfig    *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
	Qualifier     string                `json:"qualifier"`
}

type UserActivateWorkflowOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var UserActivateWorkflowOp = operations.NewOperation(
	"user-activate-workflow-op",
	semver.MustParse("1.0.0"),
	"User Activate Workflow in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UserActivateWorkflowOpInput) (UserActivateWorkflowOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.Registry.ActivateWorkflow(opts, input.WorkflowID, input.DonFamily)
			if err != nil {
				return nil, fmt.Errorf("failed to call ActivateWorkflow: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UserActivateWorkflowOpOutput{}, fmt.Errorf("failed to execute ActivateWorkflow: %w", err)
		}
		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for ActivateWorkflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully user activated workflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		}
		return UserActivateWorkflowOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

type UserDeleteWorkflowOpInput struct {
	WorkflowID [32]byte `json:"workflowID"`

	ChainSelector uint64                `json:"chainSelector"`
	MCMSConfig    *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
	Qualifier     string                `json:"qualifier"`
}

type UserDeleteWorkflowOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var UserDeleteWorkflowOp = operations.NewOperation(
	"user-delete-workflow-op",
	semver.MustParse("1.0.0"),
	"User Delete Workflow in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UserDeleteWorkflowOpInput) (UserDeleteWorkflowOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.Registry.DeleteWorkflow(opts, input.WorkflowID)
			if err != nil {
				return nil, fmt.Errorf("failed to call DeleteWorkflow: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UserDeleteWorkflowOpOutput{}, fmt.Errorf("failed to execute DeleteWorkflow: %w", err)
		}
		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for DeleteWorkflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully user deleted workflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		}
		return UserDeleteWorkflowOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

type UserLinkOwnerOpInput struct {
	ValidityTimestamp *big.Int `json:"validityTimestamp"`
	Proof             [32]byte `json:"proof"`
	Signature         []byte   `json:"signature"`

	ChainSelector uint64                `json:"chainSelector"`
	MCMSConfig    *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
	Qualifier     string                `json:"qualifier"`
}

type UserLinkOwnerOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var UserLinkOwnerOp = operations.NewOperation(
	"user-link-owner-op",
	semver.MustParse("1.0.0"),
	"User Link Owner in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UserLinkOwnerOpInput) (UserLinkOwnerOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.Registry.LinkOwner(opts, input.ValidityTimestamp, input.Proof, input.Signature)
			if err != nil {
				return nil, fmt.Errorf("failed to call LinkOwner: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UserLinkOwnerOpOutput{}, fmt.Errorf("failed to execute LinkOwner: %w", err)
		}
		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for LinkOwner on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully user linked owner on chain %d", input.ChainSelector)
		}
		return UserLinkOwnerOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

type UserUnlinkOwnerOpInput struct {
	Address           common.Address `json:"address"`
	ValidityTimestamp *big.Int       `json:"validityTimestamp"`
	Signature         []byte         `json:"signature"`

	ChainSelector uint64                `json:"chainSelector"`
	MCMSConfig    *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
	Qualifier     string                `json:"qualifier"`
}

type UserUnlinkOwnerOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var UserUnlinkOwnerOp = operations.NewOperation(
	"user-unlink-owner-op",
	semver.MustParse("1.0.0"),
	"User Unlink Owner in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UserUnlinkOwnerOpInput) (UserUnlinkOwnerOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.Registry.UnlinkOwner(opts, input.Address, input.ValidityTimestamp, input.Signature)
			if err != nil {
				return nil, fmt.Errorf("failed to call UnlinkOwner: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UserUnlinkOwnerOpOutput{}, fmt.Errorf("failed to execute UnlinkOwner: %w", err)
		}
		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for UnlinkOwner on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully user unlinked owner on chain %d", input.ChainSelector)
		}
		return UserUnlinkOwnerOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

type UserAllowlistRequestOpInput struct {
	RequestDigest   [32]byte `json:"requestDigest"`
	ExpiryTimestamp uint32   `json:"expiryTimestamp"`

	ChainSelector uint64                `json:"chainSelector"`
	MCMSConfig    *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
	Qualifier     string                `json:"qualifier"`
}

type UserAllowlistRequestOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var UserAllowlistRequestOp = operations.NewOperation(
	"user-allowlist-request-op",
	semver.MustParse("1.0.0"),
	"User Allowlist Request in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UserAllowlistRequestOpInput) (UserAllowlistRequestOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.Registry.AllowlistRequest(opts, input.RequestDigest, input.ExpiryTimestamp)
			if err != nil {
				return nil, fmt.Errorf("failed to call AllowlistRequest: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UserAllowlistRequestOpOutput{}, fmt.Errorf("failed to execute AllowlistRequest: %w", err)
		}
		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AllowlistRequest on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully user allowlisted request on chain %d", input.ChainSelector)
		}
		return UserAllowlistRequestOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)
