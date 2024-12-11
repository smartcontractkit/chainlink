package proposalutils

import (
	"fmt"
	"time"

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

type TimelockBatch struct {
	Timelock    common.Address
	ProposerMCM *gethwrappers.ManyChainMultiSig
	Operations  []mcms.Operation
}

// BuildProposalFromBatches Given batches of operations, we build the metadata and timelock addresses of those opartions
// We then return a proposal that can be executed and signed
func BuildProposalFromBatches(
	batches map[uint64]TimelockBatch,
	description string,
	minDelay time.Duration,
) (*timelock.MCMSWithTimelockProposal, error) {
	if len(batches) == 0 {
		return nil, fmt.Errorf("no operations in batch")
	}

	var chains []uint64
	proposersPerChain := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	timelocksPerChain := make(map[mcms.ChainIdentifier]common.Address)
	var chainBatches []timelock.BatchChainOperation
	for chainId, batch := range batches {
		timelocksPerChain[mcms.ChainIdentifier(chainId)] = batch.Timelock
		chainBatches = append(chainBatches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chainId),
			Batch:           batch.Operations,
		})
	}

	for chain, batch := range batches {
		chains = append(chains, chain)
		proposersPerChain[chain] = batch.ProposerMCM
	}

	mcmsMd, err := buildProposalMetadata(chains,
		proposersPerChain)
	if err != nil {
		return nil, err
	}

	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO: should be parameterized and based on current block timestamp.
		[]mcms.Signature{},
		false,
		mcmsMd,
		timelocksPerChain,
		description,
		chainBatches,
		timelock.Schedule,
		minDelay.String(),
	)
}
