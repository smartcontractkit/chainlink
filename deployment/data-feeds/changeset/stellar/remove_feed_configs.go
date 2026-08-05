package stellar

import (
	"errors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
)

// RemoveFeedConfigsRequest removes a batch of feed configs from the cache.
type RemoveFeedConfigsRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	Admin     string
	DataIDs   []string
}

var _ cldf.ChangeSetV2[*RemoveFeedConfigsRequest] = RemoveFeedConfigs{}

// RemoveFeedConfigs removes feed configs from the cache.
type RemoveFeedConfigs struct{}

type removeFeedConfigsInput struct {
	ContractID string     `json:"contract_id"`
	Admin      string     `json:"admin"`
	DataIDs    [][16]byte `json:"data_ids"`
}

func (RemoveFeedConfigs) VerifyPreconditions(env cldf.Environment, req *RemoveFeedConfigsRequest) error {
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

func (RemoveFeedConfigs) Apply(env cldf.Environment, req *RemoveFeedConfigsRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, err := resolveDeps(env, req.ChainSel, CacheContract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}

	ids, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, removeFeedConfigsOp, d.deps, removeFeedConfigsInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
		DataIDs:    ids,
	})
	return out, err
}

var removeFeedConfigsOp = operations.NewOperation(
	"df-cache:remove-feed-configs", opVersion,
	"Removes feed configs from the cache",
	func(b operations.Bundle, d StellarDeps, in removeFeedConfigsInput) (void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return void{}, c.RemoveFeedConfigs(b.GetContext(), in.Admin, in.DataIDs)
	},
)
