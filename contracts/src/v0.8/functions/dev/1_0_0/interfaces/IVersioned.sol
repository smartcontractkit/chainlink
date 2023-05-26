// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @title Chainlink versioned contract interface.
 */
interface IVersioned {
  /**
   * @notice Returns information about the contract's name and version
   * @return id Identifier for the contract
   * @return version The current version number
   */
  function idAndVersion() external view returns (string memory id, uint16 version);
}
