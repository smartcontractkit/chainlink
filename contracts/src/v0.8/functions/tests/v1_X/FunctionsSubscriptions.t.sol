// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsRouter} from "../../dev/v1_X/FunctionsRouter.sol";
import {FunctionsSubscriptions} from "../../dev/v1_X/FunctionsSubscriptions.sol";
import {FunctionsResponse} from "../../dev/v1_X/libraries/FunctionsResponse.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

import {FunctionsRouterSetup, FunctionsOwnerAcceptTermsOfServiceSetup, FunctionsClientSetup, FunctionsSubscriptionSetup, FunctionsClientRequestSetup, FunctionsFulfillmentSetup} from "./Setup.t.sol";

import "forge-std/Vm.sol";

// ================================================================
// |                    Functions Subscriptions                   |
// ================================================================

contract FunctionsSubscriptions_Constructor_Helper is FunctionsSubscriptions {
  constructor(address link) FunctionsSubscriptions(link) {}

  function getLinkToken() public view returns (IERC20) {
    return IERC20(i_linkToken);
  }

  // overrides
  function _getMaxConsumers() internal pure override returns (uint16) {
    return 0;
  }

  function _getSubscriptionDepositDetails() internal pure override returns (uint16, uint72) {
    return (0, 0);
  }

  function _onlySenderThatAcceptedToS() internal override {}

  function _onlyRouterOwner() internal override {}

  function _whenNotPaused() internal override {}
}

/// @notice #constructor
contract FunctionsSubscriptions_Constructor is BaseTest {
  FunctionsSubscriptions_Constructor_Helper s_subscriptionsHelper;
  address internal s_linkToken = 0x01BE23585060835E02B77ef475b0Cc51aA1e0709;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_subscriptionsHelper = new FunctionsSubscriptions_Constructor_Helper(s_linkToken);
  }

  function test_Constructor_Success() public {
    assertEq(address(s_linkToken), address(s_subscriptionsHelper.getLinkToken()));
  }
}

/// @notice #_markRequestInFlight
contract FunctionsSubscriptions__MarkRequestInFlight {
  // TODO: make contract internal function helper
}

/// @notice #_pay
contract FunctionsSubscriptions__Pay {
  // TODO: make contract internal function helper
}

/// @notice #ownerCancelSubscription
contract FunctionsSubscriptions_OwnerCancelSubscription is FunctionsSubscriptionSetup {
  function test_OwnerCancelSubscription_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.ownerCancelSubscription(s_subscriptionId);
  }

  function test_OwnerCancelSubscription_RevertIfNoSubscription() public {
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    uint64 invalidSubscriptionId = 123456789;
    s_functionsRouter.ownerCancelSubscription(invalidSubscriptionId);
  }

  function test_OwnerCancelSubscription_SuccessSubOwnerRefunded() public {
    uint256 subscriptionOwnerBalanceBefore = s_linkToken.balanceOf(OWNER_ADDRESS);
    s_functionsRouter.ownerCancelSubscription(s_subscriptionId);
    uint256 subscriptionOwnerBalanceAfter = s_linkToken.balanceOf(OWNER_ADDRESS);
    assertEq(subscriptionOwnerBalanceBefore + s_subscriptionInitialFunding, subscriptionOwnerBalanceAfter);
  }

  function test_OwnerCancelSubscription_SuccessWhenRequestInFlight() public {
    // send request
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets;
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);

    s_functionsClient.sendRequest(s_donId, sourceCode, secrets, args, bytesArgs, s_subscriptionId, 5500);
    s_functionsRouter.ownerCancelSubscription(s_subscriptionId);
  }

  function test_OwnerCancelSubscription_SuccessDeletesSubscription() public {
    s_functionsRouter.ownerCancelSubscription(s_subscriptionId);
    // Subscription should no longer exist
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    s_functionsRouter.getSubscription(s_subscriptionId);
  }

  event SubscriptionCanceled(uint64 indexed subscriptionId, address fundsRecipient, uint256 fundsAmount);

  function test_OwnerCancelSubscription_Success() public {
    uint256 subscriptionOwnerBalanceBefore = s_linkToken.balanceOf(OWNER_ADDRESS);

    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1SubscriptionId = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1SubscriptionId, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCanceled(s_subscriptionId, OWNER_ADDRESS, s_subscriptionInitialFunding);

    s_functionsRouter.ownerCancelSubscription(s_subscriptionId);

    uint256 subscriptionOwnerBalanceAfter = s_linkToken.balanceOf(OWNER_ADDRESS);
    assertEq(subscriptionOwnerBalanceBefore + s_subscriptionInitialFunding, subscriptionOwnerBalanceAfter);
  }
}

/// @notice #recoverFunds
contract FunctionsSubscriptions_RecoverFunds is FunctionsRouterSetup {
  event FundsRecovered(address to, uint256 amount);

  function test_RecoverFunds_Success() public {
    uint256 fundsTransferred = 1 * 1e18; // 1 LINK
    s_linkToken.transfer(address(s_functionsRouter), fundsTransferred);

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit FundsRecovered(OWNER_ADDRESS, fundsTransferred);

    uint256 subscriptionOwnerBalanceBefore = s_linkToken.balanceOf(OWNER_ADDRESS);
    s_functionsRouter.recoverFunds(OWNER_ADDRESS);
    uint256 subscriptionOwnerBalanceAfter = s_linkToken.balanceOf(OWNER_ADDRESS);
    assertEq(subscriptionOwnerBalanceBefore + fundsTransferred, subscriptionOwnerBalanceAfter);
  }

  function test_OwnerCancelSubscription_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.recoverFunds(OWNER_ADDRESS);
  }
}

