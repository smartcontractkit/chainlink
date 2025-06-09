package solana

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"
	"github.com/smartcontractkit/chainlink/deployment"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	seq "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence/operation"
)

const (
	ForwarderContract datastore.ContractType = "Forwarder"
	ForwarderState    datastore.ContractType = "ForwarderState"
)

var _ cldf.ChangeSetV2[*DeployRequest] = DeployForwarder{}

type DeployForwarder struct{}

func (cs DeployForwarder) VerifyPreconditions(env cldf.Environment, req *DeployRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return err
	}

	return nil
}

type DeployRequest = struct {
	ChainSel    uint64
	BuildConfig *helpers.BuildSolanaConfig
	Qualifier   string
	LabelSet    *datastore.LabelSet
	Version     string
}

func (cs DeployForwarder) Apply(env cldf.Environment, req *DeployRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	out.DataStore = datastore.NewMemoryDataStore()
	version := semver.MustParse(req.Version)
	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	deploySeqInput := seq.DeployForwarderSeqInput{
		ChainSel:     req.ChainSel,
		ProgramName:  deployment.KeystoneForwarderProgramName,
		Overallocate: true,
	}

	deps := operation.Deps{
		Datastore: out.DataStore.Seal(),
		Env:       env,
		Chain:     ch,
	}

	deploySeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployForwarderSeq, deps, deploySeqInput)
	if err != nil {
		return out, err
	}

	// save programID
	err = out.DataStore.Addresses().Add(
		datastore.AddressRef{
			Address:       deploySeqReport.Output.ProgramID,
			ChainSelector: req.ChainSel,
			Type:          ForwarderContract,
			Version:       version,
			Qualifier:     req.Qualifier,
			Labels:        *req.LabelSet,
		},
	)

	if err != nil {
		return out, err
	}
	// save StateID
	err = out.DataStore.Addresses().Add(
		datastore.AddressRef{
			Address:       deploySeqReport.Output.StatePubKey,
			ChainSelector: req.ChainSel,
			Type:          ForwarderState,
			Version:       version,
			Qualifier:     req.Qualifier,
			Labels:        *req.LabelSet,
		},
	)

	if err != nil {
		return out, err
	}

	return out, nil
}

type SetForwarderUpgradeAuthorityRequest = struct {
	ChainSel            uint64
	NewUpgradeAuthority solana.PublicKey
	SpillAddress        solana.PublicKey
	Qualifier           string
	Version             string
	MCMS                *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
}

var _ cldf.ChangeSetV2[*SetForwarderUpgradeAuthorityRequest] = SetForwarderUpgradeAuthority{}

type SetForwarderUpgradeAuthority struct{}

func (cs SetForwarderUpgradeAuthority) VerifyPreconditions(env cldf.Environment, req *SetForwarderUpgradeAuthorityRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return err
	}

	forwarderKey := datastore.NewAddressRefKey(req.ChainSel, ForwarderContract, version, req.Qualifier)
	_, err = env.DataStore.Addresses().Get(forwarderKey)
	if err != nil {
		return fmt.Errorf("failed to load forwarder: %w", err)
	}

	if req.MCMS != nil {
		_, err := helpers.FetchTimelockSigner(env, req.ChainSel)
		if err != nil {
			return fmt.Errorf("failed fetch timelock signer: %w")
		}
	}

	return nil
}

func (cs SetForwarderUpgradeAuthority) Apply(env cldf.Environment, req *SetForwarderUpgradeAuthorityRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)

	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	forwarderKey := datastore.NewAddressRefKey(req.ChainSel, ForwarderContract, version, req.Qualifier)
	addr, err := env.DataStore.Addresses().Get(forwarderKey)
	if err != nil {
		return out, fmt.Errorf("failed to load forwarder: %w", err)
	}

	setAuthorityInput := operation.SetUpgradeAuthorityInput{
		ChainSel:            req.ChainSel,
		NewUpgradeAuthority: req.NewUpgradeAuthority.String(),
		MCMS:                req.MCMS,
		ProgramID:           addr.Address,
	}

	deps := operation.Deps{
		Datastore: out.DataStore.Seal(),
		Env:       env,
		Chain:     ch,
	}

	execSetAuthOut, err := operations.ExecuteOperation(env.OperationsBundle, operation.SetUpgradeAuthorityOp, deps, setAuthorityInput)
	if err != nil {
		return out, err
	}

	out.MCMSTimelockProposals = execSetAuthOut.Output.Proposals

	return out, nil
}

type ConfigureForwarderRequest struct {
	WFDonName string
	// workflow don node ids in the offchain client. Used to fetch and derive the signer keys
	WFNodeIDs        []string
	RegistryChainSel uint64

	MCMS *proposalutils.TimelockConfig // if set, assumes current ownership is the timelock

	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains    map[uint64]struct{}
	Qualifier string
	Version   string
}

