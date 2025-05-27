//go:build wasip1

package main

import (
	evmcap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/capability"
	evmpb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/chain-service"
	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	pb "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
)

func RunSimpleCronWorkflow(runner sdk.DonRunner) {
	cron := &croncap.Cron{}
	cfg := &croncap.Config{
		Schedule: "*/3 * * * * *", // every three seconds
	}

	runner.Run(&sdk.WorkflowArgs[sdk.DonRuntime]{
		Handlers: []sdk.Handler[sdk.DonRuntime]{
			sdk.NewDonHandler(
				cron.Trigger(cfg),
				onTrigger,
			),
		},
	})
}

func onTrigger(runtime sdk.DonRuntime, outputs *croncap.Payload) (string, error) {
	// Add HTTP Action
	evm := evmcap.EVM{}
	evmOut := evm.CallContract(runtime, &evmpb.CallContractRequest{
		Call: &evmpb.CallMsg{
			From: &evmpb.Address{
				Address: []byte{1},
			},
			To: &evmpb.Address{
				Address: []byte{1},
			},
			Data: &evmpb.ABIPayload{
				Abi: []byte{1},
			},
		},
		BlockNumber: &pb.BigInt{
			AbsVal: []byte{1},
		},
	})
	reply, err := evmOut.Await()
	if err != nil {
		return "", err
	}
	reply.GetData()

	// consensus := consensuscap.Consensus{}
	// consensusOut := consensus.Simple(runtime, &consensusbpb.SimpleConsensusInputs{
	// 	Default:     pb.NewBoolValue(true),
	// 	Descriptors: sdk.ConsensusIdenticalAggregation[bool]().Descriptor(),
	// 	Observation: &consensusbpb.SimpleConsensusInputs_Value{
	// 		Value: pb.NewBoolValue(true),
	// 	},
	// })

	// val, err := consensusOut.Await()
	// if err != nil {
	// 	return "", err
	// }
	// val.GetBoolValue()
	return "ping", nil
}

func main() {
	RunSimpleCronWorkflow(wasm.NewDonRunner())
}
