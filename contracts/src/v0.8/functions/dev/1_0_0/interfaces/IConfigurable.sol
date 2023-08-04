// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// @title Configurable contract interface.
interface IConfigurable {
  // @notice Set the contract's configuration
  // @param config bytes containing config data
  function updateConfig(bytes calldata config) external;
}
