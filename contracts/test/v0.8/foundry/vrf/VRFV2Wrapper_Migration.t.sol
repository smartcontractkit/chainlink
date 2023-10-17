pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5} from "../../../../src/v0.8/vrf/dev/testhelpers/ExposedVRFCoordinatorV2_5.sol";
import {VRFCoordinatorV2Plus_V2Example} from "../../../../src/v0.8/vrf/dev/testhelpers/VRFCoordinatorV2Plus_V2Example.sol";
import {VRFV2PlusWrapperConsumerBase} from "../../../../src/v0.8/vrf/dev/VRFV2PlusWrapperConsumerBase.sol";
import {VRFV2PlusWrapperConsumerExample} from "../../../../src/v0.8/vrf/dev/testhelpers/VRFV2PlusWrapperConsumerExample.sol";
import {SubscriptionAPI} from "../../../../src/v0.8/vrf/dev/SubscriptionAPI.sol";
import {VRFCoordinatorV2_5} from "../../../../src/v0.8/vrf/dev/VRFCoordinatorV2_5.sol";
import {VRFV2PlusWrapper} from "../../../../src/v0.8/vrf/dev/VRFV2PlusWrapper.sol";
import {VRFV2PlusClient} from "../../../../src/v0.8/vrf/dev/libraries/VRFV2PlusClient.sol";

