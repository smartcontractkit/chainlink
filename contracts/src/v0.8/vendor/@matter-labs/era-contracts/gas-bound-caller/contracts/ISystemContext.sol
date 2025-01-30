// SPDX-License-Identifier: MIT

pragma solidity ^0.8.19;

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice Contract that stores some of the context variables, that may be either
 * block-scoped, tx-scoped or system-wide.
 */
interface ISystemContext {
  struct BlockInfo {
    uint128 timestamp;
    uint128 number;
  }

  /// @notice A structure representing the timeline for the upgrade from the batch numbers to the L2 block numbers.
  /// @dev It will be used for the L1 batch -> L2 block migration in Q3 2023 only.
  struct VirtualBlockUpgradeInfo {
    /// @notice In order to maintain consistent results for `blockhash` requests, we'll
    /// have to remember the number of the batch when the upgrade to the virtual blocks has been done.
    /// The hashes for virtual blocks before the upgrade are identical to the hashes of the corresponding batches.
    uint128 virtualBlockStartBatch;
    /// @notice L2 block when the virtual blocks have caught up with the L2 blocks. Starting from this block,
    /// all the information returned to users for block.timestamp/number, etc should be the information about the L2 blocks and
    /// not virtual blocks.
    uint128 virtualBlockFinishL2Block;
  }

  function chainId() external view returns (uint256);

  function origin() external view returns (address);

  function gasPrice() external view returns (uint256);

  function blockGasLimit() external view returns (uint256);

  function coinbase() external view returns (address);

  function difficulty() external view returns (uint256);

  function baseFee() external view returns (uint256);

  function txNumberInBlock() external view returns (uint16);

  function getBlockHashEVM(uint256 _block) external view returns (bytes32);

  function getBatchHash(uint256 _batchNumber) external view returns (bytes32 hash);

  function getBlockNumber() external view returns (uint128);

  function getBlockTimestamp() external view returns (uint128);

  function getBatchNumberAndTimestamp() external view returns (uint128 blockNumber, uint128 blockTimestamp);

  function getL2BlockNumberAndTimestamp() external view returns (uint128 blockNumber, uint128 blockTimestamp);

  function gasPerPubdataByte() external view returns (uint256 gasPerPubdataByte);

  function getCurrentPubdataSpent() external view returns (uint256 currentPubdataSpent);
}
