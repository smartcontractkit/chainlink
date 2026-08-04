package stellar

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// SetFeedFrozenRequest freezes or unfreezes a batch of feeds on an
// already-deployed DataFeedsCache contract.
type SetFeedFrozenRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	Admin     string
	DataIDs   []string
	Frozen    bool
}

var _ cldf.ChangeSetV2[*SetFeedFrozenRequest] = SetFeedFrozen{}

// SetFeedFrozen freezes or unfreezes feeds on the cache. Reads on a frozen
// feed are rejected with FeedFrozen; a feed with no recorded round cannot be
// frozen (NoFeedState).
type SetFeedFrozen struct{}

func (SetFeedFrozen) VerifyPreconditions(env cldf.Environment, req *SetFeedFrozenRequest) error {
	if err := verifyContractRef(env, req.ChainSel, CacheContract, req.Qualifier, req.Version); err != nil {
		return err
	}
	if err := validateAddress(req.Admin); err != nil {
		return err
	}
	if len(req.DataIDs) == 0 {
		return errors.New("DataIDs cannot be empty")
	}
	if _, err := dataIDsToBytes(req.DataIDs); err != nil {
		return err
	}
	return nil
}

func (SetFeedFrozen) Apply(env cldf.Environment, req *SetFeedFrozenRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	version := semver.MustParse(req.Version)
	d, _, err := resolveContractDeps(env, req.ChainSel, CacheContract, req.Qualifier, version)
	if err != nil {
		return out, err
	}

	ids, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.SetFeedFrozen, d.deps, operation.SetFeedFrozenInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
		DataIDs:    ids,
		Frozen:     req.Frozen,
	})
	return out, err
}
