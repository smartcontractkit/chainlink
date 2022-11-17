pragma solidity ^0.8.0;

import {BaseTest} from "../../BaseTest.t.sol";
import {NoCancelVRFCoordinatorV2} from "../../../../../src/v0.8/dev/special/NoCancelVRFCoordinatorV2.sol";
import {MockLinkToken as LinkToken} from "./MockLinkToken.sol";
import {ExposedNoCancelVRFCoordinatorV2} from "./ExposedNoCancelVRFCoordinatorV2.sol";

contract NoCancelVRFCoordinatorV2Setup is BaseTest {
  NoCancelVRFCoordinatorV2 s_coordinator;

  // has calculatePaymentAmount "exported" to test.
  ExposedNoCancelVRFCoordinatorV2 s_exposedCoordinator;

  LinkToken LINK;

  address internal LINK_ETH_FEED_ADDRESS = 0x2B4Ef6FAC5F6BE758508605bd87356776BfbfbA0;
  address internal BLOCKHASH_STORE_ADDRESS = 0x032477835Bc60C777321e1e09F1Fe6cFB82e111C;

  function setUp() public virtual override {
    BaseTest.setUp();

    LINK = new LinkToken();
    s_coordinator = new NoCancelVRFCoordinatorV2(address(LINK), LINK_ETH_FEED_ADDRESS, BLOCKHASH_STORE_ADDRESS);
    s_exposedCoordinator = new ExposedNoCancelVRFCoordinatorV2(address(LINK), LINK_ETH_FEED_ADDRESS, BLOCKHASH_STORE_ADDRESS);
  }
}

contract NoCancelVRFCoordinatorV2_cancelSubscription is NoCancelVRFCoordinatorV2Setup {
  address internal subscriptionOwner = 0x054684DA2dCcf0e79E72DcAb1b5db42E71b66c8A;

  event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amount);
  event SubscriptionCreated(uint64 indexed subId, address owner);
  event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance);

  function testCancelSubscription() public {
    // switch caller to subscriptionOwner so that the subscription owner and the
    // contract owner are different.
    vm.stopPrank();
    vm.startPrank(subscriptionOwner);
    vm.expectEmit(true /* subId */, false /* checkTopic2 */, false /* checkTopic3 */, true /* owner */);
    emit SubscriptionCreated(1, subscriptionOwner);
    s_coordinator.createSubscription();

    // do basic assertions on a brand new subscription
    (uint96 balance, uint64 requestCount, address owner, address[] memory consumers) = s_coordinator.getSubscription(1);
    assertEq(0 /* expected balance */, balance /* actual balance */);
    assertEq(0 /* expected count */, requestCount /* actual count */);
    assertEq(subscriptionOwner /* expected owner */, owner /* actual owner */);
    assertEq(0 /* expected length */, consumers.length /* actual length */);

    vm.expectRevert(bytes("sub cancellation not allowed"));
    s_coordinator.cancelSubscription(1, address(0)); // args don't matter, always reverts
  }

  function testOwnerCancelSubscriptionInvalidSubId() public {
    vm.expectRevert(NoCancelVRFCoordinatorV2.InvalidSubscription.selector);
    s_coordinator.ownerCancelSubscription(1); // non-existent subscription
  }

  function testOwnerCancelSubscription() public {
    // create a subscription, fund it, then cancel it.
    // assert that the funds are sent to the contract owner address.

    // switch caller to subscriptionOwner so that the subscription owner and the
    // contract owner are different.
    vm.stopPrank();
    vm.prank(subscriptionOwner);
    vm.expectEmit(true /* subId */, false /* checkTopic2 */, false /* checkTopic3 */, true /* owner */);
    emit SubscriptionCreated(1, subscriptionOwner);
    s_coordinator.createSubscription();

    // do basic assertions on a brand new subscription
    (uint96 balance, uint64 requestCount, address owner, address[] memory consumers) = s_coordinator.getSubscription(1);
    assertEq(0 /* expected balance */, balance /* actual balance */);
    assertEq(0 /* expected count */, requestCount /* actual count */);
    assertEq(subscriptionOwner /* expected owner */, owner /* actual owner */);
    assertEq(0 /* expected length */, consumers.length /* actual length */);

    // transfer and call some link to the newly created subscription
    uint256 ownerLinkBalanceBefore = LINK.balanceOf(OWNER);
    vm.expectEmit(true /* subId */, false /* checkTopic2 */, false /* checkTopic3 */, true /* oldBalance,newBalance */);
    emit SubscriptionFunded(1, 0, 1e18);
    vm.prank(OWNER);
    bool success = LINK.transferAndCall(address(s_coordinator), 1e18 /* 1 LINK */, abi.encode(uint64(1) /* sub id */));
    assertTrue(success);

    (balance, requestCount, owner, consumers) = s_coordinator.getSubscription(1);
    assertEq(1e18 /* expected balance */, balance /* actual balance */); // balance is 1 LINK

    // cancel subscription and assert OWNER link balance updated
    vm.expectEmit(true /* subId */, false /* checkTopic2 */, false /* checkTopic3 */, true /* to,amount */);
    emit SubscriptionCanceled(1, OWNER, 1e18);
    vm.prank(OWNER);
    s_coordinator.ownerCancelSubscription(1);
    uint256 ownerLinkBalanceAfter = LINK.balanceOf(OWNER);
    assertEq(ownerLinkBalanceBefore, ownerLinkBalanceAfter);
  }
}

contract NoCancelVRFCoordinatorV2_calculatePaymentAmount is NoCancelVRFCoordinatorV2Setup {
  function testFuzzCalculatePaymentAmount(uint256 gasAfterPaymentCalculation, uint32 fulfillmentFlatFeeLinkPPM, uint256 weiPerUnitGas) public {
    uint96 actualAmount = s_exposedCoordinator.calculatePaymentAmountTest(
      gasAfterPaymentCalculation,
      fulfillmentFlatFeeLinkPPM,
      weiPerUnitGas
    );
    uint96 expectedAmount = uint96(1e12 * uint256(fulfillmentFlatFeeLinkPPM));
    assertEq(expectedAmount, actualAmount);
  }
}
