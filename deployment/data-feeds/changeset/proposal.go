package changeset

import (
	"encoding/json"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func BuildMCMProposal(env deployment.Environment, description string, chainSelector uint64, contractAddress string, tx *gethTypes.Transaction, minDelay time.Duration) (*mcmslib.TimelockProposal, error) {
	state, _ := LoadOnchainState(env)
	chain := env.Chains[chainSelector]
	chainState := state.Chains[chainSelector]

	ops := &mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(chainSelector),
		Transactions: []mcmstypes.Transaction{
			{
				To:               contractAddress,
				Data:             tx.Data(),
				AdditionalFields: json.RawMessage(`{"value": 0}`),
			},
		},
	}

	timelocksPerChain := map[uint64]string{
		chainSelector: chainState.Timelock.Address().Hex(),
	}
	proposerMCMSes := map[uint64]string{
		chainSelector: chainState.ProposerMcm.Address().Hex(),
	}

	inspectorPerChain := map[uint64]sdk.Inspector{}
	inspectorPerChain[chainSelector] = evm.NewInspector(chain.Client)

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		env.GetContext(),
		timelocksPerChain,
		proposerMCMSes,
		inspectorPerChain,
		[]mcmstypes.BatchOperation{*ops},
		description,
		minDelay,
	)
	if err != nil {
		return nil, err
	}
	return proposal, err
}
