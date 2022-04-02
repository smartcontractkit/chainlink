// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./interfaces/UpkeepTranscoderInterface.sol";
import "../ConfirmedOwner.sol";

/**
 * @notice Transcoder for converting upkeep data from one keeper
 * registry version to another
 */
contract UpkeepTranscoder is UpkeepTranscoderInterface, ConfirmedOwner {
  constructor() ConfirmedOwner(msg.sender) {}

  /**
   * @notice transcodeUpkeeps transforms upkeep data from the format expected by
   * one registry to the format expected by another. It future-proofs migrations
   * by allowing keepers team to customize migration paths and set sensible defaults
   * when new fields are added
   * @param fromVersion struct version the upkeep is migrating from
   * @param toVersion struct version the upkeep is migrating to
   * @param encodedUpkeeps encoded upkeep data
   * @dev this contract & function are simple now, but should evolve as new registries
   * and migration paths are added
   */
  function transcodeUpkeeps(
    UpkeepTranscoderVersion fromVersion,
    UpkeepTranscoderVersion toVersion,
    bytes calldata encodedUpkeeps
  ) external view override returns (bytes memory) {
    require(fromVersion == UpkeepTranscoderVersion.V1 && toVersion == UpkeepTranscoderVersion.V1);
    return encodedUpkeeps;
  }
}
