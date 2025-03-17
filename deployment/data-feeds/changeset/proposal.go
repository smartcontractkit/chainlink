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

type ProposalData struct {
	contract string
	tx       *gethTypes.Transaction
}

func BuildMCMProposals(env deployment.Environment, description string, chainSelector uint64, pd []ProposalData, minDelay time.Duration) (*mcmslib.TimelockProposal, error) {
	state, _ := LoadOnchainState(env)
	chain := env.Chains[chainSelector]
	chainState := state.Chains[chainSelector]

	var transactions []mcmstypes.Transaction
	for _, proposal := range pd {
		transactions = append(transactions, mcmstypes.Transaction{
			To:               proposal.contract,
			Data:             proposal.tx.Data(),
			AdditionalFields: json.RawMessage(`{"value": 0}`),
		})
	}

	ops := &mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(chainSelector),
		Transactions:  transactions,
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
		env,
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