contract VRFV2PlusWrapperTest is BaseTest {
  address internal constant LINK_WHALE = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
  uint256 internal constant DEFAULT_NATIVE_FUNDING = 7 ether; // 7 ETH
  uint256 internal constant DEFAULT_LINK_FUNDING = 10 ether; // 10 ETH
  bytes32 vrfKeyHash = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";
  uint32 wrapperGasOverhead = 10_000;
  uint32 coordinatorGasOverhead = 20_000;

  ExposedVRFCoordinatorV2_5 s_testCoordinator;
  MockLinkToken s_linkToken;
  MockV3Aggregator s_linkNativeFeed;
  VRFV2PlusWrapper s_wrapper;
  VRFV2PlusWrapperConsumerExample s_consumer;

  VRFCoordinatorV2Plus_V2Example s_newCoordinator;

  VRFCoordinatorV2_5.FeeConfig basicFeeConfig =
    VRFCoordinatorV2_5.FeeConfig({fulfillmentFlatFeeLinkPPM: 0, fulfillmentFlatFeeNativePPM: 0});

  event CoordinatorRegistered(address coordinatorAddress);
  event MigrationCompleted(address newCoordinator, uint256 subId);
  event WrapperRequestMade(uint256 indexed requestId, uint256 paid);

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(LINK_WHALE, 10_000 ether);
    changePrank(LINK_WHALE);

    // Deploy link token and link/native feed.
    s_linkToken = new MockLinkToken();
    s_linkNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Deploy coordinator and consumer.
    s_testCoordinator = new ExposedVRFCoordinatorV2_5(address(0));
    s_wrapper = new VRFV2PlusWrapper(address(s_linkToken), address(s_linkNativeFeed), address(s_testCoordinator));
    s_consumer = new VRFV2PlusWrapperConsumerExample(address(s_linkToken), address(s_wrapper));

    // Configure the coordinator.
    s_testCoordinator.setLINKAndLINKNativeFeed(address(s_linkToken), address(s_linkNativeFeed));
    setConfigCoordinator(basicFeeConfig);
    setConfigWrapper();

    s_testCoordinator.s_config();

    // Data structures for Migrateable Wrapper
    s_newCoordinator = new VRFCoordinatorV2Plus_V2Example(address(0), address(s_testCoordinator));
    vm.expectEmit(
      false, // no first indexed topic
      false, // no second indexed topic
      false, // no third indexed topic
      true // check data (target coordinator address)
    );
    address newCoordinatorAddr = address(s_newCoordinator);
    emit CoordinatorRegistered(newCoordinatorAddr);
    s_testCoordinator.registerMigratableCoordinator(newCoordinatorAddr);
    assertTrue(s_testCoordinator.isTargetRegisteredExternal(newCoordinatorAddr));
  }

  function setConfigCoordinator(VRFCoordinatorV2_5.FeeConfig memory feeConfig) internal {
    s_testCoordinator.setConfig(
      0, // minRequestConfirmations
      2_500_000, // maxGasLimit
      1, // stalenessSeconds
      50_000, // gasAfterPaymentCalculation
      50000000000000000, // fallbackWeiPerUnitLink
      feeConfig
    );
  }

  function setConfigWrapper() internal {
    s_wrapper.setConfig(
      wrapperGasOverhead, // wrapper gas overhead
      coordinatorGasOverhead, // coordinator gas overhead
      0, // premium percentage
      vrfKeyHash, // keyHash
      10, // max number of words,
      1, // stalenessSeconds
      50000000000000000, // fallbackWeiPerUnitLink
      0, // fulfillmentFlatFeeLinkPPM
      0 // fulfillmentFlatFeeNativePPM
    );
    (
      ,
      ,
      ,
      ,
      uint32 _wrapperGasOverhead,
      uint32 _coordinatorGasOverhead,
      uint8 _wrapperPremiumPercentage,
      bytes32 _keyHash,
      uint8 _maxNumWords
    ) = s_wrapper.getConfig();
    assertEq(_wrapperGasOverhead, wrapperGasOverhead);
    assertEq(_coordinatorGasOverhead, coordinatorGasOverhead);
    assertEq(0, _wrapperPremiumPercentage);
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

  function testMigrateWrapperLINKPayment() public {
    s_linkToken.transfer(address(s_consumer), DEFAULT_LINK_FUNDING);

    uint256 subID = s_wrapper.SUBSCRIPTION_ID();
    address oldCoordinatorAddr = address(s_testCoordinator);

    // Fund subscription with native and LINK payment to check
    // if funds are transferred to new subscription after call
    // migration to new coordinator
    s_linkToken.transferAndCall(oldCoordinatorAddr, DEFAULT_LINK_FUNDING, abi.encode(subID));
    s_testCoordinator.fundSubscriptionWithNative{value: DEFAULT_NATIVE_FUNDING}(subID);

    // Get type and version.
    assertEq(s_wrapper.typeAndVersion(), "VRFV2Wrapper 1.0.0");

    // subscription exists in V1 coordinator before migration

    (
      uint96 balance,
      uint96 nativeBalance,
      uint64 reqCount,
      address owner,
      address[] memory consumers
    ) = s_testCoordinator.getSubscription(subID);
    assertEq(reqCount, 0);
    assertEq(balance, DEFAULT_LINK_FUNDING);
    assertEq(nativeBalance, DEFAULT_NATIVE_FUNDING);
    assertEq(owner, address(s_wrapper));
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_wrapper));

    vm.startPrank(LINK_WHALE);

    // Update wrapper to point to the new coordinator
    vm.expectEmit(
      false, // no first indexed field
      false, // no second indexed field
      false, // no third indexed field
      true // check data fields
    );
    address newCoordinatorAddr = address(s_newCoordinator);
    emit MigrationCompleted(newCoordinatorAddr, subID);

    s_wrapper.migrate(newCoordinatorAddr);

    // subscription no longer exists in v1 coordinator after migration
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_testCoordinator.getSubscription(subID);
    assertEq(s_testCoordinator.s_totalBalance(), 0);
    assertEq(s_testCoordinator.s_totalNativeBalance(), 0);
    assertEq(s_linkToken.balanceOf(oldCoordinatorAddr), 0);
    assertEq(oldCoordinatorAddr.balance, 0);

    // subscription exists in v2 coordinator
    (balance, nativeBalance, reqCount, owner, consumers) = s_newCoordinator.getSubscription(subID);
    assertEq(owner, address(s_wrapper));
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_wrapper));
    assertEq(reqCount, 0);
    assertEq(balance, DEFAULT_LINK_FUNDING);
    assertEq(nativeBalance, DEFAULT_NATIVE_FUNDING);
    assertEq(s_newCoordinator.s_totalLinkBalance(), DEFAULT_LINK_FUNDING);
    assertEq(s_newCoordinator.s_totalNativeBalance(), DEFAULT_NATIVE_FUNDING);
    assertEq(s_linkToken.balanceOf(newCoordinatorAddr), DEFAULT_LINK_FUNDING);
    assertEq(newCoordinatorAddr.balance, DEFAULT_NATIVE_FUNDING);

    // calling migrate again on V1 coordinator should fail
    vm.expectRevert();
    s_wrapper.migrate(newCoordinatorAddr);

    // Request randomness from wrapper.
    uint32 callbackGasLimit = 1_000_000;
    vm.expectEmit(true, true, true, true);
    uint256 wrapperCost = s_wrapper.calculateRequestPrice(callbackGasLimit);
    emit WrapperRequestMade(1, wrapperCost);
    uint256 requestId = s_consumer.makeRequest(callbackGasLimit, 0, 1);
    assertEq(requestId, 1);

    (uint256 paid, bool fulfilled, bool native) = s_consumer.s_requests(requestId);
    uint32 expectedPaid = (callbackGasLimit + wrapperGasOverhead + coordinatorGasOverhead) * 2;
    uint256 wrapperCostEstimate = s_wrapper.estimateRequestPrice(callbackGasLimit, tx.gasprice);
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPrice(callbackGasLimit);
    assertEq(paid, expectedPaid); // 1_030_000 * 2 for link/native ratio
    assertEq(uint256(paid), wrapperCostEstimate);
    assertEq(wrapperCostEstimate, wrapperCostCalculation);
    assertEq(fulfilled, false);
    assertEq(native, false);
    assertEq(s_linkToken.balanceOf(address(s_consumer)), DEFAULT_LINK_FUNDING - expectedPaid);

    (, uint256 gasLimit, ) = s_wrapper.s_callbacks(requestId);
    assertEq(gasLimit, callbackGasLimit);

    vm.stopPrank();

    vm.startPrank(newCoordinatorAddr);

    uint256[] memory words = new uint256[](1);
    words[0] = 123;
    s_wrapper.rawFulfillRandomWords(requestId, words);
    (, bool nowFulfilled, uint256[] memory storedWords) = s_consumer.getRequestStatus(requestId);
    assertEq(nowFulfilled, true);
    assertEq(storedWords[0], 123);

    vm.stopPrank();

    /// Withdraw funds from wrapper.
    vm.startPrank(LINK_WHALE);
    uint256 priorWhaleBalance = s_linkToken.balanceOf(LINK_WHALE);
    s_wrapper.withdraw(LINK_WHALE, paid);
    assertEq(s_linkToken.balanceOf(LINK_WHALE), priorWhaleBalance + paid);
    assertEq(s_linkToken.balanceOf(address(s_wrapper)), 0);

    vm.stopPrank();
  }

  function testMigrateWrapperNativePayment() public {
    vm.deal(address(s_consumer), DEFAULT_NATIVE_FUNDING);

    uint256 subID = s_wrapper.SUBSCRIPTION_ID();
    address oldCoordinatorAddr = address(s_testCoordinator);

    // Fund subscription with native and LINK payment to check
    // if funds are transferred to new subscription after call
    // migration to new coordinator
    s_linkToken.transferAndCall(oldCoordinatorAddr, DEFAULT_LINK_FUNDING, abi.encode(subID));
    s_testCoordinator.fundSubscriptionWithNative{value: DEFAULT_NATIVE_FUNDING}(subID);

    // Get type and version.
    assertEq(s_wrapper.typeAndVersion(), "VRFV2Wrapper 1.0.0");

    // subscription exists in V1 coordinator before migration
    (
      uint96 balance,
      uint96 nativeBalance,
      uint64 reqCount,
      address owner,
      address[] memory consumers
    ) = s_testCoordinator.getSubscription(subID);
    assertEq(reqCount, 0);
    assertEq(balance, DEFAULT_LINK_FUNDING);
    assertEq(nativeBalance, DEFAULT_NATIVE_FUNDING);
    assertEq(owner, address(s_wrapper));
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_wrapper));

    vm.startPrank(LINK_WHALE);

    // Update wrapper to point to the new coordinator
    vm.expectEmit(
      false, // no first indexed field
      false, // no second indexed field
      false, // no third indexed field
      true // check data fields
    );
    address newCoordinatorAddr = address(s_newCoordinator);
    emit MigrationCompleted(newCoordinatorAddr, subID);

    s_wrapper.migrate(newCoordinatorAddr);

    // subscription no longer exists in v1 coordinator after migration
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_testCoordinator.getSubscription(subID);
    assertEq(s_testCoordinator.s_totalBalance(), 0);
    assertEq(s_testCoordinator.s_totalNativeBalance(), 0);
    assertEq(s_linkToken.balanceOf(oldCoordinatorAddr), 0);
    assertEq(oldCoordinatorAddr.balance, 0);

    // subscription exists in v2 coordinator
    (balance, nativeBalance, reqCount, owner, consumers) = s_newCoordinator.getSubscription(subID);
    assertEq(owner, address(s_wrapper));
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_wrapper));
    assertEq(reqCount, 0);
    assertEq(balance, DEFAULT_LINK_FUNDING);
    assertEq(nativeBalance, DEFAULT_NATIVE_FUNDING);
    assertEq(s_newCoordinator.s_totalLinkBalance(), DEFAULT_LINK_FUNDING);
    assertEq(s_newCoordinator.s_totalNativeBalance(), DEFAULT_NATIVE_FUNDING);
    assertEq(s_linkToken.balanceOf(newCoordinatorAddr), DEFAULT_LINK_FUNDING);
    assertEq(newCoordinatorAddr.balance, DEFAULT_NATIVE_FUNDING);

    // calling migrate again on V1 coordinator should fail
    vm.expectRevert();
    s_wrapper.migrate(newCoordinatorAddr);

    // Request randomness from wrapper.
    uint32 callbackGasLimit = 1_000_000;
    vm.expectEmit(true, true, true, true);
    uint256 wrapperCost = s_wrapper.calculateRequestPriceNative(callbackGasLimit);
    emit WrapperRequestMade(1, wrapperCost);
    uint256 requestId = s_consumer.makeRequestNative(callbackGasLimit, 0, 1);
    assertEq(requestId, 1);

    (uint256 paid, bool fulfilled, bool native) = s_consumer.s_requests(requestId);
    uint32 expectedPaid = callbackGasLimit + wrapperGasOverhead + coordinatorGasOverhead;
    uint256 wrapperNativeCostEstimate = s_wrapper.estimateRequestPriceNative(callbackGasLimit, tx.gasprice);
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPriceNative(callbackGasLimit);
    assertEq(paid, expectedPaid);
    assertEq(uint256(paid), wrapperNativeCostEstimate);
    assertEq(wrapperNativeCostEstimate, wrapperCostCalculation);
    assertEq(fulfilled, false);
    assertEq(native, true);
    assertEq(address(s_consumer).balance, DEFAULT_NATIVE_FUNDING - expectedPaid);

    (, uint256 gasLimit, ) = s_wrapper.s_callbacks(requestId);
    assertEq(gasLimit, callbackGasLimit);

    vm.stopPrank();

    vm.startPrank(newCoordinatorAddr);

    uint256[] memory words = new uint256[](1);
    words[0] = 123;
    s_wrapper.rawFulfillRandomWords(requestId, words);
    (, bool nowFulfilled, uint256[] memory storedWords) = s_consumer.getRequestStatus(requestId);
    assertEq(nowFulfilled, true);
    assertEq(storedWords[0], 123);

    vm.stopPrank();

    // Withdraw funds from wrapper.
    vm.startPrank(LINK_WHALE);
    uint256 priorWhaleBalance = LINK_WHALE.balance;
    s_wrapper.withdrawNative(LINK_WHALE, paid);
    assertEq(LINK_WHALE.balance, priorWhaleBalance + paid);
    assertEq(address(s_wrapper).balance, 0);

    vm.stopPrank();
  }
}
