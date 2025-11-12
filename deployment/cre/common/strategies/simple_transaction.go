package strategies

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
)

// SimpleTransaction executes a transaction directly without MCMS
type SimpleTransaction struct {
	Chain cldf_evm.Chain
}

func (s *SimpleTransaction) Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (*mcmstypes.BatchOperation, *types.Transaction, error) {
	tx, err := callFn(s.Chain.DeployerKey)
	if err != nil {
		return nil, nil, err
	}

	_, err = s.Chain.Confirm(tx)
	return nil, tx, err
}

func (s *SimpleTransaction) BuildProposal(_ []mcmstypes.BatchOperation) (*mcmslib.TimelockProposal, error) {
	return nil, nil
}
