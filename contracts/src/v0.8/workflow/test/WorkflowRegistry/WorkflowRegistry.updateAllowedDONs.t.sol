// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_updateAllowedDONs is WorkflowRegistrySetup {
  function test_RevertWhen_TheCallerIsNotTheOwner() external {
    vm.prank(s_stranger);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registry.updateAllowedDONs(new uint32[](0), true);
  }

  // whenTheCallerIsTheOwner
  function test_RevertWhen_TheRegistryIsLocked() external {
    // Lock the registry as the owner
    vm.startPrank(s_owner);
    s_registry.lockRegistry();

    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.updateAllowedDONs(new uint32[](0), true);
    vm.stopPrank();
  }

  // whenTheCallerIsTheOwner whenTheRegistryIsNotLocked
  function test_WhenTheBoolInputIsTrue() external {
    uint32[] memory donIDsToAdd = new uint32[](3);
    donIDsToAdd[0] = 2;
    donIDsToAdd[1] = 3;
    donIDsToAdd[2] = 4;

    // Check that there is one DON ID when fetching all allowed DONs to start
    uint32[] memory allowedDONs = s_registry.getAllAllowedDONs();
    assertEq(allowedDONs.length, 1);

    // Expect the event to be emitted
    vm.expectEmit();
    emit WorkflowRegistry.AllowedDONsUpdatedV1(donIDsToAdd, true);

    // Call the function as the owner
    vm.prank(s_owner);
    s_registry.updateAllowedDONs(donIDsToAdd, true);

    // Verify that the DON IDs have been added
    allowedDONs = s_registry.getAllAllowedDONs();
    assertEq(allowedDONs.length, 4);
  }

  // whenTheCallerIsTheOwner whenTheRegistryIsNotLocked
  function test_WhenTheBoolInputIsFalse() external {
    uint32[] memory donIDsToRemove = new uint32[](1);
    donIDsToRemove[0] = s_allowedDonID;

    // Check that there is one DON ID when fetching all allowed DONs to start
    uint32[] memory allowedDONs = s_registry.getAllAllowedDONs();
    assertEq(allowedDONs.length, 1);

    // Expect the event to be emitted
    vm.expectEmit();
    emit WorkflowRegistry.AllowedDONsUpdatedV1(donIDsToRemove, false);

    // Call the function as the owner
    vm.prank(s_owner);
    s_registry.updateAllowedDONs(donIDsToRemove, false);

    // Verify that the DON IDs have been removed
    allowedDONs = s_registry.getAllAllowedDONs();
    assertEq(allowedDONs.length, 0);
  }
}
