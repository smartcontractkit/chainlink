// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

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
  /// @param donId The DON instance ID. These are stored in the CapabilityRegistry.
  /// @return configuration DON's configuration for the capability.
  function getCapabilityConfiguration(uint256 donId) external view returns (bytes memory configuration);

  // Solidity does not support generic returns types, so this cannot be part of
  // the interface. However, the implementation contract MAY implement this
  // function to enable configuration decoding on-chain.
  // function decodeCapabilityConfiguration(bytes configuration) external returns (TypedCapabilityConfigStruct config)
}
