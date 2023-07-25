// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title Configurable contract interface.
 */
interface IConfigurable {
  /**
   * @notice Get the hash of the current configuration
   * @return config hash of config bytes
   */
  function getConfigHash() external returns (bytes32 config);

  /**
   * @notice Set the contract's configuration
   * @dev Only callable by the Router
   * @param config bytes containing config data
   */
  function setConfig(bytes calldata config) external;
}
