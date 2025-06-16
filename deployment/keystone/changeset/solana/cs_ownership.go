package solana

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type TransferOwnershipForwarderRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey
	Version                     string
	Qualifier                   string

	// MCMSCfg is for the accept ownership proposal
	MCMSCfg proposalutils.TimelockConfig
}

var _ cldf.ChangeSetV2[*TransferOwnershipForwarderRequest] = TransferOwnershipForwarder{}

type TransferOwnershipForwarder struct{}

func (cs TransferOwnershipForwarder) VerifyPreconditions(env cldf.Environment, req *TransferOwnershipForwarderRequest) error {
	sel := req.ChainSel

	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return err
	}

	if _, ok := env.BlockChains.SolanaChains()[sel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", sel)
	}

	forwarderKey := datastore.NewAddressRefKey(sel, ForwarderContract, version, req.Qualifier)
	_, err = env.DataStore.Addresses().Get(forwarderKey)

	if err != nil {
		return fmt.Errorf("failed get fowarder for chain selector %d: %w", sel, err)
	}
	return nil
}

func (cs TransferOwnershipForwarder) Apply(env cldf.Environment, req *TransferOwnershipForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	version := semver.MustParse(req.Version)
	forwarderStateRef := datastore.NewAddressRefKey(req.ChainSel, ForwarderState, version, req.Qualifier)
	forwarderRef := datastore.NewAddressRefKey(req.ChainSel, ForwarderContract, version, req.Qualifier)

	forwarder, _ := env.DataStore.Addresses().Get(forwarderRef)
	forwarderState, _ := env.DataStore.Addresses().Get(forwarderStateRef)

	mcmsState, err := state.MaybeLoadMCMSWithTimelockChainStateSolanaV2(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(req.ChainSel)))
	if err != nil {
		return out, err
	}

	solChain := env.BlockChains.SolanaChains()[req.ChainSel]

	execOut, err := operations.ExecuteOperation(env.OperationsBundle,
		operations.NewOperation(
			"transfer-ownership-forwarder",
			version,
			"transfers ownership of forwarder to mcms",
			commonchangeset.TransferToTimelockSolanaOp,
		),
		commonchangeset.Deps{
			Env:   env,
			State: mcmsState,
			Chain: solChain,
		},
		commonchangeset.TransferToTimelockInput{
			Contract: commonchangeset.OwnableContract{
				Type:      cldf.ContractType(ForwarderContract),
				ProgramID: solana.MustPublicKeyFromBase58(forwarder.Address),
				OwnerPDA:  solana.MustPublicKeyFromBase58(forwarderState.Address),
			},
			MCMSCfg: req.MCMSCfg,
		},
	)
	if err != nil {
		return out, err
	}

	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}

	inspectors[req.ChainSel] = mcmssolanasdk.NewInspector(solChain.Client)
	timelocks[req.ChainSel] = mcmssolanasdk.ContractAddress(mcmsState.TimelockProgram, mcmssolanasdk.PDASeed(mcmsState.TimelockSeed))
	proposers[req.ChainSel] = mcmssolanasdk.ContractAddress(mcmsState.McmProgram, mcmssolanasdk.PDASeed(mcmsState.ProposerMcmSeed))

	// create timelock proposal with accept transactions
	proposal, err := proposalutils.BuildProposalFromBatchesV2(env, timelocks, proposers, inspectors,
		execOut.Output.Batches, "proposal to transfer ownership of contracts to timelock", req.MCMSCfg)
	if err != nil {
		return out, fmt.Errorf("failed to build proposal: %w", err)
	}
	env.Logger.Debugw("created timelock proposal", "# batches", len(execOut.Output.Batches))

	out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}

	return out, nil
}
