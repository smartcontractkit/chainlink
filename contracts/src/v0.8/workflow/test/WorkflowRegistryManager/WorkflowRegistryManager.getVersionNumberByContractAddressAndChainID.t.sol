// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";
import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";

contract WorkflowRegistryManager_getVersionNumberByContractAddressAndChainID is WorkflowRegistryManagerSetup {
  function test_WhenTheContractAddressIsInvalid() external {
    // It should revert with InvalidContractAddress
    _deployMockRegistryAndAddVersion(true);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistryManager.InvalidContractAddress.selector, address(0)));
    s_registryManager.getVersionNumberByContractAddressAndChainID(address(0), s_chainID);
  }

  // whenTheContractAddressIsValid
  function test_WhenNoVersionIsRegisteredForTheContractAddressAndChainIDCombination() external {
    // It should revert with NoVersionsRegistered.
    _deployMockRegistryAndAddVersion(true);
    vm.expectRevert(WorkflowRegistryManager.NoVersionsRegistered.selector);
    s_registryManager.getVersionNumberByContractAddressAndChainID(address(s_mockWorkflowRegistryContract), 20);
  }

  // whenTheContractAddressIsValid
  function test_WhenAVersionIsRegisteredForTheContractAddressAndChainIDCombination() external {
    // It should return the correct version number.
    _deployMockRegistryAndAddVersion(true);
    uint32 versionNumber =
      s_registryManager.getVersionNumberByContractAddressAndChainID(address(s_mockWorkflowRegistryContract), s_chainID);
    assertEq(versionNumber, 1);
  }
}
