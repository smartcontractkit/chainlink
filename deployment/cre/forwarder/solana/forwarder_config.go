package solana

import (
	"encoding/binary"
	"fmt"
	"iter"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	solstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"

	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana/sequence/operation"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

// verifyForwarderChains checks that every explicitly requested chain exists, holds a deployed
// forwarder and, when the change goes through the timelock, has MCMS deployed. An empty chain set
// means all Solana chains of the environment and is only resolved at apply time.
func verifyForwarderChains(env cldf.Environment, chains map[uint64]struct{}, version *semver.Version, qualifier string, mcmsCfg *cldfproposalutils.TimelockConfig) error {
	for sel := range chains {
		if _, ok := env.BlockChains.SolanaChains()[sel]; !ok {
			return fmt.Errorf("solana chain not found for chain selector %d", sel)
		}

		forwarderKey := datastore.NewAddressRefKey(sel, ForwarderContract, version, qualifier)
		if _, err := env.DataStore.Addresses().Get(forwarderKey); err != nil {
			return fmt.Errorf("failed get forwarder for chain selector %d: %w", sel, err)
		}

		if mcmsCfg != nil {
			_, err := solstate.MaybeLoadMCMSWithTimelockChainStateV2(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(sel)))
			if err != nil {
				return fmt.Errorf("failed to load MCMS for chain selector %d: %w", sel, err)
			}
		}
	}

	return nil
}

// forwarderChains yields the Solana chains a request applies to. An empty chain set means all
// Solana chains of the environment.
func forwarderChains(env cldf.Environment, chains map[uint64]struct{}) iter.Seq[cldfsol.Chain] {
	return func(yield func(cldfsol.Chain) bool) {
		for _, chain := range env.BlockChains.SolanaChains() {
			if _, shouldInclude := chains[chain.Selector]; len(chains) > 0 && !shouldInclude {
				continue
			}

			if !yield(chain) {
				return
			}
		}
	}
}

// resolveForwarderConfigTarget loads the forwarder addresses of a chain from the datastore and
// derives the oracles config account of the given DON.
func resolveForwarderConfigTarget(
	env cldf.Environment,
	chain cldfsol.Chain,
	version *semver.Version,
	qualifier string,
	donID uint32,
	configVersion uint32,
	mcmsCfg *cldfproposalutils.TimelockConfig,
) (operation.ForwarderConfigTarget, error) {
	var target operation.ForwarderConfigTarget

	stateRef, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(chain.Selector, ForwarderState, version, qualifier))
	if err != nil {
		return target, fmt.Errorf("failed load forwarder state: %w", err)
	}
	forwarderRef, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(chain.Selector, ForwarderContract, version, qualifier))
	if err != nil {
		return target, fmt.Errorf("failed load forwarder: %w", err)
	}

	state, err := solana.PublicKeyFromBase58(stateRef.Address)
	if err != nil {
		return target, fmt.Errorf("failed parse forwarder state %q: %w", stateRef.Address, err)
	}
	programID, err := solana.PublicKeyFromBase58(forwarderRef.Address)
	if err != nil {
		return target, fmt.Errorf("failed parse forwarder program id %q: %w", forwarderRef.Address, err)
	}

	owner := chain.DeployerKey.PublicKey()
	if mcmsCfg != nil {
		owner, err = helpers.FetchTimelockSigner(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chain.Selector)))
		if err != nil {
			return target, fmt.Errorf("failed fetch timelock signer: %w", err)
		}
	}

	return operation.ForwarderConfigTarget{
		MCMS:           mcmsCfg,
		ProgramID:      programID,
		ForwarderState: state,
		ConfigPDA:      getConfigPDA(state, donID, configVersion, programID),
		Owner:          owner,
		DonID:          donID,
		ConfigVersion:  configVersion,
		Type:           cldf.ContractType(ForwarderContract),
	}, nil
}

// buildTimelockProposals wraps one batch operation per chain into a timelock proposal.
func buildTimelockProposals(
	env cldf.Environment,
	batches map[uint64]mcmsTypes.BatchOperation,
	mcmsCfg cldfproposalutils.TimelockConfig,
	description string,
) ([]mcms.TimelockProposal, error) {
	proposals := make([]mcms.TimelockProposal, 0, len(batches))

	for chainSel, batch := range batches {
		solChain := env.BlockChains.SolanaChains()[chainSel]

		addresses := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSel))
		mcmState, err := solstate.MaybeLoadMCMSWithTimelockChainStateV2(addresses)
		if err != nil {
			return nil, fmt.Errorf("failed to load MCMS for chain selector %d: %w", chainSel, err)
		}
		if mcmState.TimelockProgram.IsZero() {
			return nil, fmt.Errorf("timelock is not found for chain selector %d", chainSel)
		}

		timelocks := map[uint64]string{
			solChain.Selector: mcmsSolana.ContractAddress(mcmState.TimelockProgram, mcmsSolana.PDASeed(mcmState.TimelockSeed)),
		}
		proposers := map[uint64]string{
			solChain.Selector: mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed)),
		}
		inspectors := map[uint64]sdk.Inspector{
			solChain.Selector: mcmsSolana.NewInspector(solChain.Client),
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(
			env,
			timelocks,
			proposers,
			inspectors,
			[]mcmsTypes.BatchOperation{batch},
			description,
			mcmsCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build proposal for chain selector %d: %w", chainSel, err)
		}

		proposals = append(proposals, *proposal)
	}

	return proposals, nil
}

func getConfigPDA(statePubkey solana.PublicKey, donID uint32, configVersion uint32, programID solana.PublicKey) solana.PublicKey {
	configID := getConfigID(donID, configVersion)
	reqIDBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(reqIDBytes, configID)

	seeds := [][]byte{
		[]byte("config"),
		statePubkey.Bytes(),
		reqIDBytes,
	}

	addr, _, _ := solana.FindProgramAddress(seeds, programID)

	return addr
}

func getConfigID(donID uint32, configVersion uint32) uint64 {
	return (uint64(donID) << 32) | uint64(configVersion)
}
