// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsRouter} from "../../dev/1_0_0/FunctionsRouter.sol";
import {FunctionsSubscriptions} from "../../dev/1_0_0/FunctionsSubscriptions.sol";

import {FunctionsRouterSetup, FunctionsOwnerAcceptTermsOfService} from "./Setup.t.sol";

// ================================================================
// |                        Functions Router                      |
// ================================================================

/// @notice #constructor
contract FunctionsRouter_Constructor {

}

/// @notice #getConfig
contract FunctionsRouter_GetConfig {

}

/// @notice #updateConfig
contract FunctionsRouter_UpdateConfig {

}

/// @notice #isValidCallbackGasLimit
contract FunctionsRouter_IsValidCallbackGasLimit {

}

/// @notice #getAdminFee
contract FunctionsRouter_GetAdminFee {

}

/// @notice #getAllowListId
contract FunctionsRouter_GetAllowListId {

}

/// @notice #setAllowListId
contract FunctionsRouter_SetAllowListId {

}

/// @notice #_getMaxConsumers
contract FunctionsRouter__GetMaxConsumers {

}

/// @notice #sendRequest
contract FunctionsRouter_SendRequest {

}

/// @notice #sendRequestToProposed
contract FunctionsRouter_SendRequestToProposed {

}

/// @notice #_sendRequest
contract FunctionsRouter__SendRequest {

}

/// @notice #fulfill
contract FunctionsRouter_Fulfill {

}

/// @notice #_callback
contract FunctionsRouter__Callback {

}

/// @notice #getContractById
contract FunctionsRouter_GetContractById {

}

/// @notice #getProposedContractById
contract FunctionsRouter_GetProposedContractById {

}

/// @notice #getProposedContractSet
contract FunctionsRouter_GetProposedContractSet {

}

/// @notice #proposeContractsUpdate
contract FunctionsRouter_ProposeContractsUpdate {

}

/// @notice #updateContracts
contract FunctionsRouter_UpdateContracts {

}

/// @notice #_whenNotPaused
contract FunctionsRouter__WhenNotPaused {

}

/// @notice #_onlyRouterOwner
contract FunctionsRouter__OnlyRouterOwner {

}

/// @notice #_onlySenderThatAcceptedToS
contract FunctionsRouter__OnlySenderThatAcceptedToS {

}

/// @notice #pause
contract FunctionsRouter_Pause is FunctionsRouterSetup {
  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();
  }

  event Paused(address account);

  function test_Pause_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.pause();
  }

  function test_Pause_Success() public {
    vm.expectEmit(false, false, false, true);
    emit Paused(OWNER_ADDRESS);

    s_functionsRouter.pause();

    bool isPaused = s_functionsRouter.paused();
    assertEq(isPaused, true);
  }
}

/// @notice #unpause
contract FunctionsRouter_Unpause is FunctionsRouterSetup {
  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();
    s_functionsRouter.pause();
  }

  event Unpaused(address account);

  function test_Unpause_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.unpause();
  }

  function test_Unpause_Success() public {
    vm.expectEmit(false, false, false, true);
    emit Unpaused(OWNER_ADDRESS);

    s_functionsRouter.unpause();

    bool isPaused = s_functionsRouter.paused();
    assertEq(isPaused, false);
  }
}

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
