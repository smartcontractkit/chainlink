// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

/**
 * @notice OVM_GasPriceOracle provides an interface to estimate optimism L1 fee in wei
 */
interface OVM_GasPriceOracle {
  // @notice Get the L1 gas fee paid by the current transaction in wei
  function getL1Fee(bytes memory data) external view returns (uint256 l1CostWei);
}
