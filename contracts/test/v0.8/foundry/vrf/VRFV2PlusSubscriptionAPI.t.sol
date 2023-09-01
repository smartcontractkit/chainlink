pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {ExposedVRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/testhelpers/ExposedVRFCoordinatorV2Plus.sol";
import {SubscriptionAPI} from "../../../../src/v0.8/dev/vrf/SubscriptionAPI.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";

contract VRFV2PlusSubscriptionAPITest is BaseTest {
  event SubscriptionFunded(uint256 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionFundedWithEth(uint256 indexed subId, uint256 oldEthBalance, uint256 newEthBalance);
  event SubscriptionCanceled(uint256 indexed subId, address to, uint256 amountLink, uint256 amountEth);

  ExposedVRFCoordinatorV2Plus s_subscriptionAPI;

  function setUp() public override {
    BaseTest.setUp();
    address bhs = makeAddr("bhs");
    s_subscriptionAPI = new ExposedVRFCoordinatorV2Plus(bhs);
  }

  function testDefaultState() public {
    assertEq(address(s_subscriptionAPI.LINK()), address(0));
    assertEq(address(s_subscriptionAPI.LINK_ETH_FEED()), address(0));
    assertEq(s_subscriptionAPI.s_currentSubNonce(), 0);
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 0);
    assertEq(s_subscriptionAPI.s_totalBalance(), 0);
    assertEq(s_subscriptionAPI.s_totalEthBalance(), 0);
  }

  function testSetLINKAndLINKETHFeed() public {
    address link = makeAddr("link");
    address linkEthFeed = makeAddr("linkEthFeed");
    s_subscriptionAPI.setLINKAndLINKETHFeed(link, linkEthFeed);
    assertEq(address(s_subscriptionAPI.LINK()), link);
    assertEq(address(s_subscriptionAPI.LINK_ETH_FEED()), linkEthFeed);

    // try setting it again, should revert
    vm.expectRevert(SubscriptionAPI.LinkAlreadySet.selector);
    s_subscriptionAPI.setLINKAndLINKETHFeed(link, linkEthFeed);
  }

  function testOwnerCancelSubscriptionNoFunds() public {
    // CASE: new subscription w/ no funds at all
    // Should cancel trivially

    // Note that the link token is not set, but this should still
    // not fail in that case.

    // Create the subscription from a separate address
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint64 nonceBefore = s_subscriptionAPI.s_currentSubNonce();
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 1);

    // change back to owner and cancel the subscription
    changePrank(OWNER);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, subOwner, 0, 0);
    s_subscriptionAPI.ownerCancelSubscription(subId);

    // assert that the subscription no longer exists
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 0);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).owner, address(0));
    // no point in checking s_subscriptions because all fields are zeroed out
    // due to no balance and no requests made
  }

  function testOwnerCancelSubscriptionNativeFundsOnly() public {
    // CASE: new subscription with native funds only
    // no link funds.
    // should cancel and return the native funds

    // Create the subscription from a separate address
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint64 nonceBefore = s_subscriptionAPI.s_currentSubNonce();
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 1);

    // fund the subscription with ether
    vm.deal(subOwner, 10 ether);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionFundedWithEth(subId, 0, 5 ether);
    s_subscriptionAPI.fundSubscriptionWithEth{value: 5 ether}(subId);

    // change back to owner and cancel the subscription
    changePrank(OWNER);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, subOwner, 0 /* link balance */, 5 ether /* native balance */);
    s_subscriptionAPI.ownerCancelSubscription(subId);

    // assert that the subscription no longer exists
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 0);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).owner, address(0));
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).ethBalance, 0);

    // check the ether balance of the subOwner, should be 5 ether
    assertEq(address(subOwner).balance, 5 ether);
  }

  function testOwnerCancelSubscriptionLinkFundsOnly() public {
    // CASE: new subscription with link funds only
    // no native funds.
    // should cancel and return the link funds

    // Create link token and set the link token on the subscription api object
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKETHFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // Create the subscription from a separate address
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint64 nonceBefore = s_subscriptionAPI.s_currentSubNonce();
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 1);

    // fund the subscription with link
    // can do it from the owner acct because anyone can fund a subscription
    changePrank(OWNER);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionFunded(subId, 0, 5 ether);
    bool success = linkToken.transferAndCall(address(s_subscriptionAPI), 5 ether, abi.encode(subId));
    assertTrue(success, "failed link transfer and call");

    // change back to owner and cancel the subscription
    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, subOwner, 5 ether /* link balance */, 0 /* native balance */);
    s_subscriptionAPI.ownerCancelSubscription(subId);

    // assert that the subscription no longer exists
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 0);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).owner, address(0));
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).balance, 0);

    // check the link balance of the sub owner, should be 5 LINK
    assertEq(linkToken.balanceOf(subOwner), 5 ether);
  }

  function testOwnerCancelSubscriptionNativeAndLinkFunds() public {
    // CASE: new subscription with link and native funds
    // should cancel and return both link and native funds

    // Create link token and set the link token on the subscription api object
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKETHFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // Create the subscription from a separate address
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint64 nonceBefore = s_subscriptionAPI.s_currentSubNonce();
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 1);

    // fund the subscription with link
    changePrank(OWNER);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionFunded(subId, 0, 5 ether);
    bool success = linkToken.transferAndCall(address(s_subscriptionAPI), 5 ether, abi.encode(subId));
    assertTrue(success, "failed link transfer and call");

    // fund the subscription with ether
    vm.deal(subOwner, 10 ether);
    changePrank(subOwner);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionFundedWithEth(subId, 0, 5 ether);
    s_subscriptionAPI.fundSubscriptionWithEth{value: 5 ether}(subId);

    // change back to owner and cancel the subscription
    changePrank(OWNER);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, subOwner, 5 ether /* link balance */, 5 ether /* native balance */);
    s_subscriptionAPI.ownerCancelSubscription(subId);

    // assert that the subscription no longer exists
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 0);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).owner, address(0));
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).balance, 0);
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).ethBalance, 0);

    // check the link balance of the sub owner, should be 5 LINK
    assertEq(linkToken.balanceOf(subOwner), 5 ether);
    // check the ether balance of the sub owner, should be 5 ether
    assertEq(address(subOwner).balance, 5 ether);
  }

  function testRecoverFundsLINKNotSet() public {

  }

  function testRecoverFundsBalanceInvariantViolated() public {

  }

  function testRecoverFundsAmountToTransfer() public {

  }

  function testRecoverFundsNothingToTransfer() public {

  }

  function testRecoverEthFundsBalanceInvariantViolated() public {

  }

  function testRecoverEthFundsAmountToTransfer() public {

  }

  function testRecoverEthFundsNothingToTransfer() public {

  }

  function testOracleWithdrawInsufficientBalance() public {

  }

  function testOracleWithdrawSufficientBalanceNoLink() public {

  }

  function testOracleWithdrawSufficientBalanceLinkSet() public {

  }

  function testOracleWithdrawEthInsufficientBalance() public {

  }

  function testOracleWithdrawEthSufficientBalance() public {

  }

  function testOnTokenTransferCallerNotLink() public {

  }

  function testOnTokenTransferInvalidCalldata() public {

  }

  function testOnTokenTransferInvalidSubscriptionId() public {

  }

  function testOnTokenTransferSuccess() public {

  }

  function testFundSubscriptionWithEthInvalidSubscriptionId() public {

  }

  function testFundSubscriptionWithEth() public {

  }

  function testCreateSubscription() public {

  }

  function testCreateSubscriptionRecreate() public {

  }

  function testSubscriptionOwnershipTransfer() public {

  }

  function testAddConsumerTooManyConsumers() public {

  }

  function testAddConsumerReaddSameConsumer() public {

  }

  function testAddConsumer() public {

  }
}
