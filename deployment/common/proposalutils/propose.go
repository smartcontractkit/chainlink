package proposalutils

import (
	"errors"
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

const (
	DefaultValidUntil = 72 * time.Hour
)

type TimelockConfig struct {
	MinDelay     time.Duration // delay for timelock worker to execute the transfers.
	MCMSAction   types.TimelockAction
	OverrideRoot bool // if true, override the previous root with the new one
}

func (tc *TimelockConfig) Validate(chain deployment.Chain, s state.MCMSWithTimelockState) error {
	// if MCMSAction is not set, default to timelock.Schedule
	if tc.MCMSAction == "" {
		tc.MCMSAction = types.TimelockActionSchedule
	}
	if tc.MCMSAction != types.TimelockActionSchedule &&
		tc.MCMSAction != types.TimelockActionCancel &&
		tc.MCMSAction != types.TimelockActionBypass {
		return fmt.Errorf("invalid MCMS type %s", tc.MCMSAction)
	}
	if s.Timelock == nil {
		return fmt.Errorf("missing timelock on %s", chain)
	}
	if tc.MCMSAction == types.TimelockActionSchedule && s.ProposerMcm == nil {
		return fmt.Errorf("missing proposerMcm on %s", chain)
	}
	if tc.MCMSAction == types.TimelockActionCancel && s.CancellerMcm == nil {
		return fmt.Errorf("missing cancellerMcm on %s", chain)
	}
	if tc.MCMSAction == types.TimelockActionBypass && s.BypasserMcm == nil {
		return fmt.Errorf("missing bypasserMcm on %s", chain)
	}
	return nil
}

func BuildProposalMetadata(
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

// BuildProposalFromBatches Given batches of operations, we build the metadata and timelock addresses of those opartions
// We then return a proposal that can be executed and signed.
// You can specify multiple batches for the same chain, but the only
// usecase to do that would be you have a batch that can't fit in a single
// transaction due to gas or calldata constraints of the chain.
// The batches are specified separately because we eventually intend
// to support user-specified cross chain ordering of batch execution by the tooling itself.
// TODO: Can/should merge timelocks and proposers into a single map for the chain.
// Deprecated: Use BuildProposalFromBatchesV2 instead.
func BuildProposalFromBatches(
	timelocksPerChain map[uint64]common.Address,
	proposerMcmsesPerChain map[uint64]*gethwrappers.ManyChainMultiSig,
	batches []timelock.BatchChainOperation,
	description string,
	minDelay time.Duration,
) (*timelock.MCMSWithTimelockProposal, error) {
	if len(batches) == 0 {
		return nil, errors.New("no operations in batch")
	}

	chains := mapset.NewSet[uint64]()
	for _, op := range batches {
		chains.Add(uint64(op.ChainIdentifier))
	}

	mcmsMd, err := BuildProposalMetadata(chains.ToSlice(), proposerMcmsesPerChain)
	if err != nil {
		return nil, err
	}

	tlsPerChainId := make(map[mcms.ChainIdentifier]common.Address)
	for chainId, tl := range timelocksPerChain {
		tlsPerChainId[mcms.ChainIdentifier(chainId)] = tl
	}
	validUntil := time.Now().Unix() + int64(DefaultValidUntil.Seconds())
	return timelock.NewMCMSWithTimelockProposal(
		"1",
		uint32(validUntil),
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

// BuildProposalFromBatchesV2 uses the new MCMS library which replaces the implementation in BuildProposalFromBatches.
func BuildProposalFromBatchesV2(
	e deployment.Environment,
	timelockAddressPerChain map[uint64]string,
	mcmsAddressPerChain map[uint64]string, inspectorPerChain map[uint64]mcmssdk.Inspector,
	batches []types.BatchOperation,
	description string,
	mcmsCfg TimelockConfig,
) (*mcmslib.TimelockProposal, error) {
	// default to schedule if not set, this is to be consistent with the old implementation
	// and to avoid breaking changes
	if mcmsCfg.MCMSAction == "" {
		mcmsCfg.MCMSAction = types.TimelockActionSchedule
	}
	if len(batches) == 0 {
		return nil, errors.New("no operations in batch")
	}

	chains := mapset.NewSet[uint64]()
	for _, op := range batches {
		chains.Add(uint64(op.ChainSelector))
	}
	tlsPerChainID := make(map[types.ChainSelector]string)
	for chainID, tl := range timelockAddressPerChain {
		tlsPerChainID[types.ChainSelector(chainID)] = tl
	}
	mcmsMd, err := buildProposalMetadataV2(e, chains.ToSlice(), inspectorPerChain, mcmsAddressPerChain)
	if err != nil {
		return nil, err
	}

	validUntil := time.Now().Unix() + int64(DefaultValidUntil.Seconds())

	builder := mcmslib.NewTimelockProposalBuilder()
	builder.
		SetVersion("v1").
		SetAction(mcmsCfg.MCMSAction).
		//nolint:gosec // G115
		SetValidUntil(uint32(validUntil)).
		SetDescription(description).
		SetDelay(types.NewDuration(mcmsCfg.MinDelay)).
		SetOverridePreviousRoot(mcmsCfg.OverrideRoot).
		SetChainMetadata(mcmsMd).
		SetTimelockAddresses(tlsPerChainID).
		SetOperations(batches)

	build, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return build, nil
}

func buildProposalMetadataV2(
	env deployment.Environment,
	chainSelectors []uint64,
	inspectorPerChain map[uint64]mcmssdk.Inspector,
	mcmsPerChain map[uint64]string, // can be proposer, canceller or bypasser
) (map[types.ChainSelector]types.ChainMetadata, error) {
	metaDataPerChain := make(map[types.ChainSelector]types.ChainMetadata)
	for _, selector := range chainSelectors {
		proposerMcms, ok := mcmsPerChain[selector]
		if !ok {
			return nil, fmt.Errorf("missing proposer mcm for chain %d", selector)
		}
		chainID := types.ChainSelector(selector)
		opCount, err := inspectorPerChain[selector].GetOpCount(env.GetContext(), proposerMcms)
		if err != nil {
			return nil, fmt.Errorf("failed to get op count for chain %d: %w", selector, err)
		}
		family, err := chain_selectors.GetSelectorFamily(selector)
		if err != nil {
			return nil, fmt.Errorf("failed to get family for chain %d: %w", selector, err)
		}
		switch family {
		case chain_selectors.FamilyEVM:
			metaDataPerChain[chainID] = types.ChainMetadata{
				StartingOpCount: opCount,
				MCMAddress:      proposerMcms,
			}
		case chain_selectors.FamilySolana:
			addresses, err := env.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return nil, fmt.Errorf("failed to load addresses for chain %d: %w", selector, err)
			}
			solanaState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(env.SolChains[selector], addresses)
			if err != nil {
				return nil, fmt.Errorf("failed to load solana state: %w", err)
			}
			metaDataPerChain[chainID], err = mcmssolanasdk.NewChainMetadata(
				opCount,
				solanaState.McmProgram,
				mcmssolanasdk.PDASeed(solanaState.ProposerMcmSeed),
				solanaState.ProposerAccessControllerAccount,
				solanaState.CancellerAccessControllerAccount,
				solanaState.BypasserAccessControllerAccount)
			if err != nil {
				return nil, fmt.Errorf("failed to create chain metadata: %w", err)
			}
		}
	}

	return metaDataPerChain, nil
}
