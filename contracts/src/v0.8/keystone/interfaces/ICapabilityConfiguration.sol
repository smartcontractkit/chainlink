// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice Interface for capability configuration contract. It MUST be
/// implemented for a contract to be used as a capability configuration.
/// The contract MAY store configuration that is shared across multiple
/// DON instances and capability versions.
/// @dev This interface does not guarantee the configuration contract's
/// correctness. It is the responsibility of the contract owner to ensure
/// that the configuration contract emits the CapabilityConfigurationSet
/// event when the configuration is set.
interface ICapabilityConfiguration {
  /// @notice Emitted when a capability configuration is set.
  event CapabilityConfigurationSet();

  /// @notice Returns the capability configuration for a particular DON instance.
  /// @dev donId is required to get DON-specific configuration. It avoids a
  /// situation where configuration size grows too large.
  /// @param donId The DON instance ID. These are stored in the CapabilitiesRegistry.
  /// @return configuration DON's configuration for the capability.
  function getCapabilityConfiguration(uint32 donId) external view returns (bytes memory configuration);

  /// @notice Called by the registry prior to the config being set for a particular DON.
  /// @param nodes The nodes that the configuration is being set for.
  /// @param donCapabilityConfig The configuration being set on the capability registry.
  /// @param donCapabilityConfigCount The number of times the DON has been configured, tracked on the capability registry.
  /// @param donId The DON ID on the capability registry.
  function beforeCapabilityConfigSet(
    bytes32[] calldata nodes,
    bytes calldata donCapabilityConfig,
    uint64 donCapabilityConfigCount,
    uint32 donId
  ) external;
}
