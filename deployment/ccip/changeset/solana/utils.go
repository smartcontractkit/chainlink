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
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func ValidateMCMSConfigSolana(e deployment.Environment, chainSelector uint64, mcms *MCMSConfigSolana) error {
	if mcms != nil {
		if mcms.MCMS == nil {
			return errors.New("MCMS config is nil")
		}
		if !mcms.FeeQuoterOwnedByTimelock && !mcms.RouterOwnedByTimelock && !mcms.OffRampOwnedByTimelock {
			return errors.New("at least one of the MCMS components must be owned by the timelock")
		}
		return ValidateMCMSConfig(e, chainSelector, mcms.MCMS)
	}
	return nil
}

func ValidateMCMSConfig(e deployment.Environment, chainSelector uint64, mcms *cs.MCMSConfig) error {
	if mcms != nil {
		// If there is no timelock and mcms proposer on the chain, the transfer will fail.
		timelockID, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.RBACTimelock)
		if err != nil {
			return fmt.Errorf("timelock not present on the chain %w", err)
		}
		proposerID, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.ProposerManyChainMultisig)
		if err != nil {
			return fmt.Errorf("mcms proposer not present on the chain %w", err)
		}
		// Make sure addresses are correctly parsed. Format is: "programID.PDASeed"
		_, _, err = mcmsSolana.ParseContractAddress(timelockID)
		if err != nil {
			return fmt.Errorf("failed to parse timelock address: %w", err)
		}
		_, _, err = mcmsSolana.ParseContractAddress(proposerID)
		if err != nil {
			return fmt.Errorf("failed to parse proposer address: %w", err)
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
		minDelay)
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
