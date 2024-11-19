// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

contract WorkflowRegistryManageraddVersion {
  function test_RevertWhen_TheCallerIsNotTheOwner() external {
    // it should revert
  }

  modifier whenTheCallerIsTheOwner() {
    _;
  }

  function test_RevertWhen_TheContractAddressIsInvalid() external whenTheCallerIsTheOwner {
    // it should revert
  }

  modifier whenTheContractAddressIsValid() {
    _;
  }

  function test_WhenAutoActivateIsTrue() external whenTheCallerIsTheOwner whenTheContractAddressIsValid {
    // it should deactivate any currently active version
    // it should activate the new version
    // it should emit VersionAddedV1 after adding the version to s_versions
    // it should emit VersionActivatedV1
  }

  function test_WhenAutoActivateIsFalse() external whenTheCallerIsTheOwner whenTheContractAddressIsValid {
    // it should not activate the new version
    // it should emit VersionAddedV1 after adding the version to s_versions
  }
}
