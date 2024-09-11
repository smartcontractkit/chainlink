// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

interface IConfigurator {
  /// @notice This event is emitted whenever a new configuration is set for a feed. It triggers a new run of the offchain reporting protocol.
  event ConfigSet(
    bytes32 indexed configId,
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    address[] signers,
    bytes32[] offchainTransmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig
  );

  function setConfig(
    bytes32 configId,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external;
}
