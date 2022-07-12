package handler

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	registry11 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
)

// Withdraw takes a keeper registry address, cancels all upkeeps and withdraws the funds
func (k *Keeper) Withdraw(ctx context.Context, hexAddr string) {
	registryAddr := common.HexToAddress(hexAddr)
	isVersion12 := k.cfg.RegistryVersion == keeper.RegistryVersion_1_2
	if isVersion12 {
		keeperRegistry12, err := registry12.NewKeeperRegistry(
			registryAddr,
			k.client,
		)
		if err != nil {
			log.Fatal("Registry failed: ", err)
		}
		activeUpkeepIds := k.getActiveUpkeepIds(ctx, keeperRegistry12)

		log.Println("Canceling upkeeps...")
		if err = k.cancelAndWithdrawActiveUpkeeps(ctx, activeUpkeepIds, keeperRegistry12); err != nil {
			log.Fatal("Failed to cancel upkeeps: ", err)
		}
	} else {
		keeperRegistry11, err := registry11.NewKeeperRegistry(
			registryAddr,
			k.client,
		)
		if err != nil {
			log.Fatal("Registry failed: ", err)
		}

		upkeepCount, err := keeperRegistry11.GetUpkeepCount(&bind.CallOpts{Context: ctx})
		if err != nil {
			log.Fatal("failed to get upkeeps count: ", err)
		}

		log.Println("Canceling upkeeps...")
		if err = k.cancelAndWithdrawUpkeeps(ctx, upkeepCount, keeperRegistry11); err != nil {
			log.Fatal("Failed to cancel upkeeps: ", err)
		}
	}
	log.Println("Upkeeps successfully canceled")
}
