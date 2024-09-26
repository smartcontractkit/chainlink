// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {MockLinkToken} from "../../mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5_Arbitrum} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5_Arbitrum.sol";
import {VRFV2PlusWrapper_Arbitrum} from "../dev/VRFV2PlusWrapper_Arbitrum.sol";
import {ArbGasInfo} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";

contract VRFV2PlusWrapperArbitrumTest is BaseTest {
  /// @dev ARBGAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbGasInfo.sol#L10
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  ArbGasInfo private constant ARBGAS = ArbGasInfo(ARBGAS_ADDR);

  address internal constant DEPLOYER = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
  bytes32 private vrfKeyHash = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";
  uint256 private s_wrapperSubscriptionId;

  ExposedVRFCoordinatorV2_5_Arbitrum private s_testCoordinator;
  MockLinkToken private s_linkToken;
  MockV3Aggregator private s_linkNativeFeed;
  VRFV2PlusWrapper_Arbitrum private s_wrapper;

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(DEPLOYER, 10_000 ether);
    vm.stopPrank();
    vm.startPrank(DEPLOYER);

    // Deploy link token and link/native feed.
    s_linkToken = new MockLinkToken();
    s_linkNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Deploy coordinator.
    s_testCoordinator = new ExposedVRFCoordinatorV2_5_Arbitrum(address(0));

    // Create subscription for all future wrapper contracts.
    s_wrapperSubscriptionId = s_testCoordinator.createSubscription();

    // Deploy wrapper.
    s_wrapper = new VRFV2PlusWrapper_Arbitrum(
      address(s_linkToken),
      address(s_linkNativeFeed),
      address(s_testCoordinator),
      uint256(s_wrapperSubscriptionId)
    );

    // Configure the wrapper.
    s_wrapper.setConfig(
      100_000, // wrapper gas overhead
      200_000, // coordinator gas overhead native
      220_000, // coordinator gas overhead link
      500, // coordinator gas overhead per word
      15, // native premium percentage,
      10, // link premium percentage
      vrfKeyHash, // keyHash
      10, // max number of words,
      1, // stalenessSeconds
      50000000000000000, // fallbackWeiPerUnitLink
      500_000, // fulfillmentFlatFeeNativePPM
      100_000 // fulfillmentFlatFeeLinkDiscountPPM
    );

    // Add wrapper as a consumer to the wrapper's subscription.
    s_testCoordinator.addConsumer(uint256(s_wrapperSubscriptionId), address(s_wrapper));
  }

  function _mockArbGasGetPricesInWei() internal {
    // return gas prices in wei, assuming the specified aggregator is used
    //        (
    //            per L2 tx,
    //            per L1 calldata unit, (zero byte = 4 units, nonzero byte = 16 units)
    //            per storage allocation,
    //            per ArbGas base,
    //            per ArbGas congestion,
    //            per ArbGas total
    //        )
    vm.mockCall(
      ARBGAS_ADDR,
      abi.encodeWithSelector(ARBGAS.getPricesInWei.selector),
      abi.encode(1 gwei, 250 gwei, 1 gwei, 1 gwei, 1 gwei, 1 gwei)
    );
  }

  function test_calculateRequestPriceNativeOnArbitrumWrapper() public {
    vm.txGasPrice(1 gwei);
    _mockArbGasGetPricesInWei();
    uint32 callbackGasLimit = 1_000_000;
    uint32 numWords = 5;
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPriceNative(callbackGasLimit, numWords);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(wrapperCostCalculation, 5.01483 * 1e17, 1e15);
  }

  function test_calculateRequestPriceLinkOnArbitrumWrapper() public {
    vm.txGasPrice(1 gwei);
    _mockArbGasGetPricesInWei();
    uint32 callbackGasLimit = 1_000_000;
    uint32 numWords = 5;
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPrice(callbackGasLimit, numWords);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(wrapperCostCalculation, 8.02846 * 1e17, 1e15);
  }
}
