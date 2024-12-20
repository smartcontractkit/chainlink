// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {LockReleaseTokenPoolSetup} from "./LockReleaseTokenPoolSetup.t.sol";

contract LockReleaseTokenPool_setRebalancer is LockReleaseTokenPoolSetup {
  function test_SetRebalancer() public {
    assertEq(address(s_lockReleaseTokenPool.getRebalancer()), OWNER);
    s_lockReleaseTokenPool.setRebalancer(STRANGER);
    assertEq(address(s_lockReleaseTokenPool.getRebalancer()), STRANGER);
  }

  function test_RevertWhen_SetRebalancer() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_lockReleaseTokenPool.setRebalancer(STRANGER);
  }
}
