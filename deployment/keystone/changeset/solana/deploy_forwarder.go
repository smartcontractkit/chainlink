package solana

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"
	"github.com/smartcontractkit/chainlink/deployment"
	cdeployment "github.com/smartcontractkit/chainlink/deployment"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/shared"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

type DeployRequest = struct {
	ChainSel    uint64
	BuildConfig *helpers.BuildSolanaConfig
}

var _ cldf.ChangeSet[*DeployRequest] = DeployForwarder

func DeployForwarder(env cldf.Environment, req *DeployRequest) (cldf.ChangesetOutput, error) {
	if req.BuildConfig != nil {
		err := helpers.BuildSolana(env, *req.BuildConfig, keystoneBuildParams)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	chain := env.SolChains[req.ChainSel]
	ab := cldf.NewMemoryAddressBook()

	address, err := helpers.DeployAndMaybeSaveToAddressBook(env, chain, ab, shared.Forwarder, cdeployment.Version1_0_0, false, "")
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// initialize
	ks_forwarder.SetProgramID(address)

	stateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create random keys: %w", err)
	}

	instruction, err := ks_forwarder.NewInitializeInstruction(stateKey.PublicKey(), chain.DeployerKey.PublicKey(), solana.SystemProgramID).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	instructions := []solana.Instruction{instruction}
	if err = chain.Confirm(instructions, solanaUtils.AddSigners(stateKey)); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm ")
	}

	tv := cldf.NewTypeAndVersion(shared.ForwarderState, deployment.Version1_0_0)
	err = ab.Save(req.ChainSel, stateKey.PublicKey().String(), tv)

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save forwarder state address: %w", err)
	}
	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

type SetForwarderUpgradeAuthorityRequest = struct {
	ChainSel            uint64
	NewUpgradeAuthority solana.PublicKey
	MCMS                *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
}

var _ cldf.ChangeSet[*SetForwarderUpgradeAuthorityRequest] = SetForwarderUpgradeAuthority

