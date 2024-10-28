// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice This interface contains the only RMN-related functions that might be used on-chain by other CCIP contracts.
interface IRMN {
  /// @notice A Merkle root tagged with the address of the commit store contract it is destined for.
  struct TaggedRoot {
    address commitStore;
    bytes32 root;
  }

  /// @notice Callers MUST NOT cache the return value as a blessed tagged root could become unblessed.
  function isBlessed(
    TaggedRoot calldata taggedRoot
  ) external view returns (bool);

  /// @notice Iff there is an active global or legacy curse, this function returns true.
  function isCursed() external view returns (bool);

  /// @notice Iff there is an active global curse, or an active curse for `subject`, this function returns true.
  /// @param subject To check whether a particular chain is cursed, set to bytes16(uint128(chainSelector)).
  function isCursed(
    bytes16 subject
  ) external view returns (bool);
}
