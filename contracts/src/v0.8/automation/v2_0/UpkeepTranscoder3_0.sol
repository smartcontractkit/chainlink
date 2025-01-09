// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../../automation/interfaces/UpkeepTranscoderInterface.sol";
import "../../shared/interfaces/ITypeAndVersion.sol";
import {Upkeep as UpkeepV1} from "../../automation/v1_2/KeeperRegistry1_2.sol";
import {Upkeep as UpkeepV2} from "../../automation/v1_3/KeeperRegistryBase1_3.sol";
import {Upkeep as UpkeepV3} from "../../automation/v2_0/KeeperRegistryBase2_0.sol";
import "../../automation/UpkeepFormat.sol";

/**
 * @notice UpkeepTranscoder 3_0 allows converting upkeep data from previous keeper registry versions 1.2 and 1.3 to
 * registry 2.0
 */
contract UpkeepTranscoder3_0 is UpkeepTranscoderInterface, ITypeAndVersion {
  error InvalidTranscoding();

  /**
   * @notice versions:
   * - UpkeepTranscoder 3.0.0: version 3.0.0 works with registry 2.0; adds temporary workaround for UpkeepFormat enum bug
   */
  string public constant override typeAndVersion = "UpkeepTranscoder 3.0.0";
  uint32 internal constant UINT32_MAX = type(uint32).max;

  /**
   * @notice transcodeUpkeeps transforms upkeep data from the format expected by
   * one registry to the format expected by another. It future-proofs migrations
   * by allowing keepers team to customize migration paths and set sensible defaults
   * when new fields are added
   * @param fromVersion struct version the upkeep is migrating from
   * @param encodedUpkeeps encoded upkeep data
   * @dev this transcoder should ONLY be use for V1/V2 --> V3 migrations
   * @dev this transcoder **ignores** the toVersion param, as it assumes all migrations are
   * for the V3 version. Therefore, it is the responsibility of the deployer of this contract
   * to ensure it is not used in any other migration paths.
   */
  function transcodeUpkeeps(
    UpkeepFormat fromVersion,
    UpkeepFormat,
    bytes calldata encodedUpkeeps
  ) external view override returns (bytes memory) {
    // this transcoder only handles upkeep V1/V2 to V3, all other formats are invalid.
    if (fromVersion == UpkeepFormat.V1) {
      (uint256[] memory ids, UpkeepV1[] memory upkeepsV1, bytes[] memory checkDatas) = abi.decode(
        encodedUpkeeps,
        (uint256[], UpkeepV1[], bytes[])
      );

      if (ids.length != upkeepsV1.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }

      address[] memory admins = new address[](ids.length);
      UpkeepV3[] memory newUpkeeps = new UpkeepV3[](ids.length);
      UpkeepV1 memory upkeepV1;
      for (uint256 idx = 0; idx < ids.length; idx++) {
        upkeepV1 = upkeepsV1[idx];
        newUpkeeps[idx] = UpkeepV3({
          executeGas: upkeepV1.executeGas,
          maxValidBlocknumber: UINT32_MAX, // maxValidBlocknumber is uint64 in V1, hence a new default value is provided
          paused: false, // migrated upkeeps are not paused by default
          target: upkeepV1.target,
          amountSpent: upkeepV1.amountSpent,
          balance: upkeepV1.balance,
          lastPerformBlockNumber: 0
        });
        admins[idx] = upkeepV1.admin;
      }
      return abi.encode(ids, newUpkeeps, checkDatas, admins);
    }

    if (fromVersion == UpkeepFormat.V2) {
      (uint256[] memory ids, UpkeepV2[] memory upkeepsV2, bytes[] memory checkDatas) = abi.decode(
        encodedUpkeeps,
        (uint256[], UpkeepV2[], bytes[])
      );

      if (ids.length != upkeepsV2.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }

      address[] memory admins = new address[](ids.length);
      UpkeepV3[] memory newUpkeeps = new UpkeepV3[](ids.length);
      UpkeepV2 memory upkeepV2;
      for (uint256 idx = 0; idx < ids.length; idx++) {
        upkeepV2 = upkeepsV2[idx];
        newUpkeeps[idx] = UpkeepV3({
          executeGas: upkeepV2.executeGas,
          maxValidBlocknumber: upkeepV2.maxValidBlocknumber,
          paused: upkeepV2.paused,
          target: upkeepV2.target,
          amountSpent: upkeepV2.amountSpent,
          balance: upkeepV2.balance,
          lastPerformBlockNumber: 0
        });
        admins[idx] = upkeepV2.admin;
      }
      return abi.encode(ids, newUpkeeps, checkDatas, admins);
    }

    revert InvalidTranscoding();
  }
}
