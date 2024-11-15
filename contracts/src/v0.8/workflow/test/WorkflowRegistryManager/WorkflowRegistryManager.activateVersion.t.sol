// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";

contract WorkflowRegistryManager_activateVersion is WorkflowRegistryManagerSetup {
  function test_RevertWhen_TheCallerIsNotTheOwner() external {
    // it should revert
    vm.prank(s_nonOwner);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registryManager.activateVersion(s_versionNumber);
  }

  // whenTheCallerIsTheOwner
  function test_RevertWhen_TheVersionNumberDoesNotExist() external {
    // it should revert
  }

  // whenTheCallerIsTheOwner whenTheVersionNumberExists
  function test_RevertWhen_TheVersionNumberIsAlreadyActive() external {
    // it should revert
  }

  function test_WhenTheVersionNumberIsNotActive() external {
    // it should deactivate the current active version (if any)
    // it should activate the specified version and update s_activeVersionNumber
    // it should add the version to s_versionNumberByAddressAndChainID
    // it should emit VersionDeactivatedV1 (if a previous version was active)
    // it should emit VersionActivatedV1
  }
}
