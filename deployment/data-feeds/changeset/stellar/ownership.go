package stellar

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// OwnershipRequest drives the two-step ownership changesets for either DF
// contract (CacheContract or ProxyContract) — the generated cache and proxy
// clients expose identical ownership methods.
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
	if req.Contract != CacheContract && req.Contract != ProxyContract {
		return fmt.Errorf("unsupported contract type %q: must be %q or %q", req.Contract, CacheContract, ProxyContract)
	}
	if err := verifyContractRef(env, req.ChainSel, req.Contract, req.Qualifier, req.Version); err != nil {
		return err
	}
	return nil
}

// resolve looks up the target contract's dependencies + address for req.
func (req *OwnershipRequest) resolve(env cldf.Environment) (stellarApplyDeps, error) {
	version := semver.MustParse(req.Version)
	d, _, err := resolveContractDeps(env, req.ChainSel, req.Contract, req.Qualifier, version)
	return d, err
}

var (
	_ cldf.ChangeSetV2[*OwnershipRequest] = TransferOwnership{}
	_ cldf.ChangeSetV2[*OwnershipRequest] = AcceptOwnership{}
	_ cldf.ChangeSetV2[*OwnershipRequest] = RenounceOwnership{}
)

// TransferOwnership begins a two-step ownership transfer on the cache or proxy.
type TransferOwnership struct{}

func (TransferOwnership) VerifyPreconditions(env cldf.Environment, req *OwnershipRequest) error {
	if err := req.verifyPreconditions(env); err != nil {
		return err
	}
	if err := validateAddress(req.NewOwner); err != nil {
		return fmt.Errorf("new owner: %w", err)
	}
	if req.LiveUntilLedger == 0 {
		return fmt.Errorf("LiveUntilLedger must be nonzero")
	}
	return nil
}

func (TransferOwnership) Apply(env cldf.Environment, req *OwnershipRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := req.resolve(env)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.TransferOwnership, d.deps, operation.OwnershipInput{
		ContractID:      d.contractID,
		IsProxy:         req.Contract == ProxyContract,
		NewOwner:        req.NewOwner,
		LiveUntilLedger: req.LiveUntilLedger,
	})
	return out, err
}

// AcceptOwnership accepts a pending ownership transfer on the cache or proxy
// (the caller signing the underlying transaction must be the pending owner).
type AcceptOwnership struct{}

func (AcceptOwnership) VerifyPreconditions(env cldf.Environment, req *OwnershipRequest) error {
	return req.verifyPreconditions(env)
}

func (AcceptOwnership) Apply(env cldf.Environment, req *OwnershipRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := req.resolve(env)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.AcceptOwnership, d.deps, operation.OwnershipInput{
		ContractID: d.contractID,
		IsProxy:    req.Contract == ProxyContract,
	})
	return out, err
}

// RenounceOwnership permanently renounces ownership of the cache or proxy.
type RenounceOwnership struct{}

func (RenounceOwnership) VerifyPreconditions(env cldf.Environment, req *OwnershipRequest) error {
	return req.verifyPreconditions(env)
}

func (RenounceOwnership) Apply(env cldf.Environment, req *OwnershipRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := req.resolve(env)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.RenounceOwnership, d.deps, operation.OwnershipInput{
		ContractID: d.contractID,
		IsProxy:    req.Contract == ProxyContract,
	})
	return out, err
}
