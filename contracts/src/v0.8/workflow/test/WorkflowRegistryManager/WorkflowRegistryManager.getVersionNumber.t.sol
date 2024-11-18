// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

contract WorkflowRegistryManagergetVersionNumber {
  function test_WhenTheContractAddressIsInvalid() external {
    // it should revert with InvalidContractAddress
  }

  modifier whenTheContractAddressIsValid() {
    _;
  }

  function test_WhenNoVersionIsRegisteredForTheContractAddressAndChainIDCombination()
    external
    whenTheContractAddressIsValid
  {
    // it should revert with NoVersionsRegistered
  }

  function test_WhenAVersionIsRegisteredForTheContractAddressAndChainIDCombination()
    external
    whenTheContractAddressIsValid
  {
    // it should return the correct version number
  }
}
