package solana

import (
	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

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

	solChain, _ := env.BlockChains.SolanaChains()[req.ChainSel]

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

	out.MCMSTimelockProposals = execOut.Output.Proposals

	return out, nil
}
