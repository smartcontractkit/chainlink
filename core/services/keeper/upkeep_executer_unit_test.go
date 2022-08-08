package keeper

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
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

func TestBuildJobSpec(t *testing.T) {
	jb := job.Job{ID: 10}
	from := ethkey.EIP55Address(testutils.NewAddress().Hex())
	contract := ethkey.EIP55Address(testutils.NewAddress().Hex())
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
	gasPrice := big.NewInt(24)
	gasTipCap := big.NewInt(48)
	gasFeeCap := big.NewInt(72)
	chainID := "250"

	m := &registryGasCheckMock{}
	m.Mock.Test(t)

	m.On("KeeperRegistryPerformGasOverhead").Return(uint32(9)).Times(2)
	m.On("KeeperRegistryCheckGasOverhead").Return(uint32(6)).Times(1)

	spec := buildJobSpec(jb, upkeep, m, m, gasPrice, gasTipCap, gasFeeCap, chainID)

	expected := map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"jobID":                 int32(10),
			"fromAddress":           from.String(),
			"contractAddress":       contract.String(),
			"upkeepID":              "4",
			"prettyID":              fmt.Sprintf("UPx%064d", 4),
			"performUpkeepGasLimit": uint32(21),
			"checkUpkeepGasLimit":   uint32(38),
			"gasPrice":              gasPrice,
			"gasTipCap":             gasTipCap,
			"gasFeeCap":             gasFeeCap,
			"evmChainID":            "250",
		},
	}

	require.Equal(t, expected, spec)
}
