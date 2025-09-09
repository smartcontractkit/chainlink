package changeset

import (
	"encoding/json"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type ProposalData struct {
	contract string
	tx       *gethTypes.Transaction
}

// MultiChainProposalConfig is a map of chain selector to a list of proposals to be executed on that chain
type MultiChainProposalConfig map[uint64][]ProposalData

func BuildMultiChainProposals(env cldf.Environment, description string, proposalConfig MultiChainProposalConfig, minDelay time.Duration) (*mcmslib.TimelockProposal, error) {
	state, _ := LoadOnchainState(env)

	var timelocksPerChain = map[uint64]string{}
	var proposerMCMSes = map[uint64]string{}
	var inspectorPerChain = map[uint64]sdk.Inspector{}
	var batches []mcmstypes.BatchOperation

	for chainSelector, proposalData := range proposalConfig {
		chain := env.BlockChains.EVMChains()[chainSelector]
		chainState := state.Chains[chainSelector]

		inspectorPerChain[chainSelector] = evm.NewInspector(chain.Client)
		timelocksPerChain[chainSelector] = chainState.Timelock.Address().Hex()
		proposerMCMSes[chainSelector] = chainState.ProposerMcm.Address().Hex()

		var transactions []mcmstypes.Transaction
		for _, proposal := range proposalData {
			transactions = append(transactions, mcmstypes.Transaction{
				To:               proposal.contract,
				Data:             proposal.tx.Data(),
				AdditionalFields: json.RawMessage(`{"value": 0}`),
			})
		}
		batches = append(batches, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  transactions,
		})
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		env,
		timelocksPerChain,
		proposerMCMSes,
		inspectorPerChain,
		batches,
		description,
		proposalutils.TimelockConfig{MinDelay: minDelay},
	)
	if err != nil {
		return nil, err
	}
	return proposal, err
}
