package handler

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	registry11 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
)

const (
	defaultMaxBlocksRange = 1000
	defaultLookBackRange  = 1000
)

var (
	checkUpkeepArguments1 abi.Arguments
	checkUpkeepArguments2 abi.Arguments
	registry11ABI         = keeper.Registry1_1ABI
	registry12ABI         = keeper.Registry1_2ABI
)

type result struct {
	block          uint64
	checkUpkeep    bool
	keeperIndex    uint64
	keeperAddress  common.Address
	reason         string
	performData    string
	maxLinkPayment *big.Int
	gasLimit       *big.Int
	adjustedGasWei *big.Int
	linkEth        *big.Int
}

func init() {
	checkUpkeepArguments1 = registry11ABI.Methods["checkUpkeep"].Outputs
	checkUpkeepArguments2 = registry12ABI.Methods["checkUpkeep"].Outputs
}

// UpkeepHistory prints the checkUpkeep status and keeper responsibility for a given upkeep in a set block range
func (k *Keeper) UpkeepHistory(ctx context.Context, upkeepId *big.Int, from, to, gasPrice uint64) {
	// There must not be a large different between boundaries
	if to-from > defaultMaxBlocksRange {
		log.Fatalf("blocks range difference must not more than %d", defaultMaxBlocksRange)
	}

	var keeperRegistry11 *registry11.KeeperRegistry
	var keeperRegistry12 *registry12.KeeperRegistry
	// var keeperRegistry20 *registry20.KeeperRegistry

	switch k.cfg.RegistryVersion {
	case keeper.RegistryVersion_1_1:
		_, keeperRegistry11 = k.getRegistry11(ctx)
	case keeper.RegistryVersion_1_2:
		_, keeperRegistry12 = k.getRegistry12(ctx)
	default:
		panic("unsupported registry version")
	}

	log.Println("Preparing a batch call request")
	var reqs []rpc.BatchElem
	var results []*string
	var keeperPerBlockIndex []uint64
	var keeperPerBlockAddress []common.Address
	for block := from; block <= to; block++ {
		callOpts := &bind.CallOpts{
			Context:     ctx,
			BlockNumber: big.NewInt(0).SetUint64(block),
		}

		var keepers []common.Address
		var bcpt uint64
		var payload []byte
		var keeperIndex uint64
		var lastKeeper common.Address

		switch k.cfg.RegistryVersion {
		case keeper.RegistryVersion_1_1:
			config, err2 := keeperRegistry11.GetConfig(callOpts)
			if err2 != nil {
				log.Fatal("failed to fetch registry config: ", err2)
			}

			bcpt = config.BlockCountPerTurn.Uint64()
			keepers, err2 = keeperRegistry11.GetKeeperList(callOpts)
			if err2 != nil {
				log.Fatal("failed to fetch keepers list: ", err2)
			}

			upkeep, err2 := keeperRegistry11.GetUpkeep(callOpts, upkeepId)
			if err2 != nil {
				log.Fatal("failed to fetch the upkeep: ", err2)
			}
			lastKeeper = upkeep.LastKeeper

		case keeper.RegistryVersion_1_2:
			state, err2 := keeperRegistry12.GetState(callOpts)
			if err2 != nil {
				log.Fatal("failed to fetch registry state: ", err2)
			}
			bcpt = state.Config.BlockCountPerTurn.Uint64()
			keepers = state.Keepers

			upkeep, err2 := keeperRegistry12.GetUpkeep(callOpts, upkeepId)
			if err2 != nil {
				log.Fatal("failed to fetch the upkeep: ", err2)
			}
			lastKeeper = upkeep.LastKeeper

		default:
			panic("unsupported registry version")
		}

		turnBinary, err2 := turnBlockHashBinary(ctx, block, bcpt, defaultLookBackRange, k.client)
		if err2 != nil {
			log.Fatal("failed to calculate turn block hash: ", err2)
		}

		// least significant 32 bits of upkeep id
		lhs := keeper.LeastSignificant32(upkeepId)

		// least significant 32 bits of the turn block hash
		turnBinaryPtr, ok := math.ParseBig256(string([]byte(turnBinary)[len(turnBinary)-32:]))
		if !ok {
			log.Fatal("failed to parse turn binary ", turnBinary)
		}
		rhs := keeper.LeastSignificant32(turnBinaryPtr)

		// bitwise XOR
		turn := lhs ^ rhs

		keepersCnt := uint64(len(keepers))
		keeperIndex = turn % keepersCnt
		if keepers[keeperIndex] == lastKeeper {
			keeperIndex = (keeperIndex + keepersCnt - 1) % keepersCnt
		}

		switch k.cfg.RegistryVersion {
		case keeper.RegistryVersion_1_1:
			payload, err2 = registry11ABI.Pack("checkUpkeep", upkeepId, keepers[keeperIndex])
			if err2 != nil {
				log.Fatal("failed to pack checkUpkeep: ", err2)
			}
		case keeper.RegistryVersion_1_2:
			payload, err2 = registry12ABI.Pack("checkUpkeep", upkeepId, keepers[keeperIndex])
			if err2 != nil {
				log.Fatal("failed to pack checkUpkeep: ", err2)
			}
		default:
			panic("unsupported registry version")
		}

		args := map[string]interface{}{
			"to":   k.cfg.RegistryAddress,
			"data": hexutil.Bytes(payload),
		}
		if gasPrice > 0 {
			args["gasPrice"] = hexutil.EncodeUint64(gasPrice)
		}

		var res string
		reqs = append(reqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				args,
				// The block at which we want to inspect the upkeep state
				hexutil.EncodeUint64(block),
			},
			Result: &res,
		})

		results = append(results, &res)
		keeperPerBlockIndex = append(keeperPerBlockIndex, keeperIndex)
		keeperPerBlockAddress = append(keeperPerBlockAddress, keepers[keeperIndex])
	}

	k.batchProcess(ctx, reqs, from, keeperPerBlockIndex, keeperPerBlockAddress, results)
}

