// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {BaseTest} from "../BaseTest.t.sol";
import {ExposedLotteryConsumer} from "./ExposedLotteryConsumer.sol";
import {VRFCoordinatorV2} from "../../../../src/v0.8/VRFCoordinatorV2.sol";
import {MockLinkToken as LinkToken} from "./special/MockLinkToken.sol";

contract LotteryConsumerSetup is BaseTest {
  VRFCoordinatorV2 s_coordinator;

  ExposedLotteryConsumer s_lotteryConsumer;

  LinkToken LINK;

  address internal LINK_ETH_FEED_ADDRESS = 0x2B4Ef6FAC5F6BE758508605bd87356776BfbfbA0;
  address internal BLOCKHASH_STORE_ADDRESS = 0x032477835Bc60C777321e1e09F1Fe6cFB82e111C;

  function setUp() public virtual override {
    BaseTest.setUp();

    LINK = new LinkToken();
    s_coordinator = new VRFCoordinatorV2(address(LINK), LINK_ETH_FEED_ADDRESS, BLOCKHASH_STORE_ADDRESS);
    s_lotteryConsumer = new ExposedLotteryConsumer(address(s_coordinator));

    VRFCoordinatorV2.FeeConfig memory feeConfig;
    s_coordinator.setConfig(1, 2500000, 8600, 32000, 1000, feeConfig);
  }
}

contract LotteryConsumerTest is LotteryConsumerSetup {
  address internal subscriptionOwner = 0x054684DA2dCcf0e79E72DcAb1b5db42E71b66c8A;
  bytes32 internal keyHash = 0xc90177494b1fe7ed1baf4da57091859c1de91d286d73491fc3ed683eec0abecd;

  address internal allowedCaller = 0xb87Fe82988ce1D3324050DD3eafBf68A9c35d623;

  event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amount);
  event SubscriptionCreated(uint64 indexed subId, address owner);
  event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionConsumerAdded(uint64 indexed subId, address consumer);

  event AllowedCallerAdded(address caller);
  event AllowedCallerRemoved(address caller);

  function testRequestRandomness() public {
    // switch caller to subscriptionOwner so that the subscription owner and the
    // contract owner are different.
    uint64 expectedSubId = 1;
    vm.stopPrank();
    vm.startPrank(subscriptionOwner);
    vm.expectEmit(true /* subId */, false /* checkTopic2 */, false /* checkTopic3 */, true /* owner */);
    emit SubscriptionCreated(expectedSubId, subscriptionOwner);
    s_coordinator.createSubscription();

    // add lottery consumer to subscription
    vm.expectEmit(true /* subId */, false /* checkTopic2 */, false /* checkTopic3 */, true /* consumer */);
    emit SubscriptionConsumerAdded(expectedSubId, address(s_lotteryConsumer));
    s_coordinator.addConsumer(expectedSubId, address(s_lotteryConsumer));

    vm.stopPrank();
    vm.startPrank(OWNER); // switch to owner of lottery consumer
    uint16 expectedMinConfs = 1;
    uint32 expectedCbGasLimit = 200000;
    uint32 expectedNumWords = 3;
    s_lotteryConsumer.setRequestConfig(
      keyHash,
      expectedSubId,
      expectedMinConfs,
      expectedCbGasLimit,
      expectedNumWords
    );
    (bytes32 actualKeyHash,
    uint64 actualSubscriptionId,
    uint16 actualMinRequestConfirmations,
    uint32 actualCallbackGasLimit,
    uint32 actualNumWords) = s_lotteryConsumer.getRequestConfig();
    assertEq(keyHash /* expected keyhash */, actualKeyHash);
    assertEq(expectedSubId /* expected subscription id */, actualSubscriptionId);
    assertEq(expectedMinConfs /* expected min confs */, actualMinRequestConfirmations);
    assertEq(expectedCbGasLimit /* expected callback gas limit */, actualCallbackGasLimit);
    assertEq(expectedNumWords /* expected num words */, expectedNumWords);

    // Set allowed callers on lottery consumer
    vm.expectEmit(false /* checkTopic1 */, false /* checkTopic2 */, false /* checkTopic3 */, true /* caller */);
    emit AllowedCallerAdded(allowedCaller);
    s_lotteryConsumer.addAllowedCaller(allowedCaller);

    // Switch to allowed caller and request randomness
    vm.stopPrank();
    vm.startPrank(allowedCaller);
    bytes32 clientRequestId = 0xc90177494b1fe7ed1baf4d222222222c1de91d286d73491fc3ed683eec0abecd;
    uint8 expectedLotteryType = 1;
    uint128 vrfExternalRequestId = 50;
    s_lotteryConsumer.requestRandomness(clientRequestId, expectedLotteryType, vrfExternalRequestId);
    uint256 vrfRequestId = s_lotteryConsumer.getMostRecentVrfRequestId();

    // switch to coordinator to fulfill
    vm.stopPrank();
    vm.startPrank(address(s_coordinator));
    uint256[] memory randomWords = new uint256[](3);
    randomWords[0] = 1;
    randomWords[1] = 2;
    randomWords[2] = 3;
    s_lotteryConsumer.fulfillRandomWordsExternal(vrfRequestId, randomWords);
  }
}
