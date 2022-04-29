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
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	defaultMaxBlocksRange = 1000
)

var (
	checkUpkeepArguments abi.Arguments
)

func init() {
	checkUpkeepArguments = keeper.Registry1_1ABI.Methods["checkUpkeep"].Outputs
}

// UpkeepHistory prints the checkUpkeep status and keeper responsibility for a given upkeep in a set block range
func (k *Keeper) UpkeepHistory(ctx context.Context, upkeepId int64, from, to, gasPrice uint64) {
	// There must not be a large different between boundaries
	if to-from > defaultMaxBlocksRange {
		log.Fatalf("blocks range difference must not more than %d", defaultMaxBlocksRange)
	}

	registryAddr, registryClient := k.GetRegistry(ctx)

	// Get positioning constant of the current registry
	positioningConstant, err := keeper.CalcPositioningConstant(utils.NewBigI(upkeepId), ethkey.EIP55AddressFromAddress(registryAddr))
	if err != nil {
		log.Fatal("failed to get positioning constant: ", err)
	}
	log.Println("Calculated Positioning Constant for the registry: ", positioningConstant)

	log.Println("Preparing a batch call request")
	var reqs []rpc.BatchElem
	var results []*string
	var keeperPerBlockIndex []uint64
	var keeperPerBlockAddress []common.Address
	for block := from; block <= to; block++ {
		callOpts := bind.CallOpts{
			Context:     ctx,
			BlockNumber: big.NewInt(0).SetUint64(block),
		}

		registryConfig, err := registryClient.GetConfig(&callOpts)
		if err != nil {
			log.Fatal("failed to fetch registry config: ", err)
		}
		blockCountPerTurn := registryConfig.BlockCountPerTurn.Uint64()

		keepersList, err := registryClient.GetKeeperList(&callOpts)
		if err != nil {
			log.Fatal("failed to fetch keepers list: ", err)
		}

		keeperIndex := (uint64(positioningConstant) + ((block - (block % blockCountPerTurn)) / blockCountPerTurn)) % uint64(len(keepersList))
		payload, err := keeper.Registry1_1ABI.Pack("checkUpkeep", big.NewInt(upkeepId), keepersList[keeperIndex])
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

		var result string
		reqs = append(reqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				args,
				// The block at which we want to inspect the upkeep state
				hexutil.EncodeUint64(block),
			},
			Result: &result,
		})

		results = append(results, &result)
		keeperPerBlockIndex = append(keeperPerBlockIndex, keeperIndex)
		keeperPerBlockAddress = append(keeperPerBlockAddress, keepersList[keeperIndex])
	}

	log.Println("Doing batch call to check upkeeps")
	if err := k.rpcClient.BatchCallContext(ctx, reqs); err != nil {
		log.Fatal("failed to batch call checkUpkeep: ", err)
	}

	log.Println("Parsing batch call response")
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
	var parsedResults []result
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

		returnValues, err := checkUpkeepArguments.UnpackValues(hexutil.MustDecode(*results[i]))
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