func SetForwarderUpgradeAuthority(env cldf.Environment, req *SetForwarderUpgradeAuthorityRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.SolChains[req.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("can't get chain for chain selector %d", req.ChainSel)
	}

	state, err := loadOnchainState(env, req.ChainSel)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if state.forwarderProgramID.IsZero() {
		return cldf.ChangesetOutput{}, fmt.Errorf("forwarder not found for chain selector %d", req.ChainSel)
	}

	currentAuthority := chain.DeployerKey.PublicKey()
	if req.MCMS != nil {
		timelockSignerPDA, err := helpers.FetchTimelockSigner(env, chain.Selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}

	env.Logger.Infow("Setting upgrade authority for", state.forwarderProgramID.String(), "newUpgradeAuthority", req.NewUpgradeAuthority.String())
	mcmsTxns := make([]mcmsTypes.Transaction, 0)
	ixn := helpers.SetUpgradeAuthority(&env, state.forwarderProgramID, currentAuthority, req.NewUpgradeAuthority, false)
	if req.MCMS == nil {
		if err := chain.Confirm([]solana.Instruction{ixn}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	} else {
		tx, err := helpers.BuildMCMSTxn(
			ixn,
			solana.BPFLoaderUpgradeableProgramID.String(),
			cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()))
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		mcmsTxns = append(mcmsTxns, *tx)
	}
	if len(mcmsTxns) > 0 {
		proposal, err := helpers.BuildProposalsForTxns(
			env, req.ChainSel, "proposal to SetUpgradeAuthority in Solana", req.MCMS.MinDelay, mcmsTxns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

type ConfigureForwarderRequest struct {
	WFDonName string
	// workflow don node ids in the offchain client. Used to fetch and derive the signer keys
	WFNodeIDs        []string
	RegistryChainSel uint64

	MCMS *proposalutils.TimelockConfig // if set, assumes current ownership is the timelock

	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains map[uint64]struct{}
}

func ConfigureForwarders(env cldf.Environment, req ConfigureForwarderRequest) (cldf.ChangesetOutput, error) {
	wfDon, err := internal.NewRegisteredDon(env, internal.RegisteredDonConfig{
		NodeIDs:          req.WFNodeIDs,
		Name:             req.WFDonName,
		RegistryChainSel: req.RegistryChainSel,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create registered don: %w", err)
	}

	mcmsBatches, err := configureForwarders(env, req, wfDon)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure forwarder: %w", err)
	}

	if req.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}

	var proposals []mcms.TimelockProposal
	for chainSel, batch := range mcmsBatches {
		// get timelocks, proposers, inspectors per chain
		solChain := env.SolChains[chainSel]

		addresses, _ := env.ExistingAddresses.AddressesForChain(chainSel)
		mcmState, _ := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
		if mcmState.TimelockProgram.IsZero() {
			return cldf.ChangesetOutput{}, fmt.Errorf("timelock is not found")
		}

		timelocks := map[uint64]string{}
		proposers := map[uint64]string{}
		inspectors := map[uint64]sdk.Inspector{}
		timelocks[solChain.Selector] = mcmsSolana.ContractAddress(
			mcmState.TimelockProgram,
			mcmsSolana.PDASeed(mcmState.TimelockSeed),
		)

		proposers[solChain.Selector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
		inspectors[solChain.Selector] = mcmsSolana.NewInspector(solChain.Client)
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocks,
			proposers,
			inspectors,
			[]mcmsTypes.BatchOperation{batch},
			"proposal to transfer ownership of keystone forwarder contract to timelock",
			*req.MCMS)

		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
		proposals = append(proposals, *proposal)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: proposals,
	}, nil
}

func configureForwarders(env cldf.Environment, req ConfigureForwarderRequest,
	wfdon *internal.RegisteredDon) (map[uint64]mcmsTypes.BatchOperation, error) {
	ops := make(map[uint64]mcmsTypes.BatchOperation)
	for _, chain := range env.SolChains {
		if _, shouldInclude := req.Chains[chain.Selector]; len(req.Chains) > 0 && !shouldInclude {
			continue
		}
		st, err := loadOnchainState(env, chain.Selector)
		if err != nil {
			return nil, fmt.Errorf("failed to load onchain state for chain selector %d: %w", chain.Selector, err)
		}

		owner := chain.DeployerKey.PublicKey()
		if req.MCMS != nil {
			timelockPDA, err := helpers.FetchTimelockSigner(env, chain.Selector)
			if err != nil {
				return nil, err
			}
			owner = timelockPDA
		}

		op, err := configureForwarder(req, st, chain, wfdon, owner)
		if err != nil {
			return nil, fmt.Errorf("failed to configure forwarder for chain selector %d: %w", chain.Selector, err)
		}

		ops[chain.Selector] = op
	}

	return ops, nil
}

func configureForwarder(req ConfigureForwarderRequest, state *state, ch cldf.SolChain, wfdon *internal.RegisteredDon, owner solana.PublicKey) (mcmsTypes.BatchOperation, error) {
	// 1. derive config pda
	forwarderState := state.forwarderState
	if forwarderState.IsZero() {
		return mcmsTypes.BatchOperation{}, fmt.Errorf("forwarder state not found for chain sel %d", ch.Selector)
	}
	forwarderProgramID := state.forwarderProgramID
	if forwarderProgramID.IsZero() {
		return mcmsTypes.BatchOperation{}, fmt.Errorf("forwarder program not found for chain sel %d", ch.Selector)
	}

	configPDA := getConfigPDA(forwarderState, wfdon.Info.Id, wfdon.Info.ConfigCount, forwarderProgramID)

	// 2. check if account exists
	var oracleExists bool
	_, err := ch.Client.GetAccountInfo(context.Background(), configPDA)
	if err != nil {
		if errors.Is(err, rpc.ErrNotFound) {
			oracleExists = false
		} else {
			return mcmsTypes.BatchOperation{}, fmt.Errorf("can't confirm oracle existence: %w", err)
		}
	} else {
		oracleExists = true
	}

	// 3. build init/update instructions
	var instructions *ks_forwarder.Instruction
	if !oracleExists {
		instructions, err = ks_forwarder.NewInitOraclesConfigInstruction(
			wfdon.Info.Id,
			wfdon.Info.ConfigCount,
			wfdon.Info.F,
			toSolSigners(wfdon.Signers(chainsel.FamilySolana)),
			state.forwarderState,
			configPDA,
			owner,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return mcmsTypes.BatchOperation{}, fmt.Errorf("cant build init oracle instruction: %w", err)
		}
	} else {
		instructions, err = ks_forwarder.NewUpdateOraclesConfigInstruction(
			wfdon.Info.Id,
			wfdon.Info.ConfigCount,
			wfdon.Info.F,
			toSolSigners(wfdon.Signers(chainsel.FamilySolana)),
			state.forwarderState,
			configPDA,
			owner,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return mcmsTypes.BatchOperation{}, fmt.Errorf("cant build init oracle instruction: %w", err)
		}
	}

	if req.MCMS == nil {
		err := ch.Confirm([]solana.Instruction{instructions})
		return mcmsTypes.BatchOperation{}, err
	}

	// 4. build mcms proposal
	tx, err := helpers.BuildMCMSTxn(
		instructions,
		solana.BPFLoaderUpgradeableProgramID.String(),
		cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()))
	if err != nil {
		return mcmsTypes.BatchOperation{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return mcmsTypes.BatchOperation{
		ChainSelector: mcmsTypes.ChainSelector(ch.Selector),
		Transactions:  []mcmsTypes.Transaction{*tx},
	}, nil

}

func getConfigPDA(statePubkey solana.PublicKey, donID uint32, configVersion uint32, programID solana.PublicKey) solana.PublicKey {
	configID := getConfigID(donID, configVersion)
	configIDBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(configIDBytes, configID)

	seeds := [][]byte{
		[]byte("config"),
		statePubkey.Bytes(),
		configIDBytes,
	}

	addr, _, _ := solana.FindProgramAddress(seeds, programID)
	return addr
}

func toSolSigners(ss []common.Address) [][20]uint8 {
	ret := make([][20]uint8, len(ss))
	for _, s := range ss {
		ret = append(ret, s)
	}

	return ret
}

func getConfigID(donID uint32, configVersion uint32) uint64 {
	return (uint64(donID) << 32) | uint64(configVersion)
}
