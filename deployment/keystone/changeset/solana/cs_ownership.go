package solana

import (
	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// TransferOwnershipForwarderRequest wraps the generic request for forwarder contracts
type TransferOwnershipForwarderRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey
	Version                     string
	Qualifier                   string
	MCMSCfg                     proposalutils.TimelockConfig
}

// TransferOwnershipForwarder implementation
var _ cldf.ChangeSetV2[*TransferOwnershipForwarderRequest] = TransferOwnershipForwarder{}

type TransferOwnershipForwarder struct{}

func (cs TransferOwnershipForwarder) VerifyPreconditions(env cldf.Environment, req *TransferOwnershipForwarderRequest) error {
	return commonchangeset.GenericVerifyPreconditions(env, req.ChainSel, req.Version, req.Qualifier, "ForwarderContract")
}

func (cs TransferOwnershipForwarder) Apply(env cldf.Environment, req *TransferOwnershipForwarderRequest) (cldf.ChangesetOutput, error) {
	genericReq := &commonchangeset.TransferOwnershipRequest{
		ChainSel:      req.ChainSel,
		CurrentOwner:  req.CurrentOwner,
		ProposedOwner: req.ProposedOwner,
		Version:       req.Version,
		Qualifier:     req.Qualifier,
		MCMSCfg:       req.MCMSCfg,
		ContractConfig: commonchangeset.ContractConfig{
			ContractType: "ForwarderContract",
			StateType:    "ForwarderState",
			OperationID:  "transfer-ownership-forwarder",
			Description:  "transfers ownership of forwarder to mcms",
		},
	}
	return commonchangeset.GenericTransferOwnership(env, genericReq)
}