/// @notice #oracleWithdraw
contract FunctionsSubscriptions_OracleWithdraw is FunctionsFulfillmentSetup {
  function test_OracleWithdraw_RevertIfPaused() public {
    s_functionsRouter.pause();

    // Subscription payable balances are set to the Coordinator
    // Send as Coordinator contract
    vm.stopPrank();
    vm.startPrank(address(s_functionsCoordinator));

    vm.expectRevert("Pausable: paused");

    uint96 amountToWithdraw = 1; // more than 0
    s_functionsRouter.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, amountToWithdraw);
  }

  function test_OracleWithdraw_RevertIfNoAmount() public {
    // Subscription payable balances are set to the Coordinator
    // Send as Coordinator contract
    vm.stopPrank();
    vm.startPrank(address(s_functionsCoordinator));

    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);

    uint96 amountToWithdraw = 0;
    s_functionsRouter.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, amountToWithdraw);
  }

  function test_OracleWithdraw_RevertIfAmountMoreThanBalance() public {
    // Subscription payable balances are set to the Coordinator
    // Send as Coordinator contract
    vm.stopPrank();
    vm.startPrank(address(s_functionsCoordinator));

    vm.expectRevert(
      abi.encodeWithSelector(FunctionsSubscriptions.InsufficientBalance.selector, s_fulfillmentCoordinatorBalance)
    );

    uint96 amountToWithdraw = s_fulfillmentCoordinatorBalance + 1;
    s_functionsRouter.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, amountToWithdraw);
  }

  function test_OracleWithdraw_RevertIfBalanceInvariant() public {
    // Subscription payable balances are set to the Coordinator
    // Send as Coordinator contract
    // vm.stopPrank();
    // vm.startPrank(address(s_functionsCoordinator));
    // TODO: Use internal function helper contract to modify s_totalLinkBalance
    // uint96 amountToWithdraw = s_fulfillmentCoordinatorBalance;
    // vm.expectRevert(abi.encodeWithSelector(FunctionsSubscriptions.TotalBalanceInvariantViolated.selector, 0, amountToWithdraw));
    // s_functionsRouter.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, amountToWithdraw);
  }

  function test_OracleWithdraw_SuccessPaysRecipient() public {
    // Subscription payable balances are set to the Coordinator
    // Send as Coordinator contract
    vm.stopPrank();
    vm.startPrank(address(s_functionsCoordinator));

    uint256 transmitterBalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);

    uint96 amountToWithdraw = s_fulfillmentCoordinatorBalance;
    s_functionsRouter.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, amountToWithdraw);

    uint256 transmitterBalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitterBalanceBefore + s_fulfillmentCoordinatorBalance, transmitterBalanceAfter);
  }

  function test_OracleWithdraw_SuccessSetsBalanceToZero() public {
    // Subscription payable balances are set to the Coordinator
    // Send as Coordinator contract
    vm.stopPrank();
    vm.startPrank(address(s_functionsCoordinator));

    uint96 amountToWithdraw = s_fulfillmentCoordinatorBalance;
    s_functionsRouter.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, amountToWithdraw);

    // Attempt to withdraw 1 Juel after withdrawing full balance
    vm.expectRevert(abi.encodeWithSelector(FunctionsSubscriptions.InsufficientBalance.selector, 0));
    s_functionsRouter.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, 1);
  }
}

/// @notice #ownerWithdraw
contract FunctionsSubscriptions_OwnerWithdraw is FunctionsFulfillmentSetup {
  function test_OwnerWithdraw_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.recoverFunds(OWNER_ADDRESS);
  }

  function test_OwnerWithdraw_RevertIfAmountMoreThanBalance() public {
    vm.expectRevert(
      abi.encodeWithSelector(FunctionsSubscriptions.InsufficientBalance.selector, s_fulfillmentRouterOwnerBalance)
    );

    uint96 amountToWithdraw = s_fulfillmentRouterOwnerBalance + 1;
    s_functionsRouter.ownerWithdraw(OWNER_ADDRESS, amountToWithdraw);
  }

  function test_OwnerWithdraw_RevertIfBalanceInvariant() public {
    // TODO: Use internal function helper contract to modify s_totalLinkBalance
    // uint96 amountToWithdraw = s_fulfillmentRouterOwnerBalance;
    // vm.expectRevert(abi.encodeWithSelector(FunctionsSubscriptions.TotalBalanceInvariantViolated.selector, 0, amountToWithdraw));
    // s_functionsRouter.ownerWithdraw(OWNER_ADDRESS, amountToWithdraw);
  }

  function test_OwnerWithdraw_SuccessIfRecipientAddressZero() public {
    uint256 balanceBefore = s_linkToken.balanceOf(address(0));
    uint96 amountToWithdraw = s_fulfillmentRouterOwnerBalance;
    s_functionsRouter.ownerWithdraw(address(0), amountToWithdraw);
    uint256 balanceAfter = s_linkToken.balanceOf(address(0));
    assertEq(balanceBefore + s_fulfillmentRouterOwnerBalance, balanceAfter);
  }

  function test_OwnerWithdraw_SuccessIfNoAmount() public {
    uint256 balanceBefore = s_linkToken.balanceOf(OWNER_ADDRESS);
    uint96 amountToWithdraw = 0;
    s_functionsRouter.ownerWithdraw(OWNER_ADDRESS, amountToWithdraw);
    uint256 balanceAfter = s_linkToken.balanceOf(OWNER_ADDRESS);
    assertEq(balanceBefore + s_fulfillmentRouterOwnerBalance, balanceAfter);
  }

  function test_OwnerWithdraw_SuccessPaysRecipient() public {
    uint256 balanceBefore = s_linkToken.balanceOf(STRANGER_ADDRESS);

    uint96 amountToWithdraw = s_fulfillmentRouterOwnerBalance;
    s_functionsRouter.ownerWithdraw(STRANGER_ADDRESS, amountToWithdraw);

    uint256 balanceAfter = s_linkToken.balanceOf(STRANGER_ADDRESS);
    assertEq(balanceBefore + s_fulfillmentRouterOwnerBalance, balanceAfter);
  }

  function test_OwnerWithdraw_SuccessSetsBalanceToZero() public {
    uint96 amountToWithdraw = s_fulfillmentRouterOwnerBalance;
    s_functionsRouter.ownerWithdraw(OWNER_ADDRESS, amountToWithdraw);

    // Attempt to withdraw 1 Juel after withdrawing full balance
    vm.expectRevert(abi.encodeWithSelector(FunctionsSubscriptions.InsufficientBalance.selector, 0));
    s_functionsRouter.ownerWithdraw(OWNER_ADDRESS, 1);
  }
}

