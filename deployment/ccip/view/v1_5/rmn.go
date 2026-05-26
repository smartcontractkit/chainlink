package v1_5

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	commoncldchangesets "github.com/smartcontractkit/cld-changesets/pkg/cldfutil"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
)

type RMNView struct {
	commoncldchangesets.ContractMetaData
	ConfigDetails            rmn_contract.GetConfigDetails `json:"configDetails"`
	PermaBlessedCommitStores []common.Address              `json:"permaBlessedCommitStores"`
}

func GenerateRMNView(r *rmn_contract.RMNContract) (RMNView, error) {
	if r == nil {
		return RMNView{}, errors.New("cannot generate view for nil RMN")
	}
	meta, err := commoncldchangesets.NewContractMetaData(r, r.Address())
	if err != nil {
		return RMNView{}, fmt.Errorf("failed to generate contract metadata for RMN: %w", err)
	}
	config, err := r.GetConfigDetails(nil)
	if err != nil {
		return RMNView{}, fmt.Errorf("failed to get config details for RMN: %w", err)
	}
	cs, err := r.GetPermaBlessedCommitStores(nil)
	if err != nil {
		return RMNView{}, err
	}
	return RMNView{
		ContractMetaData:         meta,
		ConfigDetails:            config,
		PermaBlessedCommitStores: cs,
	}, nil
}
