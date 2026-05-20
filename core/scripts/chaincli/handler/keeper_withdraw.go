package handler

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	registry20 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper2_0"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
)

type canceller interface {
	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)
	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)
}

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

func (k *Keeper) cancelAndWithdrawActiveUpkeeps(ctx context.Context, activeUpkeepIds []*big.Int, canceller canceller) error {
	for i := range activeUpkeepIds {
		upkeepId := activeUpkeepIds[i]
		tx, err := canceller.CancelUpkeep(k.buildTxOpts(ctx), upkeepId)
		if err != nil {
			return fmt.Errorf("failed to cancel upkeep %s: %w", upkeepId.String(), err)
		}

		if err = k.waitTx(ctx, tx); err != nil {
			log.Fatalf("failed to cancel upkeep for upkeepId: %s, error is: %s", upkeepId.String(), err.Error())
		}

		tx, err = canceller.WithdrawFunds(k.buildTxOpts(ctx), upkeepId, k.fromAddr)
		if err != nil {
			return fmt.Errorf("failed to withdraw upkeep %s: %w", upkeepId.String(), err)
		}

		if err = k.waitTx(ctx, tx); err != nil {
			log.Fatalf("failed to withdraw upkeep for upkeepId: %s, error is: %s", upkeepId.String(), err.Error())
		}

		log.Printf("Upkeep %s successfully canceled and refunded: ", upkeepId.String())
	}

	tx, err := canceller.RecoverFunds(k.buildTxOpts(ctx))
	if err != nil {
		return fmt.Errorf("failed to recover funds: %w", err)
	}

	if err = k.waitTx(ctx, tx); err != nil {
		log.Fatalf("failed to recover funds, error is: %s", err.Error())
	}

	return nil
}
