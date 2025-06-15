//go:build wasip1

package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	httpaction "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http"
	evmcap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	consensuscap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus"
	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	evmpb "github.com/smartcontractkit/chainlink-common/pkg/chains/evm"
	pb "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"

	consensusbpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
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
	// Relevant Addresses
	toAddress := common.HexToAddress("0x779877A7B0D9E8603169DdbD7836e478b4624789")
	walletAddress := common.HexToAddress("0x437bb34CbdB6c0Eaf859FfDC2DfC424d710e4C5B")

	// balanceOf(address) selector
	methodID := []byte{0x70, 0xa0, 0x82, 0x31}

	// Pad the address to 32 bytes
	paddedAddress := common.LeftPadBytes(walletAddress.Bytes(), 32)

	// Combine method selector and padded address
	data := append(methodID, paddedAddress...)

	// Call EVM Call Contract Capability
	evm := evmcap.Client{}
	evmOut := evm.CallContract(runtime, &evmpb.CallContractRequest{
		Call: &evmpb.CallMsg{
			From: walletAddress.Bytes(),
			To:   toAddress.Bytes(),
			Data: data,
		},
		BlockNumber: &pb.BigInt{
			AbsVal: []byte{0x0B, 0xB8}, // 3000 in hex
		},
	})

	// Await EVM Call Contract Capability
	reply, err := evmOut.Await()
	if err != nil {
		return "", err
	}

	// Get Balance from EVM Call Contract Capability
	out := reply.GetData()
	balance := new(big.Int).SetBytes(out)

	// Call Consensus Capability on Balance value
	consensus := consensuscap.Consensus{}
	consensusOut := consensus.Simple(runtime, &consensusbpb.SimpleConsensusInputs{
		Default:     pb.NewBigIntValue(1, []byte{0x01}),
		Descriptors: sdk.ConsensusIdenticalAggregation[bool]().Descriptor(),
		Observation: &consensusbpb.SimpleConsensusInputs_Value{
			Value: pb.NewBigIntValue(1, balance.Bytes()),
		},
	})

	// Await Consensus Capability
	val, err := consensusOut.Await()
	if err != nil {
		return "", err
	}

	outInt := val.GetBigintValue().GetAbsVal()

	// Call HTTP Action Capability
	_, err = runtime.RunInNodeMode(func(runtime sdk.NodeRuntime) *consensusbpb.SimpleConsensusInputs {
		httpAction := httpaction.Client{}
		httpActionOut, err := httpAction.SendRequest(runtime, &httpaction.Request{
			Url: "https://api.github.com/users/octocat",
		}).Await()
		if err != nil {
			return &consensusbpb.SimpleConsensusInputs{}
		}
		return &consensusbpb.SimpleConsensusInputs{
			Default:     pb.NewBytesValue(httpActionOut.GetBody()),
			Descriptors: sdk.ConsensusIdenticalAggregation[bool]().Descriptor(),
			Observation: &consensusbpb.SimpleConsensusInputs_Value{
				Value: pb.NewBytesValue(httpActionOut.GetBody()),
			},
		}
	}).Await()
	if err != nil {
		return "", err
	}

	// evm.WriteReport(runtime, &evmpb.WriteReportRequest{
	// 	Report: ,
	// })

	return new(big.Int).SetBytes(outInt).String(), nil
}

func main() {
	RunSimpleCronWorkflow(wasm.NewDonRunner())
}
