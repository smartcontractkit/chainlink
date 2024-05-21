// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IL1MessageQueue} from "@scroll-tech/contracts/L1/rollup/IL1MessageQueue.sol";

contract MockScrollL1MessageQueue is IL1MessageQueue {
  /// @notice The start index of all pending inclusion messages.
  function pendingQueueIndex() external pure returns (uint256) {
    return 0;
  }

  /// @notice Return the index of next appended message.
  function nextCrossDomainMessageIndex() external pure returns (uint256) {
    return 0;
  }

  /// @notice Return the message of in `queueIndex`.
  function getCrossDomainMessage(uint256 /* queueIndex */) external pure returns (bytes32) {
    return "";
  }

  /// @notice Return the amount of ETH should pay for cross domain message.
  function estimateCrossDomainMessageFee(uint256 /* gasLimit */) external pure returns (uint256) {
    return 0;
  }

  /// @notice Return the amount of intrinsic gas fee should pay for cross domain message.
  function calculateIntrinsicGasFee(bytes memory /* _calldata */) external pure returns (uint256) {
    return 0;
  }

  /// @notice Return the hash of a L1 message.
  function computeTransactionHash(
    address /* sender */,
    uint256 /* queueIndex */,
    uint256 /* value */,
    address /* target */,
    uint256 /* gasLimit */,
    bytes calldata /* data */
  ) external pure returns (bytes32) {
    return 0;
  }

  /// @notice Append a L1 to L2 message into this contract.
  /// @param target The address of target contract to call in L2.
  /// @param gasLimit The maximum gas should be used for relay this message in L2.
  /// @param data The calldata passed to target contract.
  function appendCrossDomainMessage(address target, uint256 gasLimit, bytes calldata data) external {}

  /// @notice Append an enforced transaction to this contract.
  /// @dev The address of sender should be an EOA.
  /// @param sender The address of sender who will initiate this transaction in L2.
  /// @param target The address of target contract to call in L2.
  /// @param value The value passed
  /// @param gasLimit The maximum gas should be used for this transaction in L2.
  /// @param data The calldata passed to target contract.
  function appendEnforcedTransaction(
    address sender,
    address target,
    uint256 value,
    uint256 gasLimit,
    bytes calldata data
  ) external {}

  /// @notice Pop finalized messages from queue.
  ///
  /// @dev We can pop at most 256 messages each time. And if the message is not skipped,
  ///      the corresponding entry will be cleared.
  ///
  /// @param startIndex The start index to pop.
  /// @param count The number of messages to pop.
  /// @param skippedBitmap A bitmap indicates whether a message is skipped.
  function popCrossDomainMessage(uint256 startIndex, uint256 count, uint256 skippedBitmap) external {}

  /// @notice Drop a skipped message from the queue.
  function dropCrossDomainMessage(uint256 index) external {}
}