/// @notice #onTokenTransfer
contract FunctionsSubscriptions_OnTokenTransfer is FunctionsClientSetup {
  uint64 s_subscriptionId;

  function setUp() public virtual override {
    FunctionsClientSetup.setUp();

    // Create subscription, but do not fund it
    s_subscriptionId = s_functionsRouter.createSubscription();
    s_functionsRouter.addConsumer(s_subscriptionId, address(s_functionsClient));
  }

  function test_OnTokenTransfer_RevertIfPaused() public {
    // Funding amount must be less than or equal to LINK total supply
    uint256 totalSupplyJuels = 1_000_000_000 * 1e18;
    s_functionsRouter.pause();
    vm.expectRevert("Pausable: paused");
    s_linkToken.transferAndCall(address(s_functionsRouter), totalSupplyJuels, abi.encode(s_subscriptionId));
  }

  function test_OnTokenTransfer_RevertIfCallerIsNotLink() public {
    // Funding amount must be less than or equal to LINK total supply
    uint256 totalSupplyJuels = 1_000_000_000 * 1e18;
    vm.expectRevert(FunctionsSubscriptions.OnlyCallableFromLink.selector);
    s_functionsRouter.onTokenTransfer(address(s_functionsRouter), totalSupplyJuels, abi.encode(s_subscriptionId));
  }

  function test_OnTokenTransfer_RevertIfCallerIsNoCalldata() public {
    // Funding amount must be less than or equal to LINK total supply
    uint256 totalSupplyJuels = 1_000_000_000 * 1e18;
    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);
    s_linkToken.transferAndCall(address(s_functionsRouter), totalSupplyJuels, new bytes(0));
  }

  function test_OnTokenTransfer_RevertIfCallerIsNoSubscription() public {
    // Funding amount must be less than or equal to LINK total supply
    uint256 totalSupplyJuels = 1_000_000_000 * 1e18;
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    uint64 invalidSubscriptionId = 123456789;
    s_linkToken.transferAndCall(address(s_functionsRouter), totalSupplyJuels, abi.encode(invalidSubscriptionId));
  }

  function test_OnTokenTransfer_Success() public {
    // Funding amount must be less than LINK total supply
    uint256 totalSupplyJuels = 1_000_000_000 * 1e18;
    // Some of the total supply is already in the subscription account
    s_linkToken.transferAndCall(address(s_functionsRouter), totalSupplyJuels, abi.encode(s_subscriptionId));
    uint96 subscriptionBalanceAfter = s_functionsRouter.getSubscription(s_subscriptionId).balance;
    assertEq(totalSupplyJuels, subscriptionBalanceAfter);
  }
}

/// @notice #getTotalBalance
contract FunctionsSubscriptions_GetTotalBalance is FunctionsSubscriptionSetup {
  function test_GetTotalBalance_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint96 totalBalance = s_functionsRouter.getTotalBalance();
    assertEq(totalBalance, s_subscriptionInitialFunding);
  }
}

/// @notice #getSubscriptionCount
contract FunctionsSubscriptions_GetSubscriptionCount is FunctionsSubscriptionSetup {
  function test_GetSubscriptionCount_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint96 subscriptionCount = s_functionsRouter.getSubscriptionCount();
    // One subscription was made during setup
    assertEq(subscriptionCount, 1);
  }
}

/// @notice #getSubscriptionsInRange
contract FunctionsSubscriptions_GetSubscriptionsInRange is FunctionsSubscriptionSetup {
  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    // Create 2 more subscriptions
    /* uint64 subscriptionId2 = */ s_functionsRouter.createSubscription();
    uint64 subscriptionId3 = s_functionsRouter.createSubscription();

    // Give each one unique state
    // #1 subscriptionId for requests, #2 empty, #3 proposedOwner of stranger
    s_functionsRouter.proposeSubscriptionOwnerTransfer(subscriptionId3, STRANGER_ADDRESS);
  }

  function test_GetSubscriptionsInRange_RevertIfStartIsAfterEnd() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);

    s_functionsRouter.getSubscriptionsInRange(1, 0);
  }

  function test_GetSubscriptionsInRange_RevertIfEndIsAfterLastSubscription() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint64 lastSubscriptionId = s_functionsRouter.getSubscriptionCount();
    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);
    s_functionsRouter.getSubscriptionsInRange(1, lastSubscriptionId + 1);
  }

  function test_GetSubscriptionsInRange_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint64 lastSubscriptionId = s_functionsRouter.getSubscriptionCount();
    FunctionsSubscriptions.Subscription[] memory subscriptions = s_functionsRouter.getSubscriptionsInRange(
      s_subscriptionId,
      lastSubscriptionId
    );

    assertEq(subscriptions.length, 3);

    // Check subscription 1
    assertEq(subscriptions[0].balance, s_subscriptionInitialFunding);
    assertEq(subscriptions[0].owner, OWNER_ADDRESS);
    assertEq(subscriptions[0].blockedBalance, 0);
    assertEq(subscriptions[0].proposedOwner, address(0));
    assertEq(subscriptions[0].consumers[0], address(s_functionsClient));
    assertEq(subscriptions[0].flags, bytes32(0));

    // Check subscription 2
    assertEq(subscriptions[1].balance, 0);
    assertEq(subscriptions[1].owner, OWNER_ADDRESS);
    assertEq(subscriptions[1].blockedBalance, 0);
    assertEq(subscriptions[1].proposedOwner, address(0));
    assertEq(subscriptions[1].consumers.length, 0);
    assertEq(subscriptions[1].flags, bytes32(0));

    // Check subscription 3
    assertEq(subscriptions[2].balance, 0);
    assertEq(subscriptions[2].owner, OWNER_ADDRESS);
    assertEq(subscriptions[2].blockedBalance, 0);
    assertEq(subscriptions[2].proposedOwner, address(STRANGER_ADDRESS));
    assertEq(subscriptions[2].consumers.length, 0);
    assertEq(subscriptions[2].flags, bytes32(0));
  }
}

/// @notice #getSubscription
contract FunctionsSubscriptions_GetSubscription is FunctionsSubscriptionSetup {
  function test_GetSubscription_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    FunctionsSubscriptions.Subscription memory subscription = s_functionsRouter.getSubscription(s_subscriptionId);

    assertEq(subscription.balance, s_subscriptionInitialFunding);
    assertEq(subscription.owner, OWNER_ADDRESS);
    assertEq(subscription.blockedBalance, 0);
    assertEq(subscription.proposedOwner, address(0));
    assertEq(subscription.consumers[0], address(s_functionsClient));
    assertEq(subscription.flags, bytes32(0));
  }
}

