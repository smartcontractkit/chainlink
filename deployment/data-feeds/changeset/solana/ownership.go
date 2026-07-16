package solana

import (
	"github.com/gagliardetto/solana-go"

	solchangesets "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana/changesets"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
)

// TransferOwnershipCacheRequest wraps the generic request for cache contracts
type TransferOwnershipCacheRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey
	Version                     string
	Qualifier                   string
	MCMSCfg                     cldfproposalutils.TimelockConfig
}

// TransferOwnershipCache implementation
var _ cldf.ChangeSetV2[*TransferOwnershipCacheRequest] = TransferOwnershipCache{}

type TransferOwnershipCache struct{}

func (cs TransferOwnershipCache) VerifyPreconditions(env cldf.Environment, req *TransferOwnershipCacheRequest) error {
	return solchangesets.GenericVerifyPreconditions(env, req.ChainSel, req.Version, req.Qualifier, "CacheContract")
}

func (cs TransferOwnershipCache) Apply(env cldf.Environment, req *TransferOwnershipCacheRequest) (cldf.ChangesetOutput, error) {
	genericReq := &solchangesets.TransferOwnershipRequest{
		ChainSel:      req.ChainSel,
		CurrentOwner:  req.CurrentOwner,
		ProposedOwner: req.ProposedOwner,
		Version:       req.Version,
		Qualifier:     req.Qualifier,
		MCMSCfg:       req.MCMSCfg,
		ContractConfig: solchangesets.ContractConfig{
			ContractType: CacheContract,
			StateType:    CacheState,
			OperationID:  "transfer-ownership-cache",
			Description:  "transfers ownership of cache to mcms",
		},
	}
	return solchangesets.GenericTransferOwnership(env, genericReq)
}
