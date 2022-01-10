package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ec, err := ethclient.Dial("http://127.0.0.1:8545")
	panicErr(err)
	txHash := "0x537c8aadb3ae592b0c9683a49d14e61dd74c3cd835c69d055c4c70298e491a66"
	tx, _, err := ec.TransactionByHash(context.Background(), gethCommon.HexToHash(txHash))
	panicErr(err)
	re, err := ec.TransactionReceipt(context.Background(), gethCommon.HexToHash(txHash))
	panicErr(err)
	fmt.Println(re.Status, re.GasUsed, re.CumulativeGasUsed)
	requester := gethCommon.HexToAddress("9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f")
	call := ethereum.CallMsg{
		From:     requester,
		To:       tx.To(),
		Data:     tx.Data(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
	r, err := ec.CallContract(context.Background(), call, re.BlockNumber)
	fmt.Println("call contract", "r", r, "err", err)
	reason, err := evmclient.ExtractRevertReasonFromRPCError(err)
	fmt.Println("extracting revert reason", "reason", reason, "err", err)
}
