package solana

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func ValidateMCMSConfigSolana(
	e deployment.Environment,
	mcms *MCMSConfigSolana,
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	tokenAddress solana.PublicKey,
) error {
	if mcms != nil {
		if mcms.MCMS == nil {
			return errors.New("MCMS config is nil")
		}
		if !mcms.FeeQuoterOwnedByTimelock &&
			!mcms.RouterOwnedByTimelock &&
			!mcms.OffRampOwnedByTimelock &&
			!mcms.RMNRemoteOwnedByTimelock &&
			!mcms.BurnMintTokenPoolOwnedByTimelock[tokenAddress] &&
			!mcms.LockReleaseTokenPoolOwnedByTimelock[tokenAddress] {
			return errors.New("at least one of the MCMS components must be owned by the timelock")
		}
		if err := mcms.MCMS.ValidateSolana(e, chain.Selector); err != nil {
			return fmt.Errorf("failed to validate MCMS config: %w", err)
		}
	}
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, mcms != nil && mcms.FeeQuoterOwnedByTimelock, chainState.FeeQuoter, cs.FeeQuoter, tokenAddress); err != nil {
		return fmt.Errorf("failed to validate ownership for fee quoter: %w", err)
	}
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, mcms != nil && mcms.RouterOwnedByTimelock, chainState.Router, cs.Router, tokenAddress); err != nil {
		return fmt.Errorf("failed to validate ownership for router: %w", err)
	}
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, mcms != nil && mcms.OffRampOwnedByTimelock, chainState.OffRamp, cs.OffRamp, tokenAddress); err != nil {
		return fmt.Errorf("failed to validate ownership for off ramp: %w", err)
	}
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, mcms != nil && mcms.RMNRemoteOwnedByTimelock, chainState.RMNRemote, cs.RMNRemote, tokenAddress); err != nil {
		return fmt.Errorf("failed to validate ownership for rmnremote: %w", err)
	}
	if !tokenAddress.IsZero() {
		if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, mcms != nil && mcms.BurnMintTokenPoolOwnedByTimelock[tokenAddress], chainState.BurnMintTokenPool, cs.BurnMintTokenPool, tokenAddress); err != nil {
			return fmt.Errorf("failed to validate ownership for burnmint: %w", err)
		}
		if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, mcms != nil && mcms.LockReleaseTokenPoolOwnedByTimelock[tokenAddress], chainState.LockReleaseTokenPool, cs.LockReleaseTokenPool, tokenAddress); err != nil {
			return fmt.Errorf("failed to validate ownership for lockrelease: %w", err)
		}
	}

	return nil
}

func BuildProposalsForTxns(
	e deployment.Environment,
	chainSelector uint64,
	description string,
	minDelay time.Duration,
	txns []mcmsTypes.Transaction) (*mcms.TimelockProposal, error) {
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]sdk.Inspector{}
	batches := make([]mcmsTypes.BatchOperation, 0)
	chain := e.SolChains[chainSelector]
	addresses, _ := e.ExistingAddresses.AddressesForChain(chainSelector)
	mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)

	timelocks[chainSelector] = mcmsSolana.ContractAddress(
		mcmState.TimelockProgram,
		mcmsSolana.PDASeed(mcmState.TimelockSeed),
	)
	proposers[chainSelector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
	inspectors[chainSelector] = mcmsSolana.NewInspector(chain.Client)
	batches = append(batches, mcmsTypes.BatchOperation{
		ChainSelector: mcmsTypes.ChainSelector(chainSelector),
		Transactions:  txns,
	})
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		description,
		proposalutils.TimelockConfig{MinDelay: minDelay})
	if err != nil {
		return nil, fmt.Errorf("failed to build proposal: %w", err)
	}
	return proposal, nil
}

func BuildMCMSTxn(ixn solana.Instruction, programID string, contractType deployment.ContractType) (*mcmsTypes.Transaction, error) {
	data, err := ixn.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to extract data: %w", err)
	}
	for _, account := range ixn.Accounts() {
		if account.IsSigner {
			account.IsSigner = false
		}
	}
	tx, err := mcmsSolana.NewTransaction(
		programID,
		data,
		big.NewInt(0),        // e.g. value
		ixn.Accounts(),       // pass along needed accounts
		string(contractType), // some string identifying the target
		[]string{},           // any relevant metadata
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return &tx, nil
}

func FetchTimelockSigner(e deployment.Environment, chainSelector uint64) (solana.PublicKey, error) {
	addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to load addresses for chain %d: %w", chainSelector, err)
	}
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[chainSelector], addresses)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to load mcm state: %w", err)
	}
	timelockSignerPDA := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	return timelockSignerPDA, nil
}

func GetAuthorityForIxn(
	e *deployment.Environment,
	chain deployment.SolChain,
	mcms *MCMSConfigSolana,
	contractType deployment.ContractType,
	tokenAddress solana.PublicKey, // used for burnmint and lockrelease
) (solana.PublicKey, error) {
	if mcms == nil {
		return chain.DeployerKey.PublicKey(), nil
	}
	timelockSigner, err := FetchTimelockSigner(*e, chain.Selector)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	switch contractType {
	case cs.FeeQuoter:
		if mcms.FeeQuoterOwnedByTimelock {
			return timelockSigner, nil
		}
		return chain.DeployerKey.PublicKey(), nil
	case cs.Router:
		if mcms.RouterOwnedByTimelock {
			return timelockSigner, nil
		}
		return chain.DeployerKey.PublicKey(), nil
	case cs.OffRamp:
		if mcms.OffRampOwnedByTimelock {
			return timelockSigner, nil
		}
		return chain.DeployerKey.PublicKey(), nil
	case cs.BurnMintTokenPool:
		if mcms.BurnMintTokenPoolOwnedByTimelock[tokenAddress] {
			return timelockSigner, nil
		}
		return chain.DeployerKey.PublicKey(), nil
	case cs.LockReleaseTokenPool:
		if mcms.LockReleaseTokenPoolOwnedByTimelock[tokenAddress] {
			return timelockSigner, nil
		}
		return chain.DeployerKey.PublicKey(), nil
	case cs.RMNRemote:
		if mcms.RMNRemoteOwnedByTimelock {
			return timelockSigner, nil
		}
		return chain.DeployerKey.PublicKey(), nil
	default:
		return solana.PublicKey{}, fmt.Errorf("invalid contract type: %s", contractType)
	}
}