/// @notice #getConsumer
contract FunctionsSubscriptions_GetConsumer is FunctionsSubscriptionSetup {
  function test_GetConsumer_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    FunctionsSubscriptions.Consumer memory consumer = s_functionsRouter.getConsumer(
      address(s_functionsClient),
      s_subscriptionId
    );

    assertEq(consumer.allowed, true);
    assertEq(consumer.initiatedRequests, 0);
    assertEq(consumer.completedRequests, 0);
  }
}

/// @notice #_isExistingSubscription
contract FunctionsSubscriptions__IsExistingSubscription is FunctionsSubscriptionSetup {
  // TODO: make contract internal function helper
}

/// @notice #_isAllowedConsumer
contract FunctionsSubscriptions__IsAllowedConsumer {
  // TODO: make contract internal function helper
}

/// @notice #createSubscription
contract FunctionsSubscriptions_createSubscription is FunctionsOwnerAcceptTermsOfServiceSetup {
  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);

  function test_CreateSubscription_Success() public {
    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;

    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCreated(1, OWNER_ADDRESS);
    uint64 firstCallSubscriptionId = s_functionsRouter.createSubscription();
    assertEq(firstCallSubscriptionId, 1);

    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCreated(2, OWNER_ADDRESS);
    uint64 secondCallSubscriptionId = s_functionsRouter.createSubscription();
    assertEq(secondCallSubscriptionId, 2);

    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCreated(3, OWNER_ADDRESS);
    uint64 thirdCallSubscriptionId = s_functionsRouter.createSubscription();
    assertEq(thirdCallSubscriptionId, 3);
  }

  function test_CreateSubscription_RevertIfPaused() public {
    s_functionsRouter.pause();

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.createSubscription();
  }

  function test_CreateSubscription_RevertIfNotAllowedSender() public {
    // Send as stranger, who has not accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.SenderMustAcceptTermsOfService.selector, STRANGER_ADDRESS));
    s_functionsRouter.createSubscription();
  }
}

/// @notice #createSubscriptionWithConsumer
contract FunctionsSubscriptions_CreateSubscriptionWithConsumer is FunctionsClientSetup {
  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);
  event SubscriptionConsumerAdded(uint64 indexed subscriptionId, address consumer);

  function test_CreateSubscriptionWithConsumer_Success() public {
    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;

    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCreated(1, OWNER_ADDRESS);
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionConsumerAdded(1, address(s_functionsClient));
    uint64 firstCallSubscriptionId = s_functionsRouter.createSubscriptionWithConsumer(address(s_functionsClient));
    assertEq(firstCallSubscriptionId, 1);
    assertEq(s_functionsRouter.getSubscription(firstCallSubscriptionId).consumers[0], address(s_functionsClient));

    // Consumer can be address(0)
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCreated(2, OWNER_ADDRESS);
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionConsumerAdded(2, address(0));
    uint64 secondCallSubscriptionId = s_functionsRouter.createSubscriptionWithConsumer(address(0));
    assertEq(secondCallSubscriptionId, 2);
    assertEq(s_functionsRouter.getSubscription(secondCallSubscriptionId).consumers[0], address(0));

    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCreated(3, OWNER_ADDRESS);
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionConsumerAdded(3, address(s_functionsClient));
    uint64 thirdCallSubscriptionId = s_functionsRouter.createSubscriptionWithConsumer(address(s_functionsClient));
    assertEq(thirdCallSubscriptionId, 3);
    assertEq(s_functionsRouter.getSubscription(thirdCallSubscriptionId).consumers[0], address(s_functionsClient));
  }

  function test_CreateSubscriptionWithConsumer_RevertIfPaused() public {
    s_functionsRouter.pause();

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.createSubscriptionWithConsumer(address(s_functionsClient));
  }

  function test_CreateSubscriptionWithConsumer_RevertIfNotAllowedSender() public {
    // Send as stranger, who has not accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.SenderMustAcceptTermsOfService.selector, STRANGER_ADDRESS));
    s_functionsRouter.createSubscriptionWithConsumer(address(s_functionsClient));
  }
}

