// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

/**
 * @notice ArbitrumGasInterface provides an interface to estimate arbitrum L1 fee in wei
 */
interface ArbitrumGasInterface {
  // @notice Get the L1 gas fee paid by the current transaction in wei
  function getCurrentTxL1GasFees() external view returns (uint256 l1CostWei);

  // @notice Get gas prices. Uses the caller's preferred aggregator, or the default if the caller doesn't have a preferred one.
  // @return gas price in wei for L2 transaction
  // @return gas price in wei for L1 calldata
  // @return gas price in wei for storage allocation
  // @return gas price in wei for ArbGas base
  // @return gas price in wei for ArbGas congestion
  // @return gas price in wei for ArbGas total
  function getPricesInWei()
    external
    view
    returns (
      uint256,
      uint256,
      uint256,
      uint256,
      uint256,
      uint256
    );
}
