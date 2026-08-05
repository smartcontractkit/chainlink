package stellar

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
)

// FeedAdminRequest grants or revokes feed-admin rights on the cache.
type FeedAdminRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	Admin     string
}

func (req *FeedAdminRequest) verifyPreconditions(env cldf.Environment) error {
	if err := verifyContractRef(env, req.ChainSel, CacheContract, req.Qualifier, req.Version); err != nil {
		return err
	}
	if err := validateAddress(req.Admin); err != nil {
		return err
	}
	return nil
}

var (
	_ cldf.ChangeSetV2[*FeedAdminRequest] = AddFeedAdmin{}
	_ cldf.ChangeSetV2[*FeedAdminRequest] = RemoveFeedAdmin{}
)

type feedAdminInput struct {
	ContractID string `json:"contract_id"`
	Admin      string `json:"admin"`
}

// AddFeedAdmin grants feed-admin rights on the cache.
type AddFeedAdmin struct{}

func (AddFeedAdmin) VerifyPreconditions(env cldf.Environment, req *FeedAdminRequest) error {
	return req.verifyPreconditions(env)
}

func (AddFeedAdmin) Apply(env cldf.Environment, req *FeedAdminRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := resolveDeps(env, req.ChainSel, CacheContract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, addFeedAdminOp, d.deps, feedAdminInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
	})
	return out, err
}

var addFeedAdminOp = operations.NewOperation(
	"df-cache:add-feed-admin", opVersion,
	"Grants feed-admin rights on the cache",
	func(b operations.Bundle, d StellarDeps, in feedAdminInput) (void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return void{}, c.AddFeedAdmin(b.GetContext(), in.Admin)
	},
)

// RemoveFeedAdmin revokes feed-admin rights on the cache.
type RemoveFeedAdmin struct{}

func (RemoveFeedAdmin) VerifyPreconditions(env cldf.Environment, req *FeedAdminRequest) error {
	return req.verifyPreconditions(env)
}

func (RemoveFeedAdmin) Apply(env cldf.Environment, req *FeedAdminRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := resolveDeps(env, req.ChainSel, CacheContract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, removeFeedAdminOp, d.deps, feedAdminInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
	})
	return out, err
}

var removeFeedAdminOp = operations.NewOperation(
	"df-cache:remove-feed-admin", opVersion,
	"Revokes feed-admin rights on the cache",
	func(b operations.Bundle, d StellarDeps, in feedAdminInput) (void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return void{}, c.RemoveFeedAdmin(b.GetContext(), in.Admin)
	},
)
