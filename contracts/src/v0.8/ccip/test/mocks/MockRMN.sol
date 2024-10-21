// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../interfaces/IRMN.sol";

/// @notice WARNING: This contract is to be only used for testing, all methods are unprotected.
contract MockRMN is IRMN {
  error CustomError(bytes err);

  bytes private s_isCursedRevert;

  bool private s_globalCursed;
  mapping(bytes16 subject => bool cursed) private s_cursedBySubject;
  mapping(address commitStore => mapping(bytes32 root => bool blessed)) private s_blessedByRoot;

  function setTaggedRootBlessed(IRMN.TaggedRoot calldata taggedRoot, bool blessed) external {
    s_blessedByRoot[taggedRoot.commitStore][taggedRoot.root] = blessed;
  }

  function setGlobalCursed(
    bool cursed
  ) external {
    s_globalCursed = cursed;
  }

  function setChainCursed(uint64 chainSelector, bool cursed) external {
    s_cursedBySubject[bytes16(uint128(chainSelector))] = cursed;
  }

  /// @notice Setting a revert error with length of 0 will disable reverts
  /// @dev Useful to test revert handling of ARMProxy
  function setIsCursedRevert(
    bytes calldata revertErr
  ) external {
    s_isCursedRevert = revertErr;
  }

  // IRMN implementation follows

  function isCursed() external view returns (bool) {
    if (s_isCursedRevert.length > 0) {
      revert CustomError(s_isCursedRevert);
    }
    return s_globalCursed;
  }

  function isCursed(
    bytes16 subject
  ) external view returns (bool) {
    if (s_isCursedRevert.length > 0) {
      revert CustomError(s_isCursedRevert);
    }
    return s_globalCursed || s_cursedBySubject[subject];
  }

  function isBlessed(
    IRMN.TaggedRoot calldata taggedRoot
  ) external view returns (bool) {
    return s_blessedByRoot[taggedRoot.commitStore][taggedRoot.root];
  }
}
