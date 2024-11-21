// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_getAllAllowedDONs is WorkflowRegistrySetup {
  function test_WhenTheSetOfAllowedDONsIsEmpty() external {
    // Remove the allowed DON added in the setup
    _removeDONFromAllowedDONs(s_allowedDonID);
    uint32[] memory allowedDONs = s_registry.getAllAllowedDONs();
    assertEq(allowedDONs.length, 0);
  }

  function test_WhenThereIsASingleAllowedDON() external view {
    uint32[] memory allowedDONs = s_registry.getAllAllowedDONs();
    assertEq(allowedDONs.length, 1);
    assertEq(allowedDONs[0], s_allowedDonID);
  }

  function test_WhenThereAreMultipleAllowedDONs() external {
    // Add a second DON to the allowed DONs list
    uint32 allowedDonID2 = 2;
    uint32[] memory donIDsToAdd = new uint32[](1);
    donIDsToAdd[0] = allowedDonID2;

    vm.prank(s_owner);
    s_registry.updateAllowedDONs(donIDsToAdd, true);

    uint32[] memory allowedDONs = s_registry.getAllAllowedDONs();
    assertEq(allowedDONs.length, 2);
    assertEq(allowedDONs[0], s_allowedDonID);
    assertEq(allowedDONs[1], allowedDonID2);
  }

  function test_WhenTheRegistryIsLocked() external {
    // Lock the registry
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // It should behave the same as when the registry is not locked
    vm.prank(s_stranger);
    uint32[] memory allowedDONs = s_registry.getAllAllowedDONs();
    assertEq(allowedDONs.length, 1);
    assertEq(allowedDONs[0], s_allowedDonID);
  }
}
