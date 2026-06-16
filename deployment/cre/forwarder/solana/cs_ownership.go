package solana

import (
	"github.com/gagliardetto/solana-go"

	solchangesets "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana/changesets"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
)

// TransferOwnershipForwarderRequest wraps the generic request for forwarder contracts
type TransferOwnershipForwarderRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey
	Version                     string
	Qualifier                   string
	MCMSCfg                     cldfproposalutils.TimelockConfig
}

// TransferOwnershipForwarder implementation
var _ cldf.ChangeSetV2[*TransferOwnershipForwarderRequest] = TransferOwnershipForwarder{}

type TransferOwnershipForwarder struct{}

func (cs TransferOwnershipForwarder) VerifyPreconditions(env cldf.Environment, req *TransferOwnershipForwarderRequest) error {
	return solchangesets.GenericVerifyPreconditions(env, req.ChainSel, req.Version, req.Qualifier, ForwarderContract)
}

func (cs TransferOwnershipForwarder) Apply(env cldf.Environment, req *TransferOwnershipForwarderRequest) (cldf.ChangesetOutput, error) {
	genericReq := &solchangesets.TransferOwnershipRequest{
		ChainSel:      req.ChainSel,
		CurrentOwner:  req.CurrentOwner,
		ProposedOwner: req.ProposedOwner,
		Version:       req.Version,
		Qualifier:     req.Qualifier,
		MCMSCfg:       req.MCMSCfg,
		ContractConfig: solchangesets.ContractConfig{
			ContractType: ForwarderContract,
			StateType:    ForwarderState,
			OperationID:  "transfer-ownership-forwarder",
			Description:  "transfers ownership of forwarder to mcms",
		},
	}
	return solchangesets.GenericTransferOwnership(env, genericReq)
}
