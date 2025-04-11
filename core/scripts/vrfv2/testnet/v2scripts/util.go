package v2scripts

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_test_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2_wrapper_consumer_example"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func DeployBHS(e helpers.Environment) (blockhashStoreAddress common.Address) {
	_, tx, _, err := blockhash_store.DeployBlockhashStore(e.Owner, e.Ec)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func DeployBatchBHS(e helpers.Environment, bhsAddress common.Address) (batchBHSAddress common.Address) {
	_, tx, _, err := batch_blockhash_store.DeployBatchBlockhashStore(e.Owner, e.Ec, bhsAddress)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func DeployCoordinator(
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

func DeployTestCoordinator(
	e helpers.Environment,
	linkAddress string,
	bhsAddress string,
	linkEthAddress string,
) (coordinatorAddress common.Address) {
	_, tx, _, err := vrf_coordinator_test_v2.DeployVRFCoordinatorTestV2(
		e.Owner,
		e.Ec,
		common.HexToAddress(linkAddress),
		common.HexToAddress(bhsAddress),
		common.HexToAddress(linkEthAddress))
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func DeployBatchCoordinatorV2(e helpers.Environment, coordinatorAddress common.Address) (batchCoordinatorAddress common.Address) {
	_, tx, _, err := batch_vrf_coordinator_v2.DeployBatchVRFCoordinatorV2(e.Owner, e.Ec, coordinatorAddress)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func EoaAddConsumerToSub(e helpers.Environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, subID uint64, consumerAddress string) {
	txadd, err := coordinator.AddConsumer(e.Owner, subID, common.HexToAddress(consumerAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, txadd, e.ChainID)
}

func EoaCreateSub(e helpers.Environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2) {
	tx, err := coordinator.CreateSubscription(e.Owner)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func EoaDeployConsumer(e helpers.Environment, coordinatorAddress string, linkAddress string) (consumerAddress common.Address) {
	_, tx, _, err := vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(
		e.Owner,
		e.Ec,
		common.HexToAddress(coordinatorAddress),
		common.HexToAddress(linkAddress))
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func EoaFundSubscription(e helpers.Environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, linkAddress string, amount *big.Int, subID uint64) {
	linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(linkAddress), e.Ec)
	helpers.PanicErr(err)
	bal, err := linkToken.BalanceOf(nil, e.Owner.From)
	helpers.PanicErr(err)
	fmt.Println("Initial account balance:", bal, e.Owner.From.String(), "Funding amount:", amount.String())
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, subID)
	helpers.PanicErr(err)
	tx, err := linkToken.TransferAndCall(e.Owner, coordinator.Address(), amount, b)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("sub ID: %d", subID))
}

func PrintCoordinatorConfig(coordinator *vrf_coordinator_v2.VRFCoordinatorV2) {
	cfg, err := coordinator.GetConfig(nil)
	helpers.PanicErr(err)

	feeConfig, err := coordinator.GetFeeConfig(nil)
	helpers.PanicErr(err)

	fmt.Printf("Coordinator config: %+v\n", cfg)
	fmt.Printf("Coordinator fee config: %+v\n", feeConfig)
}

func SetCoordinatorConfig(
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

func RegisterCoordinatorProvingKey(e helpers.Environment, coordinator vrf_coordinator_v2.VRFCoordinatorV2, uncompressed string, oracleAddress string) {
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

func WrapperDeploy(
	e helpers.Environment,
	link, linkEthFeed, coordinator common.Address,
) (common.Address, uint64) {
	address, tx, _, err := vrfv2_wrapper.DeployVRFV2Wrapper(e.Owner, e.Ec,
		link,
		linkEthFeed,
		coordinator)
	helpers.PanicErr(err)

	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	fmt.Println("VRFV2Wrapper address:", address)

	wrapper, err := vrfv2_wrapper.NewVRFV2Wrapper(address, e.Ec)
	helpers.PanicErr(err)

	subID, err := wrapper.SUBSCRIPTIONID(nil)
	helpers.PanicErr(err)
	fmt.Println("VRFV2Wrapper subscription id:", subID)

	return address, subID
}

func WrapperConfigure(
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

func PrintWrapperConfig(wrapper *vrfv2_wrapper.VRFV2Wrapper) {
	cfg, err := wrapper.GetConfig(nil)
	helpers.PanicErr(err)
	fmt.Printf("Wrapper config: %+v\n", cfg)
	fmt.Printf("Wrapper Keyhash: %s\n", fmt.Sprintf("0x%x", cfg.KeyHash))
}

func WrapperConsumerDeploy(
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

func EoaLoadTestConsumerWithMetricsDeploy(e helpers.Environment, consumerCoordinator string) (consumerAddress common.Address) {
	_, tx, _, err := vrf_load_test_with_metrics.DeployVRFV2LoadTestWithMetrics(
		e.Owner,
		e.Ec,
		common.HexToAddress(consumerCoordinator),
	)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func ClosestBlock(e helpers.Environment, batchBHSAddress common.Address, blockMissingBlockhash uint64, batchSize uint64) (uint64, error) {
	batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(batchBHSAddress, e.Ec)
	if err != nil {
		return 0, err
	}
	startBlock := blockMissingBlockhash + 1
	endBlock := startBlock + batchSize
	for {
		latestBlock, err := e.Ec.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return 0, err
		}
		if latestBlock.Number.Uint64() < endBlock {
			return 0, errors.New("closest block with blockhash not found")
		}
		var blockRange []*big.Int
		for i := startBlock; i <= endBlock; i++ {
			blockRange = append(blockRange, big.NewInt(int64(i)))
		}
		fmt.Println("Searching range", startBlock, "-", endBlock, "inclusive")
		hashes, err := batchBHS.GetBlockhashes(nil, blockRange)
		if err != nil {
			return 0, err
		}
		for i, hash := range hashes {
			if hash != (common.Hash{}) {
				fmt.Println("found closest block:", startBlock+uint64(i), "hash:", hexutil.Encode(hash[:]))
				fmt.Println("distance from missing block:", startBlock+uint64(i)-blockMissingBlockhash)
				return startBlock + uint64(i), nil
			}
		}
		startBlock = endBlock + 1
		endBlock = startBlock + batchSize
	}
}
