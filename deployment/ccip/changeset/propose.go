package changeset

import (
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
)

func BuildProposalMetadata(state CCIPOnChainState, chains []uint64) (map[mcms.ChainIdentifier]common.Address, map[mcms.ChainIdentifier]mcms.ChainMetadata, error) {
	tlAddressMap := make(map[mcms.ChainIdentifier]common.Address)
	metaDataPerChain := make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	for _, sel := range chains {
		chainId := mcms.ChainIdentifier(sel)
		tlAddressMap[chainId] = state.Chains[sel].Timelock.Address()
		mcm := state.Chains[sel].ProposerMcm
		opCount, err := mcm.GetOpCount(nil)
		if err != nil {
			return nil, nil, err
		}
		metaDataPerChain[chainId] = mcms.ChainMetadata{
			StartingOpCount: opCount.Uint64(),
			MCMAddress:      mcm.Address(),
		}
	}
	return tlAddressMap, metaDataPerChain, nil
}

// Given batches of operations, we build the metadata and timelock addresses of those opartions
// We then return a proposal that can be executed and signed
func BuildProposalFromBatches(state CCIPOnChainState, batches []timelock.BatchChainOperation, description string, minDelay time.Duration) (*timelock.MCMSWithTimelockProposal, error) {
	if len(batches) == 0 {
		return nil, fmt.Errorf("no operations in batch")
	}

	chains := mapset.NewSet[uint64]()
	for _, op := range batches {
		chains.Add(uint64(op.ChainIdentifier))
	}

	tls, mcmsMd, err := BuildProposalMetadata(state, chains.ToSlice())
	if err != nil {
		return nil, err
	}

	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO: should be parameterized and based on current block timestamp.
		[]mcms.Signature{},
		false,
		mcmsMd,
		tls,
		description,
		batches,
		timelock.Schedule,
		minDelay.String(),
	)
}
