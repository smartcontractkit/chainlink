pragma solidity 0.8.19;

import "./BaseTest.t.sol";
import {ExposedVRFCoordinatorV2_5} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5.sol";
import {VRFV2PlusLoadTestWithMetrics} from "../dev/testhelpers/VRFV2PlusLoadTestWithMetrics.sol";
import {SubscriptionAPI} from "../dev/SubscriptionAPI.sol";
import {MockLinkToken} from "../../functions/tests/v1_X/testhelpers/MockLinkToken.sol";
import {MockV3Aggregator} from "../../shared/mocks/MockV3Aggregator.sol";
import "@openzeppelin/contracts/utils/Strings.sol"; // for Strings.toString
import {VmSafe} from "forge-std/Vm.sol";

contract VRFV2PlusSubscriptionAPITest is BaseTest {
  event SubscriptionFunded(uint256 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionFundedWithNative(uint256 indexed subId, uint256 oldNativeBalance, uint256 newNativeBalance);
  event SubscriptionCanceled(uint256 indexed subId, address to, uint256 amountLink, uint256 amountNative);
  event FundsRecovered(address to, uint256 amountLink);
  event NativeFundsRecovered(address to, uint256 amountNative);
  event SubscriptionOwnerTransferRequested(uint256 indexed subId, address from, address to);
  event SubscriptionOwnerTransferred(uint256 indexed subId, address from, address to);
  event SubscriptionConsumerAdded(uint256 indexed subId, address consumer);
  event SubscriptionConsumerRemoved(uint256 indexed subId, address consumer);

  ExposedVRFCoordinatorV2_5 s_subscriptionAPI;

  function setUp() public override {
    BaseTest.setUp();
    address bhs = makeAddr("bhs");
    s_subscriptionAPI = new ExposedVRFCoordinatorV2_5(bhs);
  }

  function testDefaultState() public {
    assertEq(address(s_subscriptionAPI.LINK()), address(0));
    assertEq(address(s_subscriptionAPI.LINK_NATIVE_FEED()), address(0));
    assertEq(s_subscriptionAPI.s_currentSubNonce(), 0);
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 0);
    assertEq(s_subscriptionAPI.s_totalBalance(), 0);
    assertEq(s_subscriptionAPI.s_totalNativeBalance(), 0);
  }

  function testSetLINKAndLINKNativeFeed() public {
    address link = makeAddr("link");
    address linkNativeFeed = makeAddr("linkNativeFeed");
    s_subscriptionAPI.setLINKAndLINKNativeFeed(link, linkNativeFeed);
    assertEq(address(s_subscriptionAPI.LINK()), link);
    assertEq(address(s_subscriptionAPI.LINK_NATIVE_FEED()), linkNativeFeed);

    // try setting it again, should revert
    vm.expectRevert(SubscriptionAPI.LinkAlreadySet.selector);
    s_subscriptionAPI.setLINKAndLINKNativeFeed(link, linkNativeFeed);
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
    emit SubscriptionFundedWithNative(subId, 0, 5 ether);
    s_subscriptionAPI.fundSubscriptionWithNative{value: 5 ether}(subId);

    // change back to owner and cancel the subscription
    changePrank(OWNER);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, subOwner, 0 /* link balance */, 5 ether /* native balance */);
    s_subscriptionAPI.ownerCancelSubscription(subId);

    // assert that the subscription no longer exists
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 0);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).owner, address(0));
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).nativeBalance, 0);

    // check the native balance of the subOwner, should be 10 ether
    assertEq(address(subOwner).balance, 10 ether);
  }

  function testOwnerCancelSubscriptionLinkFundsOnly() public {
    // CASE: new subscription with link funds only
    // no native funds.
    // should cancel and return the link funds

    // Create link token and set the link token on the subscription api object
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
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
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
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
    emit SubscriptionFundedWithNative(subId, 0, 5 ether);
    s_subscriptionAPI.fundSubscriptionWithNative{value: 5 ether}(subId);

    // change back to owner and cancel the subscription
    changePrank(OWNER);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, subOwner, 5 ether /* link balance */, 5 ether /* native balance */);
    s_subscriptionAPI.ownerCancelSubscription(subId);

    // assert that the subscription no longer exists
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 0);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).owner, address(0));
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).balance, 0);
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).nativeBalance, 0);

    // check the link balance of the sub owner, should be 5 LINK
    assertEq(linkToken.balanceOf(subOwner), 5 ether, "link balance incorrect");
    // check the ether balance of the sub owner, should be 10 ether
    assertEq(address(subOwner).balance, 10 ether, "native balance incorrect");
  }

  function testRecoverFundsLINKNotSet() public {
    // CASE: link token not set
    // should revert with error LinkNotSet

    // call recoverFunds
    vm.expectRevert(SubscriptionAPI.LinkNotSet.selector);
    s_subscriptionAPI.recoverFunds(OWNER);
  }

  function testRecoverFundsBalanceInvariantViolated() public {
    // CASE: link token set
    // and internal balance is greater than external balance

    // Create link token and set the link token on the subscription api object
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // set the total balance to be greater than the external balance
    // so that we trigger the invariant violation
    // note that this field is not modifiable in the actual contracts
    // other than through onTokenTransfer or similar functions
    s_subscriptionAPI.setTotalBalanceTestingOnlyXXX(100 ether);

    // call recoverFunds
    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.BalanceInvariantViolated.selector, 100 ether, 0));
    s_subscriptionAPI.recoverFunds(OWNER);
  }

  function testRecoverFundsAmountToTransfer() public {
    // CASE: link token set
    // and internal balance is less than external balance
    // (i.e invariant is not violated)
    // should recover funds successfully

    // Create link token and set the link token on the subscription api object
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // transfer 10 LINK to the contract to recover
    bool success = linkToken.transfer(address(s_subscriptionAPI), 10 ether);
    assertTrue(success, "failed link transfer");

    // call recoverFunds
    vm.expectEmit(true, false, false, true);
    emit FundsRecovered(OWNER, 10 ether);
    s_subscriptionAPI.recoverFunds(OWNER);
  }

  function testRecoverFundsNothingToTransfer() public {
    // CASE: link token set
    // and there is nothing to transfer
    // should do nothing at all

    // Create link token and set the link token on the subscription api object
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // create a subscription and fund it with 5 LINK
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

    // call recoverFunds, nothing should happen because external balance == internal balance
    s_subscriptionAPI.recoverFunds(OWNER);
    assertEq(linkToken.balanceOf(address(s_subscriptionAPI)), s_subscriptionAPI.s_totalBalance());
  }

  function testRecoverNativeFundsBalanceInvariantViolated() public {
    // set the total balance to be greater than the external balance
    // so that we trigger the invariant violation
    // note that this field is not modifiable in the actual contracts
    // other than through onTokenTransfer or similar functions
    s_subscriptionAPI.setTotalNativeBalanceTestingOnlyXXX(100 ether);

    // call recoverFunds
    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.BalanceInvariantViolated.selector, 100 ether, 0));
    s_subscriptionAPI.recoverNativeFunds(payable(OWNER));
  }

  function testRecoverNativeFundsAmountToTransfer() public {
    // transfer 10 LINK to the contract to recover
    vm.deal(address(s_subscriptionAPI), 10 ether);

    // call recoverFunds
    vm.expectEmit(true, false, false, true);
    emit NativeFundsRecovered(OWNER, 10 ether);
    s_subscriptionAPI.recoverNativeFunds(payable(OWNER));
  }

  function testRecoverNativeFundsNothingToTransfer() public {
    // create a subscription and fund it with 5 ether
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint64 nonceBefore = s_subscriptionAPI.s_currentSubNonce();
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 1);

    // fund the subscription with ether
    vm.deal(subOwner, 5 ether);
    changePrank(subOwner);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionFundedWithNative(subId, 0, 5 ether);
    s_subscriptionAPI.fundSubscriptionWithNative{value: 5 ether}(subId);

    // call recoverNativeFunds, nothing should happen because external balance == internal balance
    changePrank(OWNER);
    s_subscriptionAPI.recoverNativeFunds(payable(OWNER));
    assertEq(address(s_subscriptionAPI).balance, s_subscriptionAPI.s_totalNativeBalance());
  }

  function testWithdrawNoLink() public {
    // CASE: no link token set
    vm.expectRevert(SubscriptionAPI.LinkNotSet.selector);
    s_subscriptionAPI.withdraw(OWNER);
  }

  function testWithdrawInsufficientBalance() public {
    // CASE: link token set, trying to withdraw
    // more than balance
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // call withdraw
    vm.expectRevert(SubscriptionAPI.InsufficientBalance.selector);
    s_subscriptionAPI.withdraw(OWNER);
  }

  function testWithdrawSufficientBalanceLinkSet() public {
    // CASE: link token set, trying to withdraw
    // less than balance
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // transfer 10 LINK to the contract to withdraw
    bool success = linkToken.transfer(address(s_subscriptionAPI), 10 ether);
    assertTrue(success, "failed link transfer");

    // set the withdrawable tokens of the contract to be 1 ether
    s_subscriptionAPI.setWithdrawableTokensTestingOnlyXXX(1 ether);
    assertEq(s_subscriptionAPI.getWithdrawableTokensTestingOnlyXXX(), 1 ether);

    // set the total balance to be the same as the link balance for consistency
    // (this is not necessary for the test, but just to be sane)
    s_subscriptionAPI.setTotalBalanceTestingOnlyXXX(10 ether);

    // call Withdraw from owner address
    uint256 ownerBalance = linkToken.balanceOf(OWNER);
    changePrank(OWNER);
    s_subscriptionAPI.withdraw(OWNER);
    // assert link balance of owner
    assertEq(linkToken.balanceOf(OWNER) - ownerBalance, 1 ether, "owner link balance incorrect");
    // assert state of subscription api
    assertEq(s_subscriptionAPI.getWithdrawableTokensTestingOnlyXXX(), 0, "owner withdrawable tokens incorrect");
    // assert that total balance is changed by the withdrawn amount
    assertEq(s_subscriptionAPI.s_totalBalance(), 9 ether, "total balance incorrect");
  }

  function testWithdrawNativeInsufficientBalance() public {
    // CASE: trying to withdraw more than balance
    // should revert with InsufficientBalance

    // call WithdrawNative
    changePrank(OWNER);
    vm.expectRevert(SubscriptionAPI.InsufficientBalance.selector);
    s_subscriptionAPI.withdrawNative(payable(OWNER));
  }

  function testWithdrawLinkInvalidOwner() public {
    address invalidAddress = makeAddr("invalidAddress");
    changePrank(invalidAddress);
    vm.expectRevert("Only callable by owner");
    s_subscriptionAPI.withdraw(payable(OWNER));
  }

  function testWithdrawNativeInvalidOwner() public {
    address invalidAddress = makeAddr("invalidAddress");
    changePrank(invalidAddress);
    vm.expectRevert("Only callable by owner");
    s_subscriptionAPI.withdrawNative(payable(OWNER));
  }

  function testWithdrawNativeSufficientBalance() public {
    // CASE: trying to withdraw less than balance
    // should withdraw successfully

    // transfer 10 ether to the contract to withdraw
    vm.deal(address(s_subscriptionAPI), 10 ether);

    // set the withdrawable eth of the contract to be 1 ether
    s_subscriptionAPI.setWithdrawableNativeTestingOnlyXXX(1 ether);
    assertEq(s_subscriptionAPI.getWithdrawableNativeTestingOnlyXXX(), 1 ether);

    // set the total balance to be the same as the eth balance for consistency
    // (this is not necessary for the test, but just to be sane)
    s_subscriptionAPI.setTotalNativeBalanceTestingOnlyXXX(10 ether);

    // call WithdrawNative from owner address
    changePrank(OWNER);
    s_subscriptionAPI.withdrawNative(payable(OWNER));
    // assert native balance
    assertEq(address(OWNER).balance, 1 ether, "owner native balance incorrect");
    // assert state of subscription api
    assertEq(s_subscriptionAPI.getWithdrawableNativeTestingOnlyXXX(), 0, "owner withdrawable native incorrect");
    // assert that total balance is changed by the withdrawn amount
    assertEq(s_subscriptionAPI.s_totalNativeBalance(), 9 ether, "total native balance incorrect");
  }

  function testOnTokenTransferCallerNotLink() public {
    vm.expectRevert(SubscriptionAPI.OnlyCallableFromLink.selector);
    s_subscriptionAPI.onTokenTransfer(makeAddr("someaddress"), 1 ether, abi.encode(uint256(1)));
  }

  function testOnTokenTransferInvalidCalldata() public {
    // create and set link token on subscription api
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // call link.transferAndCall with invalid calldata
    vm.expectRevert(SubscriptionAPI.InvalidCalldata.selector);
    linkToken.transferAndCall(address(s_subscriptionAPI), 1 ether, abi.encode(uint256(1), address(1)));
  }

  function testOnTokenTransferInvalidSubscriptionId() public {
    // create and set link token on subscription api
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // generate bogus sub id
    uint256 subId = uint256(keccak256("idontexist"));

    // try to fund bogus sub id
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    linkToken.transferAndCall(address(s_subscriptionAPI), 1 ether, abi.encode(subId));
  }

  function testOnTokenTransferSuccess() public {
    // happy path link funding test
    // create and set link token on subscription api
    MockLinkToken linkToken = new MockLinkToken();
    s_subscriptionAPI.setLINKAndLINKNativeFeed(address(linkToken), address(0));
    assertEq(address(s_subscriptionAPI.LINK()), address(linkToken));

    // create a subscription and fund it with 5 LINK
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

    // assert that the subscription is funded
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).balance, 5 ether);
  }

  function testFundSubscriptionWithNativeInvalidSubscriptionId() public {
    // CASE: invalid subscription id
    // should revert with InvalidSubscription

    uint256 subId = uint256(keccak256("idontexist"));

    // try to fund the subscription with native, should fail
    address funder = makeAddr("funder");
    vm.deal(funder, 5 ether);
    changePrank(funder);
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_subscriptionAPI.fundSubscriptionWithNative{value: 5 ether}(subId);
  }

  function testFundSubscriptionWithNative() public {
    // happy path test
    // funding subscription with native

    // create a subscription and fund it with native
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint64 nonceBefore = s_subscriptionAPI.s_currentSubNonce();
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 1);

    // fund the subscription with native
    vm.deal(subOwner, 5 ether);
    changePrank(subOwner);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionFundedWithNative(subId, 0, 5 ether);
    s_subscriptionAPI.fundSubscriptionWithNative{value: 5 ether}(subId);

    // assert that the subscription is funded
    assertEq(s_subscriptionAPI.getSubscriptionStruct(subId).nativeBalance, 5 ether);
  }

  function testCreateSubscription() public {
    // test that the subscription is created successfully
    // and test the initial state of the subscription
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint64 nonceBefore = s_subscriptionAPI.s_currentSubNonce();
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 1);
    assertEq(s_subscriptionAPI.getActiveSubscriptionIdsLength(), 1);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).owner, subOwner);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers.length, 0);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).requestedOwner, address(0));
  }

  function testCreateSubscriptionRecreate() public {
    // create two subscriptions from the same eoa
    // they should never be the same due to nonce incrementation
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint64 nonceBefore = s_subscriptionAPI.s_currentSubNonce();
    uint256 subId1 = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 1);
    uint256 subId2 = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.s_currentSubNonce(), nonceBefore + 2);
    assertTrue(subId1 != subId2);
  }

  function testSubscriptionOwnershipTransfer() public {
    // create two eoa's, and create a subscription from one of them
    // and transfer ownership to the other
    // assert that the subscription is now owned by the other eoa
    address oldOwner = makeAddr("oldOwner");
    address newOwner = makeAddr("newOwner");

    // create sub
    changePrank(oldOwner);
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).owner, oldOwner);

    // request ownership transfer
    changePrank(oldOwner);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionOwnerTransferRequested(subId, oldOwner, newOwner);
    s_subscriptionAPI.requestSubscriptionOwnerTransfer(subId, newOwner);

    // accept ownership transfer from newOwner
    changePrank(newOwner);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionOwnerTransferred(subId, oldOwner, newOwner);
    s_subscriptionAPI.acceptSubscriptionOwnerTransfer(subId);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).requestedOwner, address(0));
  }

  function testAddConsumerTooManyConsumers() public {
    // add 100 consumers to a sub and then
    // try adding one more and see the revert
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint256 subId = s_subscriptionAPI.createSubscription();
    for (uint256 i = 0; i < 100; i++) {
      address consumer = makeAddr(Strings.toString(i));
      vm.expectEmit(true, false, false, true);
      emit SubscriptionConsumerAdded(subId, consumer);
      s_subscriptionAPI.addConsumer(subId, consumer);
    }

    // try adding one more consumer, should revert
    address lastConsumer = makeAddr("consumer");
    changePrank(subOwner);
    vm.expectRevert(SubscriptionAPI.TooManyConsumers.selector);
    s_subscriptionAPI.addConsumer(subId, lastConsumer);
  }

  function testAddConsumerReaddSameConsumer() public {
    // try adding the same consumer twice
    // should be a no-op
    // assert state is unchanged after the 2nd add
    address subOwner = makeAddr("subOwner");
    address consumer = makeAddr("consumer");
    changePrank(subOwner);
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers.length, 0);
    changePrank(subOwner);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionConsumerAdded(subId, consumer);
    s_subscriptionAPI.addConsumer(subId, consumer);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers.length, 1);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers[0], consumer);

    // add consumer again, should be no-op
    changePrank(subOwner);
    VmSafe.Log[] memory events = vm.getRecordedLogs();
    s_subscriptionAPI.addConsumer(subId, consumer);
    assertEq(events.length, 0);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers.length, 1);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers[0], consumer);

    // remove consumer
    vm.expectEmit(true, false, false, true);
    emit SubscriptionConsumerRemoved(subId, consumer);
    s_subscriptionAPI.removeConsumer(subId, consumer);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers.length, 0);

    // removing consumer twice should revert
    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.InvalidConsumer.selector, subId, address(consumer)));
    s_subscriptionAPI.removeConsumer(subId, consumer);

    //re-add consumer
    vm.expectEmit(true, false, false, true);
    emit SubscriptionConsumerAdded(subId, consumer);
    s_subscriptionAPI.addConsumer(subId, consumer);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers.length, 1);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers[0], consumer);
  }

  function testAddConsumer() public {
    // create a subscription and add a consumer
    // assert subscription state afterwards
    address subOwner = makeAddr("subOwner");
    address consumer = makeAddr("consumer");
    changePrank(subOwner);
    uint256 subId = s_subscriptionAPI.createSubscription();
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers.length, 0);

    // only subscription owner can add a consumer
    address notSubOwner = makeAddr("notSubOwner");
    changePrank(notSubOwner);
    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.MustBeSubOwner.selector, subOwner));
    s_subscriptionAPI.addConsumer(subId, consumer);

    // subscription owner is able to add a consumer
    changePrank(subOwner);
    vm.expectEmit(true, false, false, true);
    emit SubscriptionConsumerAdded(subId, consumer);
    s_subscriptionAPI.addConsumer(subId, consumer);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers.length, 1);
    assertEq(s_subscriptionAPI.getSubscriptionConfig(subId).consumers[0], consumer);
  }
}
