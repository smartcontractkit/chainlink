// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

interface IChainModule {
  // retrieve the native block number of a chain. e.g. L2 block number on Arbitrum
  function blockNumber() external view returns (uint256);

  // retrieve the native block hash of a chain.
  function blockHash(uint256) external view returns (bytes32);

  // retrieve the L1 data fee for a L2 transaction. it should return 0 for L1 chains and
  // L2 chains which don't have L1 fee component. it uses msg.data to estimate L1 data so
  // it must be used with a transaction. Return value in wei.
  function getCurrentL1Fee() external view returns (uint256);

  // retrieve the L1 data fee for a L2 simulation. it should return 0 for L1 chains and
  // L2 chains which don't have L1 fee component. Return value in wei.
  function getMaxL1Fee(uint256 dataSize) external view returns (uint256);

  // Returns an upper bound on execution gas cost for one invocation of blockNumber(),
  // one invocation of blockHash() and one invocation of getCurrentL1Fee().
  // Returns two values, first value indicates a fixed cost and the second value is
  // the cost per msg.data byte (As some chain module's getCurrentL1Fee execution cost
  // scales with calldata size)
  function getGasOverhead()
    external
    view
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead);
}
