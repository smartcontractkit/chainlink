package handler

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"

	keeper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
)

// Withdraw takes a keeper registry address cancels all upkeeps and withdraws the funds
func (k *Keeper) Withdraw(ctx context.Context, hexAddr string) {
	registryAddr := common.HexToAddress(hexAddr)
	registryInstance, err := keeper.NewKeeperRegistry(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	log.Println("Canceling upkeeps...")
	if err := k.cancelAndWithdrawUpkeeps(ctx, registryInstance); err != nil {
		log.Fatal("Failed to cancel upkeeps: ", err)
	}
	log.Println("Upkeeps successfully canceled")
}
