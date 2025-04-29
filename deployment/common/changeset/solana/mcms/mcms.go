package mcmsnew

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// DeployMCMSWithTimelockProgramsSolana deploys an MCMS program
// and initializes 3 instances for each of the timelock roles: Bypasser, ProposerMcm, Canceller on an Solana chain.
// as well as the timelock program. It's not necessarily the only way to use
// the timelock and MCMS, but its reasonable pattern.
func DeployMCMSWithTimelockProgramsSolana(
	e deployment.Environment,
	chain deployment.SolChain,
	addressBook deployment.AddressBook,
	config commontypes.MCMSWithTimelockConfigV2,
) (*state.MCMSWithTimelockStateSolana, error) {
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil && !errors.Is(err, deployment.ErrChainNotFound) {
		return nil, fmt.Errorf("failed to get addresses for chain %v from environment: %w", chain.Selector, err)
	}

	chainState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms with timelock solana chain state: %w", err)
	}

	// access controller
	err = deployAccessControllerProgram(e, chainState, chain, addressBook)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy access controller program: %w", err)
	}
	err = initAccessController(e, chainState, commontypes.ProposerAccessControllerAccount, chain, addressBook)
	if err != nil {
		return nil, fmt.Errorf("failed to init proposer access controller: %w", err)
	}
	err = initAccessController(e, chainState, commontypes.ExecutorAccessControllerAccount, chain, addressBook)
	if err != nil {
		return nil, fmt.Errorf("failed to init access controller: %w", err)
	}
	err = initAccessController(e, chainState, commontypes.CancellerAccessControllerAccount, chain, addressBook)
	if err != nil {
		return nil, fmt.Errorf("failed to init access controller: %w", err)
	}
	err = initAccessController(e, chainState, commontypes.BypasserAccessControllerAccount, chain, addressBook)
	if err != nil {
		return nil, fmt.Errorf("failed to init access controller: %w", err)
	}

	// mcm
	err = deployMCMProgram(e, chainState, chain, addressBook)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy mcm program: %w", err)
	}
	err = initMCM(e, chainState, commontypes.BypasserManyChainMultisig, chain, addressBook, &config.Bypasser)
	if err != nil {
		return nil, fmt.Errorf("failed to init bypasser mcm: %w", err)
	}
	err = initMCM(e, chainState, commontypes.CancellerManyChainMultisig, chain, addressBook, &config.Canceller)
	if err != nil {
		return nil, fmt.Errorf("failed to init canceller mcm: %w", err)
	}
	err = initMCM(e, chainState, commontypes.ProposerManyChainMultisig, chain, addressBook, &config.Proposer)
	if err != nil {
		return nil, fmt.Errorf("failed to init proposer mcm: %w", err)
	}

	// timelock
	err = deployTimelockProgram(e, chainState, chain, addressBook)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy timelock program: %w", err)
	}
	err = initTimelock(e, chainState, chain, addressBook, config.TimelockMinDelay)
	if err != nil {
		return nil, fmt.Errorf("failed to init timelock: %w", err)
	}

	err = setupRoles(chainState, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to setup roles and ownership: %w", err)
	}

	return chainState, nil
}
