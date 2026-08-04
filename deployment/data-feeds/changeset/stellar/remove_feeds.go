package stellar

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// RemoveFeedConfigsRequest removes a batch of feed configs from an
// already-deployed DataFeedsCache contract.
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
	version := semver.MustParse(req.Version)
	d, _, err := resolveContractDeps(env, req.ChainSel, CacheContract, req.Qualifier, version)
	if err != nil {
		return out, err
	}

	ids, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.RemoveFeedConfigs, d.deps, operation.RemoveFeedConfigsInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
		DataIDs:    ids,
	})
	return out, err
}
