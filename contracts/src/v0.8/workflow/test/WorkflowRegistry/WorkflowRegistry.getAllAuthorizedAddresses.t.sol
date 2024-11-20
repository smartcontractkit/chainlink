// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistrygetAllAuthorizedAddresses is WorkflowRegistrySetup {
  function test_WhenTheSetOfAuthorizedAddressesIsEmpty() external {
    // Remove the authorized address added in the setup
    _removeAddressFromAuthorizedAddresses(s_authorizedAddress);
    address[] memory authorizedAddresses = s_registry.getAllAuthorizedAddresses();
    assertEq(authorizedAddresses.length, 0);
  }

  function test_WhenThereIsASingleAuthorizedAddress() external view {
    // it should return an array with one element
    address[] memory authorizedAddresses = s_registry.getAllAuthorizedAddresses();
    assertEq(authorizedAddresses.length, 1);
    assertEq(authorizedAddresses[0], s_authorizedAddress);
  }

  function test_WhenThereAreMultipleAuthorizedAddresses() external {
    // Add a second authorized address
    _addAddressToAuthorizedAddresses(s_unauthorizedAddress);

    // it should return an array with all the authorized addresses
    address[] memory authorizedAddresses = s_registry.getAllAuthorizedAddresses();
    assertEq(authorizedAddresses.length, 2);
    assertEq(authorizedAddresses[0], s_authorizedAddress);
    assertEq(authorizedAddresses[1], s_unauthorizedAddress);
  }
}
