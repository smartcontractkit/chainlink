package proposalutils

import (
	"errors"
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	types2 "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const (
	DefaultValidUntil = 72 * time.Hour
)

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
	proposerAddressPerChain map[uint64]string,
	inspectorPerChain map[uint64]mcmssdk.Inspector,

	batches []types.BatchOperation,
	description string,
	minDelay time.Duration,
) (*mcmslib.TimelockProposal, error) {
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
	mcmsMd, err := buildProposalMetadataV2(e, chains.ToSlice(), inspectorPerChain, proposerAddressPerChain)
	if err != nil {
		return nil, err
	}

	validUntil := time.Now().Unix() + int64(DefaultValidUntil.Seconds())

	builder := mcmslib.NewTimelockProposalBuilder()
	builder.
		SetVersion("v1").
		SetAction(types.TimelockActionSchedule).
		//nolint:gosec // G115
		SetValidUntil(uint32(validUntil)).
		SetDescription(description).
		SetDelay(types.NewDuration(minDelay)).
		SetOverridePreviousRoot(false).
		SetChainMetadata(mcmsMd).
		SetTimelockAddresses(tlsPerChainID).
		SetOperations(batches)

	build, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return build, nil
}

func getSolAccessControllerFromAddress(addresses map[string]deployment.TypeAndVersion, typeAndVersion deployment.TypeAndVersion) (solana.PublicKey, error) {
	for addr, tv := range addresses {
		if tv.Type == typeAndVersion.Type && tv.Version == typeAndVersion.Version {
			pubkey, err := solana.PublicKeyFromBase58(addr)
			if err != nil {
				return solana.PublicKey{}, fmt.Errorf("failed to parse address: %w", err)
			}
			return pubkey, nil
		}
	}

	return solana.PublicKey{}, errors.New("address not found")
}
func buildProposalMetadataV2(
	env deployment.Environment,
	chainSelectors []uint64,
	inspectorPerChain map[uint64]mcmssdk.Inspector,
	proposerMcmsesPerChain map[uint64]string,
) (map[types.ChainSelector]types.ChainMetadata, error) {
	metaDataPerChain := make(map[types.ChainSelector]types.ChainMetadata)
	for _, selector := range chainSelectors {
		proposerMcms, ok := proposerMcmsesPerChain[selector]
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
			mcmsProgramID, mcmSeed, err := mcmssolanasdk.ParseContractAddress(proposerMcms)
			if err != nil {
				return nil, fmt.Errorf("failed to parse solana contract address: %w", err)
			}
			// Get access controller from state
			addresses, err := env.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return nil, fmt.Errorf("failed to load addresses for chain %d: %w", selector, err)
			}
			proposerAC, err := getSolAccessControllerFromAddress(addresses, deployment.NewTypeAndVersion(types2.ProposerAccessControllerAccount, deployment.Version1_0_0))
			if err != nil {
				return nil, fmt.Errorf("failed to get proposer access controller: %w", err)
			}
			cancellerAC, err := getSolAccessControllerFromAddress(addresses, deployment.NewTypeAndVersion(types2.CancellerAccessControllerAccount, deployment.Version1_0_0))
			if err != nil {
				return nil, fmt.Errorf("failed to get canceller access controller: %w", err)
			}
			bypasserAC, err := getSolAccessControllerFromAddress(addresses, deployment.NewTypeAndVersion(types2.BypasserAccessControllerAccount, deployment.Version1_0_0))
			if err != nil {
				return nil, fmt.Errorf("failed to get bypasser access controller: %w", err)
			}
			metadata, err := mcmssolanasdk.NewChainMetadata(
				opCount,
				mcmsProgramID,
				mcmSeed,
				proposerAC,
				cancellerAC,
				bypasserAC)
			if err != nil {
				return nil, fmt.Errorf("failed to create chain metadata: %w", err)
			}
			metaDataPerChain[chainID] = metadata
		}
	}

	return metaDataPerChain, nil
}
