package stellar

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// OwnershipRequest drives the ownership changesets for the cache or proxy.
type OwnershipRequest struct {
	ChainSel        uint64
	Qualifier       string
	Version         string
	Contract        datastore.ContractType // CacheContract or ProxyContract
	NewOwner        string                 // TransferOwnership only
	LiveUntilLedger uint32                 // TransferOwnership only; pending-transfer expiry ledger
}

func (req *OwnershipRequest) verifyPreconditions(env cldf.Environment) error {
	if _, ok := env.BlockChains.StellarChains()[req.ChainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}
	if err := validateContract(req.Contract); err != nil {
		return err
	}
	if err := verifyContractRef(env, req.ChainSel, req.Contract, req.Qualifier, req.Version); err != nil {
		return err
	}
	return nil
}

var (
	_ cldf.ChangeSetV2[*OwnershipRequest] = TransferOwnership{}
	_ cldf.ChangeSetV2[*OwnershipRequest] = AcceptOwnership{}
)

type ownershipInput struct {
	ContractID      string `json:"contract_id"`
	IsProxy         bool   `json:"is_proxy"`
	NewOwner        string `json:"new_owner"`
	LiveUntilLedger uint32 `json:"live_until_ledger"`
}

// TransferOwnership begins a two-step ownership transfer.
type TransferOwnership struct{}

func (TransferOwnership) VerifyPreconditions(env cldf.Environment, req *OwnershipRequest) error {
	if err := req.verifyPreconditions(env); err != nil {
		return err
	}
	if err := validateAddress(req.NewOwner); err != nil {
		return fmt.Errorf("new owner: %w", err)
	}
	if req.LiveUntilLedger == 0 {
		return errors.New("LiveUntilLedger must be nonzero")
	}
	return nil
}

func (TransferOwnership) Apply(env cldf.Environment, req *OwnershipRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := resolveDeps(env, req.ChainSel, req.Contract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, transferOwnershipOp, d.deps, ownershipInput{
		ContractID:      d.contractID,
		IsProxy:         req.Contract == ProxyContract,
		NewOwner:        req.NewOwner,
		LiveUntilLedger: req.LiveUntilLedger,
	})
	return out, err
}

var transferOwnershipOp = operations.NewOperation(
	"df:transfer-ownership", opVersion,
	"Begins two-step ownership transfer",
	func(b operations.Bundle, d StellarDeps, in ownershipInput) (void, error) {
		return void{}, adminClient(d, in.ContractID, in.IsProxy).TransferOwnership(b.GetContext(), in.NewOwner, in.LiveUntilLedger)
	},
)

// AcceptOwnership accepts a pending ownership transfer. The transaction
// signer must be the pending owner.
type AcceptOwnership struct{}

func (AcceptOwnership) VerifyPreconditions(env cldf.Environment, req *OwnershipRequest) error {
	return req.verifyPreconditions(env)
}

func (AcceptOwnership) Apply(env cldf.Environment, req *OwnershipRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := resolveDeps(env, req.ChainSel, req.Contract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, acceptOwnershipOp, d.deps, ownershipInput{
		ContractID: d.contractID,
		IsProxy:    req.Contract == ProxyContract,
	})
	return out, err
}

var acceptOwnershipOp = operations.NewOperation(
	"df:accept-ownership", opVersion,
	"Accepts a pending ownership transfer (caller must be the pending owner)",
	func(b operations.Bundle, d StellarDeps, in ownershipInput) (void, error) {
		return void{}, adminClient(d, in.ContractID, in.IsProxy).AcceptOwnership(b.GetContext())
	},
)