/// @notice #proposeSubscriptionOwnerTransfer
contract FunctionsSubscriptions_ProposeSubscriptionOwnerTransfer is FunctionsSubscriptionSetup {
  uint256 internal NEW_OWNER_PRIVATE_KEY_WITH_TOS = 0x3;
  address internal NEW_OWNER_ADDRESS_WITH_TOS = vm.addr(NEW_OWNER_PRIVATE_KEY_WITH_TOS);
  uint256 internal NEW_OWNER_PRIVATE_KEY_WITH_TOS2 = 0x4;
  address internal NEW_OWNER_ADDRESS_WITH_TOS2 = vm.addr(NEW_OWNER_PRIVATE_KEY_WITH_TOS2);
  uint256 internal NEW_OWNER_PRIVATE_KEY_WITHOUT_TOS = 0x5;
  address internal NEW_OWNER_ADDRESS_WITHOUT_TOS = vm.addr(NEW_OWNER_PRIVATE_KEY_WITHOUT_TOS);

  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    // Accept ToS as new owner #1
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITH_TOS);
    bytes32 message = s_termsOfServiceAllowList.getMessage(NEW_OWNER_ADDRESS_WITH_TOS, NEW_OWNER_ADDRESS_WITH_TOS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(NEW_OWNER_ADDRESS_WITH_TOS, NEW_OWNER_ADDRESS_WITH_TOS, r, s, v);

    // Accept ToS as new owner #2
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITH_TOS2);
    bytes32 message2 = s_termsOfServiceAllowList.getMessage(NEW_OWNER_ADDRESS_WITH_TOS2, NEW_OWNER_ADDRESS_WITH_TOS2);
    bytes32 prefixedMessage2 = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message2));
    (uint8 v2, bytes32 r2, bytes32 s2) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage2);
    s_termsOfServiceAllowList.acceptTermsOfService(
      NEW_OWNER_ADDRESS_WITH_TOS2,
      NEW_OWNER_ADDRESS_WITH_TOS2,
      r2,
      s2,
      v2
    );

    vm.stopPrank();
    vm.startPrank(OWNER_ADDRESS);
  }

  function test_ProposeSubscriptionOwnerTransfer_RevertIfPaused() public {
    s_functionsRouter.pause();

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
  }

  function test_ProposeSubscriptionOwnerTransfer_RevertIfNoSubscription() public {
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    uint64 invalidSubscriptionId = 123456789;
    s_functionsRouter.proposeSubscriptionOwnerTransfer(invalidSubscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
  }

  function test_ProposeSubscriptionOwnerTransfer_RevertIfNotSubscriptionOwner() public {
    // Send as non-owner, who has accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITH_TOS);

    vm.expectRevert(FunctionsSubscriptions.MustBeSubscriptionOwner.selector);
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
  }

  function test_ProposeSubscriptionOwnerTransfer_RevertIfNotAllowedSender() public {
    // Remove owner from Allow List
    s_termsOfServiceAllowList.blockSender(OWNER_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.SenderMustAcceptTermsOfService.selector, OWNER_ADDRESS));
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
  }

  function test_ProposeSubscriptionOwnerTransfer_RevertIfEmptyNewOwner() public {
    address EMPTY_ADDRESS = address(0);
    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, EMPTY_ADDRESS);
  }

  function test_ProposeSubscriptionOwnerTransfer_RevertIfInvalidNewOwner() public {
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
  }

  event SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId, address from, address to);

  function test_ProposeSubscriptionOwnerTransfer_Success() public {
    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;

    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionOwnerTransferRequested(s_subscriptionId, OWNER_ADDRESS, NEW_OWNER_ADDRESS_WITH_TOS);
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
    assertEq(s_functionsRouter.getSubscription(s_subscriptionId).proposedOwner, NEW_OWNER_ADDRESS_WITH_TOS);
  }

  function test_ProposeSubscriptionOwnerTransfer_SuccessChangeProposedOwner() public {
    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;

    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionOwnerTransferRequested(s_subscriptionId, OWNER_ADDRESS, NEW_OWNER_ADDRESS_WITH_TOS);
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
    assertEq(s_functionsRouter.getSubscription(s_subscriptionId).proposedOwner, NEW_OWNER_ADDRESS_WITH_TOS);

    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionOwnerTransferRequested(s_subscriptionId, OWNER_ADDRESS, NEW_OWNER_ADDRESS_WITH_TOS2);
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS2);
    assertEq(s_functionsRouter.getSubscription(s_subscriptionId).proposedOwner, NEW_OWNER_ADDRESS_WITH_TOS2);
  }
}

/// @notice #acceptSubscriptionOwnerTransfer
contract FunctionsSubscriptions_AcceptSubscriptionOwnerTransfer is FunctionsSubscriptionSetup {
  uint256 internal NEW_OWNER_PRIVATE_KEY_WITH_TOS = 0x3;
  address internal NEW_OWNER_ADDRESS_WITH_TOS = vm.addr(NEW_OWNER_PRIVATE_KEY_WITH_TOS);
  uint256 internal NEW_OWNER_PRIVATE_KEY_WITHOUT_TOS = 0x4;
  address internal NEW_OWNER_ADDRESS_WITHOUT_TOS = vm.addr(NEW_OWNER_PRIVATE_KEY_WITHOUT_TOS);

  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    // Accept ToS as new owner
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITH_TOS);
    bytes32 message = s_termsOfServiceAllowList.getMessage(NEW_OWNER_ADDRESS_WITH_TOS, NEW_OWNER_ADDRESS_WITH_TOS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(NEW_OWNER_ADDRESS_WITH_TOS, NEW_OWNER_ADDRESS_WITH_TOS, r, s, v);

    vm.stopPrank();
    vm.startPrank(OWNER_ADDRESS);
  }

  function test_AcceptSubscriptionOwnerTransfer_RevertIfPaused() public {
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
    s_functionsRouter.pause();

    // Send as new owner, who has accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITH_TOS);

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.acceptSubscriptionOwnerTransfer(s_subscriptionId);
  }

  function test_AcceptSubscriptionOwnerTransfer_RevertIfNotAllowedSender() public {
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITHOUT_TOS);

    // Send as new owner, who has NOT accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITHOUT_TOS);

    vm.expectRevert(
      abi.encodeWithSelector(FunctionsRouter.SenderMustAcceptTermsOfService.selector, NEW_OWNER_ADDRESS_WITHOUT_TOS)
    );
    s_functionsRouter.acceptSubscriptionOwnerTransfer(s_subscriptionId);
  }

  function test_AcceptSubscriptionOwnerTransfer_RevertIfSenderBecomesBlocked() public {
    // Propose an address that is allowed to accept ownership
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);
    bool hasAccess = s_termsOfServiceAllowList.hasAccess(NEW_OWNER_ADDRESS_WITH_TOS, new bytes(0));
    assertEq(hasAccess, true);

    // Revoke access
    s_termsOfServiceAllowList.blockSender(NEW_OWNER_ADDRESS_WITH_TOS);

    // Send as blocked address
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITH_TOS);

    vm.expectRevert(
      abi.encodeWithSelector(FunctionsRouter.SenderMustAcceptTermsOfService.selector, NEW_OWNER_ADDRESS_WITH_TOS)
    );
    s_functionsRouter.acceptSubscriptionOwnerTransfer(s_subscriptionId);
  }

  function test_AcceptSubscriptionOwnerTransfer_RevertIfSenderIsNotNewOwner() public {
    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, STRANGER_ADDRESS);

    // Send as someone who is not hte proposed new owner
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITH_TOS);

    vm.expectRevert(abi.encodeWithSelector(FunctionsSubscriptions.MustBeProposedOwner.selector, STRANGER_ADDRESS));
    s_functionsRouter.acceptSubscriptionOwnerTransfer(s_subscriptionId);
  }

  event SubscriptionOwnerTransferred(uint64 indexed subscriptionId, address from, address to);

  function test_AcceptSubscriptionOwnerTransfer_Success() public {
    // Can transfer ownership with a pending request
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets;
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);
    s_functionsClient.sendRequest(s_donId, sourceCode, secrets, args, bytesArgs, s_subscriptionId, 5500);

    s_functionsRouter.proposeSubscriptionOwnerTransfer(s_subscriptionId, NEW_OWNER_ADDRESS_WITH_TOS);

    // Send as new owner, who has accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(NEW_OWNER_ADDRESS_WITH_TOS);

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionOwnerTransferred(s_subscriptionId, OWNER_ADDRESS, NEW_OWNER_ADDRESS_WITH_TOS);

    s_functionsRouter.acceptSubscriptionOwnerTransfer(s_subscriptionId);

    FunctionsSubscriptions.Subscription memory subscription = s_functionsRouter.getSubscription(s_subscriptionId);
    assertEq(subscription.owner, NEW_OWNER_ADDRESS_WITH_TOS);
    assertEq(subscription.proposedOwner, address(0));
  }
}

