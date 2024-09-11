// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../interfaces/IRMN.sol";
import {OwnerIsCreator} from "./../../../shared/access/OwnerIsCreator.sol";

// Inlined from RMN 1.0 contract.
// solhint-disable gas-struct-packing
contract OldRMN {
  struct Voter {
    address blessVoteAddr;
    address curseVoteAddr;
    address curseUnvoteAddr;
    uint8 blessWeight;
    uint8 curseWeight;
  }

  struct Config {
    Voter[] voters;
    uint16 blessWeightThreshold;
    uint16 curseWeightThreshold;
  }

  struct VersionedConfig {
    Config config;
    uint32 configVersion;
    uint32 blockNumber;
  }

  struct UnvoteToCurseRecord {
    address curseVoteAddr;
    bytes32 cursesHash;
    bool forceUnvote;
  }
}

/// @dev Retained almost as-is from commit 88f285b94c23d0c684d337064758a5edde380fe2 for compatibility with offchain
/// tests and scripts. Internal structs of the RMN 1.0 contract that were depended on have been inlined.
/// @dev This contract should no longer be used for any new tests or scripts.
/// @notice WARNING: This contract is to be only used for testing, all methods are unprotected.
// TODO: remove this contract when tests and scripts are updated
contract MockRMN is IRMN, OwnerIsCreator {
  error CustomError(bytes err);

  bool private s_curse;
  bytes private s_err;
  OldRMN.VersionedConfig private s_versionedConfig;
  mapping(bytes16 subject => bool cursed) private s_curseBySubject;

  function isCursed() external view override returns (bool) {
    if (s_err.length != 0) {
      revert CustomError(s_err);
    }
    return s_curse;
  }

  function isCursed(bytes16 subject) external view override returns (bool) {
    if (s_err.length != 0) {
      revert CustomError(s_err);
    }
    return s_curse || s_curseBySubject[subject];
  }

  function voteToCurse(bytes32) external {
    s_curse = true;
  }

  function voteToCurse(bytes32, bytes16 subject) external {
    s_curseBySubject[subject] = true;
  }

  function ownerUnvoteToCurse(OldRMN.UnvoteToCurseRecord[] memory) external {
    s_curse = false;
  }

  function ownerUnvoteToCurse(OldRMN.UnvoteToCurseRecord[] memory, bytes16 subject) external {
    s_curseBySubject[subject] = false;
  }

  function setRevert(bytes memory err) external {
    s_err = err;
  }

  function isBlessed(IRMN.TaggedRoot calldata) external view override returns (bool) {
    return !s_curse;
  }

  function getConfigDetails() external view returns (uint32 version, uint32 blockNumber, OldRMN.Config memory config) {
    return (s_versionedConfig.configVersion, s_versionedConfig.blockNumber, s_versionedConfig.config);
  }
}
