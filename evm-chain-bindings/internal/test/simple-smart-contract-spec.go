package ccipspec

import (
	"github.com/smartcontractkit/smart-contract-spec/internal/spec"
	"math/big"
	"reflect"
)

type FeeDefaultConfig struct {
	minFeeUSDCents    uint32
	maxFeeUSDCents    uint32
	deciBps           uint16
	destGasOverhead   uint32
	destBytesOverhead uint32
	isEnabled         bool
}

type DynamicConfig struct {
	feeQuoter                               string
	permissionLessExecutionThresholdSeconds uint32
	maxTokenTransferGas                     uint32
	maxPoolReleaseOrMintGas                 uint32
	messageValidator                        string
}

type EVM2AnyMessage struct {
	Receiver      []byte
	Data          []byte
	tokenAmmounts []int
	Address       string
	extraArgs     []byte
}

func ccipSpec() spec.OnChainProductSpec {
	spec := spec.OnChainProductSpec{
		SmartContracts: []spec.SmartContract{
			{
				Name: "OnRamp",
				Functions: []spec.Function{
					{
						Name:     "getFee",
						ReadOnly: true,
						Inputs: []spec.Arg{
							{
								Name: "destChainSelector",
								Type: reflect.TypeOf((*int)(nil)).Elem(),
							},
							{
								Name: "message",
								Type: reflect.TypeOf((*EVM2AnyMessage)(nil)).Elem(),
							},
						},
						Outputs: []spec.Arg{
							{
								Name: "feeTokenAmount",
								Type: reflect.TypeOf((*big.Int)(nil)).Elem(),
							},
						},
						Description: "Calculates the fee to be paid to send a message through CCIP to another chain",
					},
					{
						Name:     "setDynamicConfig",
						ReadOnly: false,
						Inputs: []spec.Arg{
							{
								Name: "dynamicConfig",
								Type: reflect.TypeOf((*DynamicConfig)(nil)).Elem(),
							},
						},
						Outputs:     []spec.Arg{},
						Description: "updates the DynamicConfiguration of the OffRamp contract",
					},
				},
			},
		},
	}
	return spec
}

// generated code
// type OnRamp struct {
// 	chainReader ChainReader
// 	chainWriter ChainWriter
// 	address string
// }
//
// func NewOnRamp(chainReader ChainReader, chainWriter ChainWriter, address string) {
// 	return OnRamp{
// 		chainReader: chainReader,
// 		chainWriter: chainWriter,
// 		address: address,
// 	}
// }
//
// func (or OnRamp) getFee(destChainSelector int) {
// 	return or.chainReader.GetLatestValue("OnRamp", "getFee", ...)
// }
//
// func (or OnRamp) setDynamicConfig(dynamicConfig DynaDynamicConfig) {
// 	or.chainWror.chainWriter.SubmitTransaction(.....)
// }
