package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ec, err := ethclient.Dial("TODO")
	panicErr(err)
	txHash := "0xedeeecf6bd763ecc82b5dff31e073af9cc4cf8a4b47708df526ba61cf0201d25" // non-custom on goerli
	// txHash := "0x6ec8a69657600786f0b31726f36287e80196029e60f8365528d4d540a6f70763" // custom error on mainnet
	tx, _, err := ec.TransactionByHash(context.Background(), gethCommon.HexToHash(txHash))
	panicErr(err)
	re, err := ec.TransactionReceipt(context.Background(), gethCommon.HexToHash(txHash))
	panicErr(err)
	fmt.Println(re.Status, re.GasUsed, re.CumulativeGasUsed)
	requester := gethCommon.HexToAddress("0xffe4a8b862971611dce48f3ba295d4ebfeb5b2fe")
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
