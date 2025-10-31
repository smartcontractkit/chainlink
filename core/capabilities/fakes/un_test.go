package fakes

// This test demonstrates how to unmarshal WriteReportReply from Base64-encoded bytes
// and print all fields to stdout. It covers three scenarios:
// 1. Complete message with all fields (including error message)
// 2. Success case without error message
// 3. Reverted transaction with error details

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
)

func TestUnmarshalWriteReportReply_FromBase64(t *testing.T) {
	base64Encoded := "CAIQARogoOhieV2zzajkMj9YQ/w28wnf/apZCYd7ue/X4eBWNGMiCQoFFcPVxEEQASopdW5rbm93biBpc3N1ZSBleGVjdXRpb24gcmVjZWl2ZXIgY29udHJhY3Q="
	//base64Encoded := "CAIQARogoOhieV2zzajkMj9YQ/w28wnf/apZCYd7ue/X4eBWNGMiCQoFFcPVxEEQAQ=="

	// Now test the actual unmarshaling from base64
	t.Run("Unmarshal from base64 and print fields", func(t *testing.T) {
		// Decode from base64
		decodedBytes, err := base64.StdEncoding.DecodeString(base64Encoded)
		require.NoError(t, err, "Failed to decode base64 string")

		// Unmarshal the proto message
		reply := &evmcappb.WriteReportReply{}
		err = proto.Unmarshal(decodedBytes, reply)
		require.NoError(t, err, "Failed to unmarshal WriteReportReply")

		// Print all fields to stdout
		fmt.Println("=== WriteReportReply Fields ===")
		fmt.Printf("TxStatus: %v (%s)\n", reply.TxStatus, reply.TxStatus.String())

		if reply.ReceiverContractExecutionStatus != nil {
			fmt.Printf("ReceiverContractExecutionStatus: %v (%s)\n",
				*reply.ReceiverContractExecutionStatus,
				reply.ReceiverContractExecutionStatus.String())
		} else {
			fmt.Println("ReceiverContractExecutionStatus: nil")
		}

		if reply.TxHash != nil {
			fmt.Printf("TxHash: %x (length: %d bytes)\n", reply.TxHash, len(reply.TxHash))
		} else {
			fmt.Println("TxHash: nil")
		}

		if reply.TransactionFee != nil {
			// Convert BigInt proto to big.Int
			fee := new(big.Int).SetBytes(reply.TransactionFee.AbsVal)
			if reply.TransactionFee.Sign < 0 {
				fee.Neg(fee)
			}
			fmt.Printf("TransactionFee: %s wei\n", fee.String())
			fmt.Printf("TransactionFee (ETH): %s ETH\n", weiToEth(fee))
		} else {
			fmt.Println("TransactionFee: nil")
		}

		if reply.ErrorMessage != nil {
			fmt.Printf("ErrorMessage: %s\n", *reply.ErrorMessage)
		} else {
			fmt.Println("ErrorMessage: nil")
		}
		fmt.Println("===============================")
	})
}

// weiToEth converts wei (big.Int) to ETH string representation
func weiToEth(wei *big.Int) string {
	if wei == nil {
		return "0"
	}

	// 1 ETH = 10^18 wei
	ethDivisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	// Calculate ETH value with 6 decimal places
	eth := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetInt(ethDivisor)
	eth.Quo(eth, divisor)

	return eth.Text('f', 6)
}
