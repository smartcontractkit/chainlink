package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type registryGasCheckMock struct {
	mock.Mock
}

func (_m *registryGasCheckMock) KeeperRegistryCheckGasOverhead() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

func (_m *registryGasCheckMock) KeeperRegistryPerformGasOverhead() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

func (_m *registryGasCheckMock) KeeperRegistryMaxPerformDataSize() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

func TestBuildJobSpec(t *testing.T) {
	from := ethkey.EIP55Address(testutils.NewAddress().Hex())
	contract := ethkey.EIP55Address(testutils.NewAddress().Hex())
	chainID := "250"
	jb := job.Job{
		ID: 10,
		KeeperSpec: &job.KeeperSpec{
			FromAddress:     from,
			ContractAddress: contract,
		}}

	upkeepID := utils.NewBigI(4)
	upkeep := UpkeepRegistration{
		Registry: Registry{
			FromAddress:     from,
			ContractAddress: contract,
			CheckGas:        11,
		},
		UpkeepID:   upkeepID,
		ExecuteGas: 12,
	}
	gasPrice := assets.NewWeiI(24)
	gasTipCap := assets.NewWeiI(48)
	gasFeeCap := assets.NewWeiI(72)

	m := &registryGasCheckMock{}
	m.Mock.Test(t)

	m.On("KeeperRegistryPerformGasOverhead").Return(uint32(9)).Times(1)
	m.On("KeeperRegistryMaxPerformDataSize").Return(uint32(1000)).Times(1)

	spec := buildJobSpec(jb, jb.KeeperSpec.FromAddress.Address(), upkeep, m, gasPrice, gasTipCap, gasFeeCap, chainID)

	expected := map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"jobID":                  int32(10),
			"fromAddress":            from.String(),
			"effectiveKeeperAddress": jb.KeeperSpec.FromAddress.String(),
			"contractAddress":        contract.String(),
			"upkeepID":               "4",
			"prettyID":               fmt.Sprintf("UPx%064d", 4),
			"pipelineSpec": &pipeline.Spec{
				ForwardingAllowed: false,
			},
			"performUpkeepGasLimit": uint32(5_000_000 + 9),
			"maxPerformDataSize":    uint32(1000),
			"gasPrice":              gasPrice.ToInt(),
			"gasTipCap":             gasTipCap.ToInt(),
			"gasFeeCap":             gasFeeCap.ToInt(),
			"evmChainID":            "250",
		},
	}

	require.Equal(t, expected, spec)
}
