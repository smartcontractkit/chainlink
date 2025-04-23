package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"golang.org/x/time/rate"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
)

func TestDecodeInstruction(t *testing.T) {
	data := common.Hex2Bytes("6cd886bff9ea21549b3c1f221aa3f0cc2000000000000000000000000000000078b7038a2cccb52c45334016f9c3cd95219a937b20000000000000000000000000000000000000000000000000000000000000000000000100000000b2e1ba26529431b48b4d74264979bedde967e6c78121b81dd708b35530c06c470000000000000000")

	inst, err := ccip_router.DecodeInstruction(nil, data)
	if err != nil {
		t.Fatalf("unable to decode instruction: %v", err)
	}

	if data, err := json.MarshalIndent(inst, "", " "); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	} else {
		fmt.Println(string(data))
	}
}

func TestGetCCIPSendEvent(t *testing.T) {
	txsig, err := solana.SignatureFromBase58("3KeCEazfpEY3FVr3f3jw9rwxfYFgdPAGkCwwUaBBP6PTvi2YoyjhFFNC1hLcKhybhLvjwmPBgbo4TXRV4Nbuivuc")
	if err != nil {
		t.Fatalf("failed to parse transaction signature: %v", err)
	}

	rpcClient := rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(rpc.DevNet.RPC, rate.Every(time.Second), 5))

	v := uint64(0)
	result, err := rpcClient.GetTransaction(context.Background(), txsig, &rpc.GetTransactionOpts{
		Commitment:                     rpc.CommitmentConfirmed,
		MaxSupportedTransactionVersion: &v,
	})
	if err != nil {
		t.Fatalf("failed to get transaction: %v", err)
	}

	// check CCIP event
	ccipMessageSentEvent := solccip.EventCCIPMessageSent{}
	err = solcommon.ParseEvent(result.Meta.LogMessages, "CCIPMessageSent", &ccipMessageSentEvent, true)
	if err != nil {
		t.Fatalf("failed to parse CCIPMessageSent event: %v", err)
	}

	fmt.Println("CCIPMessageSent event:")
	fmt.Printf("  MessageId: %s\n", common.Bytes2Hex(ccipMessageSentEvent.Message.Header.MessageId[:]))
	fmt.Printf("  SequenceNumber: %d\n", ccipMessageSentEvent.Message.Header.SequenceNumber)
}
