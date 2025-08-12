package mcmsnew

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms/sequence"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms/sequence/operation"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// DeployMCMSWithTimelockProgramsSolana deploys an MCMS program
// and initializes 3 instances for each of the timelock roles: Bypasser, ProposerMcm, Canceller on an Solana chain.
// as well as the timelock program. It's not necessarily the only way to use
// the timelock and MCMS, but its reasonable pattern.
func DeployMCMSWithTimelockProgramsSolana(
	e cldf.Environment,
	chain cldf_solana.Chain,
	addressBook cldf.AddressBook,
	config commontypes.MCMSWithTimelockConfigV2,
) (*state.MCMSWithTimelockStateSolana, error) {
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil && !errors.Is(err, cldf.ErrChainNotFound) {
		return nil, fmt.Errorf("failed to get addresses for chain %v from environment: %w", chain.Selector, err)
	}

	chainState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms with timelock solana chain state: %w", err)
	}

	// access controller
	err = deployAccessControllerProgram(e, chainState, chain, addressBook)
	err = waitForProgramDeployment(e.GetContext(), chain.Client, chainState.AccessControllerProgram, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("access controller program not ready: %w", err)
	}
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
	err = waitForProgramDeployment(e.GetContext(), chain.Client, chainState.AccessControllerProgram, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("access controller program not ready: %w", err)
	}
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
	err = waitForProgramDeployment(e.GetContext(), chain.Client, chainState.AccessControllerProgram, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("access controller program not ready: %w", err)
	}
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

// DeployMCMSWithTimelockProgramsSolanaV2 deploys an MCMS program using Operations API
// saves addresses to datastore
func DeployMCMSWithTimelockProgramsSolanaV2(
	e cldf.Environment,
	ds datastore.MutableDataStore,
	chain cldf_solana.Chain,
	config commontypes.MCMSWithTimelockConfigV2) (*state.MCMSWithTimelockStateSolana, error) {
	chainstate, err := state.MaybeLoadMCMSWithTimelockChainStateSolanaV2(e.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chain.Selector)))
	if err != nil {
		return nil, err
	}

	deps := operation.Deps{
		State:     chainstate,
		Chain:     chain,
		Datastore: ds,
	}

	_, err = operations.ExecuteSequence(e.OperationsBundle, sequence.DeployMCMSWithTimelockSeq, deps, sequence.DeployMCMSWithTimelockInput{
		MCMConfig:        config,
		TimelockMinDelay: config.TimelockMinDelay,
	})
	if err != nil {
		return nil, err
	}

	return chainstate, nil
}

func waitForProgramDeployment(ctx context.Context, client *rpc.Client, programID solana.PublicKey, maxWait time.Duration) error {
	timeout := time.After(maxWait)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out waiting for program %s to be deployed", programID.String())
		case <-ticker.C:
			resp, err := client.GetAccountInfo(ctx, programID)
			if err != nil {
				continue // Retry on RPC errors
			}
			if resp != nil && resp.Value != nil && resp.Value.Executable {
				return nil // Ready
			}
		}
	}
}
