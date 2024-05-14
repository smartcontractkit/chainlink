// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IMultiCommitStore {
  /// @notice Source chain specific config
  struct SourceChainConfig {
    bool isEnabled; // ────╮ Whether the source chain is enabled
    uint64 minSeqNr; //    | The min sequence number expected for future messages
    address onRamp; // ────╯ The onRamp address on the source chain
  }

  /// @notice Returns timestamp of when root was accepted or 0 if verification fails.
  /// @dev This method uses a merkle tree within a merkle tree, with the hashedLeaves,
  /// proofs and proofFlagBits being used to get the root of the inner tree.
  /// This root is then used as the singular leaf of the outer tree.
  function verify(
    uint64 sourceChainSelector,
    bytes32[] calldata hashedLeaves,
    bytes32[] calldata proofs,
    uint256 proofFlagBits
  ) external view returns (uint256 timestamp);

  /// @notice Returns the source chain config for a given source chain selector.
  function getSourceChainConfig(uint64 sourceChainSelector) external view returns (SourceChainConfig memory);
}
