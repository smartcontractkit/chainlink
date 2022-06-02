package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_external_sub_owner_example"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func confirmTXMined(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64, txInfo ...string) {
	fmt.Println("Executing TX", helpers.ExplorerLink(chainID, transaction.Hash()), txInfo)
	receipt, err := bind.WaitMined(context, client, transaction)
	helpers.PanicErr(err)
	fmt.Println("TX", receipt.TxHash, "mined. \nBlock Number:", receipt.BlockNumber, "\nGas Used: ", receipt.GasUsed)
}

func confirmContractDeployed(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64) (address common.Address) {
	fmt.Println("Executing contract deployment, TX:", helpers.ExplorerLink(chainID, transaction.Hash()))
	contractAddress, err := bind.WaitDeployed(context, client, transaction)
	helpers.PanicErr(err)
	fmt.Println("Contract Address:", contractAddress.String())
	return contractAddress
}

func deployBHS(e environment) (blockhashStoreAddress common.Address) {
	_, tx, _, err := blockhash_store.DeployBlockhashStore(e.owner, e.ec)
	helpers.PanicErr(err)
	return confirmContractDeployed(context.Background(), e.ec, tx, e.chainID)
}

func deployCoordinator(
	e environment,
	linkAddress string,
	bhsAddress string,
	linkEthAddress string,
) (coordinatorAddress common.Address) {
	_, tx, _, err := vrf_coordinator_v2.DeployVRFCoordinatorV2(
		e.owner,
		e.ec,
		common.HexToAddress(linkAddress),
		common.HexToAddress(bhsAddress),
		common.HexToAddress(linkEthAddress))
	helpers.PanicErr(err)
	return confirmContractDeployed(context.Background(), e.ec, tx, e.chainID)
}

func eoaAddConsumerToSub(e environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, subID uint64, consumerAddress string) {
	txadd, err := coordinator.AddConsumer(e.owner, subID, common.HexToAddress(consumerAddress))
	helpers.PanicErr(err)
	confirmTXMined(context.Background(), e.ec, txadd, e.chainID)
}

func eoaCreateSub(e environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2) {
	tx, err := coordinator.CreateSubscription(e.owner)
	helpers.PanicErr(err)
	confirmTXMined(context.Background(), e.ec, tx, e.chainID)
}

func eoaDeployConsumer(e environment, coordinatorAddress string, linkAddress string) (consumerAddress common.Address) {
	_, tx, _, err := vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(
		e.owner,
		e.ec,
		common.HexToAddress(coordinatorAddress),
		common.HexToAddress(linkAddress))
	helpers.PanicErr(err)
	return confirmContractDeployed(context.Background(), e.ec, tx, e.chainID)
}

func eoaFundSubscription(e environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, linkAddress string, amount *big.Int, subID uint64) {
	linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(linkAddress), e.ec)
	helpers.PanicErr(err)
	bal, err := linkToken.BalanceOf(nil, e.owner.From)
	helpers.PanicErr(err)
	fmt.Println("Initial account balance:", bal, e.owner.From.String(), "Funding amount:", amount.String())
	b, err := utils.GenericEncode([]string{"uint64"}, subID)
	helpers.PanicErr(err)
	e.owner.GasLimit = 500000
	tx, err := linkToken.TransferAndCall(e.owner, coordinator.Address(), amount, b)
	helpers.PanicErr(err)
	confirmTXMined(context.Background(), e.ec, tx, e.chainID, fmt.Sprintf("sub ID: %d", subID))
}

func printCoordinatorConfig(e environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2) {
	cfg, err := coordinator.GetConfig(nil)
	helpers.PanicErr(err)

	feeConfig, err := coordinator.GetFeeConfig(nil)
	helpers.PanicErr(err)

	fmt.Printf("Coordinator config: %+v\n", cfg)
	fmt.Printf("Coordinator fee config: %+v\n", feeConfig)
}

func setCoordinatorConfig(
	e environment,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2,
	minConfs uint16,
	maxGasLimit uint32,
	stalenessSeconds uint32,
	gasAfterPayment uint32,
	fallbackWeiPerUnitLink *big.Int,
	feeConfig vrf_coordinator_v2.VRFCoordinatorV2FeeConfig,
) {
	tx, err := coordinator.SetConfig(
		e.owner,
		minConfs,               // minRequestConfirmations
		maxGasLimit,            // max gas limit
		stalenessSeconds,       // stalenessSeconds
		gasAfterPayment,        // gasAfterPaymentCalculation
		fallbackWeiPerUnitLink, // 0.01 eth per link fallbackLinkPrice
		feeConfig,
	)
	helpers.PanicErr(err)
	confirmTXMined(context.Background(), e.ec, tx, e.chainID)
}

