// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.19;

import {UpkeepTranscoderInterfaceV2} from "../interfaces/UpkeepTranscoderInterfaceV2.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

enum RegistryVersion {
  V12,
  V13,
  V20,
  V21,
  V23
}

/**
 * @notice UpkeepTranscoder is a contract that allows converting upkeep data from previous registry versions to newer versions
 * @dev it currently only supports 2.3 -> 2.3 migrations
 */
contract UpkeepTranscoder5_0 is UpkeepTranscoderInterfaceV2, ITypeAndVersion {
  error InvalidTranscoding();

  string public constant override typeAndVersion = "UpkeepTranscoder 5.0.0";

  /**
   * @notice transcodeUpkeeps transforms upkeep data from the format expected by
   * one registry to the format expected by another. It future-proofs migrations
   * by allowing automation team to customize migration paths and set sensible defaults
   * when new fields are added
   * @param fromVersion version the upkeep is migrating from
   * @param toVersion version the upkeep is migrating to
   * @param encodedUpkeeps encoded upkeep data
   * @dev this transcoder should ONLY be use for V23->V23 migrations for now
   */
  function transcodeUpkeeps(
    uint8 fromVersion,
    uint8 toVersion,
    bytes calldata encodedUpkeeps
  ) external view override returns (bytes memory) {
    if (toVersion == uint8(RegistryVersion.V23) && fromVersion == uint8(RegistryVersion.V23)) {
      return encodedUpkeeps;
    }

    revert InvalidTranscoding();
  }
}
