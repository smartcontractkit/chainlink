// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_updateAuthorizedAddresses is WorkflowRegistrySetup {
  function test_RevertWhen_TheCallerIsNotTheOwner() external {
    vm.prank(s_stranger);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registry.updateAuthorizedAddresses(new address[](0), true);
  }

  // whenTheCallerIsTheOwner
  function test_RevertWhen_TheRegistryIsLocked() external {
    // Lock the registry as the owner
    vm.startPrank(s_owner);
    s_registry.lockRegistry();

    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.updateAuthorizedAddresses(new address[](0), true);
    vm.stopPrank();
  }

  // whenTheCallerIsTheOwner whenTheRegistryIsNotLocked
  function test_WhenTheBoolInputIsTrue() external {
    address[] memory addressesToAdd = new address[](3);
    addressesToAdd[0] = makeAddr("1");
    addressesToAdd[1] = makeAddr("2");
    addressesToAdd[2] = makeAddr("3");

    // Check that there is one authorized address when fetching all authorized addresses to start
    address[] memory authorizedAddresses = s_registry.getAllAuthorizedAddresses();
    assertEq(authorizedAddresses.length, 1);

    // Expect the event to be emitted
    vm.expectEmit();
    emit WorkflowRegistry.AuthorizedAddressesUpdatedV1(addressesToAdd, true);

    // Call the function as the owner
    vm.prank(s_owner);
    s_registry.updateAuthorizedAddresses(addressesToAdd, true);

    // Verify that the addresses have been added
    authorizedAddresses = s_registry.getAllAuthorizedAddresses();
    assertEq(authorizedAddresses.length, 4);
  }

  // whenTheCallerIsTheOwner whenTheRegistryIsNotLocked
  function test_WhenTheBoolInputIsFalse() external {
    address[] memory addressesToRemove = new address[](1);
    addressesToRemove[0] = s_authorizedAddress;

    // Check that there is one authorized address when fetching all authorized addresses to start
    address[] memory authorizedAddresses = s_registry.getAllAuthorizedAddresses();
    assertEq(authorizedAddresses.length, 1);

    // Expect the event to be emitted
    vm.expectEmit();
    emit WorkflowRegistry.AuthorizedAddressesUpdatedV1(addressesToRemove, false);

    // Call the function as the owner
    vm.prank(s_owner);
    s_registry.updateAuthorizedAddresses(addressesToRemove, false);

    // Verify that the addresses have been removed
    authorizedAddresses = s_registry.getAllAuthorizedAddresses();
    assertEq(authorizedAddresses.length, 0);
  }
}