func registerCoordinatorProvingKey(e environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, uncompressed string, oracleAddress string) {
	pubBytes, err := hex.DecodeString(uncompressed)
	helpers.PanicErr(err)
	pk, err := crypto.UnmarshalPubkey(pubBytes)
	helpers.PanicErr(err)
	tx, err := coordinator.RegisterProvingKey(e.owner,
		common.HexToAddress(oracleAddress),
		[2]*big.Int{pk.X, pk.Y})
	helpers.PanicErr(err)
	confirmTXMined(
		context.Background(),
		e.ec,
		tx,
		e.chainID,
		fmt.Sprintf("Uncompressed public key: %s,", uncompressed),
		fmt.Sprintf("Oracle address: %s,", oracleAddress),
	)
}

func parseIntSlice(arg string) (ret []*big.Int) {
	parts := strings.Split(arg, ",")
	ret = []*big.Int{}
	for _, part := range parts {
		ret = append(ret, decimal.RequireFromString(part).BigInt())
	}
	return ret
}

func parseAddressSlice(arg string) (ret []common.Address) {
	parts := strings.Split(arg, ",")
	ret = []common.Address{}
	for _, part := range parts {
		ret = append(ret, common.HexToAddress(part))
	}
	return
}

func parseHashSlice(arg string) (ret []common.Hash) {
	parts := strings.Split(arg, ",")
	ret = []common.Hash{}
	for _, part := range parts {
		ret = append(ret, common.HexToHash(part))
	}
	return
}

// decreasingBlockRange creates a continugous block range starting with
// block `start` and ending at block `end`.
func decreasingBlockRange(start, end *big.Int) (ret []*big.Int, err error) {
	if start.Cmp(end) == -1 {
		return nil, fmt.Errorf("start (%s) must be greater than end (%s)", start.String(), end.String())
	}
	ret = []*big.Int{}
	for i := new(big.Int).Set(start); i.Cmp(end) >= 0; i.Sub(i, big.NewInt(1)) {
		ret = append(ret, new(big.Int).Set(i))
	}
	return
}

func getRlpHeaders(ec *ethclient.Client, blockNumbers []*big.Int) (headers [][]byte, err error) {
	headers = [][]byte{}
	for _, blockNum := range blockNumbers {
		// Get child block since it's the one that has the parent hash in it's header.
		h, err := ec.HeaderByNumber(
			context.Background(),
			new(big.Int).Set(blockNum).Add(blockNum, big.NewInt(1)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get header: %+v", err)
		}
		rlpHeader, err := rlp.EncodeToBytes(h)
		if err != nil {
			return nil, fmt.Errorf("failed to encode rlp: %+v", err)
		}
		// Uncomment in case storeVerifyHeader calls are reverting, there may be an issue with the RLP
		// encoding.
		// h2, err := ec.HeaderByNumber(context.Background(), blockNum)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to get header: %v", err)
		// }
		// fmt.Println("block number:", blockNum, "blockhash:", h2.Hash(), "encoded header of next block:", common.Bytes2Hex(rlpHeader))
		headers = append(headers, rlpHeader)
	}
	return
}

// binarySearch finds the highest value within the range bottom-top at which the test function is
// true.
func binarySearch(top, bottom *big.Int, test func(amount *big.Int) bool) *big.Int {
	var runs int
	// While the difference between top and bottom is > 1
	for new(big.Int).Sub(top, bottom).Cmp(big.NewInt(1)) > 0 {
		// Calculate midpoint between top and bottom
		midpoint := new(big.Int).Sub(top, bottom)
		midpoint.Div(midpoint, big.NewInt(2))
		midpoint.Add(midpoint, bottom)

		// Check if the midpoint amount is withdrawable
		if test(midpoint) {
			bottom = midpoint
		} else {
			top = midpoint
		}

		runs++
		if runs%10 == 0 {
			fmt.Printf("Searching... current range %s-%s\n", bottom.String(), top.String())
		}
	}

	return bottom
}