/// @notice #removeConsumer
contract FunctionsSubscriptions_RemoveConsumer is FunctionsSubscriptionSetup {
  function test_RemoveConsumer_RevertIfPaused() public {
    s_functionsRouter.pause();

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.removeConsumer(s_subscriptionId, address(s_functionsClient));
  }

  function test_RemoveConsumer_RevertIfNoSubscription() public {
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    uint64 invalidSubscriptionId = 123456789;
    s_functionsRouter.removeConsumer(invalidSubscriptionId, address(s_functionsClient));
  }

  function test_RemoveConsumer_RevertIfNotSubscriptionOwner() public {
    // Accept ToS as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);
    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, STRANGER_ADDRESS, r, s, v);

    // Send as non-subscription owner, who has accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(FunctionsSubscriptions.MustBeSubscriptionOwner.selector);
    s_functionsRouter.removeConsumer(s_subscriptionId, address(s_functionsClient));
  }

  function test_RemoveConsumer_RevertIfNotAllowedSender() public {
    // Remove owner from Allow List
    s_termsOfServiceAllowList.blockSender(OWNER_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.SenderMustAcceptTermsOfService.selector, OWNER_ADDRESS));
    s_functionsRouter.removeConsumer(s_subscriptionId, address(s_functionsClient));
  }

  function test_RemoveConsumer_RevertIfInvalidConsumer() public {
    vm.expectRevert(FunctionsSubscriptions.InvalidConsumer.selector);
    s_functionsRouter.removeConsumer(s_subscriptionId, address(0));
  }

  function test_RemoveConsumer_RevertIfPendingRequests() public {
    // Send a minimal request
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets;
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);

    s_functionsClient.sendRequest(s_donId, sourceCode, secrets, args, bytesArgs, s_subscriptionId, 5000);

    vm.expectRevert(FunctionsSubscriptions.CannotRemoveWithPendingRequests.selector);
    s_functionsRouter.removeConsumer(s_subscriptionId, address(s_functionsClient));
  }

  event SubscriptionConsumerRemoved(uint64 indexed subscriptionId, address consumer);

  function test_RemoveConsumer_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionConsumerRemoved(s_subscriptionId, address(s_functionsClient));
    s_functionsRouter.removeConsumer(s_subscriptionId, address(s_functionsClient));

    FunctionsSubscriptions.Subscription memory subscription = s_functionsRouter.getSubscription(s_subscriptionId);
    assertEq(subscription.consumers, new address[](0));
  }
}

/// @notice #_getMaxConsumers
contract FunctionsSubscriptions__GetMaxConsumers is FunctionsRouterSetup {
  // TODO: make contract internal function helper
}

/// @notice #addConsumer
contract FunctionsSubscriptions_AddConsumer is FunctionsSubscriptionSetup {
  function test_AddConsumer_RevertIfPaused() public {
    s_functionsRouter.pause();

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.addConsumer(s_subscriptionId, address(1));
  }

  function test_AddConsumer_RevertIfNoSubscription() public {
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    uint64 invalidSubscriptionId = 123456789;
    s_functionsRouter.addConsumer(invalidSubscriptionId, address(1));
  }

  function test_AddConsumer_RevertIfNotSubscriptionOwner() public {
    // Accept ToS as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);
    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, STRANGER_ADDRESS, r, s, v);

    // Send as non-subscription owner, who has accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(FunctionsSubscriptions.MustBeSubscriptionOwner.selector);
    s_functionsRouter.addConsumer(s_subscriptionId, address(1));
  }

  function test_AddConsumer_RevertIfNotAllowedSender() public {
    // Remove owner from Allow List
    s_termsOfServiceAllowList.blockSender(OWNER_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.SenderMustAcceptTermsOfService.selector, OWNER_ADDRESS));
    s_functionsRouter.addConsumer(s_subscriptionId, address(1));
  }

  function test_AddConsumer_RevertIfMaximumConsumers() public {
    // Fill Consumers to s_maxConsumersPerSubscription
    // Already has one from setup
    s_functionsRouter.addConsumer(s_subscriptionId, address(1));
    s_functionsRouter.addConsumer(s_subscriptionId, address(2));

    vm.expectRevert(
      abi.encodeWithSelector(FunctionsSubscriptions.TooManyConsumers.selector, s_maxConsumersPerSubscription)
    );
    s_functionsRouter.addConsumer(s_subscriptionId, address(3));
  }

  function test_AddConsumer_RevertIfMaximumConsumersAfterConfigUpdate() public {
    // Fill Consumers to s_maxConsumersPerSubscription
    // Already has one from setup
    s_functionsRouter.addConsumer(s_subscriptionId, address(1));
    s_functionsRouter.addConsumer(s_subscriptionId, address(2));

    // Lower maxConsumersPerSubscription
    s_maxConsumersPerSubscription = 1;
    FunctionsRouter.Config memory newRouterConfig = getRouterConfig();
    s_functionsRouter.updateConfig(newRouterConfig);

    // .AddConsumer should still revert
    vm.expectRevert(
      abi.encodeWithSelector(FunctionsSubscriptions.TooManyConsumers.selector, s_maxConsumersPerSubscription)
    );
    s_functionsRouter.addConsumer(s_subscriptionId, address(3));
  }

  event SubscriptionConsumerAdded(uint64 indexed subscriptionId, address consumer);

  function test_AddConsumer_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionConsumerAdded(s_subscriptionId, address(1));
    s_functionsRouter.addConsumer(s_subscriptionId, address(1));

    FunctionsSubscriptions.Subscription memory subscription = s_functionsRouter.getSubscription(s_subscriptionId);
    assertEq(subscription.consumers[1], address(1));
    FunctionsSubscriptions.Consumer memory consumer = s_functionsRouter.getConsumer(address(1), s_subscriptionId);
    assertEq(consumer.allowed, true);
  }
}

