// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

// Future versions of VRFCoordinatorV2Plus must implement IVRFCoordinatorV2PlusMigration
// to support migrations from previous versions
interface IVRFCoordinatorV2PlusMigration {
  /**
   * @notice called by older versions of coordinator for migration.
   * @notice only callable by older versions of coordinator
   * @notice supports transfer of native currency
   * @param encodedData - user data from older version of coordinator
   */
  function onMigration(bytes calldata encodedData) external payable;
}
