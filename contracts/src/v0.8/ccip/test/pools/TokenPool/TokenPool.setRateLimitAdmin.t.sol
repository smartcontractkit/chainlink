// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_setRateLimitAdmin is TokenPoolSetup {
  function test_SetRateLimitAdmin() public {
    assertEq(address(0), s_tokenPool.getRateLimitAdmin());
    vm.expectEmit();
    emit TokenPool.RateLimitAdminSet(OWNER);
    s_tokenPool.setRateLimitAdmin(OWNER);
    assertEq(OWNER, s_tokenPool.getRateLimitAdmin());
  }

  // Reverts

  function test_RevertWhen_SetRateLimitAdmin() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_tokenPool.setRateLimitAdmin(STRANGER);
  }
}
