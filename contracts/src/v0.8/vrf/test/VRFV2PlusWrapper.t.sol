// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {MockLinkToken} from "../../mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5.sol";
import {SubscriptionAPI} from "../dev/SubscriptionAPI.sol";
import {VRFV2PlusWrapperConsumerExample} from "../dev/testhelpers/VRFV2PlusWrapperConsumerExample.sol";
import {VRFCoordinatorV2_5} from "../dev/VRFCoordinatorV2_5.sol";
import {VRFConsumerBaseV2Plus} from "../dev/VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusWrapper} from "../dev/VRFV2PlusWrapper.sol";
import {VRFV2PlusClient} from "../dev/libraries/VRFV2PlusClient.sol";

contract VRFV2PlusWrapperTest is BaseTest {
  address internal constant LINK_WHALE = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
  bytes32 private vrfKeyHash = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";
  uint32 private wrapperGasOverhead = 10_000;
  uint32 private coordinatorGasOverhead = 20_000;
  uint256 private s_wrapperSubscriptionId;

  ExposedVRFCoordinatorV2_5 private s_testCoordinator;
  MockLinkToken private s_linkToken;
  MockV3Aggregator private s_linkNativeFeed;
  VRFV2PlusWrapper private s_wrapper;
  VRFV2PlusWrapperConsumerExample private s_consumer;

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(LINK_WHALE, 10_000 ether);
    vm.stopPrank();
    vm.startPrank(LINK_WHALE);

    // Deploy link token and link/native feed.
    s_linkToken = new MockLinkToken();
    s_linkNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Deploy coordinator.
    s_testCoordinator = new ExposedVRFCoordinatorV2_5(address(0));

    // Create subscription for all future wrapper contracts.
    s_wrapperSubscriptionId = s_testCoordinator.createSubscription();

    // Deploy wrapper.
    s_wrapper = new VRFV2PlusWrapper(
      address(s_linkToken),
      address(s_linkNativeFeed),
      address(s_testCoordinator),
      uint256(s_wrapperSubscriptionId)
    );
    assertEq(address(s_linkToken), address(s_wrapper.link()));
    assertEq(address(s_linkNativeFeed), address(s_wrapper.linkNativeFeed()));

    // Add wrapper as a consumer to the wrapper's subscription.
    s_testCoordinator.addConsumer(uint256(s_wrapperSubscriptionId), address(s_wrapper));

    // Deploy consumer.
    s_consumer = new VRFV2PlusWrapperConsumerExample(address(s_wrapper));

    // Configure the coordinator.
    s_testCoordinator.setLINKAndLINKNativeFeed(address(s_linkToken), address(s_linkNativeFeed));
    setConfigCoordinator();
    setConfigWrapper();

    s_testCoordinator.s_config();
  }

  function setConfigCoordinator() internal {
    s_testCoordinator.setConfig(
      0, // minRequestConfirmations
      2_500_000, // maxGasLimit
      1, // stalenessSeconds
      50_000, // gasAfterPaymentCalculation
      50000000000000000, // fallbackWeiPerUnitLink
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );
  }

  function setConfigWrapper() internal {
    vm.expectEmit(false, false, false, true, address(s_wrapper));
    emit ConfigSet(wrapperGasOverhead, coordinatorGasOverhead, 0, 0, vrfKeyHash, 10, 1, 50000000000000000, 0, 0);
    s_wrapper.setConfig(
      wrapperGasOverhead, // wrapper gas overhead
      coordinatorGasOverhead, // coordinator gas overhead
      0, // native premium percentage,
      0, // link premium percentage
      vrfKeyHash, // keyHash
      10, // max number of words,
      1, // stalenessSeconds
      50000000000000000, // fallbackWeiPerUnitLink
      0, // fulfillmentFlatFeeNativePPM
      0 // fulfillmentFlatFeeLinkDiscountPPM
    );
    (
      ,
      ,
      ,
      ,
      uint32 _wrapperGasOverhead,
      uint32 _coordinatorGasOverhead,
      uint8 _coordinatorNativePremiumPercentage,
      uint8 _coordinatorLinkPremiumPercentage,
      bytes32 _keyHash,
      uint8 _maxNumWords
    ) = s_wrapper.getConfig();
    assertEq(_wrapperGasOverhead, wrapperGasOverhead);
    assertEq(_coordinatorGasOverhead, coordinatorGasOverhead);
    assertEq(0, _coordinatorNativePremiumPercentage);
    assertEq(0, _coordinatorLinkPremiumPercentage);
    assertEq(vrfKeyHash, _keyHash);
    assertEq(10, _maxNumWords);
  }

  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint256 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    bytes extraArgs,
    address indexed sender
  );

  // IVRFV2PlusWrapper events
  event LinkNativeFeedSet(address linkNativeFeed);
  event FulfillmentTxSizeSet(uint32 size);
  event ConfigSet(
    uint32 wrapperGasOverhead,
    uint32 coordinatorGasOverhead,
    uint8 coordinatorNativePremiumPercentage,
    uint8 coordinatorLinkPremiumPercentage,
    bytes32 keyHash,
    uint8 maxNumWords,
    uint32 stalenessSeconds,
    int256 fallbackWeiPerUnitLink,
    uint32 fulfillmentFlatFeeNativePPM,
    uint32 fulfillmentFlatFeeLinkDiscountPPM
  );
  event FallbackWeiPerUnitLinkUsed(uint256 requestId, int256 fallbackWeiPerUnitLink);
  event Withdrawn(address indexed to, uint256 amount);
  event NativeWithdrawn(address indexed to, uint256 amount);
  event Enabled();
  event Disabled();

  // VRFV2PlusWrapperConsumerBase events
  event LinkTokenSet(address link);

  // SubscriptionAPI events
  event SubscriptionConsumerAdded(uint256 indexed subId, address consumer);

  function testVRFV2PlusWrapperZeroAddress() public {
    vm.expectRevert(VRFConsumerBaseV2Plus.ZeroAddress.selector);
    new VRFV2PlusWrapper(address(s_linkToken), address(s_linkNativeFeed), address(0), uint256(0));
  }

  function testCreationOfANewVRFV2PlusWrapper() public {
    // second wrapper contract will simply add itself to the same subscription
    VRFV2PlusWrapper nextWrapper = new VRFV2PlusWrapper(
      address(s_linkToken),
      address(s_linkNativeFeed),
      address(s_testCoordinator),
      s_wrapperSubscriptionId
    );
    assertEq(s_wrapperSubscriptionId, nextWrapper.SUBSCRIPTION_ID());
  }

  function testVRFV2PlusWrapperWithZeroSubscriptionId() public {
    vm.expectRevert(VRFV2PlusWrapper.SubscriptionIdMissing.selector);
    new VRFV2PlusWrapper(address(s_linkToken), address(s_linkNativeFeed), address(s_testCoordinator), uint256(0));
  }

  function testVRFV2PlusWrapperWithInvalidSubscriptionId() public {
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    new VRFV2PlusWrapper(address(s_linkToken), address(s_linkNativeFeed), address(s_testCoordinator), uint256(123456));
  }

  function testSetFulfillmentTxSize() public {
    uint32 fulfillmentTxSize = 100_000;
    vm.expectEmit(false, false, false, true, address(s_wrapper));
    emit FulfillmentTxSizeSet(fulfillmentTxSize);
    s_wrapper.setFulfillmentTxSize(fulfillmentTxSize);
    assertEq(s_wrapper.s_fulfillmentTxSizeBytes(), fulfillmentTxSize);
  }

  function testSetCoordinatorZeroAddress() public {
    vm.expectRevert(VRFConsumerBaseV2Plus.ZeroAddress.selector);
    s_wrapper.setCoordinator(address(0));
  }

  function testRequestAndFulfillRandomWordsNativeWrapper() public {
    // Fund subscription.
    s_testCoordinator.fundSubscriptionWithNative{value: 10 ether}(s_wrapper.SUBSCRIPTION_ID());
    vm.deal(address(s_consumer), 10 ether);

    // Get type and version.
    assertEq(s_wrapper.typeAndVersion(), "VRFV2PlusWrapper 1.0.0");

    // Cannot make request while disabled.
    vm.expectEmit(false, false, false, true, address(s_wrapper));
    emit Disabled();
    s_wrapper.disable();
    vm.expectRevert("wrapper is disabled");
    s_consumer.makeRequestNative(500_000, 0, 1);
    vm.expectEmit(false, false, false, true, address(s_wrapper));
    emit Enabled();
    s_wrapper.enable();

    // Request randomness from wrapper.
    uint32 callbackGasLimit = 1_000_000;
    vm.expectEmit(true, true, true, true);
    (uint256 requestId, uint256 preSeed) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(s_wrapper),
      s_wrapper.SUBSCRIPTION_ID(),
      1
    );
    uint32 EIP150Overhead = callbackGasLimit / 63 + 1;
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId,
      preSeed,
      s_wrapper.SUBSCRIPTION_ID(), // subId
      0, // minConfirmations
      callbackGasLimit + EIP150Overhead + wrapperGasOverhead, // callbackGasLimit - accounts for EIP 150
      1, // numWords
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true})), // extraArgs
      address(s_wrapper) // requester
    );
    requestId = s_consumer.makeRequestNative(callbackGasLimit, 0, 1);

    (uint256 paid, bool fulfilled, bool native) = s_consumer.s_requests(requestId);
    uint32 expectedPaid = callbackGasLimit + wrapperGasOverhead + coordinatorGasOverhead;
    uint256 wrapperNativeCostEstimate = s_wrapper.estimateRequestPriceNative(callbackGasLimit, tx.gasprice);
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPriceNative(callbackGasLimit);
    assertEq(paid, expectedPaid);
    assertEq(uint256(paid), wrapperNativeCostEstimate);
    assertEq(wrapperNativeCostEstimate, wrapperCostCalculation);
    assertEq(fulfilled, false);
    assertEq(native, true);
    assertEq(address(s_consumer).balance, 10 ether - expectedPaid);

    (, uint256 gasLimit, ) = s_wrapper.s_callbacks(requestId);
    assertEq(gasLimit, callbackGasLimit);

    changePrank(address(s_testCoordinator));
    uint256[] memory words = new uint256[](1);
    words[0] = 123;
    s_wrapper.rawFulfillRandomWords(requestId, words);
    (, bool nowFulfilled, uint256[] memory storedWords) = s_consumer.getRequestStatus(requestId);
    assertEq(nowFulfilled, true);
    assertEq(storedWords[0], 123);

    // Withdraw funds from wrapper.
    changePrank(LINK_WHALE);
    uint256 priorWhaleBalance = LINK_WHALE.balance;
    vm.expectEmit(true, false, false, true, address(s_wrapper));
    emit NativeWithdrawn(LINK_WHALE, paid);
    s_wrapper.withdrawNative(LINK_WHALE);
    assertEq(LINK_WHALE.balance, priorWhaleBalance + paid);
    assertEq(address(s_wrapper).balance, 0);
  }

  function testSetConfigFulfillmentFlatFee_LinkDiscountTooHigh() public {
    // Test that setting link discount flat fee higher than native flat fee reverts
    vm.expectRevert(abi.encodeWithSelector(VRFV2PlusWrapper.LinkDiscountTooHigh.selector, uint32(501), uint32(500)));
    s_wrapper.setConfig(
      wrapperGasOverhead, // wrapper gas overhead
      coordinatorGasOverhead, // coordinator gas overhead
      0, // native premium percentage,
      0, // link premium percentage
      vrfKeyHash, // keyHash
      10, // max number of words,
      1, // stalenessSeconds
      50000000000000000, // fallbackWeiPerUnitLink
      500, // fulfillmentFlatFeeNativePPM
      501 // fulfillmentFlatFeeLinkDiscountPPM
    );
  }

  function testSetConfigFulfillmentFlatFee_LinkDiscountEqualsNative() public {
    // Test that setting link discount flat fee equal to native flat fee does not revert
    s_wrapper.setConfig(
      wrapperGasOverhead, // wrapper gas overhead
      coordinatorGasOverhead, // coordinator gas overhead
      0, // native premium percentage,
      0, // link premium percentage
      vrfKeyHash, // keyHash
      10, // max number of words,
      1, // stalenessSeconds
      50000000000000000, // fallbackWeiPerUnitLink
      450, // fulfillmentFlatFeeNativePPM
      450 // fulfillmentFlatFeeLinkDiscountPPM
    );
  }

  function testSetConfigNativePremiumPercentageInvalidPremiumPercentage() public {
    // Test that setting native premium percentage higher than 155 will revert
    vm.expectRevert(
      abi.encodeWithSelector(VRFCoordinatorV2_5.InvalidPremiumPercentage.selector, uint8(156), uint8(155))
    );
    s_wrapper.setConfig(
      wrapperGasOverhead, // wrapper gas overhead
      coordinatorGasOverhead, // coordinator gas overhead
      156, // native premium percentage,
      0, // link premium percentage
      vrfKeyHash, // keyHash
      10, // max number of words,
      1, // stalenessSeconds
      50000000000000000, // fallbackWeiPerUnitLink
      0, // fulfillmentFlatFeeNativePPM
      0 // fulfillmentFlatFeeLinkDiscountPPM
    );
  }

  function testSetConfigLinkPremiumPercentageInvalidPremiumPercentage() public {
    // Test that setting LINK premium percentage higher than 155 will revert
    vm.expectRevert(
      abi.encodeWithSelector(VRFCoordinatorV2_5.InvalidPremiumPercentage.selector, uint8(202), uint8(155))
    );
    s_wrapper.setConfig(
      wrapperGasOverhead, // wrapper gas overhead
      coordinatorGasOverhead, // coordinator gas overhead
      15, // native premium percentage,
      202, // link premium percentage
      vrfKeyHash, // keyHash
      10, // max number of words,
      1, // stalenessSeconds
      50000000000000000, // fallbackWeiPerUnitLink
      0, // fulfillmentFlatFeeNativePPM
      0 // fulfillmentFlatFeeLinkDiscountPPM
    );
  }

  function testRequestAndFulfillRandomWordsLINKWrapper() public {
    // Fund subscription.
    s_linkToken.transferAndCall(address(s_testCoordinator), 10 ether, abi.encode(s_wrapper.SUBSCRIPTION_ID()));
    s_linkToken.transfer(address(s_consumer), 10 ether);

    // Request randomness from wrapper.
    uint32 callbackGasLimit = 1_000_000;
    vm.expectEmit(true, true, true, true);
    (uint256 requestId, uint256 preSeed) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(s_wrapper),
      s_wrapper.SUBSCRIPTION_ID(),
      1
    );
    uint32 EIP150Overhead = callbackGasLimit / 63 + 1;
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId,
      preSeed,
      s_wrapper.SUBSCRIPTION_ID(), // subId
      0, // minConfirmations
      callbackGasLimit + EIP150Overhead + wrapperGasOverhead, // callbackGasLimit - accounts for EIP 150
      1, // numWords
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false})), // extraArgs
      address(s_wrapper) // requester
    );
    s_consumer.makeRequest(callbackGasLimit, 0, 1);

    // Assert that the request was made correctly.
    (uint256 paid, bool fulfilled, bool native) = s_consumer.s_requests(requestId);
    uint32 expectedPaid = (callbackGasLimit + wrapperGasOverhead + coordinatorGasOverhead) * 2;
    uint256 wrapperCostEstimate = s_wrapper.estimateRequestPrice(callbackGasLimit, tx.gasprice);
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPrice(callbackGasLimit);
    assertEq(paid, expectedPaid); // 1_030_000 * 2 for link/native ratio
    assertEq(uint256(paid), wrapperCostEstimate);
    assertEq(wrapperCostEstimate, wrapperCostCalculation);
    assertEq(fulfilled, false);
    assertEq(native, false);
    assertEq(s_linkToken.balanceOf(address(s_consumer)), 10 ether - expectedPaid);
    (, uint256 gasLimit, ) = s_wrapper.s_callbacks(requestId);
    assertEq(gasLimit, callbackGasLimit);

    // Fulfill the request.
    changePrank(address(s_testCoordinator));
    uint256[] memory words = new uint256[](1);
    words[0] = 456;
    s_wrapper.rawFulfillRandomWords(requestId, words);
    (, bool nowFulfilled, uint256[] memory storedWords) = s_consumer.getRequestStatus(requestId);
    assertEq(nowFulfilled, true);
    assertEq(storedWords[0], 456);

    // Withdraw funds from wrapper.
    changePrank(LINK_WHALE);
    uint256 priorWhaleBalance = s_linkToken.balanceOf(LINK_WHALE);
    vm.expectEmit(true, false, false, true, address(s_wrapper));
    emit Withdrawn(LINK_WHALE, paid);
    s_wrapper.withdraw(LINK_WHALE);
    assertEq(s_linkToken.balanceOf(LINK_WHALE), priorWhaleBalance + paid);
    assertEq(s_linkToken.balanceOf(address(s_wrapper)), 0);
  }

  function testRequestRandomWordsLINKWrapperFallbackWeiPerUnitLinkUsed() public {
    // Fund subscription.
    s_linkToken.transferAndCall(address(s_testCoordinator), 10 ether, abi.encode(s_wrapper.SUBSCRIPTION_ID()));
    s_linkToken.transfer(address(s_consumer), 10 ether);

    // Set the link feed to be stale.
    (, , , uint32 stalenessSeconds, , , , , ) = s_testCoordinator.s_config();
    int256 fallbackWeiPerUnitLink = s_testCoordinator.s_fallbackWeiPerUnitLink();
    (uint80 roundId, int256 answer, uint256 startedAt, , ) = s_linkNativeFeed.latestRoundData();
    uint256 timestamp = block.timestamp - stalenessSeconds - 1;
    s_linkNativeFeed.updateRoundData(roundId, answer, timestamp, startedAt);

    // Request randomness from wrapper.
    uint32 callbackGasLimit = 1_000_000;
    (uint256 requestId, uint256 preSeed) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(s_wrapper),
      s_wrapper.SUBSCRIPTION_ID(),
      1
    );
    uint32 EIP150Overhead = callbackGasLimit / 63 + 1;
    vm.expectEmit(true, true, true, true);
    emit FallbackWeiPerUnitLinkUsed(requestId, fallbackWeiPerUnitLink);
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId,
      preSeed,
      s_wrapper.SUBSCRIPTION_ID(), // subId
      0, // minConfirmations
      callbackGasLimit + EIP150Overhead + wrapperGasOverhead, // callbackGasLimit - accounts for EIP 150
      1, // numWords
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false})), // extraArgs
      address(s_wrapper) // requester
    );
    s_consumer.makeRequest(callbackGasLimit, 0, 1);
  }

  function testRequestRandomWordsInNativeNotConfigured() public {
    VRFV2PlusWrapper wrapper = new VRFV2PlusWrapper(
      address(s_linkToken),
      address(s_linkNativeFeed),
      address(s_testCoordinator),
      uint256(s_wrapperSubscriptionId)
    );

    vm.expectRevert("wrapper is not configured");
    wrapper.requestRandomWordsInNative(500_000, 0, 1, "");
  }

  function testRequestRandomWordsInNativeDisabled() public {
    s_wrapper.disable();

    vm.expectRevert("wrapper is disabled");
    s_wrapper.requestRandomWordsInNative(500_000, 0, 1, "");
  }
}
