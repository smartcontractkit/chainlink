// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

/**
 * @notice ArbitrumGasInterface provides an interface to estimate arbitrum L1 fee in wei
 */
interface ArbitrumGasInterface {
  // @notice Get the L1 gas fee paid by the current transaction in wei
  function getCurrentTxL1GasFees() external view returns (uint256 l1CostWei);
}