/// @notice #cancelSubscription
contract FunctionsSubscriptions_CancelSubscription is FunctionsSubscriptionSetup {
  function test_CancelSubscription_RevertIfPaused() public {
    s_functionsRouter.pause();

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.cancelSubscription(s_subscriptionId, OWNER_ADDRESS);
  }

  function test_CancelSubscription_RevertIfNoSubscription() public {
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    uint64 invalidSubscriptionId = 123456789;
    s_functionsRouter.cancelSubscription(invalidSubscriptionId, OWNER_ADDRESS);
  }

  function test_CancelSubscription_RevertIfNotSubscriptionOwner() public {
    // Accept ToS as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);
    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, STRANGER_ADDRESS, r, s, v);

    // Send as non-subscription owner, who has accepted Terms of Service
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(FunctionsSubscriptions.MustBeSubscriptionOwner.selector);
    s_functionsRouter.cancelSubscription(s_subscriptionId, OWNER_ADDRESS);
  }

  function test_CancelSubscription_RevertIfNotAllowedSender() public {
    // Remove owner from Allow List
    s_termsOfServiceAllowList.blockSender(OWNER_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.SenderMustAcceptTermsOfService.selector, OWNER_ADDRESS));
    s_functionsRouter.cancelSubscription(s_subscriptionId, OWNER_ADDRESS);
  }

  function test_CancelSubscription_RevertIfPendingRequests() public {
    // Send a minimal request
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets;
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);

    s_functionsClient.sendRequest(s_donId, sourceCode, secrets, args, bytesArgs, s_subscriptionId, 5000);

    vm.expectRevert(FunctionsSubscriptions.CannotRemoveWithPendingRequests.selector);
    s_functionsRouter.cancelSubscription(s_subscriptionId, OWNER_ADDRESS);
  }

  event SubscriptionCanceled(uint64 indexed subscriptionId, address fundsRecipient, uint256 fundsAmount);

  function test_CancelSubscription_SuccessForfeitAllBalanceAsDeposit() public {
    // No requests have been completed
    assertEq(s_functionsRouter.getConsumer(address(s_functionsClient), s_subscriptionId).completedRequests, 0);
    // Subscription balance is less than deposit amount
    assertLe(s_functionsRouter.getSubscription(s_subscriptionId).balance, s_subscriptionDepositJuels);

    uint256 subscriptionOwnerBalanceBefore = s_linkToken.balanceOf(OWNER_ADDRESS);

    uint96 expectedRefund = 0;

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCanceled(s_subscriptionId, OWNER_ADDRESS, expectedRefund);

    s_functionsRouter.cancelSubscription(s_subscriptionId, OWNER_ADDRESS);

    uint256 subscriptionOwnerBalanceAfter = s_linkToken.balanceOf(OWNER_ADDRESS);
    assertEq(subscriptionOwnerBalanceBefore + expectedRefund, subscriptionOwnerBalanceAfter);

    // Subscription should no longer exist
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    s_functionsRouter.getSubscription(s_subscriptionId);

    // Router owner should have expectedDepositWithheld to withdraw
    uint96 expectedDepositWithheld = s_subscriptionInitialFunding;
    uint256 balanceBeforeWithdraw = s_linkToken.balanceOf(STRANGER_ADDRESS);
    s_functionsRouter.ownerWithdraw(STRANGER_ADDRESS, 0);
    uint256 balanceAfterWithdraw = s_linkToken.balanceOf(STRANGER_ADDRESS);
    assertEq(balanceBeforeWithdraw + expectedDepositWithheld, balanceAfterWithdraw);
  }

  function test_CancelSubscription_SuccessForfeitSomeBalanceAsDeposit() public {
    // No requests have been completed
    assertEq(s_functionsRouter.getConsumer(address(s_functionsClient), s_subscriptionId).completedRequests, 0);
    // Subscription balance is more than deposit amount, double fund the subscription
    s_linkToken.transferAndCall(address(s_functionsRouter), s_subscriptionInitialFunding, abi.encode(s_subscriptionId));
    assertGe(s_functionsRouter.getSubscription(s_subscriptionId).balance, s_subscriptionDepositJuels);

    uint256 subscriptionOwnerBalanceBefore = s_linkToken.balanceOf(OWNER_ADDRESS);

    uint96 expectedRefund = (s_subscriptionInitialFunding * 2) - s_subscriptionDepositJuels;
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCanceled(s_subscriptionId, OWNER_ADDRESS, expectedRefund);

    s_functionsRouter.cancelSubscription(s_subscriptionId, OWNER_ADDRESS);

    uint256 subscriptionOwnerBalanceAfter = s_linkToken.balanceOf(OWNER_ADDRESS);
    assertEq(subscriptionOwnerBalanceBefore + expectedRefund, subscriptionOwnerBalanceAfter);

    // Subscription should no longer exist
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    s_functionsRouter.getSubscription(s_subscriptionId);

    // Router owner should have expectedDepositWithheld to withdraw
    uint96 expectedDepositWithheld = s_subscriptionDepositJuels;
    uint256 balanceBeforeWithdraw = s_linkToken.balanceOf(STRANGER_ADDRESS);
    s_functionsRouter.ownerWithdraw(STRANGER_ADDRESS, 0);
    uint256 balanceAfterWithdraw = s_linkToken.balanceOf(STRANGER_ADDRESS);
    assertEq(balanceBeforeWithdraw + expectedDepositWithheld, balanceAfterWithdraw);
  }
}

