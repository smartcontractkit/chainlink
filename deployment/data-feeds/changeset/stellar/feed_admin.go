package stellar

import (
	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// FeedAdminRequest grants or revokes feed-admin rights on an already-deployed
// DataFeedsCache contract.
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

// resolve looks up the cache contract's dependencies + address for req.
func (req *FeedAdminRequest) resolve(env cldf.Environment) (stellarApplyDeps, error) {
	version := semver.MustParse(req.Version)
	d, _, err := resolveContractDeps(env, req.ChainSel, CacheContract, req.Qualifier, version)
	return d, err
}

var (
	_ cldf.ChangeSetV2[*FeedAdminRequest] = AddFeedAdmin{}
	_ cldf.ChangeSetV2[*FeedAdminRequest] = RemoveFeedAdmin{}
)

// AddFeedAdmin grants feed-admin rights on the cache.
type AddFeedAdmin struct{}

func (AddFeedAdmin) VerifyPreconditions(env cldf.Environment, req *FeedAdminRequest) error {
	return req.verifyPreconditions(env)
}

func (AddFeedAdmin) Apply(env cldf.Environment, req *FeedAdminRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := req.resolve(env)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.AddFeedAdmin, d.deps, operation.FeedAdminInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
	})
	return out, err
}

// RemoveFeedAdmin revokes feed-admin rights on the cache.
type RemoveFeedAdmin struct{}

func (RemoveFeedAdmin) VerifyPreconditions(env cldf.Environment, req *FeedAdminRequest) error {
	return req.verifyPreconditions(env)
}

func (RemoveFeedAdmin) Apply(env cldf.Environment, req *FeedAdminRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := req.resolve(env)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.RemoveFeedAdmin, d.deps, operation.FeedAdminInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
	})
	return out, err
}
