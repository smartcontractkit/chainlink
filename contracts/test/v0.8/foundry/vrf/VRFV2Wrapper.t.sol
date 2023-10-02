pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5} from "../../../../src/v0.8/dev/vrf/testhelpers/ExposedVRFCoordinatorV2_5.sol";
import {VRFV2PlusWrapperConsumerBase} from "../../../../src/v0.8/dev/vrf/VRFV2PlusWrapperConsumerBase.sol";
import {VRFV2PlusWrapperConsumerExample} from "../../../../src/v0.8/dev/vrf/testhelpers/VRFV2PlusWrapperConsumerExample.sol";
import {VRFCoordinatorV2_5} from "../../../../src/v0.8/dev/vrf/VRFCoordinatorV2_5.sol";
import {VRFV2PlusWrapper} from "../../../../src/v0.8/dev/vrf/VRFV2PlusWrapper.sol";
import {VRFV2PlusClient} from "../../../../src/v0.8/dev/vrf/libraries/VRFV2PlusClient.sol";
import {console} from "forge-std/console.sol";

contract VRFV2PlusWrapperTest is BaseTest {
  address internal constant LINK_WHALE = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
  bytes32 vrfKeyHash = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";
  uint32 wrapperGasOverhead = 10_000;
  uint32 coordinatorGasOverhead = 20_000;

  ExposedVRFCoordinatorV2_5 s_testCoordinator;
  MockLinkToken s_linkToken;
  MockV3Aggregator s_linkNativeFeed;
  VRFV2PlusWrapper s_wrapper;
  VRFV2PlusWrapperConsumerExample s_consumer;

  VRFCoordinatorV2_5.FeeConfig basicFeeConfig =
    VRFCoordinatorV2_5.FeeConfig({fulfillmentFlatFeeLinkPPM: 0, fulfillmentFlatFeeNativePPM: 0});

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

  function testSetLinkAndLinkNativeFeed() public {
    VRFV2PlusWrapper wrapper = new VRFV2PlusWrapper(address(0), address(0), address(s_testCoordinator));

    // Set LINK and LINK/Native feed on wrapper.
    wrapper.setLINK(address(s_linkToken));
    wrapper.setLinkNativeFeed(address(s_linkNativeFeed));
    assertEq(address(wrapper.s_link()), address(s_linkToken));
    assertEq(address(wrapper.s_linkNativeFeed()), address(s_linkNativeFeed));

    // Revert for subsequent assignment.
    vm.expectRevert(VRFV2PlusWrapper.LinkAlreadySet.selector);
    wrapper.setLINK(address(s_linkToken));

    // Consumer can set LINK token.
    VRFV2PlusWrapperConsumerExample consumer = new VRFV2PlusWrapperConsumerExample(address(0), address(wrapper));
    consumer.setLinkToken(address(s_linkToken));

    // Revert for subsequent assignment.
    vm.expectRevert(VRFV2PlusWrapperConsumerBase.LINKAlreadySet.selector);
    consumer.setLinkToken(address(s_linkToken));
  }

  function testRequestAndFulfillRandomWordsNativeWrapper() public {
    // Fund subscription.
    s_testCoordinator.fundSubscriptionWithNative{value: 10 ether}(s_wrapper.SUBSCRIPTION_ID());
    vm.deal(address(s_consumer), 10 ether);

    // Get type and version.
    assertEq(s_wrapper.typeAndVersion(), "VRFV2Wrapper 1.0.0");

    // Cannot make request while disabled.
    s_wrapper.disable();
    vm.expectRevert("wrapper is disabled");
    s_consumer.makeRequestNative(500_000, 0, 1);
    s_wrapper.enable();

    // Request randomness from wrapper.
    uint32 callbackGasLimit = 1_000_000;
    vm.expectEmit(true, true, true, true);
    (uint256 requestId, uint256 preSeed) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(s_wrapper),
      s_wrapper.SUBSCRIPTION_ID(),
      2
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
    s_wrapper.withdrawNative(LINK_WHALE, paid);
    assertEq(LINK_WHALE.balance, priorWhaleBalance + paid);
    assertEq(address(s_wrapper).balance, 0);
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
      2
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
    s_wrapper.withdraw(LINK_WHALE, paid);
    assertEq(s_linkToken.balanceOf(LINK_WHALE), priorWhaleBalance + paid);
    assertEq(s_linkToken.balanceOf(address(s_wrapper)), 0);
  }
}
