package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
)

type CCIPConfigView struct {
	types.ContractMetaData
	CapabilityRegistry common.Address `json:"capabilityRegistry,omitempty"`
}

func GenerateCCIPConfigView(ccipConfig *ccip_config.CCIPConfig) (CCIPConfigView, error) {
	view := CCIPConfigView{}
	var err error
	view.ContractMetaData, err = types.NewContractMetaData(ccipConfig, ccipConfig.Address())
	if err != nil {
		return CCIPConfigView{}, fmt.Errorf("metadata error for CCIPConfig: %w", err)
	}
	view.CapabilityRegistry, err = ccipConfig.GetCapabilityRegistry(nil)
	if err != nil {
		return CCIPConfigView{}, fmt.Errorf("failed to get CapabilityRegistry: %w", err)
	}
	return view, nil
}
