// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_unlockRegistry is WorkflowRegistrySetup {
  function test_RevertWhen_TheCallerIsNotTheContractOwner() external {
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registry.unlockRegistry();
  }

  function test_WhenTheCallerIsTheContractOwner() external {
    // Lock the registry first
    vm.startPrank(s_owner);
    s_registry.lockRegistry();

    assertTrue(s_registry.isRegistryLocked());

    // Unlock the registry
    vm.expectEmit(true, true, false, false);
    emit WorkflowRegistry.RegistryUnlockedV1(s_owner);

    s_registry.unlockRegistry();

    assertFalse(s_registry.isRegistryLocked());
    vm.stopPrank();
  }
}
