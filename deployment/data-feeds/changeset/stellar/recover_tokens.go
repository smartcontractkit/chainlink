package stellar

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// RecoverTokensRequest recovers tokens accidentally sent to the cache or proxy.
type RecoverTokensRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	Contract  datastore.ContractType // CacheContract or ProxyContract
	Token     string
	To        string
	Amount    int64
}

var _ cldf.ChangeSetV2[*RecoverTokensRequest] = RecoverTokens{}

// RecoverTokens recovers tokens accidentally sent to the cache or proxy.
type RecoverTokens struct{}

type recoverTokensInput struct {
	ContractID string `json:"contract_id"`
	IsProxy    bool   `json:"is_proxy"`
	Token      string `json:"token"`
	To         string `json:"to"`
	Amount     int64  `json:"amount"`
}

func (RecoverTokens) VerifyPreconditions(env cldf.Environment, req *RecoverTokensRequest) error {
	if err := validateContract(req.Contract); err != nil {
		return err
	}
	if err := verifyContractRef(env, req.ChainSel, req.Contract, req.Qualifier, req.Version); err != nil {
		return err
	}
	if err := validateAddress(req.Token); err != nil {
		return fmt.Errorf("token: %w", err)
	}
	if err := validateAddress(req.To); err != nil {
		return fmt.Errorf("to: %w", err)
	}
	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	return nil
}

func (RecoverTokens) Apply(env cldf.Environment, req *RecoverTokensRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	d, _, err := resolveContractDeps(env, req.ChainSel, req.Contract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}
	_, err = operations.ExecuteOperation(env.OperationsBundle, recoverTokensOp, d.deps, recoverTokensInput{
		ContractID: d.contractID,
		IsProxy:    req.Contract == ProxyContract,
		Token:      req.Token,
		To:         req.To,
		Amount:     req.Amount,
	})
	return out, err
}

var recoverTokensOp = operations.NewOperation(
	"df:recover-tokens", opVersion,
	"Recovers tokens accidentally sent to the contract",
	func(b operations.Bundle, d StellarDeps, in recoverTokensInput) (void, error) {
		return void{}, adminClient(d, in.ContractID, in.IsProxy).RecoverTokens(b.GetContext(), in.Token, in.To, in.Amount)
	},
)
