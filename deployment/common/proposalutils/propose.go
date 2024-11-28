package proposalutils

import (
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
)

func buildProposalMetadata(
	chainSelectors []uint64,
	proposerMcmsesPerChain map[uint64]*gethwrappers.ManyChainMultiSig,
) (map[mcms.ChainIdentifier]mcms.ChainMetadata, error) {
	metaDataPerChain := make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	for _, selector := range chainSelectors {
		proposerMcms, ok := proposerMcmsesPerChain[selector]
		if !ok {
			return nil, fmt.Errorf("missing proposer mcm for chain %d", selector)
		}
		chainId := mcms.ChainIdentifier(selector)
		opCount, err := proposerMcms.GetOpCount(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get op count for chain %d: %w", selector, err)
		}
		metaDataPerChain[chainId] = mcms.ChainMetadata{
			StartingOpCount: opCount.Uint64(),
			MCMAddress:      proposerMcms.Address(),
		}
	}
	return metaDataPerChain, nil
}

// Given batches of operations, we build the metadata and timelock addresses of those opartions
// We then return a proposal that can be executed and signed
func BuildProposalFromBatches(
	timelocksPerChain map[uint64]common.Address,
	proposerMcmsesPerChain map[uint64]*gethwrappers.ManyChainMultiSig,
	batches []timelock.BatchChainOperation,
	description string,
	minDelay time.Duration,
) (*timelock.MCMSWithTimelockProposal, error) {
	if len(batches) == 0 {
		return nil, fmt.Errorf("no operations in batch")
	}

	chains := mapset.NewSet[uint64]()
	for _, op := range batches {
		chains.Add(uint64(op.ChainIdentifier))
	}

	mcmsMd, err := buildProposalMetadata(chains.ToSlice(), proposerMcmsesPerChain)
	if err != nil {
		return nil, err
	}

	tlsPerChainId := make(map[mcms.ChainIdentifier]common.Address)
	for chainId, tl := range timelocksPerChain {
		tlsPerChainId[mcms.ChainIdentifier(chainId)] = tl
	}

	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO: should be parameterized and based on current block timestamp.
		[]mcms.Signature{},
		false,
		mcmsMd,
		tlsPerChainId,
		description,
		batches,
		timelock.Schedule,
		minDelay.String(),
	)
}
