package v1_6

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
)

type CCIPHome struct {
	types.ContractMetaData
}

func GenerateCCIPHomeView(c *ccip_home.CCIPHome) (CCIPHome, error) {
	if c == nil {
		return CCIPHome{}, fmt.Errorf("cannot generate view for nil CCIPHome")
	}
	meta, err := types.NewContractMetaData(c, c.Address())
	if err != nil {
		return CCIPHome{}, fmt.Errorf("failed to generate contract metadata for CCIPHome: %w", err)
	}
	return CCIPHome{
		ContractMetaData: meta,
	}, nil
}
