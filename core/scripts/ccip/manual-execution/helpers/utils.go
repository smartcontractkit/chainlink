package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func VerifyAddress(addr string) error {
	if addr == "" {
		return errors.New("address is blank")
	}
	if !common.IsHexAddress(addr) {
		return fmt.Errorf("address %s is invalid", addr)
	}
	return nil
}

func WaitForSuccessfulTxReceipt(client ethereum.TransactionReader, hash common.Hash) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	for {
		select {
		case <-ticker.C:
			log.Println("[MINING] waiting for tx to be mined...")
			receipt, _ := client.TransactionReceipt(context.Background(), hash)
			if receipt != nil {
				if receipt.Status == types.ReceiptStatusFailed {
					return fmt.Errorf("[MINING] ERROR tx reverted %s", hash.Hex())
				}
				if receipt.Status == types.ReceiptStatusSuccessful {
					log.Printf("[MINING] tx mined %s successful\n", hash.Hex())
					return nil
				}
			}
		case <-ctx.Done():
			return errors.New("tx not confirmed within time")
		}
	}
}
