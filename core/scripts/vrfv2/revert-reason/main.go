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
	ec, err := ethclient.Dial("https://eth.getblock.io/goerli/608a6708-d787-4ba8-bcbd-493001ea7fd9/")
	panicErr(err)
	txHash := "0x329f997acdffc0d1ce712526030bc34f5e3ca00c816a802f873b6d87118c02de" // non-custom on goerli
	//txHash := "0x6ec8a69657600786f0b31726f36287e80196029e60f8365528d4d540a6f70763" // custom error on mainnet
	tx, _, err := ec.TransactionByHash(context.Background(), gethCommon.HexToHash(txHash))
	panicErr(err)
	re, err := ec.TransactionReceipt(context.Background(), gethCommon.HexToHash(txHash))
	panicErr(err)
	fmt.Println(re.Status, re.GasUsed, re.CumulativeGasUsed)
	requester := gethCommon.HexToAddress("0xeFF41C8725be95e66F6B10489B6bF34b08055853")
	call := ethereum.CallMsg{
		From:     requester,
		To:       tx.To(),
		Data:     tx.Data(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
	r, err := ec.CallContract(context.Background(), call, re.BlockNumber)
	fmt.Println("call contract", "r", r, "err", err)
	rpcError, err := evmclient.ExtractRPCError(err)
	fmt.Println("extracting rpc error", rpcError.String(), err)
}
