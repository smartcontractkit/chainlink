package v1_2_0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_2_0"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

func TestGetRouter(t *testing.T) {
	routerAddr := utils.RandomAddress()

	mockOffRamp := mock_contracts.NewEVM2EVMOffRampInterface(t)
	mockOffRamp.On("GetDynamicConfig", mock.Anything).Return(evm_2_evm_offramp_1_2_0.EVM2EVMOffRampDynamicConfig{
		Router: routerAddr,
	}, nil)

	offRamp := OffRamp{
		offRampV120: mockOffRamp,
	}

	ctx := testutils.Context(t)
	gotRouterAddr, err := offRamp.GetRouter(ctx)
	require.NoError(t, err)

	gotRouterEvmAddr, err := ccipcalc.GenericAddrToEvm(gotRouterAddr)
	require.NoError(t, err)
	assert.Equal(t, routerAddr, gotRouterEvmAddr)
}
