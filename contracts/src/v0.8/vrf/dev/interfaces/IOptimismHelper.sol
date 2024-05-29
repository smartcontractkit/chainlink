// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

// IOptimismHelper is a helper contract that provides a function to get the L1 gas fees for a transaction.
interface IOptimismHelper {
  function getTxL1GasFees(
    bytes memory data
  ) external view returns (uint256);
}
