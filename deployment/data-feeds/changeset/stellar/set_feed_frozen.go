package stellar

import (
	"errors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
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

type setFeedFrozenInput struct {
	ContractID string     `json:"contract_id"`
	Admin      string     `json:"admin"`
	DataIDs    [][16]byte `json:"data_ids"`
	Frozen     bool       `json:"frozen"`
}

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
	d, err := resolveDeps(env, req.ChainSel, CacheContract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}

	ids, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, setFeedFrozenOp, d.deps, setFeedFrozenInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
		DataIDs:    ids,
		Frozen:     req.Frozen,
	})
	return out, err
}

var setFeedFrozenOp = operations.NewOperation(
	"df-cache:set-feed-frozen", opVersion,
	"Freezes or unfreezes feeds on the cache",
	func(b operations.Bundle, d StellarDeps, in setFeedFrozenInput) (void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return void{}, c.SetFeedFrozen(b.GetContext(), in.Admin, in.DataIDs, in.Frozen)
	},
)
