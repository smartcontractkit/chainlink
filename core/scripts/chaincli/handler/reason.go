package handler

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RevertReason attempts to fetch more info on failed TX
func (h *baseHandler) RevertReason(hash string) {
	txHash := common.HexToHash(hash)
	// Get transaction object
	tx, isPending, err := h.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal("Transaction not found")
	}
	if isPending {
		log.Fatal("Transaction is still pending")
	}
	// Get transaction receipt
	receipt, err := h.client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal("Failed to retrieve receipt: " + err.Error())
	}

	if receipt.Status == 1 {
		log.Println("Transaction was successful")
		return
	}

	// Get failure reason
	reason := getFailureReason(h.client, h.fromAddr, tx, receipt.BlockNumber)
	log.Println("Revert reason: " + reason)
}

func getFailureReason(client *ethclient.Client, from common.Address, tx *types.Transaction, blockNumber *big.Int) string {
	code, err := client.CallContract(context.Background(), createCallMsgFromTransaction(from, tx), blockNumber)
	if err != nil {
		log.Println("Cannot not get revert reason: " + err.Error())
		return "not found"
	}
	if len(code) == 0 {
		return "no error message or out of gas"
	}
	return string(code)
}

func createCallMsgFromTransaction(from common.Address, tx *types.Transaction) ethereum.CallMsg {
	return ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
}
