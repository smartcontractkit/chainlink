// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.16;

import {UpkeepTranscoderInterfaceV2} from "../interfaces/UpkeepTranscoderInterfaceV2.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {KeeperRegistryBase2_1 as R21} from "./KeeperRegistryBase2_1.sol";
import {IAutomationForwarder} from "../interfaces/IAutomationForwarder.sol";

enum RegistryVersion {
  V12,
  V13,
  V20,
  V21
}

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
  uint32 lastPerformedBlockNumber;
}

/**
 * @notice UpkeepTranscoder allows converting upkeep data from previous keeper registry versions 1.2, 1.3, and
 * 2.0 to registry 2.1
 */
contract UpkeepTranscoder4_0 is UpkeepTranscoderInterfaceV2, TypeAndVersionInterface {
  error InvalidTranscoding();

  /**
   * @notice versions:
   * - UpkeepTranscoder 4.0.0: adds support for registry 2.1; adds support for offchainConfigs
   * - UpkeepTranscoder 3.0.0: works with registry 2.0; adds temporary workaround for UpkeepFormat enum bug
   */
  string public constant override typeAndVersion = "UpkeepTranscoder 4.0.0";
  uint32 internal constant UINT32_MAX = type(uint32).max;
  IAutomationForwarder internal constant ZERO_FORWARDER = IAutomationForwarder(address(0));

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
    uint8 fromVersion,
    uint8,
    bytes calldata encodedUpkeeps
  ) external view override returns (bytes memory) {
    // v1.2 => v2.1
    if (fromVersion == uint8(RegistryVersion.V12)) {
      (uint256[] memory ids, UpkeepV12[] memory upkeepsV12, bytes[] memory checkDatas) = abi.decode(
        encodedUpkeeps,
        (uint256[], UpkeepV12[], bytes[])
      );
      if (ids.length != upkeepsV12.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }
      address[] memory targets = new address[](ids.length);
      address[] memory admins = new address[](ids.length);
      R21.Upkeep[] memory newUpkeeps = new R21.Upkeep[](ids.length);
      UpkeepV12 memory upkeepV12;
      for (uint256 idx = 0; idx < ids.length; idx++) {
        upkeepV12 = upkeepsV12[idx];
        newUpkeeps[idx] = R21.Upkeep({
          performGas: upkeepV12.executeGas,
          maxValidBlocknumber: UINT32_MAX, // maxValidBlocknumber is uint64 in V1, hence a new default value is provided
          paused: false, // migrated upkeeps are not paused by default
          forwarder: ZERO_FORWARDER,
          amountSpent: upkeepV12.amountSpent,
          balance: upkeepV12.balance,
          lastPerformedBlockNumber: 0
        });
        targets[idx] = upkeepV12.target;
        admins[idx] = upkeepV12.admin;
      }
      return abi.encode(ids, newUpkeeps, targets, admins, checkDatas, new bytes[](ids.length), new bytes[](ids.length));
    }
    // v1.3 => v2.1
    if (fromVersion == uint8(RegistryVersion.V13)) {
      (uint256[] memory ids, UpkeepV13[] memory upkeepsV13, bytes[] memory checkDatas) = abi.decode(
        encodedUpkeeps,
        (uint256[], UpkeepV13[], bytes[])
      );
      if (ids.length != upkeepsV13.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }
      address[] memory targets = new address[](ids.length);
      address[] memory admins = new address[](ids.length);
      R21.Upkeep[] memory newUpkeeps = new R21.Upkeep[](ids.length);
      UpkeepV13 memory upkeepV13;
      for (uint256 idx = 0; idx < ids.length; idx++) {
        upkeepV13 = upkeepsV13[idx];
        newUpkeeps[idx] = R21.Upkeep({
          performGas: upkeepV13.executeGas,
          maxValidBlocknumber: upkeepV13.maxValidBlocknumber,
          paused: upkeepV13.paused,
          forwarder: ZERO_FORWARDER,
          amountSpent: upkeepV13.amountSpent,
          balance: upkeepV13.balance,
          lastPerformedBlockNumber: 0
        });
        targets[idx] = upkeepV13.target;
        admins[idx] = upkeepV13.admin;
      }
      return abi.encode(ids, newUpkeeps, targets, admins, checkDatas, new bytes[](ids.length), new bytes[](ids.length));
    }
    // v2.0 => v2.1
    if (fromVersion == uint8(RegistryVersion.V20)) {
      (uint256[] memory ids, UpkeepV20[] memory upkeepsV20, bytes[] memory checkDatas, address[] memory admins) = abi
        .decode(encodedUpkeeps, (uint256[], UpkeepV20[], bytes[], address[]));
      if (ids.length != upkeepsV20.length || ids.length != checkDatas.length) {
        revert InvalidTranscoding();
      }
      // bit of a hack - transcodeUpkeeps should be a pure function
      R21.Upkeep[] memory newUpkeeps = new R21.Upkeep[](ids.length);
      bytes[] memory emptyBytes = new bytes[](ids.length);
      address[] memory targets = new address[](ids.length);
      UpkeepV20 memory upkeepV20;
      for (uint256 idx = 0; idx < ids.length; idx++) {
        upkeepV20 = upkeepsV20[idx];
        newUpkeeps[idx] = R21.Upkeep({
          performGas: upkeepV20.executeGas,
          maxValidBlocknumber: upkeepV20.maxValidBlocknumber,
          paused: upkeepV20.paused,
          forwarder: ZERO_FORWARDER,
          amountSpent: upkeepV20.amountSpent,
          balance: upkeepV20.balance,
          lastPerformedBlockNumber: 0
        });
        targets[idx] = upkeepV20.target;
      }
      return abi.encode(ids, newUpkeeps, targets, admins, checkDatas, emptyBytes, emptyBytes);
    }
    // v2.1 => v2.1
    if (fromVersion == uint8(RegistryVersion.V21)) {
      return encodedUpkeeps;
    }

    revert InvalidTranscoding();
  }
}
