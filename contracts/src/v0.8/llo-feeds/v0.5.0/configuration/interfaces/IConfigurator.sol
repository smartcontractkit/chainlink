// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

interface IConfigurator {
  /// @notice This event is emitted whenever a new production configuration is set for a feed. It triggers a new run of the offchain reporting protocol.
  event ProductionConfigSet(
    bytes32 indexed configId,
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    bytes[] signers,
    bytes32[] offchainTransmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig,
    bool isGreenProduction
  );

  /// @notice This event is emitted whenever a new staging configuration is set for a feed. It triggers a new run of the offchain reporting protocol.
  event StagingConfigSet(
    bytes32 indexed configId,
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    bytes[] signers,
    bytes32[] offchainTransmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig,
    bool isGreenProduction
  );

  event PromoteStagingConfig(bytes32 indexed configId, bytes32 indexed retiredConfigDigest, bool isGreenProduction);

  /// @notice Promotes the staging configuration to production
  // currentState must match the current state for the given configId (prevents
  // accidentally double-flipping if same transaction is sent twice)
  function promoteStagingConfig(bytes32 configId, bool currentState) external;

  function setProductionConfig(
    bytes32 configId,
    bytes[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external;

  function setStagingConfig(
    bytes32 configId,
    bytes[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external;
}
