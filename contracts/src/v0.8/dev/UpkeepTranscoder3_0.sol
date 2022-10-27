// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./../interfaces/UpkeepTranscoderInterface.sol";
import "./../interfaces/TypeAndVersionInterface.sol";
import {Upkeep as UpkeepV1} from "./../KeeperRegistry1_2.sol";
import {Upkeep as UpkeepV2} from "./KeeperRegistryBase.sol";
import {Upkeep as UpkeepV3} from "./keeper2_0/KeeperRegistryBase2_0.sol";
import "./../UpkeepFormat.sol";

/**
 * @notice UpkeepTranscoder 3_0 allows converting upkeep data from previous keeper registry versions 1.2 and 1.3 to
 * registry 2.0
 */
contract UpkeepTranscoder3_0 is UpkeepTranscoderInterface, TypeAndVersionInterface {
  error InvalidTranscoding();

  /**
   * @notice versions:
   * - UpkeepTranscoder 3.0.0: version 3.0.0 works with registry 2.0
   */
  string public constant override typeAndVersion = "UpkeepTranscoder 3.0.0";
  uint32 internal constant UINT32_MAX = type(uint32).max;

  /**
   * @notice transcodeUpkeeps transforms upkeep data from the format expected by
   * one registry to the format expected by another. It future-proofs migrations
   * by allowing keepers team to customize migration paths and set sensible defaults
   * when new fields are added
   * @param fromVersion struct version the upkeep is migrating from
   * @param toVersion struct version the upkeep is migrating to
   * @param encodedUpkeeps encoded upkeep data
   */
  function transcodeUpkeeps(
    UpkeepFormat fromVersion,
    UpkeepFormat toVersion,
    bytes calldata encodedUpkeeps
  ) external view override returns (bytes memory) {
    if (fromVersion == toVersion) {
      return encodedUpkeeps;
    }

    // this transcoder only handles upkeep V1/V2 to V3, all other formats are invalid.
    if (fromVersion == UpkeepFormat.V1 && toVersion == UpkeepFormat.V3) {
      (uint256[] memory ids, UpkeepV1[] memory upkeeps, bytes[] memory checkDatas) = abi.decode(
        encodedUpkeeps,
        (uint256[], UpkeepV1[], bytes[])
      );

      if (ids.length != upkeeps.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }

      address[] memory admins = new address[](ids.length);
      UpkeepV3[] memory newUpkeeps = new UpkeepV3[](ids.length);
      for (uint256 idx = 0; idx < ids.length; idx++) {
        UpkeepV1 memory upkeep = upkeeps[idx];
        uint32 maxValidBlock;
        if (upkeep.maxValidBlocknumber >= UINT32_MAX) {
          maxValidBlock = UINT32_MAX;
        } else {
          maxValidBlock = uint32(upkeep.maxValidBlocknumber);
        }
        newUpkeeps[idx] = UpkeepV3({
          executeGas: upkeep.executeGas,
          maxValidBlocknumber: maxValidBlock,
          paused: false,
          target: upkeep.target,
          amountSpent: upkeep.amountSpent,
          balance: upkeep.balance,
          lastPerformBlockNumber: 0
        });
        admins[idx] = upkeep.admin;
      }
      return abi.encode(ids, newUpkeeps, checkDatas, admins);
    }

    if (fromVersion == UpkeepFormat.V2 && toVersion == UpkeepFormat.V3) {
      (uint256[] memory ids, UpkeepV2[] memory upkeeps, bytes[] memory checkDatas) = abi.decode(
        encodedUpkeeps,
        (uint256[], UpkeepV2[], bytes[])
      );

      if (ids.length != upkeeps.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }

      address[] memory admins = new address[](ids.length);
      UpkeepV3[] memory newUpkeeps = new UpkeepV3[](ids.length);
      for (uint256 idx = 0; idx < ids.length; idx++) {
        UpkeepV2 memory upkeep = upkeeps[idx];
        newUpkeeps[idx] = UpkeepV3({
          executeGas: upkeep.executeGas,
          maxValidBlocknumber: upkeep.maxValidBlocknumber,
          paused: upkeep.paused,
          target: upkeep.target,
          amountSpent: upkeep.amountSpent,
          balance: upkeep.balance,
          lastPerformBlockNumber: 0
        });
        admins[idx] = upkeep.admin;
      }
      return abi.encode(ids, newUpkeeps, checkDatas, admins);
    }

    revert InvalidTranscoding();
  }
}
