// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

interface IChainSpecific {
  // retrieve the native block number of a chain. e.g. L2 block number on Arbitrum
  function _blockNumber() external view returns (uint256);

  // retrieve the native block hash of a chain.
  function _blockHash(uint256) external view returns (bytes32);

  // retrieve the L1 data fee for a L2 transaction. it should return 0 for L1 chains and
  // L2 chains which don't have L1 fee component.
  function _getL1FeeForTransaction(bytes calldata txCallData) external view returns (uint256);

  // retrieve the L1 data fee for a L2 simulation. it should return 0 for L1 chains and
  // L2 chains which don't have L1 fee component.
  function _getL1FeeForSimulation(bytes calldata txCallData) external view returns (uint256);
}