/// @notice #cancelSubscription
contract FunctionsSubscriptions_CancelSubscription_ReceiveDeposit is FunctionsFulfillmentSetup {
  event SubscriptionCanceled(uint64 indexed subscriptionId, address fundsRecipient, uint256 fundsAmount);

  function test_CancelSubscription_SuccessRecieveDeposit() public {
    uint96 totalCostJuels = s_fulfillmentRouterOwnerBalance + s_fulfillmentCoordinatorBalance;

    uint256 subscriptionOwnerBalanceBefore = s_linkToken.balanceOf(OWNER_ADDRESS);

    uint96 expectedRefund = s_subscriptionInitialFunding - totalCostJuels;
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit SubscriptionCanceled(s_subscriptionId, OWNER_ADDRESS, expectedRefund);

    s_functionsRouter.cancelSubscription(s_subscriptionId, OWNER_ADDRESS);

    uint256 subscriptionOwnerBalanceAfter = s_linkToken.balanceOf(OWNER_ADDRESS);
    assertEq(subscriptionOwnerBalanceBefore + expectedRefund, subscriptionOwnerBalanceAfter);

    // Subscription should no longer exist
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    s_functionsRouter.getSubscription(s_subscriptionId);
  }
}

/// @notice #_cancelSubscriptionHelper
contract FunctionsSubscriptions__CancelSubscriptionHelper {
  // TODO: make contract internal function helper
}

/// @notice #pendingRequestExists
contract FunctionsSubscriptions_PendingRequestExists is FunctionsFulfillmentSetup {
  function test_PendingRequestExists_SuccessFalseIfNoPendingRequests() public {
    bool hasPendingRequests = s_functionsRouter.pendingRequestExists(s_subscriptionId);
    assertEq(hasPendingRequests, false);
  }

  function test_PendingRequestExists_SuccessTrueIfPendingRequests() public {
    // Send a minimal request
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets;
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);

    s_functionsClient.sendRequest(s_donId, sourceCode, secrets, args, bytesArgs, s_subscriptionId, 5000);

    bool hasPendingRequests = s_functionsRouter.pendingRequestExists(s_subscriptionId);
    assertEq(hasPendingRequests, true);
  }
}

/// @notice #setFlags
contract FunctionsSubscriptions_SetFlags is FunctionsSubscriptionSetup {
  function test_SetFlags_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    bytes32 flagsToSet = bytes32("1");
    s_functionsRouter.setFlags(s_subscriptionId, flagsToSet);
  }

  function test_SetFlags_RevertIfNoSubscription() public {
    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    uint64 invalidSubscriptionId = 123456789;
    bytes32 flagsToSet = bytes32("1");
    s_functionsRouter.setFlags(invalidSubscriptionId, flagsToSet);
  }

  function test_SetFlags_Success() public {
    bytes32 flagsToSet = bytes32("1");
    s_functionsRouter.setFlags(s_subscriptionId, flagsToSet);
    bytes32 flags = s_functionsRouter.getFlags(s_subscriptionId);
    assertEq(flags, flagsToSet);
  }
}

/// @notice #getFlags
contract FunctionsSubscriptions_GetFlags is FunctionsSubscriptionSetup {
  function test_GetFlags_SuccessInvalidSubscription() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint64 invalidSubscriptionId = 999999;

    bytes32 flags = s_functionsRouter.getFlags(invalidSubscriptionId);
    assertEq(flags, bytes32(0));
  }

  function test_GetFlags_SuccessValidSubscription() public {
    // Set flags
    bytes32 flagsToSet = bytes32("1");
    s_functionsRouter.setFlags(s_subscriptionId, flagsToSet);

    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 flags = s_functionsRouter.getFlags(s_subscriptionId);
    assertEq(flags, flagsToSet);
  }
}

/// @notice #timeoutRequests
contract FunctionsSubscriptions_TimeoutRequests is FunctionsClientRequestSetup {
  function test_TimeoutRequests_RevertIfPaused() public {
    s_functionsRouter.pause();

    vm.expectRevert("Pausable: paused");
    FunctionsResponse.Commitment[] memory commitments = new FunctionsResponse.Commitment[](1);
    commitments[0] = s_requests[1].commitmentOnchain;
    s_functionsRouter.timeoutRequests(commitments);
  }

  function test_TimeoutRequests_RevertInvalidRequest() public {
    // Modify the commitment so that it doesn't match
    s_requests[1].commitmentOnchain.donFee = 123456789;
    FunctionsResponse.Commitment[] memory commitments = new FunctionsResponse.Commitment[](1);
    commitments[0] = s_requests[1].commitmentOnchain;
    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);
    s_functionsRouter.timeoutRequests(commitments);
  }

  function test_TimeoutRequests_RevertIfTimeoutNotExceeded() public {
    vm.expectRevert(FunctionsSubscriptions.TimeoutNotExceeded.selector);
    FunctionsResponse.Commitment[] memory commitments = new FunctionsResponse.Commitment[](1);
    commitments[0] = s_requests[1].commitmentOnchain;
    s_functionsRouter.timeoutRequests(commitments);
  }

  event RequestTimedOut(bytes32 indexed requestId);

  function test_TimeoutRequests_Success() public {
    uint64 consumerCompletedRequestsBefore = s_functionsRouter
      .getConsumer(address(s_functionsClient), s_subscriptionId)
      .completedRequests;

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit RequestTimedOut(s_requests[1].requestId);

    // Jump ahead in time past timeout timestamp
    vm.warp(s_requests[1].commitmentOnchain.timeoutTimestamp + 1);

    FunctionsResponse.Commitment[] memory commitments = new FunctionsResponse.Commitment[](1);
    commitments[0] = s_requests[1].commitmentOnchain;
    s_functionsRouter.timeoutRequests(commitments);

    // Releases blocked balance and increments completed requests
    uint96 subscriptionBlockedBalanceAfter = s_functionsRouter.getSubscription(s_subscriptionId).blockedBalance;
    assertEq(0, subscriptionBlockedBalanceAfter);
    uint64 consumerCompletedRequestsAfter = s_functionsRouter
      .getConsumer(address(s_functionsClient), s_subscriptionId)
      .completedRequests;
    assertEq(consumerCompletedRequestsBefore + 1, consumerCompletedRequestsAfter);
  }
}

// @notice #_onlySubscriptionOwner
contract FunctionsSubscriptions__OnlySubscriptionOwner {
  // TODO: make contract internal function helper
}
