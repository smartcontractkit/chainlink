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
   * @param fromRegistry registry the upkeep is migrating from
   * @param toRegistry registry the upkeep is migrating to
   * @param encodedUpkeeps encoded upkeep data
   * @dev this contract & function are simple now, but should evolve as new registries
   * and migration paths are added
   */
  function transcodeUpkeeps(
    address fromRegistry,
    address toRegistry,
    bytes calldata encodedUpkeeps
  ) external view override returns (bytes memory) {
    return encodedUpkeeps;
  }
}
