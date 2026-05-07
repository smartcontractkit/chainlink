package handler

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	registry20 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper2_0"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
)

// Withdraw takes a keeper registry address, cancels all upkeeps and withdraws the funds
func (k *Keeper) Withdraw(ctx context.Context, hexAddr string) {
	registryAddr := common.HexToAddress(hexAddr)
	switch k.cfg.RegistryVersion {
	case config.RegistryVersion2_0:
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
