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
    // topic0 (always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    vm.expectEmit(false, false, false, true);
    emit Paused(OWNER_ADDRESS);

    s_functionsRouter.pause();

    bool isPaused = s_functionsRouter.paused();
    assertEq(isPaused, true);

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.createSubscription();
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
    // topic0 (always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    vm.expectEmit(false, false, false, true);
    emit Unpaused(OWNER_ADDRESS);

    s_functionsRouter.unpause();

    bool isPaused = s_functionsRouter.paused();
    assertEq(isPaused, false);

    s_functionsRouter.createSubscription();
  }
}
