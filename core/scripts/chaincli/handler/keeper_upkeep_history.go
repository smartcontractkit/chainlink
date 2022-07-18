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

	registry11 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
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

	isVersion12 := k.cfg.RegistryVersion == keeper.RegistryVersion_1_2
	var registryAddr common.Address
	var keeperRegistry11 *registry11.KeeperRegistry
	var keeperRegistry12 *registry12.KeeperRegistry

	if isVersion12 {
		registryAddr, keeperRegistry12 = k.getRegistry2(ctx)
	} else {
		registryAddr, keeperRegistry11 = k.getRegistry1(ctx)
	}

	// Get positioning constant of the current registry
	var positioningConstant int32
	var err error
	if !isVersion12 {
		positioningConstant, err = keeper.CalcPositioningConstant(utils.NewBig(upkeepId), ethkey.EIP55AddressFromAddress(registryAddr))
		if err != nil {
			log.Fatal("failed to get positioning constant: ", err)
		}
		log.Println("Calculated Positioning Constant for the registry: ", positioningConstant)
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

		if isVersion12 {
			state, err := keeperRegistry12.GetState(callOpts)
			if err != nil {
				log.Fatal("failed to fetch registry state: ", err)
			}
			bcpt = state.Config.BlockCountPerTurn.Uint64()
			keepers = state.Keepers
			keepersCnt := uint64(len(state.Keepers))

			upkeep, err := keeperRegistry12.GetUpkeep(callOpts, upkeepId)
			if err != nil {
				log.Fatal("failed to fetch the upkeep: ", err)
			}

			turnBinary, err := turnBlockHashBinary(block, bcpt, defaultLookBackRange, k.client)
			if err != nil {
				log.Fatal("failed to calculate turn block hash: ", err)
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

			keeperIndex = turn % keepersCnt
			if keepers[keeperIndex] == upkeep.LastKeeper {
				keeperIndex = (keeperIndex + keepersCnt - 1) % keepersCnt
			}
			payload, err = registry12ABI.Pack("checkUpkeep", upkeepId, keepers[keeperIndex])
		} else {
			config, err := keeperRegistry11.GetConfig(callOpts)
			if err != nil {
				log.Fatal("failed to fetch registry config: ", err)
			}

			bcpt = config.BlockCountPerTurn.Uint64()
			keepers, err = keeperRegistry11.GetKeeperList(callOpts)
			if err != nil {
				log.Fatal("failed to fetch keepers list: ", err)
			}

			keeperIndex = (uint64(positioningConstant) + ((block - (block % bcpt)) / bcpt)) % uint64(len(keepers))
			payload, err = registry11ABI.Pack("checkUpkeep", upkeepId, keepers[keeperIndex])
		}
		if err != nil {
			log.Fatal("failed to pack checkUpkeep: ", err)
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

func turnBlockHashBinary(blockNum, bcpt, lookback uint64, ethClient *ethclient.Client) (string, error) {
	turnBlock := blockNum - (blockNum % bcpt) - lookback
	block, err := ethClient.BlockByNumber(context.Background(), big.NewInt(int64(turnBlock)))
	if err != nil {
		return "", err
	}
	hashAtHeight := block.Hash()
	binaryString := fmt.Sprintf("%b", hashAtHeight.Big())
	return binaryString, nil
}
