package handler

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	iregistry21 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	registry20 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
)

// Keeper is the keepers commands handler
type Keeper struct {
	*baseHandler
}

// NewKeeper creates new instance of Keeper
func NewKeeper(cfg *config.Config) *Keeper {
	return &Keeper{
		baseHandler: NewBaseHandler(cfg),
	}
}

func (k *Keeper) VerifyContract(params ...string) {
	if err := k.changeToContractsDirectory(); err != nil {
		log.Fatalf("failed to change to directory where the hardhat.config.ts file is located: %v", err)
	}

	commandArgs := append([]string{}, params...)
	command := fmt.Sprintf(
		"NODE_HTTP_URL='%s' EXPLORER_API_KEY='%s' NETWORK_NAME='%s' pnpm hardhat verify --network env %s",
		k.cfg.NodeHttpURL,
		k.cfg.ExplorerAPIKey,
		k.cfg.NetworkName,
		strings.Join(commandArgs, " "),
	)

	fmt.Println("Running command to verify contract: ", command)
	if err := k.runCommand(command); err != nil {
		log.Println("Contract verification on Explorer failed: ", err)
	}
}

// UpdateRegistry attaches to an existing registry and possibly updates registry config
func (k *Keeper) UpdateRegistry(ctx context.Context) {
	var registryAddr common.Address
	switch k.cfg.RegistryVersion {
	case config.RegistryVersion2_0:
		registryAddr, _ = k.getRegistry20(ctx)
	case config.RegistryVersion2_1:
		registryAddr, _ = k.getRegistry21(ctx)
	default:
		panic("unexpected registry address")
	}
	log.Println("KeeperRegistry at:", registryAddr)
}

func (k *Keeper) getRegistry20(ctx context.Context) (common.Address, *registry20.KeeperRegistry) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry20, err := registry20.NewKeeperRegistry(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	if k.cfg.RegistryConfigUpdate {
		panic("KeeperRegistry2.0 could not be updated")
	}
	log.Println("KeeperRegistry2.0 config not updated: KEEPER_CONFIG_UPDATE=false")
	return registryAddr, keeperRegistry20
}

func (k *Keeper) getRegistry21(ctx context.Context) (common.Address, *iregistry21.IKeeperRegistryMaster) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry21, err := iregistry21.NewIKeeperRegistryMaster(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	if k.cfg.RegistryConfigUpdate {
		panic("KeeperRegistry2.1 could not be updated")
	}
	log.Println("KeeperRegistry2.1 config not updated: KEEPER_CONFIG_UPDATE=false")
	return registryAddr, keeperRegistry21
}

type activeUpkeepGetter interface {
	Address() common.Address
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
}

func (k *Keeper) getActiveUpkeepIds(ctx context.Context, registry activeUpkeepGetter, from, to *big.Int) []*big.Int {
	activeUpkeepIds, _ := registry.GetActiveUpkeepIDs(&bind.CallOpts{
		Pending: false,
		From:    k.fromAddr,
		Context: ctx,
	}, from, to)
	return activeUpkeepIds
}
