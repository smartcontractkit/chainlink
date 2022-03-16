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
)

const (
	defaultMaxBlocksRange = 1000
)

// UpkeepHistory prints the checkUpkeep status and keeper responsibility for a given upkeep in a set block range
func (k *Keeper) UpkeepHistory(ctx context.Context, upkeepId int64, from, to uint64) {
	// There must not be a large different between boundaries
	if to-from > defaultMaxBlocksRange {
		log.Fatal("blocks range difference must not more than 1000")
	}

	registryAddr, registryClient := k.GetRegistry(ctx)

	// Get positioning constant of the current registry
	positioningConstant, err := keeper.CalcPositioningConstant(upkeepId, ethkey.EIP55AddressFromAddress(registryAddr))
	if err != nil {
		log.Fatal("failed to get positioning constant: ", err)
	}
	log.Println("Calculated Positioning Constant for the registry: ", positioningConstant)

	callOpts := bind.CallOpts{Context: ctx}

	log.Println("Fetching registry config")
	registryConfig, err := registryClient.GetConfig(&callOpts)
	if err != nil {
		log.Fatal("failed to fetch registry config: ", err)
	}
	blockCountPerTurn := registryConfig.BlockCountPerTurn.Uint64()

	log.Println("Fetching registry keepers list")
	keepersList, err := registryClient.GetKeeperList(&callOpts)
	if err != nil {
		log.Fatal("failed to fetch keepers list: ", err)
	}

	log.Println("Preparing a batch call request")
	var reqs []rpc.BatchElem
	var results []*string
	var keeperPerBlock []uint64
	for block := from; block <= to; block++ {
		keeperIndex := (uint64(positioningConstant) + ((block - (block % blockCountPerTurn)) / blockCountPerTurn)) % uint64(len(keepersList))
		payload, err := keeper.RegistryABI.Pack("checkUpkeep", big.NewInt(upkeepId), keepersList[keeperIndex])
		if err != nil {
			log.Fatal("failed to pack checkUpkeep: ", err)
		}

		var result string
		reqs = append(reqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   k.cfg.RegistryAddress,
					"data": hexutil.Bytes(payload),
				},
				// The block at which we want to inspect the upkeep state
				hexutil.EncodeUint64(block),
			},
			Result: &result,
		})

		results = append(results, &result)
		keeperPerBlock = append(keeperPerBlock, keeperIndex)
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
				keeperIndex:   keeperPerBlock[i],
				keeperAddress: keepersList[keeperPerBlock[i]],
				reason:        strings.TrimPrefix(req.Error.Error(), "execution reverted: "),
			})
			continue
		}

		returnValues, err := keeper.RegistryABI.Methods["checkUpkeep"].
			Outputs.UnpackValues(hexutil.MustDecode(*results[i]))
		if err != nil {
			log.Fatal("unpack checkUpkeep return: ", err, *results[i])
		}

		parsedResults = append(parsedResults, result{
			block:          uint64(i) + from,
			checkUpkeep:    true,
			keeperIndex:    keeperPerBlock[i],
			keeperAddress:  keepersList[keeperPerBlock[i]],
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

func (k *Keeper) getPositioningConstant(upkeepId int64) (int32, error) {
	registryAddress, err := ethkey.NewEIP55Address(k.cfg.RegistryAddress)
	if err != nil {
		log.Println("failed to parse registry address: ", err)
		return 0, err
	}

	positioningConstant, err := keeper.CalcPositioningConstant(upkeepId, registryAddress)
	if err != nil {
		log.Println("failed to Calc Positioning Constant: ", err)
		return 0, err
	}

	return positioningConstant, nil
}
