// SPDX-License-Identifier: MIT

pragma solidity 0.8.16;

import "../../../interfaces/automation/UpkeepTranscoderInterface.sol";
import "../../../interfaces/TypeAndVersionInterface.sol";
import {KeeperRegistryBase2_1} from "./KeeperRegistryBase2_1.sol";
import "../../../automation/UpkeepFormat.sol";

/**
 * @dev structs copied directly from source (can't import without changing the contract version)
 */
struct UpkeepV12 {
  uint96 balance;
  address lastKeeper;
  uint32 executeGas;
  uint64 maxValidBlocknumber;
  address target;
  uint96 amountSpent;
  address admin;
}

struct UpkeepV13 {
  uint96 balance;
  address lastKeeper;
  uint96 amountSpent;
  address admin;
  uint32 executeGas;
  uint32 maxValidBlocknumber;
  address target;
  bool paused;
}

struct UpkeepV20 {
  uint32 executeGas;
  uint32 maxValidBlocknumber;
  bool paused;
  address target;
  uint96 amountSpent;
  uint96 balance;
  uint32 lastPerformBlockNumber;
}

/**
 * @notice UpkeepTranscoder allows converting upkeep data from previous keeper registry versions 1.2, 1.3, and
 * 2.0 to registry 2.1
 */
contract UpkeepTranscoder4_0 is UpkeepTranscoderInterface, TypeAndVersionInterface {
  error InvalidTranscoding();

  /**
   * @notice versions:
   * - UpkeepTranscoder 4.0.0: adds support for registry 2.1
   * - UpkeepTranscoder 3.0.0: works with registry 2.0; adds temporary workaround for UpkeepFormat enum bug
   */
  string public constant override typeAndVersion = "UpkeepTranscoder 4.0.0";
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
      (uint256[] memory ids, UpkeepV12[] memory upkeepsV1, bytes[] memory checkDatas) = abi.decode(
        encodedUpkeeps,
        (uint256[], UpkeepV12[], bytes[])
      );

      if (ids.length != upkeepsV1.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }

      address[] memory admins = new address[](ids.length);
      UpkeepV20[] memory newUpkeeps = new UpkeepV20[](ids.length);
      UpkeepV12 memory upkeepV1;
      for (uint256 idx = 0; idx < ids.length; idx++) {
        upkeepV1 = upkeepsV1[idx];
        newUpkeeps[idx] = UpkeepV20({
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
      (uint256[] memory ids, UpkeepV13[] memory upkeepsV2, bytes[] memory checkDatas) = abi.decode(
        encodedUpkeeps,
        (uint256[], UpkeepV13[], bytes[])
      );

      if (ids.length != upkeepsV2.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }

      address[] memory admins = new address[](ids.length);
      UpkeepV20[] memory newUpkeeps = new UpkeepV20[](ids.length);
      UpkeepV13 memory upkeepV2;
      for (uint256 idx = 0; idx < ids.length; idx++) {
        upkeepV2 = upkeepsV2[idx];
        newUpkeeps[idx] = UpkeepV20({
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
