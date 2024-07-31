// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ICommitStore {
  /// @notice Returns timestamp of when root was accepted or 0 if verification fails.
  /// @dev This method uses a merkle tree within a merkle tree, with the hashedLeaves,
  /// proofs and proofFlagBits being used to get the root of the inner tree.
  /// This root is then used as the singular leaf of the outer tree.
  function verify(
    bytes32[] calldata hashedLeaves,
    bytes32[] calldata proofs,
    uint256 proofFlagBits
  ) external view returns (uint256 timestamp);

  /// @notice Returns the expected next sequence number
  function getExpectedNextSequenceNumber() external view returns (uint64 sequenceNumber);
}
