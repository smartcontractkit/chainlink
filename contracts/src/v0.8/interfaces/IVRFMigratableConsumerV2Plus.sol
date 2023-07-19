// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice The IVRFMigratableConsumerV2Plus interface defines the
/// @notice method required to be implemented by all V2Plus consumers.
/// @dev This interface is designed to be used in VRFConsumerBaseV2Plus.
interface IVRFMigratableConsumerV2Plus {
  /// @notice Set the VRF Coordinator address for the consumer.
  /// @notice This method is should only be callable by the subscription admin.
  function setVRFCoordinator(address vrfCoordinator) external;
}
