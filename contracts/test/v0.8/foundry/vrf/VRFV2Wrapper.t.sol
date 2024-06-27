pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/testhelpers/ExposedVRFCoordinatorV2Plus.sol";
import {VRFV2PlusWrapperConsumerExample} from "../../../../src/v0.8/dev/vrf/testhelpers/VRFV2PlusWrapperConsumerExample.sol";
import {VRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/VRFCoordinatorV2Plus.sol";
import {VRFV2PlusWrapper} from "../../../../src/v0.8/dev/vrf/VRFV2PlusWrapper.sol";
import {VRFV2PlusClient} from "../../../../src/v0.8/dev/vrf/libraries/VRFV2PlusClient.sol";
import {console} from "forge-std/console.sol";

contract VRFV2PlusWrapperTest is BaseTest {
  address internal constant LINK_WHALE = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
  bytes32 vrfKeyHash = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";
  uint32 wrapperGasOverhead = 10_000;
  uint32 coordinatorGasOverhead = 20_000;

  ExposedVRFCoordinatorV2Plus s_testCoordinator;
  MockLinkToken s_linkToken;
  MockV3Aggregator s_linkEthFeed;
  VRFV2PlusWrapper s_wrapper;
  VRFV2PlusWrapperConsumerExample s_consumer;

  VRFCoordinatorV2Plus.FeeConfig basicFeeConfig =
    VRFCoordinatorV2Plus.FeeConfig({fulfillmentFlatFeeLinkPPM: 0, fulfillmentFlatFeeEthPPM: 0});

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(LINK_WHALE, 10_000 ether);
    changePrank(LINK_WHALE);

    // Deploy link token and link/eth feed.
    s_linkToken = new MockLinkToken();
    s_linkEthFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Deploy coordinator and consumer.
    s_testCoordinator = new ExposedVRFCoordinatorV2Plus(address(0));
    s_wrapper = new VRFV2PlusWrapper(address(s_linkToken), address(s_linkEthFeed), address(s_testCoordinator));
    s_consumer = new VRFV2PlusWrapperConsumerExample(address(s_linkToken), address(s_wrapper));

    // Configure the coordinator.
    s_testCoordinator.setLINKAndLINKETHFeed(address(s_linkToken), address(s_linkEthFeed));
    setConfigCoordinator(basicFeeConfig);
    setConfigWrapper();

    s_testCoordinator.s_config();
  }

  function setConfigCoordinator(VRFCoordinatorV2Plus.FeeConfig memory feeConfig) internal {
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
      10 // max number of words
    );
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

  function testRequestAndFulfillRandomWordsNativeWrapper() public {
    // Fund subscription.
    s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(s_wrapper.SUBSCRIPTION_ID());
    vm.deal(address(s_consumer), 10 ether);

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
    assertEq(paid, expectedPaid);
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
    assertEq(paid, expectedPaid); // 1_030_000 * 2 for link/eth ratio
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
  }
}
