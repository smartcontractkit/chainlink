package v2plusscripts

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_sub_owner"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper_consumer_example"
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
	_, tx, _, err := vrf_coordinator_v2_5.DeployVRFCoordinatorV25(
		e.Owner,
		e.Ec,
		common.HexToAddress(bhsAddress))
	helpers.PanicErr(err)
	coordinatorAddress = helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)

	// Set LINK and LINK ETH
	coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(coordinatorAddress, e.Ec)
	helpers.PanicErr(err)

	if linkAddress != "" && linkEthAddress != "" {
		linkTx, err := coordinator.SetLINKAndLINKNativeFeed(e.Owner,
			common.HexToAddress(linkAddress), common.HexToAddress(linkEthAddress))
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, linkTx, e.ChainID)
	}
	return coordinatorAddress
}

func DeployBatchCoordinatorV2(e helpers.Environment, coordinatorAddress common.Address) (batchCoordinatorAddress common.Address) {
	_, tx, _, err := batch_vrf_coordinator_v2plus.DeployBatchVRFCoordinatorV2Plus(e.Owner, e.Ec, coordinatorAddress)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func EoaAddConsumerToSub(
	e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25,
	subID *big.Int,
	consumerAddress string,
) {
	txadd, err := coordinator.AddConsumer(e.Owner, subID, common.HexToAddress(consumerAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, txadd, e.ChainID)
}

func EoaCreateSub(e helpers.Environment, coordinator vrf_coordinator_v2_5.VRFCoordinatorV25) {
	tx, err := coordinator.CreateSubscription(e.Owner)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

// returns subscription ID that belongs to the given owner. Returns result found first
func FindSubscriptionID(e helpers.Environment, coordinator *vrf_coordinator_v2_5.VRFCoordinatorV25) *big.Int {
	// Use most recent 500 blocks as search window.
	head, err := e.Ec.BlockNumber(context.Background())
	helpers.PanicErr(err)
	fopts := &bind.FilterOpts{
		Start: head - 500,
	}

	subscriptionIterator, err := coordinator.FilterSubscriptionCreated(fopts, nil)
	helpers.PanicErr(err)

	if !subscriptionIterator.Next() {
		helpers.PanicErr(fmt.Errorf("expected at least 1 subID for the given owner %s", e.Owner.From.Hex()))
	}
	return subscriptionIterator.Event.SubId
}

func EoaDeployConsumer(e helpers.Environment,
	coordinatorAddress string,
	linkAddress string) (
	consumerAddress common.Address) {
	_, tx, _, err := vrf_v2plus_sub_owner.DeployVRFV2PlusExternalSubOwnerExample(
		e.Owner,
		e.Ec,
		common.HexToAddress(coordinatorAddress),
		common.HexToAddress(linkAddress))
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func EoaFundSubWithLink(
	e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25,
	linkAddress string, amount,
	subID *big.Int,
) {
	linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(linkAddress), e.Ec)
	helpers.PanicErr(err)
	bal, err := linkToken.BalanceOf(nil, e.Owner.From)
	helpers.PanicErr(err)
	fmt.Println("Initial account balance (Juels):", bal, e.Owner.From.String(), "Funding amount:", amount.String())
	b, err := utils.ABIEncode(`[{"type":"uint256"}]`, subID)
	helpers.PanicErr(err)
	tx, err := linkToken.TransferAndCall(e.Owner, coordinator.Address(), amount, b)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("sub ID: %d", subID))
}

func EoaFundSubWithNative(e helpers.Environment, coordinatorAddress common.Address, subID *big.Int, amount *big.Int) {
	coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(coordinatorAddress, e.Ec)
	helpers.PanicErr(err)
	e.Owner.Value = amount
	tx, err := coordinator.FundSubscriptionWithNative(e.Owner, subID)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func PrintCoordinatorConfig(coordinator *vrf_coordinator_v2_5.VRFCoordinatorV25) {
	cfg, err := coordinator.SConfig(nil)
	helpers.PanicErr(err)

	fmt.Printf("Coordinator config: %+v\n", cfg)
}

func SetCoordinatorConfig(
	e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25,
	minConfs uint16,
	maxGasLimit uint32,
	stalenessSeconds uint32,
	gasAfterPayment uint32,
	fallbackWeiPerUnitLink *big.Int,
	fulfillmentFlatFeeNativePPM uint32,
	fulfillmentFlatFeeLinkDiscountPPM uint32,
	nativePremiumPercentage uint8,
	linkPremiumPercentage uint8,
) {
	tx, err := coordinator.SetConfig(
		e.Owner,
		minConfs,               // minRequestConfirmations
		maxGasLimit,            // max gas limit
		stalenessSeconds,       // stalenessSeconds
		gasAfterPayment,        // gasAfterPaymentCalculation
		fallbackWeiPerUnitLink, // 0.01 eth per link fallbackLinkPrice
		fulfillmentFlatFeeNativePPM,
		fulfillmentFlatFeeLinkDiscountPPM,
		nativePremiumPercentage,
		linkPremiumPercentage,
	)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func RegisterCoordinatorProvingKey(e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25, uncompressed string, gasLaneMaxGas uint64) {
	pubBytes, err := hex.DecodeString(uncompressed)
	helpers.PanicErr(err)
	pk, err := crypto.UnmarshalPubkey(pubBytes)
	helpers.PanicErr(err)
	tx, err := coordinator.RegisterProvingKey(e.Owner,
		[2]*big.Int{pk.X, pk.Y}, gasLaneMaxGas)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(
		context.Background(),
		e.Ec,
		tx,
		e.ChainID,
		fmt.Sprintf("Uncompressed public key: %s,", uncompressed),
		fmt.Sprintf("Gas Lane Max Gas: %d,", gasLaneMaxGas),
	)
}

func RegisterMigratableCoordinator(
	e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25,
	coordinatorMigrateToAddress common.Address,
) {
	tx, err := coordinator.RegisterMigratableCoordinator(e.Owner, coordinatorMigrateToAddress)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(
		context.Background(),
		e.Ec,
		tx,
		e.ChainID,
		fmt.Sprintf("Coordinator %s registered migratable coordinator %s", coordinator.Address().String(), coordinatorMigrateToAddress.String()),
	)
}

func MigrateSub(
	e helpers.Environment,
	coordinatorMigrateSubFrom vrf_coordinator_v2_5.VRFCoordinatorV25,
	coordinatorMigrateSubTo common.Address,
	subID *big.Int,
) {
	tx, err := coordinatorMigrateSubFrom.Migrate(e.Owner, subID, coordinatorMigrateSubTo)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(
		context.Background(),
		e.Ec,
		tx,
		e.ChainID,
		fmt.Sprintf("Sub Migrated from Coordinator: %s,", coordinatorMigrateSubFrom.Address().String()),
		fmt.Sprintf("Sub Migrated TO Coordinator: %s,", coordinatorMigrateSubTo.String()),
		fmt.Sprintf("Sub ID which was migrated: %s,", subID.String()),
	)
}

func WrapperDeploy(
	e helpers.Environment,
	link, linkEthFeed, coordinator common.Address, subID *big.Int,
) common.Address {
	address, tx, _, err := vrfv2plus_wrapper.DeployVRFV2PlusWrapper(e.Owner, e.Ec,
		link,
		linkEthFeed,
		coordinator,
		subID)
	helpers.PanicErr(err)

	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	fmt.Println("VRFV2Wrapper address:", address)

	return address
}

func WrapperConfigure(
	e helpers.Environment,
	wrapperAddress common.Address,
	wrapperGasOverhead uint,
	coordinatorGasOverheadNative, coordinatorGasOverheadLink uint,
	coordinatorGasOverheadPerWord uint,
	nativePremiumPercentage, linkPremiumPercentage uint,
	keyHash string,
	maxNumWords uint,
	fallbackWeiPerUnitLink *big.Int,
	stalenessSeconds uint32,
	fulfillmentFlatFeeNativePPM uint32,
	fulfillmentFlatFeeLinkDiscountPPM uint32,
) {
	wrapper, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(wrapperAddress, e.Ec)
	helpers.PanicErr(err)

	tx, err := wrapper.SetConfig(
		e.Owner,
		uint32(wrapperGasOverhead),
		uint32(coordinatorGasOverheadNative),
		uint32(coordinatorGasOverheadLink),
		uint16(coordinatorGasOverheadPerWord),
		uint8(nativePremiumPercentage),
		uint8(linkPremiumPercentage),
		common.HexToHash(keyHash),
		uint8(maxNumWords),
		stalenessSeconds,
		fallbackWeiPerUnitLink,
		fulfillmentFlatFeeNativePPM,
		fulfillmentFlatFeeLinkDiscountPPM,
	)

	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func PrintWrapperConfig(wrapper *vrfv2plus_wrapper.VRFV2PlusWrapper) {
	cfg, err := wrapper.GetConfig(nil)
	helpers.PanicErr(err)
	fmt.Printf("Wrapper config: %+v\n", cfg)
	fmt.Printf("Wrapper Keyhash: %s\n", fmt.Sprintf("0x%x", cfg.KeyHash))
}

func WrapperConsumerDeploy(
	e helpers.Environment,
	link, wrapper common.Address,
) common.Address {
	address, tx, _, err := vrfv2plus_wrapper_consumer_example.DeployVRFV2PlusWrapperConsumerExample(e.Owner, e.Ec,
		wrapper)
	helpers.PanicErr(err)

	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	fmt.Printf("VRFV2WrapperConsumerExample address: %s\n", address)
	return address
}

func EoaV2PlusLoadTestConsumerWithMetricsDeploy(e helpers.Environment, consumerCoordinator string) (consumerAddress common.Address) {
	_, tx, _, err := vrf_v2plus_load_test_with_metrics.DeployVRFV2PlusLoadTestWithMetrics(
		e.Owner,
		e.Ec,
		common.HexToAddress(consumerCoordinator),
	)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}
