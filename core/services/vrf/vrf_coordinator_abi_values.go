package vrf

import (
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"chainlink/core/services/vrf/generated/solidity_vrf_coordinator_interface"
)

// CoordinatorABI returns the ABI for the VRFCoordinator contract
func CoordinatorABI() abi.ABI {
	return coordinatorABIValues().coordinatorABI
}

// FulfillSelector returns the signature of the fulfillRandomnessRequest method
// on the VRFCoordinator contract
func FulfillSelector() string {
	return coordinatorABIValues().fulfillSelector
}

// RandomnessRequestLogTopic returns the signature of the RandomnessRequest log
// emitted by the VRFCoordinator contract
func RandomnessRequestLogTopic() common.Hash {
	return coordinatorABIValues().randomnessRequestLogTopic
}

// randomnessRequestRawDataArgs returns a list of the arguments to the
// RandomnessRequest log emitted by the VRFCoordinator contract
func randomnessRequestRawDataArgs() abi.Arguments {
	return coordinatorABIValues().randomnessRequestRawDataArgs
}

var fulfillMethodName = "fulfillRandomnessRequest"

// abiValues is a singleton carrying information parsed once from the
// VRFCoordinator abi string
type abiValues struct {
	// CoordinatorABI is the ABI of the VRFCoordinator
	coordinatorABI  abi.ABI
	fulfillSelector string
	// RandomnessRequestLogTopic is the signature of the RandomnessRequest log
	randomnessRequestLogTopic    common.Hash
	randomnessRequestRawDataArgs abi.Arguments
}

var dontUseThisUseGetterFunctionsAbove abiValues
var parseABIOnce sync.Once

func coordinatorABIValues() *abiValues {
	parseABIOnce.Do(readCoordinatorABI)
	return &dontUseThisUseGetterFunctionsAbove
}

func readCoordinatorABI() {
	v := &dontUseThisUseGetterFunctionsAbove
	var err error
	v.coordinatorABI, err = abi.JSON(strings.NewReader(
		solidity_vrf_coordinator_interface.VRFCoordinatorABI))
	if err != nil {
		panic(err)
	}
	var seen bool
	for methodName, method := range v.coordinatorABI.Methods {
		if methodName == fulfillMethodName {
			v.fulfillSelector = hexutil.Encode(method.ID())
			seen = true
		}
	}
	if !seen {
		panic("failed to find fulfill method")
	}
	randomnessRequestABI := v.coordinatorABI.Events["RandomnessRequest"]
	v.randomnessRequestLogTopic = randomnessRequestABI.ID()
	for _, arg := range randomnessRequestABI.Inputs {
		if !arg.Indexed {
			v.randomnessRequestRawDataArgs = append(v.randomnessRequestRawDataArgs, arg)
		}
	}
}
