// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/v1.7.0/packages/contracts-bedrock/src/L2/L2ToL1MessagePasser.sol
pragma solidity ^0.8.0;

interface IOptimismL2ToL1MessagePasser {
  /// @notice Emitted any time a withdrawal is initiated.
  /// @param nonce          Unique value corresponding to each withdrawal.
  /// @param sender         The L2 account address which initiated the withdrawal.
  /// @param target         The L1 account address the call will be send to.
  /// @param value          The ETH value submitted for withdrawal, to be forwarded to the target.
  /// @param gasLimit       The minimum amount of gas that must be provided when withdrawing.
  /// @param data           The data to be forwarded to the target on L1.
  /// @param withdrawalHash The hash of the withdrawal.
  event MessagePassed(
    uint256 indexed nonce,
    address indexed sender,
    address indexed target,
    uint256 value,
    uint256 gasLimit,
    bytes data,
    bytes32 withdrawalHash
  );
}