func (k *Keeper) batchProcess(ctx context.Context, reqs []rpc.BatchElem, from uint64, keeperPerBlockIndex []uint64, keeperPerBlockAddress []common.Address, results []*string) {
	log.Println("Doing batch call to check upkeeps")
	if err := k.rpcClient.BatchCallContext(ctx, reqs); err != nil {
		log.Fatal("failed to batch call checkUpkeep: ", err)
	}

	log.Println("Parsing batch call response")
	var parsedResults []result
	isVersion12 := k.cfg.RegistryVersion == keeper.RegistryVersion_1_2
	for i, req := range reqs {
		if req.Error != nil {
			parsedResults = append(parsedResults, result{
				block:         uint64(i) + from,
				checkUpkeep:   false,
				keeperIndex:   keeperPerBlockIndex[i],
				keeperAddress: keeperPerBlockAddress[i],
				reason:        strings.TrimPrefix(req.Error.Error(), "execution reverted: "),
			})
			continue
		}

		var returnValues []interface{}
		var err error
		if isVersion12 {
			returnValues, err = checkUpkeepArguments2.UnpackValues(hexutil.MustDecode(*results[i]))
		} else {
			returnValues, err = checkUpkeepArguments1.UnpackValues(hexutil.MustDecode(*results[i]))
		}
		if err != nil {
			log.Fatal("unpack checkUpkeep return: ", err, *results[i])
		}

		parsedResults = append(parsedResults, result{
			block:          uint64(i) + from,
			checkUpkeep:    true,
			keeperIndex:    keeperPerBlockIndex[i],
			keeperAddress:  keeperPerBlockAddress[i],
			performData:    "0x" + hex.EncodeToString(*abi.ConvertType(returnValues[0], new([]byte)).(*[]byte)),
			maxLinkPayment: *abi.ConvertType(returnValues[1], new(*big.Int)).(**big.Int),
			gasLimit:       *abi.ConvertType(returnValues[2], new(*big.Int)).(**big.Int),
			adjustedGasWei: *abi.ConvertType(returnValues[3], new(*big.Int)).(**big.Int),
			linkEth:        *abi.ConvertType(returnValues[4], new(*big.Int)).(**big.Int),
		})
	}

	printResultsToConsole(parsedResults)
}

// printResultsToConsole writes parsed results to the console
func printResultsToConsole(parsedResults []result) {
	writer := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	defer writer.Flush()

	fmt.Fprintf(writer, "\n %s\t\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t", "Block", "checkUpkeep", "Keeper Index", "Keeper Address", "Max LINK Payment", "Gas Limit", "Adjusted Gas", "LINK ETH", "Perform Data", "Reason")
	fmt.Fprintf(writer, "\n %s\t\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t", "----", "----", "----", "----", "----", "----", "----", "----", "----", "----")
	for _, res := range parsedResults {
		fmt.Fprintf(writer, "\n %d\t\t%t\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t",
			res.block,
			res.checkUpkeep,
			res.keeperIndex,
			res.keeperAddress,
			res.maxLinkPayment,
			res.gasLimit,
			res.adjustedGasWei,
			res.linkEth,
			res.performData,
			res.reason,
		)
	}
	fmt.Fprintf(writer, "\n %s\t\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n", "----", "----", "----", "----", "----", "----", "----", "----", "----", "----")
}

func turnBlockHashBinary(ctx context.Context, blockNum, bcpt, lookback uint64, ethClient *ethclient.Client) (string, error) {
	turnBlock := blockNum - (blockNum % bcpt) - lookback
	block, err := ethClient.BlockByNumber(ctx, big.NewInt(int64(turnBlock)))
	if err != nil {
		return "", err
	}
	hashAtHeight := block.Hash()
	binaryString := fmt.Sprintf("%b", hashAtHeight.Big())
	return binaryString, nil
}
