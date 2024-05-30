// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

// IGasModule is an interface for the standalone contract that
// calculates the gas fees for the given chain
interface IGasModule {
  function getTxL1GasFees(
    bytes memory data
  ) external view returns (uint256);
}
