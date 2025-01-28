// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_lockRegistry is WorkflowRegistrySetup {
  function test_RevertWhen_TheCallerIsNotTheContractOwner() external {
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registry.lockRegistry();
  }

  function test_WhenTheCallerIsTheContractOwner() external {
    vm.expectEmit(true, true, false, false);
    emit WorkflowRegistry.RegistryLockedV1(s_owner);

    vm.prank(s_owner);
    s_registry.lockRegistry();

    assertTrue(s_registry.isRegistryLocked());
  }
}
