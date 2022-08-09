// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../interfaces/UpkeepTranscoderInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import {Upkeep} from "./interfaces/KeeperRegistryInterfaceDev.sol";

/**
 * @notice Transcoder for converting upkeep data from one keeper
 * registry version to another
 */
contract UpkeepTranscoderDev is UpkeepTranscoderInterface, TypeAndVersionInterface {
  struct UpkeepV1 {
    uint96 balance;
    address lastKeeper; // 1 full evm word
    uint32 executeGas;
    uint64 maxValidBlocknumber;
    address target; // 2 full evm words
    uint96 amountSpent;
    address admin; // 3 full evm words
  }

  error InvalidTranscoding();

  /**
   * @notice versions:
   * - UpkeepTranscoder 1.3.0: placeholder to allow new formats in the future
   */
  string public constant override typeAndVersion = "UpkeepTranscoder 1.3.0";

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
    if (toVersion == fromVersion) {
      return encodedUpkeeps;
    }
    if (fromVersion != UpkeepFormat.V1 || toVersion != UpkeepFormat.V2) {
      revert InvalidTranscoding();
    }

    // the only possible case here is V1 to V2
    (uint256[] memory ids, UpkeepV1[] memory upkeeps, bytes[] memory checkDatas) = abi.decode(
      encodedUpkeeps,
      (uint256[], UpkeepV1[], bytes[])
    );

    UpkeepV1 memory fromUpkeep;
    Upkeep[] memory newUpkeeps = new Upkeep[](upkeeps.length);
    for (uint256 idx = 0; idx < ids.length; idx++) {
      fromUpkeep = upkeeps[idx];
      Upkeep memory upkeep = Upkeep({
        balance: fromUpkeep.balance,
        lastKeeper: fromUpkeep.lastKeeper,
        executeGas: fromUpkeep.executeGas,
        maxValidBlocknumber: fromUpkeep.maxValidBlocknumber,
        target: fromUpkeep.target,
        amountSpent: fromUpkeep.amountSpent,
        admin: fromUpkeep.admin,
        // there is no pause notion for Upkeep V1, so all migrated upkeeps are not paused
        paused: false
      });
      newUpkeeps[idx] = upkeep;
    }

    return abi.encode(ids, newUpkeeps, checkDatas);
  }
}
