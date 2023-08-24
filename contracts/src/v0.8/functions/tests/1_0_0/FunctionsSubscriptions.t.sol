// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsRouter} from "../../dev/1_0_0/FunctionsRouter.sol";
import {FunctionsSubscriptions} from "../../dev/1_0_0/FunctionsSubscriptions.sol";

import {FunctionsRouterSetup, FunctionsOwnerAcceptTermsOfService} from "./Setup.t.sol";

// ================================================================
// |                    Functions Subscriptions                   |
// ================================================================

/// @notice #constructor
contract FunctionsSubscriptions_Constructor {

}

/// @notice #_markRequestInFlight
contract FunctionsSubscriptions__MarkRequestInFlight {

}

/// @notice #_pay
contract FunctionsSubscriptions__Pay {

}

/// @notice #ownerCancelSubscription
contract FunctionsSubscriptions_OwnerCancelSubscription {

}

/// @notice #recoverFunds
contract FunctionsSubscriptions_RecoverFunds {

}

/// @notice #oracleWithdraw
contract FunctionsSubscriptions_OracleWithdraw {

}

/// @notice #ownerWithdraw
contract FunctionsSubscriptions_OwnerWithdraw {

}

/// @notice #onTokenTransfer
contract FunctionsSubscriptions_OnTokenTransfer {

}

/// @notice #getTotalBalance
contract FunctionsSubscriptions_GetTotalBalance {

}

/// @notice #getSubscriptionCount
contract FunctionsSubscriptions_GetSubscriptionCount {

}

/// @notice #getSubscription
contract FunctionsSubscriptions_GetSubscription {

}

/// @notice #getConsumer
contract FunctionsSubscriptions_GetConsumer {

}

/// @notice #_isExistingSubscription
contract FunctionsSubscriptions__IsExistingSubscription {

}

/// @notice #_isAllowedConsumer
contract FunctionsSubscriptions__IsAllowedConsumer {

}

/// @notice #createSubscription
contract FunctionsSubscriptions_createSubscription is FunctionsOwnerAcceptTermsOfService {
  function setUp() public virtual override {
    FunctionsOwnerAcceptTermsOfService.setUp();
  }

  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);

  function test_CreateSubscription_Success() public {
    vm.expectEmit(true, false, false, true);
    emit SubscriptionCreated(1, OWNER_ADDRESS);
    uint64 firstCallSubscriptionId = s_functionsRouter.createSubscription();
    assertEq(firstCallSubscriptionId, 1);

    vm.expectEmit(true, false, false, true);
    emit SubscriptionCreated(2, OWNER_ADDRESS);
    uint64 secondCallSubscriptionId = s_functionsRouter.createSubscription();
    assertEq(secondCallSubscriptionId, 2);

    vm.expectEmit(true, false, false, true);
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
contract FunctionsSubscriptions_CreateSubscriptionWithConsumer {

}

/// @notice #proposeSubscriptionOwnerTransfer
contract FunctionsSubscriptions_ProposeSubscriptionOwnerTransfer {

}

/// @notice #acceptSubscriptionOwnerTransfer
contract FunctionsSubscriptions_AcceptSubscriptionOwnerTransfer {

}

/// @notice #removeConsumer
contract FunctionsSubscriptions_RemoveConsumer {

}

/// @notice #_getMaxConsumers
contract FunctionsSubscriptions__GetMaxConsumers {

}

/// @notice #addConsumer
contract FunctionsSubscriptions_AddConsumer {

}

/// @notice #cancelSubscription
contract FunctionsSubscriptions_CancelSubscription {

}

/// @notice #_cancelSubscriptionHelper
contract FunctionsSubscriptions__CancelSubscriptionHelper {

}

/// @notice #pendingRequestExists
contract FunctionsSubscriptions_PendingRequestExists {

}

/// @notice #setFlags
contract FunctionsSubscriptions_SetFlags {

}

/// @notice #getFlags
contract FunctionsSubscriptions_GetFlags {

}

/// @notice #timeoutRequests
contract FunctionsSubscriptions_TimeoutRequests {

}

/// @notice #_onlySubscriptionOwner
contract FunctionsSubscriptions__OnlySubscriptionOwner {

}
