package changeset

import (
	"fmt"
	"math/big"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
)

func GenerateAcceptOwnershipProposal(
	state CCIPOnChainState,
	homeChain uint64,
	chains []uint64,
) (*timelock.MCMSWithTimelockProposal, error) {
	// TODO: Accept rest of contracts
	var batches []timelock.BatchChainOperation
	for _, sel := range chains {
		chain, _ := chainsel.ChainBySelector(sel)
		acceptOnRamp, err := state.Chains[sel].OnRamp.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return nil, err
		}
		acceptFeeQuoter, err := state.Chains[sel].FeeQuoter.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return nil, err
		}
		chainSel := mcms.ChainIdentifier(chain.Selector)
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: chainSel,
			Batch: []mcms.Operation{
				{
					To:    state.Chains[sel].OnRamp.Address(),
					Data:  acceptOnRamp.Data(),
					Value: big.NewInt(0),
				},
				{
					To:    state.Chains[sel].FeeQuoter.Address(),
					Data:  acceptFeeQuoter.Data(),
					Value: big.NewInt(0),
				},
			},
		})
	}

	acceptCR, err := state.Chains[homeChain].CapabilityRegistry.AcceptOwnership(deployment.SimTransactOpts())
	if err != nil {
		return nil, err
	}
	acceptCCIPConfig, err := state.Chains[homeChain].CCIPHome.AcceptOwnership(deployment.SimTransactOpts())
	if err != nil {
		return nil, err
	}
	homeChainID := mcms.ChainIdentifier(homeChain)
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: homeChainID,
		Batch: []mcms.Operation{
			{
				To:    state.Chains[homeChain].CapabilityRegistry.Address(),
				Data:  acceptCR.Data(),
				Value: big.NewInt(0),
			},
			{
				To:    state.Chains[homeChain].CCIPHome.Address(),
				Data:  acceptCCIPConfig.Data(),
				Value: big.NewInt(0),
			},
		},
	})

	return BuildProposalFromBatches(state, batches, "accept ownership operations", 0)
}

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
