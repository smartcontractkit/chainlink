// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice This interface is implemented by all VRF V2+ coordinators that can
/// @notice migrate subscription data to new coordinators.
interface IVRFV2PlusMigrate {
  /**
   * @notice migrate the provided subscription ID to the provided VRF coordinator
   * @notice msg.sender must be the subscription owner and newCoordinator must
   * @notice implement IVRFCoordinatorV2PlusMigration.
   * @param subId the subscription ID to migrate
   * @param newCoordinator the vrf coordinator to migrate to
   */
  function migrate(uint256 subId, address newCoordinator) external;
}
