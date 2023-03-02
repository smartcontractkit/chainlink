package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/nocancel_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrfv2_wrapper"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrfv2_wrapper_consumer_example"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func deployBHS(e helpers.Environment) (blockhashStoreAddress common.Address) {
	_, tx, _, err := blockhash_store.DeployBlockhashStore(e.Owner, e.Ec)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployBatchBHS(e helpers.Environment, bhsAddress common.Address) (batchBHSAddress common.Address) {
	_, tx, _, err := batch_blockhash_store.DeployBatchBlockhashStore(e.Owner, e.Ec, bhsAddress)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployCoordinator(
	e helpers.Environment,
	linkAddress string,
	bhsAddress string,
	linkEthAddress string,
) (coordinatorAddress common.Address) {
	_, tx, _, err := vrf_coordinator_v2.DeployVRFCoordinatorV2(
		e.Owner,
		e.Ec,
		common.HexToAddress(linkAddress),
		common.HexToAddress(bhsAddress),
		common.HexToAddress(linkEthAddress))
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployNoCancelCoordinator(
	e helpers.Environment,
	linkAddress string,
	bhsAddress string,
	linkEthAddress string,
) (coordinatorAddress common.Address) {
	addr, _, _, err := nocancel_vrf_coordinator_v2.DeployNoCancelVRFCoordinatorV2(
		e.Owner,
		e.Ec,
		common.HexToAddress(linkAddress),
		common.HexToAddress(bhsAddress),
		common.HexToAddress(linkEthAddress))
	helpers.PanicErr(err)
	helpers.ConfirmCodeAt(context.Background(), e.Ec, addr, e.ChainID)
	return addr
}

func deployBatchCoordinatorV2(e helpers.Environment, coordinatorAddress common.Address) (batchCoordinatorAddress common.Address) {
	_, tx, _, err := batch_vrf_coordinator_v2.DeployBatchVRFCoordinatorV2(e.Owner, e.Ec, coordinatorAddress)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func eoaAddConsumerToSub(e helpers.Environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, subID uint64, consumerAddress string) {
	txadd, err := coordinator.AddConsumer(e.Owner, subID, common.HexToAddress(consumerAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, txadd, e.ChainID)
}

func eoaCreateSub(e helpers.Environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2) {
	tx, err := coordinator.CreateSubscription(e.Owner)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func eoaDeployConsumer(e helpers.Environment, coordinatorAddress string, linkAddress string) (consumerAddress common.Address) {
	_, tx, _, err := vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(
		e.Owner,
		e.Ec,
		common.HexToAddress(coordinatorAddress),
		common.HexToAddress(linkAddress))
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func eoaFundSubscription(e helpers.Environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, linkAddress string, amount *big.Int, subID uint64) {
	linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(linkAddress), e.Ec)
	helpers.PanicErr(err)
	bal, err := linkToken.BalanceOf(nil, e.Owner.From)
	helpers.PanicErr(err)
	fmt.Println("Initial account balance:", bal, e.Owner.From.String(), "Funding amount:", amount.String())
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, subID)
	helpers.PanicErr(err)
	e.Owner.GasLimit = 500000
	tx, err := linkToken.TransferAndCall(e.Owner, coordinator.Address(), amount, b)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("sub ID: %d", subID))
}

func printCoordinatorConfig(coordinator *vrf_coordinator_v2.VRFCoordinatorV2) {
	cfg, err := coordinator.GetConfig(nil)
	helpers.PanicErr(err)

	feeConfig, err := coordinator.GetFeeConfig(nil)
	helpers.PanicErr(err)

	fmt.Printf("Coordinator config: %+v\n", cfg)
	fmt.Printf("Coordinator fee config: %+v\n", feeConfig)
}

func setCoordinatorConfig(
	e helpers.Environment,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2,
	minConfs uint16,
	maxGasLimit uint32,
	stalenessSeconds uint32,
	gasAfterPayment uint32,
	fallbackWeiPerUnitLink *big.Int,
	feeConfig vrf_coordinator_v2.VRFCoordinatorV2FeeConfig,
) {
	tx, err := coordinator.SetConfig(
		e.Owner,
		minConfs,               // minRequestConfirmations
		maxGasLimit,            // max gas limit
		stalenessSeconds,       // stalenessSeconds
		gasAfterPayment,        // gasAfterPaymentCalculation
		fallbackWeiPerUnitLink, // 0.01 eth per link fallbackLinkPrice
		feeConfig,
	)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func registerCoordinatorProvingKey(e helpers.Environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, uncompressed string, oracleAddress string) {
	pubBytes, err := hex.DecodeString(uncompressed)
	helpers.PanicErr(err)
	pk, err := crypto.UnmarshalPubkey(pubBytes)
	helpers.PanicErr(err)
	tx, err := coordinator.RegisterProvingKey(e.Owner,
		common.HexToAddress(oracleAddress),
		[2]*big.Int{pk.X, pk.Y})
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(
		context.Background(),
		e.Ec,
		tx,
		e.ChainID,
		fmt.Sprintf("Uncompressed public key: %s,", uncompressed),
		fmt.Sprintf("Oracle address: %s,", oracleAddress),
	)
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

func getRlpHeaders(env helpers.Environment, blockNumbers []*big.Int) (headers [][]byte, err error) {
	headers = [][]byte{}
	for _, blockNum := range blockNumbers {
		// Avalanche block headers are special, handle them by using the avalanche rpc client
		// rather than the regular go-ethereum ethclient.
		if helpers.IsAvaxNetwork(env.ChainID) {
			// Get child block since it's the one that has the parent hash in its header.
			h, err := env.AvaxEc.HeaderByNumber(
				context.Background(),
				new(big.Int).Set(blockNum).Add(blockNum, big.NewInt(1)),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get header: %+v", err)
			}
			// We can still use vanilla go-ethereum rlp.EncodeToBytes, see e.g
			// https://github.com/ava-labs/coreth/blob/e3ca41bf5295a9a7ca1aeaf29d541fcbb94f79b1/core/types/hashing.go#L49-L57.
			rlpHeader, err := rlp.EncodeToBytes(h)
			if err != nil {
				return nil, fmt.Errorf("failed to encode rlp: %+v", err)
			}

			// Sanity check - can be un-commented if storeVerifyHeader is failing due to unexpected
			// blockhash.
			//bh := crypto.Keccak256Hash(rlpHeader)
			//fmt.Println("Calculated BH:", bh.String(),
			//	"fetched BH:", h.Hash(),
			//	"block number:", new(big.Int).Set(blockNum).Add(blockNum, big.NewInt(1)).String())

			headers = append(headers, rlpHeader)
		} else {
			// Get child block since it's the one that has the parent hash in its header.
			h, err := env.Ec.HeaderByNumber(
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

			headers = append(headers, rlpHeader)
		}
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

func wrapperDeploy(
	e helpers.Environment,
	link, linkEthFeed, coordinator common.Address,
) (common.Address, uint64) {
	address, tx, _, err := vrfv2_wrapper.DeployVRFV2Wrapper(e.Owner, e.Ec,
		link,
		linkEthFeed,
		coordinator)
	helpers.PanicErr(err)

	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	fmt.Printf("VRFV2Wrapper address: %s\n", address)

	wrapper, err := vrfv2_wrapper.NewVRFV2Wrapper(address, e.Ec)
	helpers.PanicErr(err)

	subID, err := wrapper.SUBSCRIPTIONID(nil)
	helpers.PanicErr(err)

	return address, subID
}

func wrapperConfigure(
	e helpers.Environment,
	wrapperAddress common.Address,
	wrapperGasOverhead, coordinatorGasOverhead, premiumPercentage uint,
	keyHash string,
	maxNumWords uint,
) {
	wrapper, err := vrfv2_wrapper.NewVRFV2Wrapper(wrapperAddress, e.Ec)
	helpers.PanicErr(err)

	tx, err := wrapper.SetConfig(
		e.Owner,
		uint32(wrapperGasOverhead),
		uint32(coordinatorGasOverhead),
		uint8(premiumPercentage),
		common.HexToHash(keyHash),
		uint8(maxNumWords))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func wrapperConsumerDeploy(
	e helpers.Environment,
	link, wrapper common.Address,
) common.Address {
	address, tx, _, err := vrfv2_wrapper_consumer_example.DeployVRFV2WrapperConsumerExample(e.Owner, e.Ec,
		link,
		wrapper)
	helpers.PanicErr(err)

	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	fmt.Printf("VRFV2WrapperConsumerExample address: %s\n", address)
	return address
}
