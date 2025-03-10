package handler

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	registry11 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
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
	default:
		panic("unexpected registry version")
	}
	log.Println("Upkeeps successfully canceled")
}
