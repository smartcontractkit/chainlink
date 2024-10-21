// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {Internal} from "../libraries/Internal.sol";

/// @notice This interface contains the only RMN-related functions that might be used on-chain by other CCIP contracts.
interface IRMNRemote {
  /// @notice signature components from RMN nodes
  struct Signature {
    bytes32 r;
    bytes32 s;
  }

  /// @notice Verifies signatures of RMN nodes, on dest lane updates as provided in the CommitReport
  /// @param offRampAddress is not inferred by msg.sender, in case the call is made through ARMProxy
  /// @param merkleRoots must be well formed, and is a representation of the CommitReport received from the oracles
  /// @param signatures rmnNodes ECDSA sigs, only r & s, must be sorted in ascending order by signer address
  /// @param rawVs rmnNodes ECDSA sigs, part v bitmap
  /// @dev Will revert if verification fails
  function verify(
    address offRampAddress,
    Internal.MerkleRoot[] memory merkleRoots,
    Signature[] memory signatures,
    uint256 rawVs
  ) external view;

  /// @notice gets the current set of cursed subjects
  /// @return subjects the list of cursed subjects
  function getCursedSubjects() external view returns (bytes16[] memory subjects);

  /// @notice If there is an active global or legacy curse, this function returns true.
  /// @return bool true if there is an active global curse
  function isCursed() external view returns (bool);

  /// @notice If there is an active global curse, or an active curse for `subject`, this function returns true.
  /// @param subject To check whether a particular chain is cursed, set to bytes16(uint128(chainSelector)).
  /// @return bool true if the provided subject is cured *or* if there is an active global curse
  function isCursed(
    bytes16 subject
  ) external view returns (bool);
}
