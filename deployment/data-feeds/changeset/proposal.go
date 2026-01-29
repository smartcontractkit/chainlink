package changeset

import (
	"encoding/json"
	"fmt"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type ProposalData struct {
	contract          string
	tx                *gethTypes.Transaction
	timeLockQualifier string
}

// MultiChainProposalConfig is a map of chain selector to a list of proposals to be executed on that chain
type MultiChainProposalConfig map[uint64][]ProposalData

func BuildMultiChainProposals(env cldf.Environment, description string, proposalConfig MultiChainProposalConfig, minDelay time.Duration) (*mcmslib.TimelockProposal, error) {
	var timelocksPerChain = map[uint64]string{}
	var proposerMCMSes = map[uint64]string{}
	var inspectorPerChain = map[uint64]sdk.Inspector{}
	var batches []mcmstypes.BatchOperation

	for chainSelector, proposalData := range proposalConfig {
		var transactions []mcmstypes.Transaction
		for _, proposal := range proposalData {
			mcmsChainState, err := commonchangeset.MaybeLoadMCMSWithTimelockStateWithQualifier(env, []uint64{chainSelector}, proposal.timeLockQualifier)
			if err != nil {
				return nil, fmt.Errorf("failed to load MCMS contracts for chain %d: %w", chainSelector, err)
			}

			chain := env.BlockChains.EVMChains()[chainSelector]

			inspectorPerChain[chainSelector] = evm.NewInspector(chain.Client)
			timelocksPerChain[chainSelector] = mcmsChainState[chainSelector].Timelock.Address().Hex()
			proposerMCMSes[chainSelector] = mcmsChainState[chainSelector].ProposerMcm.Address().Hex()

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
