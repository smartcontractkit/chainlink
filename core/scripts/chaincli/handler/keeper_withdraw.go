package handler

import (
	"context"
	"log"
	"math/big"

	registry20 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	registry11 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
)

// Withdraw takes a keeper registry address, cancels all upkeeps and withdraws the funds
func (k *Keeper) Withdraw(ctx context.Context, hexAddr string) {
	registryAddr := common.HexToAddress(hexAddr)
	switch k.cfg.RegistryVersion {
	case keeper.RegistryVersion_1_1:
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
	case keeper.RegistryVersion_1_2:
		keeperRegistry12, err := registry12.NewKeeperRegistry(
			registryAddr,
			k.client,
		)
		if err != nil {
			log.Fatal("Registry failed: ", err)
		}

		activeUpkeepIds := k.getActiveUpkeepIds(ctx, keeperRegistry12, big.NewInt(0), big.NewInt(0))

		log.Println("Canceling upkeeps...")
		if err = k.cancelAndWithdrawActiveUpkeeps(ctx, activeUpkeepIds, keeperRegistry12); err != nil {
			log.Fatal("Failed to cancel upkeeps: ", err)
		}
	case keeper.RegistryVersion_2_0:
		keeperRegistry20, err := registry20.NewKeeperRegistry(
			registryAddr,
			k.client,
		)
		if err != nil {
			log.Fatal("Registry failed: ", err)
		}

		activeUpkeepIds := k.getActiveUpkeepIds(ctx, keeperRegistry20, big.NewInt(0), big.NewInt(0))

		log.Println("Canceling upkeeps...")
		if err = k.cancelAndWithdrawActiveUpkeeps(ctx, activeUpkeepIds, keeperRegistry20); err != nil {
			log.Fatal("Failed to cancel upkeeps: ", err)
		}
	default:
		panic("unexpected registry version")
	}
	log.Println("Upkeeps successfully canceled")
}
