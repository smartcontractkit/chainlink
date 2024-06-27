// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import {UpkeepTranscoderInterface} from "./interfaces/UpkeepTranscoderInterface.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {UpkeepFormat} from "./UpkeepFormat.sol";

/**
 * @notice Transcoder for converting upkeep data from one keeper
 * registry version to another
 */
contract UpkeepTranscoder is UpkeepTranscoderInterface, TypeAndVersionInterface {
  error InvalidTranscoding();

  /**
   * @notice versions:
   * - UpkeepTranscoder 1.0.0: placeholder to allow new formats in the future
   */
  string public constant override typeAndVersion = "UpkeepTranscoder 1.0.0";

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
    UpkeepFormat fromVersion,
    UpkeepFormat toVersion,
    bytes calldata encodedUpkeeps
  ) external view override returns (bytes memory) {
    if (fromVersion != toVersion) {
      revert InvalidTranscoding();
    }

    return encodedUpkeeps;
  }
}
