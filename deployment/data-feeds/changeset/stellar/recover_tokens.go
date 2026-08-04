package stellar

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// RecoverTokensRequest recovers tokens accidentally sent to an already-deployed
// DataFeedsCache contract.
type RecoverTokensRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	Token     string
	To        string
	Amount    int64
}

var _ cldf.ChangeSetV2[*RecoverTokensRequest] = RecoverTokens{}

// RecoverTokens recovers tokens accidentally sent to the cache.
type RecoverTokens struct{}

func (RecoverTokens) VerifyPreconditions(env cldf.Environment, req *RecoverTokensRequest) error {
	if err := verifyContractRef(env, req.ChainSel, CacheContract, req.Qualifier, req.Version); err != nil {
		return err
	}
	if err := validateAddress(req.Token); err != nil {
		return fmt.Errorf("token: %w", err)
	}
	if err := validateAddress(req.To); err != nil {
		return fmt.Errorf("to: %w", err)
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	return nil
}

func (RecoverTokens) Apply(env cldf.Environment, req *RecoverTokensRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	version := semver.MustParse(req.Version)
	d, _, err := resolveContractDeps(env, req.ChainSel, CacheContract, req.Qualifier, version)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.RecoverTokens, d.deps, operation.RecoverTokensInput{
		ContractID: d.contractID,
		Token:      req.Token,
		To:         req.To,
		Amount:     req.Amount,
	})
	return out, err
}