var _ cldf.ChangeSetV2[*ConfigureForwarderRequest] = ConfigureForwarders{}

type ConfigureForwarders struct{}

func (cs ConfigureForwarders) VerifyPreconditions(env cldf.Environment, req *ConfigureForwarderRequest) error {
	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return err
	}

	if req.Chains != nil {
		for sel := range req.Chains {
			if _, ok := env.BlockChains.SolanaChains()[sel]; !ok {
				return fmt.Errorf("solana chain not found for chain selector %d", sel)
			}
			forwarderKey := datastore.NewAddressRefKey(sel, ForwarderContract, version, req.Qualifier)
			_, err := env.DataStore.Addresses().Get(forwarderKey)

			if err != nil {
				return fmt.Errorf("failed get fowarder for chain selector %d: %w", sel, err)
			}
			if req.MCMS != nil {
				_, err := helpers.FetchTimelockSigner(env, sel)
				if err != nil {
					return fmt.Errorf("failed fetch timelock signer for chain selector %d: %w", sel, err)
				}
			}
		}
	}

	if _, err := internal.NewRegisteredDon(env, internal.RegisteredDonConfig{
		NodeIDs:          req.WFNodeIDs,
		Name:             req.WFDonName,
		RegistryChainSel: req.RegistryChainSel}); err != nil {
		return fmt.Errorf("failed to create registered don: %w", err)
	}

	return nil
}

func (cs ConfigureForwarders) Apply(env cldf.Environment, req *ConfigureForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	// TODO
	return out, nil
}

func ConfigureForwarders_(env cldf.Environment, req ConfigureForwarderRequest) (cldf.ChangesetOutput, error) {
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
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to requre forwarder: %w", err)
	}

	if req.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}

	var proposals []mcms.TimelockProposal
	for chainSel, batch := range mcmsBatches {
		// get timelocks, proposers, inspectors per chain
		solChain := env.BlockChains.SolanaChains()[chainSel]

		addresses, _ := env.ExistingAddresses.AddressesForChain(chainSel)
		mcmState, _ := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
		if mcmState.TimelockProgram.IsZero() {
			return cldf.ChangesetOutput{}, errors.New("timelock is not found")
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
	for _, chain := range env.BlockChains.SolanaChains() {
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
			return nil, fmt.Errorf("failed to requre forwarder for chain selector %d: %w", chain.Selector, err)
		}

		ops[chain.Selector] = op
	}

	return ops, nil
}

func configureForwarder(req ConfigureForwarderRequest, state *state, ch cldfsol.Chain, wfdon *internal.RegisteredDon, owner solana.PublicKey) (mcmsTypes.BatchOperation, error) {
	// 1. derive req pda
	forwarderState := state.forwarderState
	if forwarderState.IsZero() {
		return mcmsTypes.BatchOperation{}, fmt.Errorf("forwarder state not found for chain sel %d", ch.Selector)
	}
	forwarderProgramID := state.forwarderProgramID
	if forwarderProgramID.IsZero() {
		return mcmsTypes.BatchOperation{}, fmt.Errorf("forwarder program not found for chain sel %d", ch.Selector)
	}

	reqPDA := getConfigPDA(forwarderState, wfdon.Info.Id, wfdon.Info.ConfigCount, forwarderProgramID)

	// 2. check if account exists
	var oracleExists bool
	_, err := ch.Client.GetAccountInfo(context.Background(), reqPDA)
	if err != nil {
		if !errors.Is(err, rpc.ErrNotFound) {
			return mcmsTypes.BatchOperation{}, fmt.Errorf("can't confirm oracle existence: %w", err)
		}
		oracleExists = false
	} else {
		oracleExists = true
	}

	// 3. build init/update instructions
	var instructions *ks_forwarder.Instruction
	signers := toSolSigners(wfdon.Signers(chainsel.FamilySolana))
	if !oracleExists {
		instructions, err = ks_forwarder.NewInitOraclesConfigInstruction(
			wfdon.Info.Id,
			wfdon.Info.ConfigCount,
			wfdon.Info.F,
			signers,
			state.forwarderState,
			reqPDA,
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
			signers,
			state.forwarderState,
			reqPDA,
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

func toSolSigners(ss []common.Address) [][20]uint8 {
	ret := make([][20]uint8, 0, len(ss))
	slices.SortFunc(ss, func(a, b common.Address) int {
		return slices.Compare(a.Bytes(), b.Bytes())
	})
	for _, s := range ss {
		ret = append(ret, s)
	}

	return ret
}

func getConfigID(donID uint32, configVersion uint32) uint64 {
	return (uint64(donID) << 32) | uint64(configVersion)
}
