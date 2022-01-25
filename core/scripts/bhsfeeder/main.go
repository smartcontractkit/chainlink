package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/scripts/common"
)

// config holds environment variables necessary for this script to run successfully.
type config struct {
	ethURL  string
	chainID int64
	account *bind.TransactOpts
}

func main() {
	config, err := getConfig()
	common.PanicErr(err)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "not enough arguments to bhsfeeder command")
		os.Exit(1)
	}

	ethClient, err := ethclient.Dial(config.ethURL)
	common.PanicErr(err)

	chainID, err := ethClient.ChainID(context.Background())
	common.PanicErr(err)

	switch os.Args[1] {
	case "bhs-deploy":
		bhsAddress, tx, _, err := blockhash_store.DeployBlockhashStore(config.account, ethClient)
		common.PanicErr(err)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		fmt.Println("BlockhashStore", bhsAddress.String(), "hash", tx.Hash())
		fmt.Println("Waiting for tx to get mined...")
		_, err = bind.WaitDeployed(ctx, ethClient, tx)
		common.PanicErr(err)
		fmt.Println("BHS Deploy Tx mined")
	case "forward":
		cmd := flag.NewFlagSet("forward", flag.ExitOnError)
		bhsAddress := cmd.String("bhs-address", "", "The address of the deployed blockhash store contract to feed")
		gasMultiplier := cmd.Float64("gas-multiplier", 1.15, "Gas multiplier to multiply with the gas price to get the final sending gas price")

		common.ParseArgs(cmd, os.Args[2:], "bhs-address", "gas-multiplier")

		bhs, err := blockhash_store.NewBlockhashStore(gethcommon.HexToAddress(*bhsAddress), ethClient)
		common.PanicErr(err)

		bhsAbi, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
		common.PanicErr(err)

		feeder := &forwardFeeder{
			account:   config.account,
			ethClient: ethClient,
			chainID:   chainID,
			bhs:       bhs,
			bhsABI:    bhsAbi,
			gasEstimator: &gasEstimator{
				client:        ethClient,
				gasMultiplier: *gasMultiplier,
			},
		}
		err = feeder.feed()
		common.PanicErr(err)
	case "backward":
		cmd := flag.NewFlagSet("backward", flag.ExitOnError)
		bhsAddress := cmd.String("bhs-address", "", "The address of the deployed blockhash store contract to feed")
		gasMultiplier := cmd.Float64("gas-multiplier", 1.15, "Gas multiplier to multiply with the gas price to get the final sending gas price")

		common.ParseArgs(cmd, os.Args[2:], "bhs-address", "gas-multiplier")

		bhs, err := blockhash_store.NewBlockhashStore(gethcommon.HexToAddress(*bhsAddress), ethClient)
		common.PanicErr(err)

		bhsAbi, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
		common.PanicErr(err)

		feeder := &backwardFeeder{
			account:   config.account,
			ethClient: ethClient,
			chainID:   chainID,
			bhs:       bhs,
			bhsABI:    bhsAbi,
			gasEstimator: &gasEstimator{
				client:        ethClient,
				gasMultiplier: *gasMultiplier,
			},
		}
		err = feeder.feed()
		common.PanicErr(err)
	case "store-earliest":
		cmd := flag.NewFlagSet("store-earliest", flag.ExitOnError)
		bhsAddress := cmd.String("bhs-address", "", "The address of the deployed blockhash store contract to feed")
		gasMultiplier := cmd.Float64("gas-multiplier", 1.15, "Gas multiplier to multiply with the gas price to get the final sending gas price")

		common.ParseArgs(cmd, os.Args[2:], "bhs-address", "gas-multiplier")

		bhs, err := blockhash_store.NewBlockhashStore(gethcommon.HexToAddress(*bhsAddress), ethClient)
		common.PanicErr(err)

		bhsAbi, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
		common.PanicErr(err)

		feeder := &backwardFeeder{
			account:   config.account,
			ethClient: ethClient,
			chainID:   chainID,
			bhs:       bhs,
			bhsABI:    bhsAbi,
			gasEstimator: &gasEstimator{
				client:        ethClient,
				gasMultiplier: *gasMultiplier,
			},
		}

		_, err = feeder.storeEarliest()
		common.PanicErr(err)

	case "blockheader":
		cmd := flag.NewFlagSet("blockheader", flag.ExitOnError)
		blockNumber := cmd.Int64("blocknumber", -1, "The blocknumber to get the block header for")

		common.ParseArgs(cmd, os.Args[2:], "blocknumber")

		header, err := serializedBlockHeader(ethClient, big.NewInt(*blockNumber), chainID)
		common.PanicErr(err)

		fmt.Println("Serialized header:", hex.EncodeToString(header))
	default:
		fmt.Println("Please provide a subcommand, one of 'backward', 'forward', and 'missed'")
	}
}

func hashes(txs []*types.Transaction) []gethcommon.Hash {
	hs := []gethcommon.Hash{}
	for _, tx := range txs {
		hs = append(hs, tx.Hash())
	}
	return hs
}
